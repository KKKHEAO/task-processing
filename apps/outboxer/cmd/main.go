package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/KKKHEAO/task-processing/apps/outboxer/internal/outbox"
	"github.com/KKKHEAO/task-processing/packages/config"
	"github.com/KKKHEAO/task-processing/packages/logger"
	"github.com/KKKHEAO/task-processing/packages/postgres"
	"github.com/KKKHEAO/task-processing/packages/repository"
	"go.uber.org/zap"
)

// App представляет основное приложение outbox worker
type App struct {
	config    *config.Config
	db        *sql.DB
	publisher *outbox.Publisher
	worker    *outbox.Worker
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	log       *zap.Logger
}

// NewApp создает новое приложение
func NewApp(cfg *config.Config, log *zap.Logger) (*App, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Форматируем retry topics для красивого вывода
	retryTopicsStr := ""
	for i, rt := range cfg.Kafka.RetryTopics {
		if i > 0 {
			retryTopicsStr += ", "
		}
		retryTopicsStr += fmt.Sprintf("%s (delay: %v, max: %d)", rt.Name, rt.Delay, rt.MaxRetry)
	}

	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	taskRepo := repository.NewPostgresRepo(psqlDB)
	publisher := outbox.NewPublisher(&cfg.Kafka, log)
	worker := outbox.NewWorker(taskRepo, publisher, &cfg.Kafka, log)

	return &App{
		config:    cfg,
		db:        psqlDB,
		publisher: publisher,
		worker:    worker,
		log:       log,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	// Запускаем основной worker
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.worker.Start(ctx)
		a.log.Info("Worker stopped")
	}()

	a.log.Info("Outbox worker started successfully")
	a.log.Info("Press Ctrl+C to stop")

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		a.log.Info("Received signal", zap.String("signal", sig.String()))
		return a.Shutdown(10 * time.Second)
	case <-ctx.Done():
		a.log.Info("Context cancelled")
		return nil
	}
}

// Shutdown gracefully останавливает приложение
func (a *App) Shutdown(timeout time.Duration) error {
	a.log.Info("Shutting down gracefully...")

	// Отменяем контекст
	if a.cancel != nil {
		a.cancel()
	}

	// Ожидаем завершения всех worker'ов с таймаутом
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		a.log.Info("All workers stopped")
	case <-time.After(timeout):
		a.log.Info("Timeout waiting for workers to stop")
	}

	// Закрываем соединения
	var errs []error

	if a.publisher != nil {
		if err := a.publisher.Close(); err != nil {
			errs = append(errs, fmt.Errorf("publisher close: %w", err))
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("database close: %w", err))
		}
	}

	return errors.Join(errs...)
}

func main() {
	cfg := config.NewConfig()
	log, _ := logger.NewLogger(cfg)
	app, err := NewApp(cfg, log)
	if err != nil {
		log.Fatal("Error new app: ", zap.Error(err))
	}

	if err := app.Run(); err != nil {
		log.Fatal("Error start app: ", zap.Error(err))
	}

	log.Info("Application stopped")
}

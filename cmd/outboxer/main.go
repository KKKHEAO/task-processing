package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"task-processing/config"
	"task-processing/internal/api/outbox"
	"task-processing/internal/api/repository"
	"task-processing/pkg/postgres"
)

// App представляет основное приложение outbox worker
type App struct {
	config    *config.Config
	db        *sql.DB
	publisher *outbox.Publisher
	worker    *outbox.Worker
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewApp создает новое приложение
func NewApp() (*App, error) {
	cfg := config.NewConfig()
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
	publisher := outbox.NewPublisher(&cfg.Kafka)
	worker := outbox.NewWorker(taskRepo, publisher, &cfg.Kafka)

	return &App{
		config:    cfg,
		db:        psqlDB,
		publisher: publisher,
		worker:    worker,
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
		log.Println("Worker stopped")
	}()

	log.Println("Outbox worker started successfully")
	log.Println("Press Ctrl+C to stop")

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("Received signal: %s", sig)
		return a.Shutdown(10 * time.Second)
	case <-ctx.Done():
		log.Println("Context cancelled")
		return nil
	}
}

// Shutdown gracefully останавливает приложение
func (a *App) Shutdown(timeout time.Duration) error {
	log.Println("Shutting down gracefully...")

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
		log.Println("All workers stopped")
	case <-time.After(timeout):
		log.Println("Timeout waiting for workers to stop")
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
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Application stopped")
}

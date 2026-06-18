package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KKKHEAO/task-processing/apps/api/internal/service"
	"github.com/KKKHEAO/task-processing/apps/api/internal/transport/grpc"
	"github.com/KKKHEAO/task-processing/packages/config"
	"github.com/KKKHEAO/task-processing/packages/logger"
	"github.com/KKKHEAO/task-processing/packages/postgres"
	"github.com/KKKHEAO/task-processing/packages/repository"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("invalid config: %v", err))
	}

	log, err := logger.NewLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot init logger: %v", err))
	}
	defer log.Sync()

	if err := run(cfg, log); err != nil {
		log.Error("server stopped with error", zap.Error(err))
		os.Exit(1)
	}
	log.Info("Server stopped")
}

func run(cfg *config.Config, log *zap.Logger) error {
	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}
	defer psqlDB.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	taskRepo := repository.NewPostgresRepo(psqlDB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := grpc.NewTaskHandler(taskService)

	errChan := make(chan error, 1)
	go func() {
		errChan <- grpc.RunServer(ctx, taskHandler, cfg.Server.GRPCPort, log)
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutting down server gracefully...")
		stop()

		select {
		case err := <-errChan:
			if err != nil && !errors.Is(err, context.Canceled) {
				return fmt.Errorf("server error during shutdown: %w", err)
			}
		case <-time.After(cfg.Server.ShutdownTimeout):
			return errors.New("shutdown timeout")
		}
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	}

	return nil
}

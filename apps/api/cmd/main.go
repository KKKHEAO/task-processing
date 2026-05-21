package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	log, _ := logger.NewLogger(cfg)
	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		log.Fatal("Ошибка при инициализации postgres", zap.Error(err))
	}
	defer psqlDB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskRepo := repository.NewPostgresRepo(psqlDB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := grpc.NewTaskHandler(taskService)

	errChan := make(chan error, 1)
	go func() {
		errChan <- grpc.RunServer(ctx, taskHandler, "50051")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Info("Shutting down server gracefully...")
		cancel()
		err := <-errChan
		if err != nil && err != context.Canceled {
			log.Fatal("Server error during shutdown: ", zap.Error(err))
		}
		log.Info("Server stopped")
	case err := <-errChan:
		if err != nil {
			log.Fatal("Server error: ", zap.Error(err))
		}
		log.Info("Server stopped")
	}
}

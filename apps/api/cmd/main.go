package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KKKHEAO/task-processing/packages/repository"
	"github.com/KKKHEAO/task-processing/apps/api/internal/service"
	"github.com/KKKHEAO/task-processing/apps/api/internal/transport/grpc"
	"github.com/KKKHEAO/task-processing/packages/config"
	"github.com/KKKHEAO/task-processing/packages/postgres"
)

func main() {
	log.Println("Starting gRPC api server...")
	cfg := config.NewConfig()
	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		log.Fatal(err)
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
		log.Println("Shutting down server gracefully...")
		cancel()
		err := <-errChan
		if err != nil && err != context.Canceled {
			log.Fatalf("Server error during shutdown: %v", err)
		}
		log.Println("Server stopped")
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
		log.Println("Server stopped")
	}
}

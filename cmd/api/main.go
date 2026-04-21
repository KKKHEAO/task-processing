package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-processing/config"
	"task-processing/internal/api/repository"
	"task-processing/internal/api/service"
	"task-processing/internal/api/transport/grpc"
	"task-processing/pkg/postgres"
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

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-processing/config"
	"task-processing/internal/outbox"
	"task-processing/internal/repository"
	"task-processing/internal/service"
	"task-processing/internal/transport/grpc"
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

	go func() {
		if err := grpc.RunServer(taskHandler, "50051"); err != nil {
			log.Fatal(err)
		}
	}()

	publisher := outbox.NewPublisher("kafka:9092")
	worker := outbox.NewWorker(taskRepo, publisher)
	go worker.Start(ctx)

	// черновой вариант для синхронизации
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
}

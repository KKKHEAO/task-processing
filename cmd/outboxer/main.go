package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-processing/config"
	"task-processing/internal/api/outbox"
	"task-processing/internal/api/repository"
	"task-processing/pkg/postgres"
	"time"
)

func main() {
	log.Println("Starting outbox worker...")
	cfg := config.NewConfig()
	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer psqlDB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskRepo := repository.NewPostgresRepo(psqlDB)
	publisher := outbox.NewPublisher("kafka:9092")
	defer publisher.Close()

	worker := outbox.NewWorker(taskRepo, publisher)

	done := make(chan struct{})
	go func() {
		worker.Start(ctx)
		close(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down gracefully...")
		cancel()

		select {
		case <-done:
			log.Println("Worker stopped")
		case <-time.After(5 * time.Second):
			log.Println("Timeout waiting for worker to stop")
		}
	case <-done:
		log.Println("Worker stopped unexpectedly")
		cancel()
	}
}

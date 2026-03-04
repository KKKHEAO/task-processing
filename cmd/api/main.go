package main

import (
	"log"
	"task-processing/config"
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
	taskRepo := repository.NewPostgresRepo(psqlDB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := grpc.NewTaskHandler(taskService)
	if err := grpc.RunServer(taskHandler, "50051"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"log"
	"task-processing/internal/kafka"
	"task-processing/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewConsumer(
		"kafka:9092",
		"tasks.created",
		"task-workers",
	)

	pool := worker.NewPool(100)

	for {

		msg, err := consumer.Read(ctx)
		if err != nil {
			log.Println(err)
			continue
		}

		pool.Submit(worker.Job{
			Payload: msg.Value,
		})
	}
}

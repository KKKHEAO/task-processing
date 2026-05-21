package main

import (
	"context"
	"log"
	"github.com/KKKHEAO/task-processing/packages/kafka"
	"github.com/KKKHEAO/task-processing/apps/worker/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	consumer := kafka.NewConsumer(
		"kafka:9092",
		"tasks.created",
		"task-workers",
	)

	retryProducer := kafka.NewProducer("kafka:9092", "tasks.retry")
	dlqProducer := kafka.NewProducer("kafka:9092", "tasks.dlq")

	pool := worker.NewPool(5, retryProducer, dlqProducer)

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

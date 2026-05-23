package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/KKKHEAO/task-processing/apps/worker/internal/worker"
	"github.com/KKKHEAO/task-processing/packages/config"
	"github.com/KKKHEAO/task-processing/packages/kafka"
	"github.com/KKKHEAO/task-processing/packages/logger"
	sgkafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	log, _ := logger.NewLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	topics := []string{cfg.Kafka.MainTopic}
	for _, rt := range cfg.Kafka.RetryTopics {
		topics = append(topics, rt.Name)
	}

	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers[0],
		topics,
		"task-workers",
	)

	retryProducer := kafka.NewProducer(cfg.Kafka.Brokers[0], cfg.Kafka.RetryTopics[0].Name)
	dlqProducer := kafka.NewProducer(cfg.Kafka.Brokers[0], cfg.Kafka.DLQTopic)

	pool := worker.NewPool(ctx, 5, retryProducer, dlqProducer, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	msgCh := make(chan struct {
		msg sgkafka.Message
		err error
	})

	go func() {
		for {
			msg, err := consumer.Read(ctx)
			msgCh <- struct {
				msg sgkafka.Message
				err error
			}{msg, err}
		}
	}()

	for {
		select {
		case <-quit:
			log.Info("Shutting down...")
			cancel()
			return
		case job := <-msgCh:
			if job.err != nil {
				log.Error("Consumer read error", zap.Error(job.err))
				continue
			}
			pool.Submit(worker.Job{
				Payload: job.msg.Value,
			})
		}
	}
}

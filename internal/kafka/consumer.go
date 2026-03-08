package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker string, topic string, groupId string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: groupId,
		Topic:   topic,
	})

	return &Consumer{
		reader: reader,
	}
}

func (c *Consumer) Read(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

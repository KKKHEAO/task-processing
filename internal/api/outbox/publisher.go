package outbox

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(broker string) *Publisher {

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Balancer: &kafka.LeastBytes{},
	}

	return &Publisher{
		writer: writer,
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, key string, payload []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}

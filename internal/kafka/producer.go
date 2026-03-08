package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker, topic string) *Producer {

	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: topic,
		},
	}
}

func (p *Producer) Send(ctx context.Context, payload []byte) error {

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
}

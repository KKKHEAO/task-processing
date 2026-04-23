package outbox

import (
	"context"
	"time"

	"task-processing/config"
	"task-processing/internal/domain"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
	config *config.KafkaConfig
}

func NewPublisher(cfg *config.KafkaConfig) *Publisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Publisher{
		writer: writer,
		config: cfg,
	}
}

// PublishEvent публикует событие с учетом retry логики
// Возвращает target topic и error
func (p *Publisher) PublishEvent(ctx context.Context, event *domain.OutboxEvent) (string, error) {
	// Определяем целевой топик на основе retry count
	targetTopic := p.selectTargetTopic(event.RetryCount)

	// Создаем сообщение с минимальными заголовками
	msg := kafka.Message{
		Topic: targetTopic,
		Key:   []byte(event.Key),
		Value: event.Payload,
		Headers: []kafka.Header{
			{Key: "X-Event-ID", Value: []byte(event.Id.String())},
			{Key: "X-Original-Topic", Value: []byte(event.Topic)},
		},
	}

	// Устанавливаем timestamp для retry сообщений
	if event.RetryCount > 0 && event.NextRetryAt != nil {
		msg.Time = *event.NextRetryAt
	}

	err := p.writer.WriteMessages(ctx, msg)
	return targetTopic, err
}

// selectTargetTopic выбирает целевой топик на основе номера попытки
func (p *Publisher) selectTargetTopic(retryCount int) string {
	// Если это первая попытка, публикуем в основной топик
	if retryCount == 0 {
		return p.config.MainTopic
	}

	// Если превышено максимальное количество попыток, отправляем в DLQ
	if retryCount >= p.config.MaxRetries {
		return p.config.DLQTopic
	}

	// Выбираем подходящий retry топик на основе счетчика попыток
	for _, retryTopic := range p.config.RetryTopics {
		if retryCount <= retryTopic.MaxRetry {
			return retryTopic.Name
		}
	}

	// Если не нашли подходящий retry топик, отправляем в DLQ
	return p.config.DLQTopic
}

// CalculateNextRetryTime вычисляет время следующей попытки на основе retry count
func (p *Publisher) CalculateNextRetryTime(retryCount int) time.Time {
	// Для первой retry попытки используем минимальную задержку
	if retryCount == 1 {
		if len(p.config.RetryTopics) > 0 {
			return time.Now().Add(p.config.RetryTopics[0].Delay)
		}
		return time.Now().Add(1 * time.Minute) // fallback
	}

	// Ищем подходящий retry топик для текущей попытки
	for _, retryTopic := range p.config.RetryTopics {
		if retryCount <= retryTopic.MaxRetry {
			return time.Now().Add(retryTopic.Delay)
		}
	}

	// Если не нашли, используем максимальную задержку
	if len(p.config.RetryTopics) > 0 {
		lastTopic := p.config.RetryTopics[len(p.config.RetryTopics)-1]
		return time.Now().Add(lastTopic.Delay)
	}

	return time.Now().Add(30 * time.Minute) // fallback
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}

// Publish - упрощенный метод для обратной совместимости
func (p *Publisher) Publish(ctx context.Context, topic string, key string, payload []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}

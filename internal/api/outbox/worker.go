package outbox

import (
	"context"
	"log"
	"time"

	"task-processing/config"
	"task-processing/internal/api/repository"
)

type Worker struct {
	repo      repository.TaskRepository
	publisher *Publisher
	config    *config.KafkaConfig
}

func NewWorker(repo repository.TaskRepository, pub *Publisher, cfg *config.KafkaConfig) *Worker {
	return &Worker{
		repo:      repo,
		publisher: pub,
		config:    cfg,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.config.PollInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) {
	events, err := w.repo.FetchOutboxBatch(ctx, w.config.BatchSize)
	if err != nil {
		log.Println("fetch error:", err)
		return
	}

	log.Printf("Fetched outbox batch: %d events", len(events))

	for _, e := range events {
		targetTopic, err := w.publisher.PublishEvent(ctx, e)
		if err != nil {
			log.Printf("Failed to publish event %s to topic %s: %v", e.Id, targetTopic, err)

			// Обновляем retry информацию в БД
			newRetryCount := e.RetryCount + 1
			now := time.Now()
			nextRetryAt := w.publisher.CalculateNextRetryTime(newRetryCount)
			errorMsg := err.Error()

			// Если превышено максимальное количество попыток, помечаем как обработанное (отправлено в DLQ)
			if newRetryCount >= w.config.MaxRetries {
				log.Printf("Event %s reached max retries (%d), marking as processed", e.Id, w.config.MaxRetries)
				if markErr := w.repo.MarkOutboxProcessed(ctx, e.Id); markErr != nil {
					log.Printf("Failed to mark event %s as processed: %v", e.Id, markErr)
				}
			} else {
				// Обновляем retry информацию для следующей попытки
				if updateErr := w.repo.UpdateOutboxRetry(ctx, e.Id, newRetryCount, &now, &nextRetryAt, &errorMsg); updateErr != nil {
					log.Printf("Failed to update retry info for event %s: %v", e.Id, updateErr)
				} else {
					log.Printf("Updated retry info for event %s: retry %d, next retry at %v",
						e.Id, newRetryCount, nextRetryAt.Format(time.RFC3339))
				}
			}
			continue
		}

		// Если публикация успешна, помечаем как обработанное
		log.Printf("Successfully published event %s to topic %s", e.Id, targetTopic)
		if err := w.repo.MarkOutboxProcessed(ctx, e.Id); err != nil {
			log.Printf("Failed to mark event %s as processed: %v", e.Id, err)
		}
	}
}

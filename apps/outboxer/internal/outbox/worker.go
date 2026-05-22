package outbox

import (
	"context"
	"time"

	"github.com/KKKHEAO/task-processing/packages/config"
	"github.com/KKKHEAO/task-processing/packages/domain"
	"go.uber.org/zap"
)

type Worker struct {
	repo      domain.TaskRepository
	publisher *Publisher
	config    *config.KafkaConfig
	log       *zap.Logger
}

func NewWorker(repo domain.TaskRepository, pub *Publisher, cfg *config.KafkaConfig, log *zap.Logger) *Worker {
	return &Worker{
		repo:      repo,
		publisher: pub,
		config:    cfg,
		log:       log,
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
		w.log.Error("fetch error:", zap.Error(err))
		return
	}

	w.log.Info("Fetched outbox batch len: ", zap.Int("length", len(events)))

	for _, e := range events {
		targetTopic, err := w.publisher.PublishEvent(ctx, e)
		if err != nil {
			w.log.Error("Failed to publish event ", zap.Error(err))

			// Обновляем retry информацию в БД
			newRetryCount := e.RetryCount + 1
			now := time.Now()
			nextRetryAt := w.publisher.CalculateNextRetryTime(newRetryCount)
			errorMsg := err.Error()

			// Если превышено максимальное количество попыток, помечаем как обработанное (отправлено в DLQ)
			if newRetryCount >= w.config.MaxRetries {
				w.log.Info("event reached max retries",
					zap.String("event_id", e.Id.String()),
					zap.Int("max_retries", w.config.MaxRetries),
				)
				if markErr := w.repo.MarkOutboxProcessed(ctx, e.Id); markErr != nil {
					w.log.Error("failed to mark event as processed after max retries",
						zap.String("event_id", e.Id.String()),
						zap.Error(markErr),
					)
				}
			} else {
				// Обновляем retry информацию для следующей попытки
				if updateErr := w.repo.UpdateOutboxRetry(ctx, e.Id, newRetryCount, &now, &nextRetryAt, &errorMsg); updateErr != nil {
					w.log.Error("failed to update retry info",
						zap.String("event_id", e.Id.String()),
						zap.Error(updateErr),
					)
				} else {
					w.log.Info("updated retry info",
						zap.String("event_id", e.Id.String()),
						zap.Int("retry", newRetryCount),
						zap.Time("next_retry_at", nextRetryAt),
					)
				}
			}
			continue
		}

		// Если публикация успешна, помечаем как обработанное
		w.log.Info("successfully published event",
			zap.String("event_id", e.Id.String()),
			zap.String("topic", targetTopic),
		)
		if err := w.repo.MarkOutboxProcessed(ctx, e.Id); err != nil {
			w.log.Error("failed to mark event as processed",
				zap.String("event_id", e.Id.String()),
				zap.Error(err),
			)
		}
	}
}

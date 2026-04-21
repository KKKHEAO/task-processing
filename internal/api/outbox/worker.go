package outbox

import (
	"context"
	"log"
	"task-processing/internal/api/repository"
	"time"
)

type Worker struct {
	repo      repository.TaskRepository
	publisher *Publisher
}

func NewWorker(repo repository.TaskRepository, pub *Publisher) *Worker {
	return &Worker{
		repo:      repo,
		publisher: pub,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)

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

	events, err := w.repo.FetchOutboxBatch(ctx, 100)
	if err != nil {
		log.Println("fetch error:", err)
		return
	}
	log.Println("Fetched outbox batch:", len(events))
	for _, e := range events {
		err := w.publisher.Publish(
			ctx,
			e.Topic,
			e.Key,
			e.Payload,
		)

		if err != nil {
			log.Println("kafka publish error:", err)
			continue
		}

		err = w.repo.MarkOutboxProcessed(ctx, e.Id)
		if err != nil {
			log.Println("mark processed error:", err)
		}
	}
}

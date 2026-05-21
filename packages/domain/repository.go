package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TaskRepository — порт (интерфейс) для хранения задач
type TaskRepository interface {
	Create(ctx context.Context, task *Task, event *OutboxEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	FetchOutboxBatch(ctx context.Context, limit int) ([]*OutboxEvent, error)
	MarkOutboxProcessed(ctx context.Context, id uuid.UUID) error
	UpdateOutboxRetry(ctx context.Context, id uuid.UUID, retryCount int, lastRetryAt *time.Time, nextRetryAt *time.Time, errorMessage *string) error
}

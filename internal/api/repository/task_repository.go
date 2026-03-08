package repository

import (
	"context"
	"task-processing/internal/domain"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task, event *domain.OutboxEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	FetchOutboxBatch(ctx context.Context, limit int) ([]*domain.OutboxEvent, error)
	MarkOutboxProcessed(ctx context.Context, id uuid.UUID) error
}

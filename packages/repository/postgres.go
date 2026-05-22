package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KKKHEAO/task-processing/packages/domain"

	"github.com/google/uuid"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.TaskRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, task *domain.Task, event *domain.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(
		ctx,
		createTaskQuery,
		task.Id,
		task.Type,
		task.Payload,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert task error: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		createOutBoxQuery,
		event.Id,
		event.Topic,
		event.Key,
		event.Payload,
		event.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert outbox event error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	//TODO: добавить спаны для opentelemetry
	t := domain.Task{}
	if err := r.db.QueryRowContext(ctx, getByIdQuery, id).Scan(
		&t.Id,
		&t.Type,
		&t.Payload,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("get task %s error: %w", id, err)
	}
	return &t, nil
}

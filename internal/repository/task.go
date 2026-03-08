package repository

import (
	"context"
	"database/sql"
	"task-processing/internal/domain"

	"github.com/google/uuid"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) TaskRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, task *domain.Task, event *domain.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		createTaskQuery,
		task.Id,
		task.Type,
		task.Payload,
		task.Status,
		task.Retries,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return err
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
		return err
	}
	return tx.Commit()
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	//TODO: добавить спаны для opentelemetry
	t := domain.Task{}
	if err := r.db.QueryRowContext(ctx, getByIdQuery, id).Scan(
		&t.Id,
		&t.Type,
		&t.Payload,
		&t.Status,
		&t.Retries,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

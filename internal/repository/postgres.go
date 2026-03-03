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

func NewPostgresRepo(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, task *domain.Task) error {
	if err := r.db.QueryRowContext(ctx, "", task).Err(); err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	//TODO: добавить спаны для opentelemetry
	t := &domain.Task{}
	if err := r.db.QueryRowContext(ctx, "", id).Scan(t); err != nil {
		return nil, err
	}
	return t, nil
}

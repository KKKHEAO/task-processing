package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/KKKHEAO/task-processing/packages/domain"

	"github.com/google/uuid"
)

// пока все в одном репе, но потом разнести по разным

func (r *postgresRepository) FetchOutboxBatch(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx, fetchOutboxBatch, limit)
	if err != nil {
		return nil, fmt.Errorf("query outbox batch: %w", err)
	}

	defer rows.Close()

	events := make([]*domain.OutboxEvent, 0)

	for rows.Next() {
		var e domain.OutboxEvent
		var lastRetryAt, nextRetryAt *time.Time
		var errorMessage *string

		err := rows.Scan(
			&e.Id,
			&e.Topic,
			&e.Key,
			&e.Payload,
			&e.CreatedAt,
			&e.RetryCount,
			&lastRetryAt,
			&nextRetryAt,
			&errorMessage,
		)

		if err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}

		e.LastRetryAt = lastRetryAt
		e.NextRetryAt = nextRetryAt
		e.ErrorMessage = errorMessage

		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox rows: %w", err)
	}

	return events, nil
}

func (r *postgresRepository) MarkOutboxProcessed(
	ctx context.Context,
	id uuid.UUID,
) error {

	_, err := r.db.ExecContext(ctx, updateOutBoxQuery, id)

	return err
}

func (r *postgresRepository) UpdateOutboxRetry(
	ctx context.Context,
	id uuid.UUID,
	retryCount int,
	lastRetryAt *time.Time,
	nextRetryAt *time.Time,
	errorMessage *string,
) error {

	_, err := r.db.ExecContext(ctx, updateOutboxRetry,
		id,
		retryCount,
		lastRetryAt,
		nextRetryAt,
		errorMessage,
	)

	return err
}

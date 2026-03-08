package repository

import (
	"context"
	"task-processing/internal/domain"

	"github.com/google/uuid"
)

// пока все в одном репе, но потом разнести по разным

func (r *postgresRepository) FetchOutboxBatch(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx, fetchOutboxBatch, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := make([]*domain.OutboxEvent, 0)

	for rows.Next() {

		var e domain.OutboxEvent

		err := rows.Scan(
			&e.Id,
			&e.Topic,
			&e.Key,
			&e.Payload,
			&e.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, &e)
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

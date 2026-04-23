package domain

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	Id           uuid.UUID
	Topic        string
	Key          string
	Payload      []byte
	CreatedAt    time.Time
	RetryCount   int
	LastRetryAt  *time.Time
	NextRetryAt  *time.Time
	ErrorMessage *string
}

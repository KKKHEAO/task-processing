package domain

import "time"

import "github.com/google/uuid"

type OutboxEvent struct {
	Id        uuid.UUID
	Topic     string
	Key       string
	Payload   []byte
	CreatedAt time.Time
}

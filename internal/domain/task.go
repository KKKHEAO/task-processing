package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusDone       Status = "DONE"
	StatusFailed     Status = "FAILED"
)

type Task struct {
	Id        uuid.UUID
	Type      string
	Payload   []byte
	Status    Status
	Retries   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

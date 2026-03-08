package domain

import "time"

type TaskCreatedEvent struct {
	TaskId    string    `json:"task_id"`
	Retry     int       `json:"retry"`
	NextRetry time.Time `json:"next_retry"`
}

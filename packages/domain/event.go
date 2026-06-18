package domain

import "time"

type TaskCreatedEvent struct {
	Id        string    `json:"id"`
	Retry     int       `json:"retry"`
	NextRetry time.Time `json:"next_retry"`
}

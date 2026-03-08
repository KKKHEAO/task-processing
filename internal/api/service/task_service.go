package service

import (
	"context"
	"encoding/json"
	"task-processing/internal/api/repository"
	"task-processing/internal/domain"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	repository repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) *TaskService {
	return &TaskService{
		repository: repository,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, taskType string, payload []byte) (uuid.UUID, error) {
	task := &domain.Task{
		Id:        uuid.New(),
		Type:      taskType,
		Payload:   payload,
		Status:    domain.StatusPending,
		Retries:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	eventPayload, err := json.Marshal(task)
	if err != nil {
		return uuid.Nil, err
	}

	event := &domain.OutboxEvent{
		Id:        uuid.New(),
		Topic:     "tasks.created",
		Key:       task.Id.String(),
		Payload:   eventPayload,
		CreatedAt: time.Now(),
	}

	if err := s.repository.Create(ctx, task, event); err != nil {
		return uuid.Nil, err
	}

	return task.Id, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.repository.GetByID(ctx, id)
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KKKHEAO/task-processing/packages/domain"
	"github.com/google/uuid"
)

type TaskService struct {
	repository domain.TaskRepository
}

func NewTaskService(repository domain.TaskRepository) *TaskService {
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
		return uuid.Nil, fmt.Errorf("create task %s: %w", task.Id, err)
	}

	return task.Id, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task %s: %w", id, err)
	}
	return task, nil
}

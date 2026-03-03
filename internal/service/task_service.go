package service

import (
	"context"
	"task-processing/internal/domain"
	"task-processing/internal/repository"
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

	if err := s.repository.Create(ctx, task); err != nil {
		return uuid.Nil, err
	}

	return task.Id, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.repository.GetByID(ctx, id)
}

package grpc

import (
	"context"
	"task-processing/internal/api/service"
	taskpb "task-processing/proto"

	"github.com/google/uuid"
)

type TaskHandler struct {
	taskpb.UnimplementedTaskServiceServer
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	id, err := h.service.CreateTask(ctx, req.Type, req.Payload)
	if err != nil {
		return nil, err
	}

	return &taskpb.CreateTaskResponse{
		Id: id.String(),
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	task, err := h.service.GetTask(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &taskpb.GetTaskResponse{
		Id:      task.Id.String(),
		Type:    task.Type,
		Status:  string(task.Status),
		Payload: task.Payload,
	}, nil
}

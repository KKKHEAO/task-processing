package grpc

import (
	"context"

	taskpb "task-processing/proto"

	"github.com/KKKHEAO/task-processing/apps/api/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}

	if len(req.Payload) == 0 {
		return nil, status.Error(codes.InvalidArgument, "payload is required")
	}

	if len(req.Payload) > 1000000 {
		return nil, status.Error(codes.InvalidArgument, "payload too large")
	}

	id, err := h.service.CreateTask(ctx, req.Type, req.Payload)
	if err != nil {
		return nil, err
	}

	return &taskpb.CreateTaskResponse{
		Id: id.String(),
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

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

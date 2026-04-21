package grpc

import (
	"context"
	"net"

	taskpb "task-processing/proto"

	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, handler *TaskHandler, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(server, handler)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(lis)
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		server.GracefulStop()
		return nil
	}
}

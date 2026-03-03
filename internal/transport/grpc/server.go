package grpc

import (
	"net"

	taskpb "task-processing/proto"

	"google.golang.org/grpc"
)

func RunServer(handler *TaskHandler, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(server, handler)

	return server.Serve(lis)
}

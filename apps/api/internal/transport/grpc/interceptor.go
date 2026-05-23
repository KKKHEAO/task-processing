package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryLogging(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			log.Error("gRPC call",
				zap.String("method", info.FullMethod),
				zap.Duration("latency", time.Since(start)),
				zap.Error(err),
			)
		} else {
			log.Info("gRPC call",
				zap.String("method", info.FullMethod),
				zap.Duration("latency", time.Since(start)),
			)
		}

		return resp, err
	}
}

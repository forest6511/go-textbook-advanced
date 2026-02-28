package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// loggingInterceptor is a gRPC unary server interceptor for logging
func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	logger.Info("gRPC call", slog.String("method", info.FullMethod))
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error("gRPC error", slog.String("method", info.FullMethod), slog.Any("err", err))
	}
	return resp, err
}

// notFoundError demonstrates gRPC status codes
func notFoundError(resource string, id int64) error {
	return status.Errorf(codes.NotFound, "%s with id %d not found", resource, id)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	logger.Info("gRPC server starting", slog.String("addr", ":50051"))

	// Demonstrate status error creation
	err = notFoundError("user", 42)
	st, _ := status.FromError(err)
	logger.Info("demo error",
		slog.String("code", st.Code().String()),
		slog.String("message", st.Message()),
	)

	_ = srv
	_ = lis
	fmt.Println("gRPC server configured (not starting in demo mode)")
}

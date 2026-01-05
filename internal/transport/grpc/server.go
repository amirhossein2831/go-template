package grpc

import (
	"context"
	configs "event-collector/internal/config"
	"fmt"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewGRPCServer creates and manages the lifecycle of the gRPC server.
func NewGRPCServer(lc fx.Lifecycle, cfg *configs.Config) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	// Use fx.Lifecycle to handle graceful start and stop.
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Printf("Starting gRPC server on port %d ✅", cfg.Server.GRPC.Port)
			go s.Serve(lis)
			return nil
		},
		OnStop: func(context.Context) error {
			log.Println("Gracefully stopping gRPC server... ✅")
			s.GracefulStop()
			return nil
		},
	})

	return s
}

package grpc

import (
	"context"
	"event-collector/internal/firstapp/config"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

// NewGRPCServer creates and manages the lifecycle of the gRPC server.
func NewGRPCServer(lc fx.Lifecycle, cfg *config.Config) *grpc.Server {
	grpcPort := cfg.GetEnv(cfg.Server.GRPC.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	// Use fx.Lifecycle to handle graceful start and stop.
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Printf("Starting gRPC server on port %s ✅", grpcPort)
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

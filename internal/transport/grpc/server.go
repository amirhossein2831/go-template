package grpc

import (
	"context"
	"event-collector/internal/config"
	"event-collector/internal/service"
	providev1 "event-collector/pkg/grpc/provide/v1"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

// server implements the gRPC server interface defined in the .proto file.
type server struct {
	providev1.UnimplementedGreetingServiceServer
	greetingService *service.GreetingService
}

// NewGRPCServer creates and manages the lifecycle of the gRPC server.
func NewGRPCServer(lc fx.Lifecycle, cfg *config.Config, greetService *service.GreetingService) *grpc.Server {
	grpcPort := cfg.GetEnv(cfg.Server.GRPC.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	grpcServer := &server{
		greetingService: greetService,
	}
	// Register the implementation with the gRPC server.
	providev1.RegisterGreetingServiceServer(s, grpcServer)

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

// SayGreeting is the implementation of our gRPC RPC.
func (s *server) SayGreeting(ctx context.Context, req *providev1.SayGreetingRequest) (*providev1.SayGreetingResponse, error) {
	log.Printf("gRPC transport: Received SayGreeting request: Name=%s", req.Name)

	// Call the shared business logic. The transport layer's only job is translation.
	greeting, err := s.greetingService.GenerateGreetingLogic(ctx, req.Name)
	if err != nil {
		// In a real app, convert this to a proper gRPC status error.
		return nil, err
	}

	return &providev1.SayGreetingResponse{GreetingMessage: greeting}, nil
}

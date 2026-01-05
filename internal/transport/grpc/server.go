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

func NewGRPCServer(lc fx.Lifecycle, cfg *configs.Config) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPC.Port))
	if err != nil {
		return nil, fmt.Errorf("grpc listen failed: %w", err)
	}

	s := grpc.NewServer()
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Printf("Starting gRPC server on port %d ✅", cfg.Server.GRPC.Port)
			go func() {
				if err = s.Serve(lis); err != nil {
					log.Printf("gRPC server stopped with error %d ✅", cfg.Server.GRPC.Port)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			log.Println("Gracefully stopping gRPC server... ✅")
			s.GracefulStop()
			return nil
		},
	})

	return s, nil
}

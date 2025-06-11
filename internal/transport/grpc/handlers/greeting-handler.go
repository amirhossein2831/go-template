package handlers

import (
	"context"
	"event-collector/internal/services"
	providev1 "event-collector/pkg/grpc/provide/v1"
	"log"
)

// GreetingHandler implements the gRPC server interface for the GreetingService.
type GreetingHandler struct {
	providev1.UnimplementedGreetingServiceServer
	greetingService *services.GreetingService
}

// NewGreetingHandler is the constructor for our gRPC handler.
func NewGreetingHandler(gs *services.GreetingService) *GreetingHandler {
	return &GreetingHandler{
		greetingService: gs,
	}
}

// SayGreeting is the implementation of our gRPC RPC.
func (h *GreetingHandler) SayGreeting(ctx context.Context, req *providev1.SayGreetingRequest) (*providev1.SayGreetingResponse, error) {
	log.Printf("gRPC transport: Received SayGreeting request: Name=%s", req.Name)

	greeting, err := h.greetingService.GenerateGreetingLogic(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &providev1.SayGreetingResponse{GreetingMessage: greeting}, nil
}

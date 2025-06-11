package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// GreetingService encapsulates the business logic for greetings.
type GreetingService struct {
	db *mongo.Client
}

// NewGreetingService is the constructor for GreetingService. FX will provide the dependencies.
func NewGreetingService(db *mongo.Client) *GreetingService {
	return &GreetingService{db: db}
}

// GenerateGreetingLogic contains the actual logic for creating a greeting.
func (s *GreetingService) GenerateGreetingLogic(ctx context.Context, name string) (string, error) {
	log.Printf("SERVICE: Generating greeting for '%s'", name)
	return fmt.Sprintf("Hello, %s!", name), nil
}

package services

import (
	"context"
	"event-collector/internal/firstapp/logger"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// GreetingService encapsulates the business logic for greetings.
type GreetingService struct {
	db  *mongo.Client
	log logger.Logger
}

// NewGreetingService is the constructor for GreetingService. FX will provide the dependencies.
func NewGreetingService(db *mongo.Client, log logger.Logger) *GreetingService {
	return &GreetingService{db: db, log: log}
}

// GenerateGreetingLogic contains the actual logic for creating a greeting.
func (s *GreetingService) GenerateGreetingLogic(ctx context.Context, name string) (string, error) {
	s.log.Info("GenerateGreetingLogic::::::::", zap.String("name", name))
	return fmt.Sprintf("Hello, %s!", name), nil
}

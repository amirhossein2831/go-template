package mongo

import (
	"context"
	"event-collector/internal/config"
	"fmt"
	"time"

	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
)

// DB holds the client and implements the singleton logic.
type DB struct {
	Client *mongo.Client
}

// NewMongo returns a thread-safe, singleton instance of our DB client.
func NewMongo(cfg *configs.Config, lc fx.Lifecycle) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		return nil, fmt.Errorf("mongo connect failed: %w", err)
	}

	// Verify connection
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err = client.Disconnect(ctx); err != nil {
				log.Printf("Failed to disconnect MongoDB err: %v", err)
				return err
			}
			log.Println("Connection to MongoDB closed. ✅")
			return nil
		},
	})

	log.Println("Successfully connected to MongoDB! ✅")
	return client, nil
}

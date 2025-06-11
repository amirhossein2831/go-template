package database

import (
	"context"
	"event-collector/internal/firstapp/config"
	"go.uber.org/fx"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB holds the client and implements the singleton logic.
type MongoDB struct {
	Client *mongo.Client
}

// NewMongo returns a thread-safe, singleton instance of our MongoDB client.
func NewMongo(cfg *config.Config, lc fx.Lifecycle) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetEnv(cfg.Database.URI)))
	if err != nil {
		log.Fatalf("FATAL: Failed to create MongoDB client: %v", err)
		return nil, err
	}

	// lc.Append is a hook that registers a function to be executed on shutdown.
	// This is the idiomatic way to handle cleanup in FX.
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			dbError := client.Disconnect(ctx)
			if dbError != nil {
				log.Fatalf("Failed to disconnect from MongoDB: %v", err)
			}
			log.Println("Connection to MongoDB closed. ✅")
			return dbError
		},
	})

	// Use a ping to verify the connection is alive.
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("FATAL: Could not ping MongoDB: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to MongoDB! ✅")
	return client, nil
}

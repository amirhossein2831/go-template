// my-go-app/internal/database/mongo.go
package database

import (
	"context"
	"event-collector/internal/config"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	instance *MongoDB  // *MongoDB	is the single-ton instance
	once     sync.Once // sync.Once will ensure the initialization code is executed only once.
)

// MongoDB holds the client and implements the singleton logic.
type MongoDB struct {
	Client *mongo.Client
}

// GetMongo returns a thread-safe, singleton instance of our MongoDB client.
func GetMongo(ctx context.Context) *MongoDB {
	once.Do(func() {
		client := connect(ctx)

		instance = &MongoDB{Client: client}
	})

	return instance
}

func connect(ctx context.Context) *mongo.Client {
	// This function will only be executed on the first call to GetInstance

	appConfig := config.GetConfig()
	mongoURI := appConfig.Database.URI

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("FATAL: Failed to create MongoDB client: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Fatalf("FATAL: Could not ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB! ✅")
	return client
}

// Close gracefully disconnects the client from MongoDB.
// It's essential to call this on application shutdown.
func (db *MongoDB) Close(ctx context.Context) {
	if db.Client == nil {
		return
	}
	if err := db.Client.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	log.Println("Connection to MongoDB closed. ✅")
}

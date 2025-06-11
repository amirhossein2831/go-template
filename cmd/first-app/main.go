package main

import (
	"event-collector/internal/config"
	"event-collector/internal/database"
	"event-collector/internal/service"
	"event-collector/internal/transport/grpc"
	"event-collector/internal/transport/http"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Provide all of our constructors. FX will call them in the
		// right order and handle their lifecycle.
		fx.Provide(
			// Core Components
			config.NewConfig,
			database.NewMongo,

			// Shared Business Logic Service
			service.NewGreetingService,

			// Transport Layers (both are provided)
			http.NewHTTPServer,
			grpc.NewGRPCServer,
		),
		// Invoke is used for functions that are needed for their side effects,
		// but don't provide any new types. This is our main application logic.
		fx.Invoke(func(*fiber.App) {}, runApplication),
	)

	// Run the application. This call is blocking.
	// It will only return when the application is shutting down.
	app.Run()
}

func runApplication() {
	fmt.Println("Starting application...")
}

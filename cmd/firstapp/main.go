package main

import (
	"event-collector/internal/config"
	"event-collector/internal/database"
	"event-collector/internal/services"
	"event-collector/internal/transport/grpc"
	handlers2 "event-collector/internal/transport/grpc/handlers"
	"event-collector/internal/transport/http"
	"event-collector/internal/transport/http/handlers"
	"event-collector/internal/transport/http/route"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
	grpc2 "google.golang.org/grpc"
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
			services.NewGreetingService,

			// Http handlers
			handlers.NewGreetingHandler,

			// Grpc handlers
			handlers2.NewGreetingHandler,

			// HTTP and GRPC Servers
			http.NewHTTPServer,
			grpc.NewGRPCServer,
		),
		// Invoke is used for functions that are needed for their side effects,
		// but don't provide any new types. This is our main application logic.
		fx.Invoke(
			func(*fiber.App) {},
			route.RegisterRoutes,
			func(server *grpc2.Server) {},
			grpc.RegisterServices,
			runApplication,
		),
	)

	// Run the application. This call is blocking.
	// It will only return when the application is shutting down.
	app.Run()
}

func runApplication() {
	fmt.Println("Starting application...")
}

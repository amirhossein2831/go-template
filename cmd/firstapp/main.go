package main

import (
	"event-collector/internal/firstapp/config"
	"event-collector/internal/firstapp/database"
	"event-collector/internal/firstapp/services"
	"event-collector/internal/firstapp/transport/grpc"
	handlers2 "event-collector/internal/firstapp/transport/grpc/handlers"
	"event-collector/internal/firstapp/transport/http"
	"event-collector/internal/firstapp/transport/http/handlers"
	"event-collector/internal/firstapp/transport/http/route"
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
			// Instantiate http server and register routes
			func(*fiber.App) {},
			route.RegisterRoutes,

			// Instantiate grpc server and register services
			func(server *grpc2.Server) {},
			grpc.RegisterServices,

			// run the app
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

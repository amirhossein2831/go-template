package http

import (
	"context"
	"event-collector/internal/config"
	"event-collector/pkg/validation"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/fx"
	"log"
)

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config) *fiber.App {
	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		StructValidator: &validation.StructValidator{Validator: validator.New()},
		AppName:         cfg.GetEnv(cfg.APP.Name),
	})

	// append middleware here
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	lc.Append(fx.Hook{
		// OnStart is called when the application starts.
		OnStart: func(ctx context.Context) error {
			go func() {
				port := cfg.GetEnv(cfg.Server.HTTP.Port)
				host := cfg.GetEnv(cfg.Server.HTTP.Host)
				log.Printf("Starting Fiber server on %s:%s ✅", host, port)
				if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
					log.Printf("FATAL: Failed to start Fiber server: %v", err)
				}
			}()
			return nil
		},
		// OnStop is called when the application is shutting down.
		OnStop: func(ctx context.Context) error {
			fmt.Println("Gracefully shutting down Fiber server ✅")
			return app.Shutdown()
		},
	})

	return app
}

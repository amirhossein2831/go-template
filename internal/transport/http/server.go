package http

import (
	"context"
	configs "event-collector/internal/config"
	"event-collector/pkg/validation"
	"fmt"
	"net"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, cfg *configs.Config) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		StructValidator: &validation.StructValidator{Validator: validator.New()},
		AppName:         cfg.APP.Name,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("http listen failed: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting Fiber server on %v ✅", addr)
			go func() {
				if err := app.Listener(ln); err != nil {
					log.Printf("Fiber server start failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Gracefully shutting down Fiber server ✅")
			app.Shutdown()
			ln.Close()
			return nil
		},
	})

	return app, nil
}

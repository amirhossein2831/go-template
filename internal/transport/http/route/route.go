package route

import (
	"event-collector/internal/transport/http/handlers"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

// RouterParams defines the dependencies for the route registration.
// Using 'fx.In' makes it easy to add more handlers later.
type RouterParams struct {
	fx.In
	App             *fiber.App
	GreetingHandler *handlers.GreetingHandler
	// Add other handlers here as dependencies, e.g.,
	// UserHandler *handlers.UserHandler
}

// RegisterRoutes sets up the application's specific routes.
func RegisterRoutes(p RouterParams) {
	api := p.App.Group("/api/v1")
	{
		// Register routes for the greeting service
		api.Post("/greeting", p.GreetingHandler.SayGreeting)

		// Register routes for other services here
		// api.Get("/users", p.UserHandler.ListUsers)
	}
}

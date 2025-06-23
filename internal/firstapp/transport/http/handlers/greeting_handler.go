package handlers

import (
	"context"
	"event-collector/internal/firstapp/services"
	"event-collector/internal/firstapp/transport/http/exceptions"
	"event-collector/internal/firstapp/transport/http/requests"
	"event-collector/internal/firstapp/transport/http/responses"
	"github.com/gofiber/fiber/v3"
)

// GreetingHandler implements the http server interface for the GreetingService.
type GreetingHandler struct {
	greetingService *services.GreetingService
}

// NewGreetingHandler is the constructor for our http handler.
func NewGreetingHandler(gs *services.GreetingService) *GreetingHandler {
	return &GreetingHandler{
		greetingService: gs,
	}
}

// SayGreeting is the implementation for http greeting
func (controller *GreetingHandler) SayGreeting(c fiber.Ctx) error {
	ctx := context.Background()
	req := new(requests.GreetingRequest)

	if err := c.Bind().Body(req); err != nil {
		return exceptions.HandleGreetingException(c, err)
	}

	res, err := controller.greetingService.GenerateGreetingLogic(ctx, req.Name)
	if err != nil {
		return exceptions.HandleGreetingException(c, err)
	}

	return responses.NewResponse(c).Status(fiber.StatusOK).Success(res)
}

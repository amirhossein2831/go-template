package handlers

import (
	"context"
	"event-collector/internal/services"
	"event-collector/internal/transport/http/requests"
	"event-collector/pkg/validation"
	"github.com/gofiber/fiber/v3"
)

type GreetingHandler struct {
	greetingService *services.GreetingService
}

func (controller *GreetingHandler) List(c fiber.Ctx) error {
	ctx := context.Background()
	req := new(requests.GreetingRequest)

	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(validation.ValidateStruct(req))
	}

	res, err := controller.greetingService.GenerateGreetingLogic(ctx, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"greeting_message": res,
	})
}

package exceptions

import (
	"errors"
	"event-collector/internal/firstapp/transport/http/errs"
	"event-collector/internal/firstapp/transport/http/responses"
	"github.com/gofiber/fiber/v3"
)

func HandleGreetingException(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errs.GreetingError):
		return responses.NewResponse(c).Status(fiber.StatusBadRequest).Error(err.Error())
	default:
		return HandleCommonException(c, err)
	}
}

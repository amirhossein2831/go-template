package exceptions

import (
	"errors"
	"event-collector/internal/firstapp/transport/http/errs"
	"event-collector/internal/firstapp/transport/http/responses"
	"event-collector/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func HandleCommonException(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errs.MongoInternalServerError):
		return responses.NewResponse(c).Status(fiber.StatusInternalServerError).Error(errs.DefaultServerError.Error())
	case errors.As(err, &validator.ValidationErrors{}):
		return responses.NewResponse(c).Status(fiber.StatusBadRequest).Error(validation.ValidateStruct(err))
	default:
		return responses.NewResponse(c).Status(fiber.StatusInternalServerError).Error(errs.DefaultServerError.Error())
	}
}

package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	Validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

type ErrorResponse struct {
	Field       string `json:"field"`           // The field name (e.g., "name", "redirect_uris[0]")
	Tag         string `json:"tag"`             // The validation tag (e.g., "required", "min", "url")
	Value       string `json:"value,omitempty"` // The invalid value, if applicable
	Message     string `json:"message"`         // A user-friendly error message
	Param       string `json:"param,omitempty"` // The parameter for the tag (e.g., "3" for min=3)
	Namespace   string `json:"namespace"`       // Full path to the field (e.g., "CreateClientRequest.Name")
	StructField string `json:"struct_field"`    // Original struct field name (e.g., "Name")
}

func ValidateStruct(err error) []ErrorResponse {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var formattedErrors []ErrorResponse
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			if fieldErr.Field() == "" {
				fieldName = strings.ToLower(fieldErr.StructField())
			}

			formattedErrors = append(formattedErrors, ErrorResponse{
				Field:       fieldName,
				Tag:         fieldErr.Tag(),
				Value:       fmt.Sprintf("%v", fieldErr.Value()),
				Message:     getValidationErrorMessage(fieldErr),
				Param:       fieldErr.Param(),
				Namespace:   fieldErr.Namespace(),
				StructField: fieldErr.StructField(),
			})
		}
		return formattedErrors
	}
	return nil
}

// getValidationErrorMessage provides a user-friendly message for common validation errors.
func getValidationErrorMessage(err validator.FieldError) string {
	fieldName := strings.ReplaceAll(err.Field(), "_", " ")
	fieldName = strings.Title(fieldName)

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", fieldName)
	case "min":
		if err.Kind().String() == "int" || err.Kind().String() == "float64" {
			return fmt.Sprintf("%s must be at least %s.", fieldName, err.Param())
		}
		return fmt.Sprintf("%s must have a minimum length of %s.", fieldName, err.Param())
	case "max":
		if err.Kind().String() == "int" || err.Kind().String() == "float64" {
			return fmt.Sprintf("%s must be at most %s.", fieldName, err.Param())
		}
		return fmt.Sprintf("%s must have a maximum length of %s.", fieldName, err.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID.", fieldName)
	case "url":
		return fmt.Sprintf("%s must be a valid URL.", fieldName)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI.", fieldName)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s.", fieldName, strings.ReplaceAll(err.Param(), " ", ", "))
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters.", fieldName)
	case "dive":
		return fmt.Sprintf("Invalid entry in %s.", fieldName)
	case "hexadecimal":
		return fmt.Sprintf("%s must be a valid hexadecimal string.", fieldName)
	case "len":
		return fmt.Sprintf("%s must have a length of %s.", fieldName, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address.", fieldName)
	default:
		return fmt.Sprintf("Validation failed for %s on the %s tag.", fieldName, err.Tag())
	}
}

package validation

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

// Validate Validator needs to implement the Validate method
func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(req interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = validationErr.StructNamespace()
			element.Tag = validationErr.Tag()
			element.Value = validationErr.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

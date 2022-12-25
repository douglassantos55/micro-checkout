package pkg

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	errors map[string]string
}

func (ValidationError) StatusCode() int {
	return 400
}

func (e ValidationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]string{"err": e.errors})
}

func (e ValidationError) Error() string {
	bytes, err := json.Marshal(e.errors)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type Validator interface {
	Validate(data any) error
}

type validate struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	return &validate{validator.New()}
}

func (v *validate) Validate(data any) error {
	err := v.validator.Struct(data)
	return makeValidationError(err.(validator.ValidationErrors))
}

func makeValidationError(err validator.ValidationErrors) ValidationError {
	errors := make(map[string]string)
	for _, error := range err {
		errors[error.Field()] = getErrorMessage(error.Tag())
	}
	return ValidationError{errors}
}

func getErrorMessage(tag string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "this field must be a valid email"
	default:
		return "something wrong is not right with this field"
	}
}

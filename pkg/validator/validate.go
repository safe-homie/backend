package validator

import "github.com/go-playground/validator/v10"

type Validator interface {
	Validate(i interface{}) error
}

type cValidator struct {
	validator *validator.Validate
}

func New() Validator {
	return &cValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (cv *cValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

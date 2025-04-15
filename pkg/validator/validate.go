package validator

import "github.com/go-playground/validator/v10"

type Validator interface {
	Validate(i interface{}) error
}

type cValidator struct {
	validator *validator.Validate
}

func New() Validator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("sensor_type", ValidateSensorType)
	v.RegisterValidation("device_type", ValidateDeviceType)
	v.RegisterValidation("device_action", ValidateDeviceAction)
	v.RegisterValidation("schedule_repeat", ValidateScheduleRepeat)

	return &cValidator{validator: v}
}

func (cv *cValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

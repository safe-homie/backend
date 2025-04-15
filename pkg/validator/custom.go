package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/safe-homie/backend/internal/domain"
)

func ValidateSensorType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return domain.IsValidSensorType(value)
}

func ValidateDeviceType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return domain.IsValidDeviceType(value)
}

func ValidateDeviceAction(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return domain.IsValidDeviceAction(value)
}

func ValidateScheduleRepeat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return domain.IsValidScheduleRepeat(value)
}

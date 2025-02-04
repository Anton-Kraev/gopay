package validator

import (
	"github.com/go-playground/validator/v10"

	"github.com/Anton-Kraev/gopay"
)

func ValidateID(fl validator.FieldLevel) bool {
	return gopay.ID(fl.Field().String()).Validate()
}

func ValidateStatus(fl validator.FieldLevel) bool {
	return gopay.Status(fl.Field().String()).Validate()
}

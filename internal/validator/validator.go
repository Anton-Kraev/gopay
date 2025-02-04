package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() (*Validator, error) {
	validate := validator.New()

	if err := validate.RegisterValidation("id", ValidateID); err != nil {
		return nil, err
	}

	if err := validate.RegisterValidation("status", ValidateStatus); err != nil {
		return nil, err
	}

	return &Validator{validator: validate}, nil
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

package app

import "github.com/go-playground/validator/v10"

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}

func (v *Validator) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		return NewValidationFailedError(err.Error())
	}

	return nil
}

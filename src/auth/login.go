package auth

import (
	"errors"
	"strings"

	"github.com/go-playground/validator"
)

type Login struct {
	Phone    string `json:"Phone"  validate:"required"`
	Password string `json:"Password"  validate:"required"`
}

func (c *Login) ValidateLogin() error {
	validate := validator.New()

	// Validate the struct
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("validation failed")
		}

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field(), err.Tag())
		}
		return errors.New("validation errors: " + strings.Join(validationErrors, ", "))
	}
	return nil
}

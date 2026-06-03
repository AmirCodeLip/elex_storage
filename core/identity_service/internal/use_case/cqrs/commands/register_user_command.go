package commands

import (
	"elex_storage/pkg/shared_kernel/models"

	"github.com/go-playground/validator/v10"
)

type RegisterUserCommand struct {
	Name     string `validate:"required,min=3"`
	Email    *string
	Phone    *string
	Password string `validate:"required,min=4"`
}

func (command RegisterUserCommand) Validate() error {
	validate := validator.New()
	err := validate.Struct(command)
	if err != nil {
		errors := models.NewCommonError(nil)
		for _, e := range err.(validator.ValidationErrors) {
			errors.AppendTag(e.Field(), e.Tag(), e.Param())
		}
		return errors
	}
	return nil
}

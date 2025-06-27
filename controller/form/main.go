package controller

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type FormValidationController struct{}

func NewFormValidationController() *FormValidationController {
	return &FormValidationController{}
}

// ValidateStruct validates a struct based on its tags
func (c *FormValidationController) ValidateStruct(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}

// ReturnFirstInvalidField returns the first invalid field in case of validation errors
func (c *FormValidationController) ReturnFirstInvalidField(err error) string {
	var field string
	for _, err := range err.(validator.ValidationErrors) {
		f := err.StructNamespace()
		field = strings.Split(f, ".")[1]
		break
	}
	return field
}

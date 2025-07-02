package controller

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
)

type formValidationControllerImpl struct{}

func NewFormValidationController() interfaces.FormValidationController {
	return &formValidationControllerImpl{}
}

// ValidateStruct validates a struct based on its tags
func (c *formValidationControllerImpl) ValidateStruct(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}

// ReturnFirstInvalidField returns the first invalid field in case of validation errors
func (c *formValidationControllerImpl) ReturnFirstInvalidField(err error) string {
	for _, err := range err.(validator.ValidationErrors) {
		ns := err.StructNamespace()
		if dotIndex := strings.Index(ns, "."); dotIndex != -1 {
			return ns[dotIndex+1:]
		}
		return ns
	}
	return ""
}
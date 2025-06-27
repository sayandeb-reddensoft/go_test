package services

import (
	controller "github.com/nelsonin-research-org/clenz-auth/controller/form"
)

type FormValidationServiceImpl struct {
	Controller controller.FormValidationController
}

func NewFormValidationService() *FormValidationServiceImpl {
	return &FormValidationServiceImpl{}
}

func (service *FormValidationServiceImpl) ValidateStruct(s interface{}) error {
	return service.Controller.ValidateStruct(s)
}

func (service *FormValidationServiceImpl) ReturnFirstInvalidField(err error) string {
	return service.Controller.ReturnFirstInvalidField(err)
}

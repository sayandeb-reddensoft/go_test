package services

import (
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
)

type formValidateServiceImpl struct {
	formValidateController interfaces.FormValidationController
}

func NewFormValidationService(c interfaces.FormValidationController) interfaces.FormValidateService {
	return &formValidateServiceImpl{formValidateController: c}
}

func (service *formValidateServiceImpl) ValidateStruct(s interface{}) error {
	return service.formValidateController.ValidateStruct(s)
}

func (service *formValidateServiceImpl) ReturnFirstInvalidField(err error) string {
	return service.formValidateController.ReturnFirstInvalidField(err)
}

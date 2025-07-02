package services

import (
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	emailModel "github.com/nelsonin-research-org/cdc-auth/models/email"
)

type emailServiceImpl struct {
	emailController interfaces.EmailController
}

func NewEmailService(c interfaces.EmailController) interfaces.EmailService {
	return &emailServiceImpl{emailController: c}
}

func (service *emailServiceImpl) SendResetPasswordOTP(data *emailModel.ResetPasswordMailContent, emails ...string) (bool, error) {
	return service.emailController.SendResetPasswordOTP(data, emails...)
}

func (service *emailServiceImpl) SendWelcomeOTP(data *emailModel.WelcomeAccountMailContent, emails ...string) (bool, error) {
	return service.emailController.SendWelcomeOTP(data, emails...)
}

func (service *emailServiceImpl) SendDeleteAccountOTP(data *emailModel.DeleteAccountMailContent, emails ...string) (bool, error) {
	return service.emailController.SendDeleteAccountOTP(data, emails...)
}

package services

import (
	controller "github.com/nelsonin-research-org/clenz-auth/controller/email"
	emailModel "github.com/nelsonin-research-org/clenz-auth/models/email"
)

type EmailServiceImpl struct {
	Controller controller.EmailController
}

func NewEmailService() *EmailServiceImpl {
	return &EmailServiceImpl{}
}

func (service *EmailServiceImpl) SendResetPasswordOTP(data *emailModel.ResetPasswordMailContent, email string) (bool, error) {
	return service.Controller.SendResetPasswordOTP(data, email)
}

func (service *EmailServiceImpl) SendWelcomeOTP(data *emailModel.WelcomeAccountMailContent, email string) (bool, error) {
	return service.Controller.SendWelcomeOTP(data, email)
}

func (service *EmailServiceImpl) ParseTemplate(templateName string, data interface{}) ([]byte, error) {
	return service.Controller.ParseTemplate(templateName, data)
}

func (service *EmailServiceImpl) SendEmail(recp []string, subject string, body []byte) error {
	return service.Controller.SendEmail(recp, subject, body)
}

func (service *EmailServiceImpl) DevSendEmail(recp []string, subject string, body []byte) error {
	return service.Controller.DevSendEmail(recp, subject, body)
}

func (service *EmailServiceImpl) SendDeleteAccountOTP(data *emailModel.DeleteAccountMailContent, email string) (bool, error) {
	return service.Controller.SendDeleteAccountOTP(data, email)
}

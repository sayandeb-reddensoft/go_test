package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/nelsonin-research-org/clenz-auth/globals"
	model "github.com/nelsonin-research-org/clenz-auth/models/email"
	"github.com/nelsonin-research-org/clenz-auth/utils"
)

type EmailController struct{}

func NewEmailController() *EmailController {
	return &EmailController{}
}

func (c *EmailController) SendResetPasswordOTP(data *model.ResetPasswordMailContent, emails ...string) (bool, error) {
	var MailContent model.ResetPasswordMailContent
	MailContent.Code = data.Code

	bodyData, err := c.ParseTemplate("forgotpassword.html", MailContent)
	if err != nil {
		return false, err
	}

	err = c.SendEmail(emails, "Reset Password", bodyData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *EmailController) SendWelcomeOTP(data *model.WelcomeAccountMailContent, emails ...string) (bool, error) {
	var MailContent model.WelcomeAccountMailContent
	MailContent.Code = data.Code
	MailContent.Name = data.Name

	bodyData, err := c.ParseTemplate("welcome.html", MailContent)
	if err != nil {
		return false, err
	}

	err = c.SendEmail(emails, "Welcome to Clenz!", bodyData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *EmailController) SendDeleteAccountOTP(data *model.DeleteAccountMailContent, emails ...string) (bool, error) {
	var MailContent model.DeleteAccountMailContent
	MailContent.Code = data.Code
	MailContent.Name = data.Name
	MailContent.Message = data.Message

	bodyData, err := c.ParseTemplate("delete-account.html", MailContent)
	if err != nil {
		return false, err
	}

	err = c.SendEmail(emails, "Account Deletion Requested", bodyData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *EmailController) ParseTemplate(templateName string, data interface{}) ([]byte, error) {
	templateFileName := os.Getenv("MAIL_TEMPLATE_PATH") + templateName
	tpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *EmailController) SendEmail(recp []string, subject string, body []byte) error {
	if utils.IsDevelopment() {
		err := c.DevSendEmail(recp, subject, body)
		if err != nil {
			return err
		}

		return nil
	} else {
		svc := ses.New(globals.AWSSesSession)

		from := os.Getenv("MAIL_FROM")
		to := recp
		subjectStr := subject
		bodyStr := string(body)

		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				ToAddresses: aws.StringSlice(to),
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(bodyStr),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subjectStr),
				},
			},
			Source: aws.String(from),
		}

		_, err := svc.SendEmail(input)
		if err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	}

	return nil
}

func (c *EmailController) DevSendEmail(recp []string, subject string, body []byte) error {
	if os.Getenv("TESTING") == "true" || os.Getenv("APP_STAGING") == "true" {
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
		from := "From: " + os.Getenv("MAIL_FROM") + "\n"
		to := "To: " + strings.Join(recp, ",") + "\n"
		sub := "Subject: " + subject + "\n"
		headers := from + to + sub + mime + "\n"

		msg := []byte(headers + string(body))

		username := os.Getenv("MAIL_USERNAME")
		password := os.Getenv("MAIL_PASSWORD")
		host := os.Getenv("MAIL_HOST")
		portStr := os.Getenv("MAIL_PORT")

		port, err := utils.StringToInt(portStr)
		if err != nil {
			return err
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		auth := smtp.PlainAuth("", username, password, host)

		err = smtp.SendMail(addr, auth, os.Getenv("MAIL_FROM"), recp, msg)
		if err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	}
	return nil
}

package controller

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	constants "github.com/nelsonin-research-org/cdc-auth/const"
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	model "github.com/nelsonin-research-org/cdc-auth/models/email"
	"github.com/nelsonin-research-org/cdc-auth/utils"
	env "github.com/nelsonin-research-org/cdc-auth/utils"
)

type emailControllerImpl struct{
	templateDir     string
	mailHost        string
	mailUserName    string
	mailPassword    string
	mailFrom        string
	mailPort        string
	awsSession      *session.Session
}

func NewEmailController(td, host, username, password, from, port string, sesSession *session.Session) interfaces.EmailController {
	return &emailControllerImpl{
		templateDir: td,
		mailHost: host,
		mailUserName: username,
		mailPassword: password,
		mailFrom: from,
		mailPort: port,
		awsSession: sesSession,
	}
}

func (c *emailControllerImpl) SendResetPasswordOTP(data *model.ResetPasswordMailContent, emails ...string) (bool, error) {
	var MailContent model.ResetPasswordMailContent
	MailContent.Code = data.Code

	bodyData, err := c.parseTemplate(constants.EMAIL_FORGOT_PASSWORD_ACCOUNT_TEMPLATE, MailContent)
	if err != nil {
		return false, err
	}

	err = c.sendEmail(emails, constants.EMAIL_FORGOT_PASSWORD_ACCOUNT_SUBJECT, bodyData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *emailControllerImpl) SendWelcomeOTP(data *model.WelcomeAccountMailContent, emails ...string) (bool, error) {
	var MailContent model.WelcomeAccountMailContent
	MailContent.Code = data.Code
	MailContent.Name = data.Name

	bodyData, err := c.parseTemplate(constants.EMAIL_ONBOARD_ACCOUNT_TEMPLATE, MailContent)
	if err != nil {
		return false, err
	}

	err = c.sendEmail(emails, constants.EMAIL_ONBOARD_ACCOUNT_SUBJECT, bodyData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *emailControllerImpl) SendDeleteAccountOTP(data *model.DeleteAccountMailContent, emails ...string) (bool, error) {
	var MailContent model.DeleteAccountMailContent
	MailContent.Code = data.Code
	MailContent.Name = data.Name
	MailContent.Message = data.Message

	bodyData, err := c.parseTemplate(constants.EMAIL_DELETE_ACCOUNT_TEMPLATE, MailContent)
	if err != nil {
		return false, err
	}

	err = c.sendEmail(emails, constants.EMAIL_DELETE_ACCOUNT_SUBJECT, bodyData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *emailControllerImpl) parseTemplate(templateName string, data interface{}) ([]byte, error) {
	templateFileName := c.templateDir + templateName
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

func (c *emailControllerImpl) sendEmail(recp []string, subject string, body []byte) error {
	if env.IsDevelopment() || env.IsStage() {
		err := c.sendEmailTest(recp, subject, body)
		if err != nil {
			return err
		}

		return nil
	} else {
		svc := ses.New(c.awsSession)

		from := c.mailFrom
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
			return errors.New("failed to send email:" + err.Error())
		}
	}

	return nil
}

func (c *emailControllerImpl) sendEmailTest(recp []string, subject string, body []byte) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	from := "From: " + c.mailFrom + "\n"
	to := "To: " + strings.Join(recp, ",") + "\n"
	sub := "Subject: " + subject + "\n"
	headers := from + to + sub + mime + "\n"

	msg := []byte(headers + string(body))

	username := c.mailUserName
	password := c.mailPassword
	host := c.mailHost
	portStr := c.mailPort

	port, err := utils.StringToInt(portStr)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", username, password, host)

	err = smtp.SendMail(addr, auth, c.mailFrom, recp, msg)
	if err != nil {
		return errors.New("failed to send email:" + err.Error())
	}

	return nil
}

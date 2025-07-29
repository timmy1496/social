package mailer

import (
	"bytes"
	"errors"
	gomail "gopkg.in/mail.v2"
	"text/template"
)

type mailTrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailTrapClient, error) {
	if apiKey == "" {
		return mailTrapClient{}, errors.New("api key is required")
	}

	return mailTrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailTrapClient) Send(templateFile, username, email string, data any, isSandbox bool) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)

	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

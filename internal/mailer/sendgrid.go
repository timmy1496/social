package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"text/template"
	"time"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetries; i++ {
		response, retryErr := m.client.Send(message)
		if retryErr != nil {

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts, err: %s", maxRetries, retryErr)
}

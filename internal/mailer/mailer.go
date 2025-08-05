package mailer

import "embed"

const (
	FromName            string = "GopherSocial"
	maxRetries          int    = 3
	UserWelcomeTemplate string = "user_invitations.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

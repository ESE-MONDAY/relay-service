package sender

import (
	"context"

	"github.com/ESE-MONDAY/relay-service/internal/models"
	mail "github.com/wneessen/go-mail"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPSender(
	host string,
	port int,
	username string,
	password string,
	from string,
) *SMTPSender {

	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) Send(
	ctx context.Context,
	email *models.Email,
) error {

	client, err := mail.NewClient(
		s.host,
		mail.WithPort(s.port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.username),
		mail.WithPassword(s.password),
	)

	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	if err := msg.From(s.from); err != nil {
		return err
	}

	if err := msg.To(email.To); err != nil {
		return err
	}

	msg.Subject(email.Subject)

	// Plain-text body
	if email.Text != "" {
		msg.SetBodyString(
			mail.TypeTextPlain,
			email.Text,
		)
	}

	// HTML alternative (multipart/alternative)
	if email.HTML != "" {
		msg.AddAlternativeString(
			mail.TypeTextHTML,
			email.HTML,
		)
	}

	return client.DialAndSendWithContext(ctx, msg)
}

package mailer

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/core/env"
)

type sendgridMailer struct {
	apiKey string
}

func NewSendgridMailer(config env.Config) (Mailer, error) {
	return &sendgridMailer{
		apiKey: config.SendgridAPIKey,
	}, nil
}

func (mailer *sendgridMailer) Send(e Mail) error {
	// create new *SGMailV3
	m := mail.NewV3Mail()

	from := mail.NewEmail(e.From.Name, e.From.Email)
	m.SetFrom(from)

	if e.PlainText != "" {
		content := mail.NewContent("text/plain", e.PlainText)
		m.AddContent(content)
	}

	if e.HTML != "" {
		content := mail.NewContent("text/html", e.HTML)
		m.AddContent(content)
	}

	// create new *Personalization
	personalization := mail.NewPersonalization()
	personalization.Subject = e.Subject

	// populate `personalization` with data
	var tos, ccs, bccs []*mail.Email
	for _, to := range e.To {
		tos = append(tos, mail.NewEmail(to.Name, to.Email))
	}
	personalization.AddTos(tos...)

	for _, cc := range e.Cc {
		ccs = append(ccs, mail.NewEmail(cc.Name, cc.Email))
	}
	personalization.AddCCs(ccs...)

	for _, bcc := range e.Bcc {
		bccs = append(bccs, mail.NewEmail(bcc.Name, bcc.Email))
	}
	personalization.AddBCCs(bccs...)

	// add `personalization` to `m`
	m.AddPersonalizations(personalization)

	request := sendgrid.GetRequest(mailer.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		return err
	}

	// TODO: Since response is a pointer, maybe it can be nil... Need to test.
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	return core.Unexpected(
		core.WithCode(fmt.Sprint(response.StatusCode)),
		core.WithMessage(response.Body),
	)
}

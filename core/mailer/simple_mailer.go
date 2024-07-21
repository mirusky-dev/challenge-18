package mailer

import (
	"strconv"
	"sync"
	"time"

	"github.com/mirusky-dev/challenge-18/core/env"
	simplemail "github.com/xhit/go-simple-mail/v2"
)

type simpleMailMailer struct {
	server *simplemail.SMTPServer
	client *simplemail.SMTPClient

	mu *sync.Mutex
}

func NewSimpleMailMailer(config env.Config) (Mailer, error) {
	server := simplemail.NewSMTPClient()

	port, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		return nil, err
	}

	server.Host = config.SMTPURL
	server.Port = port
	server.KeepAlive = true
	server.Username = config.SMTPClientID
	server.Password = config.SMTPClientSecret
	server.Encryption = simplemail.EncryptionSTARTTLS
	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}
	smtpClient.KeepAlive = true

	mu := sync.Mutex{}

	// https://github.com/xhit/go-simple-mail/issues/23#issuecomment-752015639
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			mu.Lock()
			if smtpClient.Noop() != nil {
				smtpClient, _ = server.Connect()
			}
			mu.Unlock()
		}
	}()

	return &simpleMailMailer{
		server: server,
		client: smtpClient,
		mu:     &mu,
	}, nil
}

func (mailer *simpleMailMailer) Send(email Mail) error {

	var tos, cc, bcc []string
	for _, v := range email.To {
		tos = append(tos, v.Email)
	}
	for _, v := range email.Cc {
		cc = append(cc, v.Email)
	}
	for _, v := range email.Bcc {
		bcc = append(bcc, v.Email)
	}

	mail := simplemail.NewMSG()
	mail.SetFrom(email.From.Email).
		AddTo(tos...).
		AddCc(cc...).
		AddBcc(bcc...).
		SetSubject(email.Subject)

	if email.PlainText != "" {
		mail.SetBody(simplemail.TextPlain, email.PlainText)
	}

	if email.HTML != "" {
		mail.SetBody(simplemail.TextHTML, email.HTML)
	}

	if mail.Error != nil {
		return mail.Error
	}

	// Just have sure we have an active connection and the pointer will not change while we are using
	// https://github.com/xhit/go-simple-mail/issues/23#issuecomment-752015639
	mailer.mu.Lock()
	defer mailer.mu.Unlock()

	if mailer.client.Noop() != nil {
		smtpClient, err := mailer.server.Connect()
		if err != nil {
			return err
		}

		mailer.client = smtpClient
	}

	err := mail.Send(mailer.client)
	return err
}

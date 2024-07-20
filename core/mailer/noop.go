package mailer

import (
	"fmt"

	"github.com/gobp/gobp/core/env"
)

type noopMailer struct{}

func NewNoopMailer(config env.Config) (Mailer, error) {
	return &noopMailer{}, nil
}

func (mailer *noopMailer) Send(email Mail) error {

	fmt.Printf("%+v\n", email)

	return nil
}

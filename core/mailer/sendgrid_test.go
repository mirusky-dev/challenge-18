package mailer

import (
	"path"
	"testing"

	"github.com/gobp/gobp/core/env"
)

func Test_sendgridMailer_Send(t *testing.T) {
	config, err := env.Load(path.Join("..", "..", "configs"))
	if err != nil || config.SendgridAPIKey == "" {
		t.Error("precondition failed: env.Load() error = ", err, " or config.SendgridAPIKey is empty")
		return
	}

	mailer, err := NewSendgridMailer(config)
	if err != nil {
		t.Error("precondition failed: NewSendgridMailer() error = ", err)
		return
	}

	type args struct {
		e Mail
	}
	tests := []struct {
		name    string
		mailer  Mailer
		args    args
		wantErr bool
	}{
		{
			name:   "Sendgrid integration test",
			mailer: mailer,
			args: args{
				e: Mail{
					From: EmailInfo{
						Name:  config.EmailSenderName,
						Email: config.EmailSender,
					},
					To: []EmailInfo{
						{
							Name:  "Example User",
							Email: "test@example.com",
						},
					},
					Subject:   "Sending with SendGrid is Fun",
					PlainText: "and easy to do anywhere, even with Go",
					HTML:      "<strong>and easy to do anywhere, even with Go</strong>",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mailer.Send(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("sendgridMailer.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

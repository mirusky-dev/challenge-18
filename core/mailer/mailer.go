package mailer

type EmailInfo struct {
	Name  string
	Email string
}

type Mail struct {
	From      EmailInfo
	To        []EmailInfo
	Cc        []EmailInfo
	Bcc       []EmailInfo
	Subject   string
	PlainText string
	HTML      string
}

type Mailer interface {
	Send(Mail) error
}

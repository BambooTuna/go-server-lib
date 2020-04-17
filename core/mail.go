package core

type Mailer interface {
	Send(toAddress []string, subject, message string) error
}

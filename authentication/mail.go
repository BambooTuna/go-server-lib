package authentication

type ActivateMailer interface {
	Send(toAddress []string, subject, message string) error
	SendActivateCode(code, mailAddress string) error
}

package services

type MailService interface {
	SendEmail(to string, subject string, body string) error
}

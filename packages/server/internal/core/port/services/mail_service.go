//go:generate mockgen -source=mail_service.go -destination=../../service/mocks/mail_service_mock.go -package=mocks -typed

package services

type MailService interface {
	SendEmail(to string, subject string, body string) error
}

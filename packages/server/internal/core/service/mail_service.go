package service

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	SMTPServer   string
	SMTPPort     int
	FromMail     string
	FromPassword string
}

func NewMailService(smtpServer string, smtpPort int, fromMail string, fromPassword string) *MailService {
	return &MailService{
		SMTPServer:   smtpServer,
		SMTPPort:     smtpPort,
		FromMail:     fromMail,
		FromPassword: fromPassword,
	}
}

func (ms *MailService) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ms.FromMail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(ms.SMTPServer, ms.SMTPPort, ms.FromMail, ms.FromPassword)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return err
	}
	return nil
}

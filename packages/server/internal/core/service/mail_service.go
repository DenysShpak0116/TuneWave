package service

import (
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"gopkg.in/gomail.v2"
)

type MailService struct {
	SMTPServer   string
	SMTPPort     int
	FromMail     string
	FromPassword string
}

func NewMailService(cfg *config.Config) services.MailService {
	return &MailService{
		SMTPServer:   cfg.Mail.StmpServer,
		SMTPPort:     cfg.Mail.SmtpPort,
		FromMail:     cfg.Mail.FromMail,
		FromPassword: cfg.Mail.FromPassword,
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

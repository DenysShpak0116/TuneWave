package service

import (
	"log/slog"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"gopkg.in/gomail.v2"
)

type MailService struct {
	SMTPServer   string
	SMTPPort     int
	FromMail     string
	FromPassword string
	logger  *slog.Logger
}

func NewMailService(cfg *config.Config, logger *slog.Logger) services.MailService {
	return &MailService{
		SMTPServer:   cfg.Mail.StmpServer,
		SMTPPort:     cfg.Mail.SmtpPort,
		FromMail:     cfg.Mail.FromMail,
		FromPassword: cfg.Mail.FromPassword,
		logger: logger,
	}
}

func (ms *MailService) SendEmail(to string, subject string, body string) error {
	const op = "core.service.MailService.SendEmail"
	logger := ms.logger.With(
		slog.String("op", op),
	)

	m := gomail.NewMessage()
	m.SetHeader("From", ms.FromMail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(ms.SMTPServer, ms.SMTPPort, ms.FromMail, ms.FromPassword)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		logger.Error("Failed to send mail", "err", err.Error())
		return err
	}

	logger.Info("Message sent")
	return nil
}

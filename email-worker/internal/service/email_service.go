package service

import (
	"net/smtp"

	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/config"
	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/model"
)

func SendEmail(cfg config.Config, event model.NotificationCreatedEvent) error {
	auth := smtp.PlainAuth(
		"",
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPHost,
	)

	message := []byte(
		"From: " + cfg.SMTPUsername + "\r\n" +
			"To: " + event.Recipient + "\r\n" +
			"Subject: " + event.Subject + "\r\n" +
			"\r\n" +
			event.Body,
	)

	return smtp.SendMail(
		cfg.SMTPHost+":"+cfg.SMTPPort,
		auth,
		cfg.SMTPUsername,
		[]string{event.Recipient},
		message,
	)
}

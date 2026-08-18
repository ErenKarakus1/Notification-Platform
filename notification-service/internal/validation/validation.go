package validation

import (
	"errors"
	"net/mail"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
)

func ValidateNotificationRequest(req model.NotificationRequest) error {
	if req.Body == "" {
		return errors.New("body is required")
	}
	if len(req.Body) > 5000 {
		return errors.New("body must be at most 5000 characters")
	}

	if req.Recipient == "" {
		return errors.New("recipient is required")
	}

	switch req.Channel {
	case "email":
		_, err := mail.ParseAddress(req.Recipient)
		if err != nil {
			return errors.New("invalid email")
		}
	default:
		return errors.New("invalid channel")
	}

	if len(req.Subject) > 500 {
		return errors.New("subject must be at most 500 characters")
	}
	return nil
}

package validation

import (
	"net/mail"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
)

func ValidateNotificationRequest(req model.NotificationRequest) bool {
	if req.Body == "" || len(req.Body) > 1000 {
		return false
	}

	if req.Recipient == "" {
		return false
	}

	switch req.Channel {
	case "email":
		_, err := mail.ParseAddress(req.Recipient)
		if err != nil {
			return false
		}
	default:
		return false
	}

	if len(req.Subject) > 500 {
		return false
	}
	return true
}

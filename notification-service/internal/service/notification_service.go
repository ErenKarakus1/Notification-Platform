package service

import (
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
	"github.com/google/uuid"
)

func CreateNotification(req model.NotificationRequest, customerID uuid.UUID) model.Notification {
	return model.Notification{
		ID:         uuid.New(),
		CustomerID: customerID,
		Recipient:  req.Recipient,
		Channel:    req.Channel,
		Subject:    req.Subject,
		Body:       req.Body,
		Status:     model.StatusQueued,
	}
}

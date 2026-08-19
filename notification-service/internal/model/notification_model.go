package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	StatusQueued = "queued"
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type Notification struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Recipient  string    `json:"recipient"`
	Channel    string    `json:"channel"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type NotificationRequest struct {
	Recipient string `json:"recipient"`
	Channel   string `json:"channel"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (r *NotificationRequest) Normalize() {
	r.Recipient = strings.ToLower(strings.TrimSpace(r.Recipient))
	r.Channel = strings.ToLower(strings.TrimSpace(r.Channel))
	r.Subject = strings.TrimSpace(r.Subject)
	r.Body = strings.TrimSpace(r.Body)
}

type CreateNotificationResponse struct {
	ID     uuid.UUID
	Status string
}

type NotificationCreatedEvent struct {
	NotificationID uuid.UUID `json:"notification_id"`
	CustomerID     uuid.UUID `json:"customer_id"`
	Recipient      string    `json:"recipient"`
	Channel        string    `json:"channel"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
}

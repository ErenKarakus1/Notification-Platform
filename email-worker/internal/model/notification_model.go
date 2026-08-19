package model

type NotificationCreatedEvent struct {
	NotificationID string `json:"notification_id"`
	CustomerID     string `json:"customer_id"`
	Recipient      string `json:"recipient"`
	Channel        string `json:"channel"`
	Subject        string `json:"subject"`
	Body           string `json:"body"`
}

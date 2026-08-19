package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "notification.created",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Close() {
	p.writer.Close()
}

func (p *Producer) PublishNotificationCreated(ctx context.Context, event model.NotificationCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.New("invalid paylaod")
	}
	return p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(event.NotificationID.String()),
			Value: payload,
		},
	)
}

package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/model"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishNotificationSent(ctx context.Context, event model.NotificationSentEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.New("invalid payload")
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.NotificationID),
		Value: payload,
	})
}

func (p *Producer) PublishNotificationFailed(ctx context.Context, event model.NotificationFailedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.New("invalid payload")
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.NotificationID),
		Value: payload,
	})
}

func (p *Producer) Close() {
	p.writer.Close()
}

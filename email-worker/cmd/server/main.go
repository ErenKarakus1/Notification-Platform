package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/config"
	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/kafka"
	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/model"
	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	consumer := kafka.NewConsumer([]string{"localhost:9092"})
	defer consumer.Close()

	sentProducer := kafka.NewProducer([]string{"localhost:9092"}, "notification.sent")
	defer sentProducer.Close()

	failedProducer := kafka.NewProducer([]string{"localhost:9092"}, "notification.failed")
	defer failedProducer.Close()

	log.Println("email worker started")

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}

		var event model.NotificationCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println(err)
			continue
		}
		log.Printf("notification recieved: %+v", event)

		if err := service.SendEmail(cfg, event); err != nil {
			log.Print("failed to send email")
			failedEvent := model.NotificationFailedEvent{NotificationID: event.NotificationID}
			if err := failedProducer.PublishNotificationFailed(context.Background(), failedEvent); err != nil {
				log.Println("failed to publish notification.failed")
				continue
			}
			log.Println("notification.failed published")
		} else {
			log.Println("email sent")
			sentEvent := model.NotificationSentEvent{NotificationID: event.NotificationID}
			if err := sentProducer.PublishNotificationSent(context.Background(), sentEvent); err != nil {
				log.Println("failed to publish notification.sent")
				continue
			}
			log.Println("notification.sent published")
		}

	}
}

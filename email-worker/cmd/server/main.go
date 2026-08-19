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

	producer := kafka.NewProducer([]string{"localhost:9092"})
	defer producer.Close()

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
			continue
		} else {
			log.Println("email sent")
			sentEvent := model.NotificationSentEvent{NotificationID: event.NotificationID}
			if err := producer.PublishNotificationSent(context.Background(), sentEvent); err != nil {
				log.Println("failed to publish notification.sent")
				continue
			}
			log.Println("notification.sent published")
		}

	}
}

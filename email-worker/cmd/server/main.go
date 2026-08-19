package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/kafka"
	"github.com/ErenKarakus1/Notification-Platform/email-worker/internal/model"
)

func main() {
	consumer := kafka.NewConsumer([]string{"localhost:9092"})
	defer consumer.Close()
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
	}
}

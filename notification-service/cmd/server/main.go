package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/config"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/db"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/handler"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/kafka"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/middleware"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("couldnt connect to postgres")
	}
	defer pool.Close()
	log.Println("connected to postgres")

	producer := kafka.NewProducer([]string{"localhost:9092"})
	defer producer.Close()
	log.Println("kafka producer initialized")

	sentConsumer := kafka.NewConsumer([]string{"localhost:9092"}, "notification.sent", "notification-sent")
	defer sentConsumer.Close()

	failedConsumer := kafka.NewConsumer([]string{"localhost:9092"}, "notification.failed", "notification-failed")
	defer failedConsumer.Close()

	go runNotificationConsumer(sentConsumer, pool, model.StatusSent)
	go runNotificationConsumer(failedConsumer, pool, model.StatusFailed)
	log.Println("kafka consumer initialized")

	router := gin.Default()
	router.POST("/notifications", middleware.CustomerAuthMiddleware(), handler.CreateNotificationHandler(pool, producer))
	router.GET("/notifications/:id", middleware.CustomerAuthMiddleware(), handler.GetNotificationByIDHandler(pool))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("couldnt run router")
	}
}

func runNotificationConsumer(consumer *kafka.Consumer, pool *pgxpool.Pool, status string) {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Println("failed to read message: ", err)
			continue
		}
		var event model.NotificationStatusEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("failed to unmarshal: ", err)
			continue
		}
		notificationID, err := uuid.Parse(event.NotificationID)
		if err != nil {
			log.Println("couldnt parse notification id:", err)
			continue
		}
		err = repository.UpdateNotificationStatus(context.Background(), pool, notificationID, status)
		if err != nil {
			if errors.Is(err, repository.ErrNotificationNotFound) {
				log.Println("notification not found: ", notificationID)
				continue
			}
			log.Println("couldnt update notification status: ", err)
			continue
		}
		log.Printf("marked notification as %s: %s", status, event.NotificationID)
	}
}

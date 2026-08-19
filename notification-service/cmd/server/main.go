package main

import (
	"log"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/config"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/db"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/handler"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/kafka"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/middleware"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	router.POST("/notifications", middleware.CustomerAuthMiddleware(), handler.CreateNotificationHandler(pool, producer))
	router.GET("/notifications/:id", middleware.CustomerAuthMiddleware(), handler.GetNotificationByIDHandler(pool))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("couldnt run router")
	}
}

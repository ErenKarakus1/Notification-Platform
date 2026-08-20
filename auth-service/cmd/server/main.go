package main

import (
	"log"

	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/config"
	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/db"
	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to postgres")

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterHandler(pool))
	router.POST("/auth/login", handler.LoginHandler(pool, cfg.JWTSecret))

	if err := router.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}

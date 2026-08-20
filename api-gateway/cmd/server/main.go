package main

import (
	"log"
	"time"

	"github.com/ErenKarakus1/Notification-Platform/api-gateway/internal/config"
	"github.com/ErenKarakus1/Notification-Platform/api-gateway/internal/middleware"
	"github.com/ErenKarakus1/Notification-Platform/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	authProxy := proxy.NewProxy("http://localhost:8083")
	notificationProxy := proxy.NewProxy("http://localhost:8082")

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	router := gin.Default()

	router.POST("/auth/register", authProxy)
	router.POST("/auth/login", authProxy)

	protected := router.Group("/")
	protected.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RateLimiter(redisClient, 100, time.Minute),
	)

	protected.POST("/notifications", middleware.AuthMiddleware(cfg.JWTSecret), notificationProxy)
	protected.GET("/notifications/:id", middleware.AuthMiddleware(cfg.JWTSecret), notificationProxy)

	if err := router.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}

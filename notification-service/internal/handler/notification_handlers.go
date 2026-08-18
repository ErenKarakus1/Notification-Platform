package handler

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/model"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/repository"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/service"
	"github.com/ErenKarakus1/Notification-Platform/notification-service/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNotificationHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerIDStr, _ := ctx.Get("customer_id")
		parsedCustomerID := customerIDStr.(uuid.UUID)
		var req model.NotificationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		req.Normalize()
		if err := validation.ValidateNotificationRequest(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		notification := service.CreateNotification(req, parsedCustomerID)
		if err := repository.CreateNotification(ctx.Request.Context(), pool, notification); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusCreated, model.CreateNotificationResponse{
			ID:     notification.ID,
			Status: notification.Status,
		})
	}
}

func GetNotificationByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerIDStr, _ := ctx.Get("customer_id")
		parsedCustomerID := customerIDStr.(uuid.UUID)
		notificationIDStr := ctx.Param("id")
		notificationID, err := uuid.Parse(notificationIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
			return
		}
		notification, err := repository.GetNotificationByID(ctx.Request.Context(), pool, notificationID, parsedCustomerID)
		if err != nil {
			if errors.Is(err, repository.ErrNotificationNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, notification)
	}
}

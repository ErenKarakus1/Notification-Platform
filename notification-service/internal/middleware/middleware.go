package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CustomerAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerIDStr := ctx.GetHeader("X-Customer-ID")
		if customerIDStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "customer id required"})
			return
		}
		parsedCustomerID, err := uuid.Parse(customerIDStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid customer id"})
			return
		}
		ctx.Set("customer_id", parsedCustomerID)
		ctx.Next()
	}
}

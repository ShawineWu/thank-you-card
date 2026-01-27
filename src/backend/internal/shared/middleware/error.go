package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler handles panics and errors in a consistent format
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				log.Printf("[%s] Panic recovered: %v", requestID, err)

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error": gin.H{
						"code":      "INTERNAL_SERVER_ERROR",
						"message":   "An unexpected error occurred",
						"requestId": requestID,
					},
				})
			}
		}()

		c.Next()
	}
}

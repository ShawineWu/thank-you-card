package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware logs HTTP requests with request ID and performance metrics
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set("requestId", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Start timer
		startTime := time.Now()

		// Log request
		log.Printf("[%s] %s %s - Started", requestID, c.Request.Method, c.Request.URL.Path)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Log response
		log.Printf(
			"[%s] %s %s - Completed in %v - Status: %d",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			duration,
			c.Writer.Status(),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("[%s] Error: %s", requestID, err.Error())
			}
		}
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("requestId")
	if !exists {
		return ""
	}

	str, _ := requestID.(string)
	return str
}

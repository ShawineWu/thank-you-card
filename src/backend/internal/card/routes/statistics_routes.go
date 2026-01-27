package routes

import (
	"github.com/castlery/thank-you-card/internal/card/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterStatisticsRoutes registers all statistics-related routes
func RegisterStatisticsRoutes(rg *gin.RouterGroup, h *handlers.StatisticsHandler) {
	stats := rg.Group("/statistics")
	{
		stats.GET("/user/:id", h.GetUserStats)
		stats.GET("/top10", h.GetTopRecipients)
	}
}

package routes

import (
	"github.com/castlery/thank-you-card/internal/analytics/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterAnalyticsRoutes registers all analytics-related routes
func RegisterAnalyticsRoutes(rg *gin.RouterGroup, h *handlers.AnalyticsHandler) {
	analytics := rg.Group("/analytics")
	{
		analytics.GET("/dashboard", h.GetDashboard)
		analytics.GET("/recognizers/top", h.GetTopRecognizers)
		analytics.GET("/teams", h.GetTeamAnalytics)
		analytics.GET("/values/distribution", h.GetValueDistribution)
		analytics.POST("/export", h.ExportData)
	}
}

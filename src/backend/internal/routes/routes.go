package routes

import (
	"github.com/castlery/thank-you-card/internal/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(
	r *gin.Engine,
	cardHandler *handlers.CardHandler,
	valueHandler *handlers.ValueHandler,
	statsHandler *handlers.StatisticsHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	teamsHandler *handlers.TeamsHandler,
) {
	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "thank-you-card-service",
		})
	})

	// Protected API routes
	api := r.Group("/api/v1")
	{
		// Card routes
		cards := api.Group("/cards")
		{
			cards.POST("", cardHandler.CreateCard)
			cards.GET("", cardHandler.GetCards)
			cards.GET("/received", cardHandler.GetReceivedCards)
			cards.GET("/sent", cardHandler.GetSentCards)
			cards.GET("/:id", cardHandler.GetCardByID)
		}

		// Employee routes
		employees := api.Group("/employees")
		{
			employees.GET("/search", cardHandler.SearchEmployees)
		}

		// Value routes
		values := api.Group("/values")
		{
			values.GET("", valueHandler.GetValues)
		}

		// Statistics routes
		stats := api.Group("/statistics")
		{
			stats.GET("/user/:id", statsHandler.GetUserStats)
			stats.GET("/top10", statsHandler.GetTopRecipients)
		}

		// Analytics routes
		analytics := api.Group("/analytics")
		{
			analytics.GET("/dashboard", analyticsHandler.GetDashboard)
			analytics.GET("/recognizers/top", analyticsHandler.GetTopRecognizers)
			analytics.GET("/teams", analyticsHandler.GetTeamAnalytics)
			analytics.GET("/values/distribution", analyticsHandler.GetValueDistribution)
			analytics.POST("/export", analyticsHandler.ExportData)
		}
	}

	// Teams Bot routes (unprotected, Teams uses its own auth)
	teams := r.Group("/api/teams")
	{
		teams.POST("/messages", teamsHandler.HandleMessages)
	}
}

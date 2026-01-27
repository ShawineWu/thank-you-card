package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/company/thank-you-card/internal/presentation/handlers"
	"github.com/company/thank-you-card/internal/presentation/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	cardHandler *handlers.CardHandler,
	employeeHandler *handlers.EmployeeHandler,
	milestoneHandler *handlers.MilestoneHandler,
) {
	// Global middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "thank-you-card-api",
			"timestamp": "2026-01-27T10:30:00Z",
		})
	})

	// API v1 routes with authentication
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		// Card routes
		cards := v1.Group("/cards")
		{
			cards.POST("", cardHandler.CreateCard)                    // Create card
			cards.GET("", cardHandler.GetCardFeed)                    // Get company feed
			cards.GET("/sent", cardHandler.GetSentCards)              // Get sent cards
			cards.GET("/received", cardHandler.GetReceivedCards)      // Get received cards
			cards.GET("/search", cardHandler.SearchCards)             // Search cards
			cards.GET("/:id", cardHandler.GetCard)                    // Get specific card
			cards.POST("/:id/share", cardHandler.ShareCard)           // Share card
			cards.DELETE("/:id", cardHandler.DeleteCard)              // Delete card
		}

		// Employee routes
		employees := v1.Group("/employees")
		{
			employees.GET("/search", employeeHandler.SearchEmployees)    // Search employees
			employees.GET("/:id", employeeHandler.GetEmployee)           // Get employee
			employees.POST("/validate", employeeHandler.ValidateEmployees) // Validate employees
		}

		// Company values routes
		v1.GET("/values", cardHandler.GetCompanyValues) // Get company values

		// Statistics routes
		statistics := v1.Group("/statistics")
		{
			statistics.GET("/personal", cardHandler.GetPersonalStatistics) // Personal stats
			statistics.GET("/top10", cardHandler.GetTop10)                 // Top 10 recognized
		}

		// Milestone routes
		milestones := v1.Group("/milestones")
		{
			milestones.GET("/me", milestoneHandler.GetMyMilestones)           // My milestones
			milestones.GET("/recent", milestoneHandler.GetRecentMilestones)   // Recent milestones
			milestones.GET("/type/:type", milestoneHandler.GetMilestonesByType) // Milestones by type
			milestones.GET("/:userId", milestoneHandler.GetUserMilestones)    // User milestones
		}

		// HR-only routes
		hr := v1.Group("/hr")
		hr.Use(middleware.RequireRole("HR"))
		{
			// HR analytics endpoints can be added here
			// hr.GET("/analytics/dashboard", hrHandler.GetAnalyticsDashboard)
			// hr.GET("/analytics/export", hrHandler.ExportData)
		}
	}
}
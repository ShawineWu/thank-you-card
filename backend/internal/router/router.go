package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"thank-you-card-backend/internal/handlers"
	"thank-you-card-backend/internal/middleware"
	"thank-you-card-backend/internal/repositories"
	"thank-you-card-backend/internal/services"
)

// SetupRouter wires repositories, services, and handlers, and returns a configured Gin engine.
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// global middleware
	r.Use(middleware.MockUserMiddleware())

	// Repositories
	cardRepo := repositories.NewCardRepository(db)
	emojiRepo := repositories.NewEmojiReactionRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	companyValueRepo := repositories.NewCompanyValueRepository(db)

	// Services
	cardService := services.NewCardService(cardRepo, emojiRepo)
	statsService := services.NewStatsService(cardRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, cardRepo, employeeRepo)
	companyValueService := services.NewCompanyValueService(companyValueRepo)

	// Handlers
	cardHandler := handlers.NewCardHandler(cardService)
	statsHandler := handlers.NewStatsHandler(statsService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	teamsHandler := handlers.NewTeamsHandler(cardService)
	companyValueHandler := handlers.NewCompanyValueHandler(companyValueService)

	api := r.Group("/api")
	{
		// Company Values
		api.GET("/company-values", companyValueHandler.GetByType)
		api.GET("/company-values/:id", companyValueHandler.GetByID)

		// Cards
		api.POST("/cards", cardHandler.CreateCard)
		api.GET("/cards/feed", cardHandler.GetFeed)
		api.GET("/cards/me/received", cardHandler.GetMyReceived)
		api.GET("/cards/me/sent", cardHandler.GetMySent)

		// Personal Stats
		api.GET("/cards/me/stats", statsHandler.GetPersonalStats)
		// Top recipients is available to all employees (not just HR)
		api.GET("/cards/top-recipients", analyticsHandler.GetTopRecognizedEmployees)

		// Emoji reactions
		api.POST("/cards/:id/reactions", cardHandler.React)
		api.DELETE("/cards/:id/reactions", cardHandler.RemoveReaction)

		// Teams Bot endpoints (simplified)
		teams := api.Group("/teams")
		{
			teams.GET("/feed", teamsHandler.GetFeed)
			teams.GET("/cards/:id", teamsHandler.GetCardDetail)
		}

		// HR Analytics (requires HR admin)
		analytics := api.Group("/analytics")
		analytics.Use(middleware.RequireHRAdmin(db))
		{
			analytics.GET("/dashboard", analyticsHandler.GetDashboard)
			analytics.GET("/recognizers", analyticsHandler.GetMostActiveRecognizers)
			analytics.GET("/teams", analyticsHandler.GetTeamPatterns)
			analytics.GET("/values", analyticsHandler.GetValuesDistribution)
			analytics.GET("/export", analyticsHandler.ExportCards)
		}
	}

	return r
}

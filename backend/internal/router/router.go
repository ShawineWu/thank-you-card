package router

import (
	"os"

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
	env := os.Getenv("ENV")
	if env == "production" {
		// Production: only allow specific origins
		config.AllowOrigins = []string{"http://localhost:8099", "http://localhost:5173"}
		config.AllowCredentials = true
	} else {
		// Development: allow all origins (including local network IPs)
		// This allows access from other devices on the same network
		config.AllowAllOrigins = true
		// Note: AllowAllOrigins and AllowCredentials cannot be used together
		// For development, this is acceptable
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	r.Use(cors.New(config))

	// Store db in context for middleware
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Repositories
	cardRepo := repositories.NewCardRepository(db)
	emojiRepo := repositories.NewEmojiReactionRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	companyValueRepo := repositories.NewCompanyValueRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Services
	cardService := services.NewCardService(cardRepo, emojiRepo)
	statsService := services.NewStatsService(cardRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, cardRepo, employeeRepo)
	companyValueService := services.NewCompanyValueService(companyValueRepo)
	authService := services.NewAuthService(userRepo, employeeRepo)

	// Handlers
	cardHandler := handlers.NewCardHandler(cardService)
	statsHandler := handlers.NewStatsHandler(statsService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	teamsHandler := handlers.NewTeamsHandler(cardService)
	companyValueHandler := handlers.NewCompanyValueHandler(companyValueService)
	authHandler := handlers.NewAuthHandler(authService)

	api := r.Group("/api")
	{
		// Auth (public)
		api.POST("/auth/login", authHandler.Login)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService, userRepo))
		{
			// Company Values
			protected.GET("/company-values", companyValueHandler.GetByType)
			protected.GET("/company-values/:id", companyValueHandler.GetByID)

			// Cards
			protected.POST("/cards", cardHandler.CreateCard)
			protected.GET("/cards/feed", cardHandler.GetFeed)
			protected.GET("/cards/me/received", cardHandler.GetMyReceived)
			protected.GET("/cards/me/sent", cardHandler.GetMySent)

			// Personal Stats
			protected.GET("/cards/me/stats", statsHandler.GetPersonalStats)
			// Top recipients is available to all employees (not just HR)
			protected.GET("/cards/top-recipients", analyticsHandler.GetTopRecognizedEmployees)

			// Emoji reactions
			protected.POST("/cards/:id/reactions", cardHandler.React)
			protected.DELETE("/cards/:id/reactions", cardHandler.RemoveReaction)

			// Teams Bot endpoints (simplified)
			teams := protected.Group("/teams")
			{
				teams.GET("/feed", teamsHandler.GetFeed)
				teams.GET("/cards/:id", teamsHandler.GetCardDetail)
			}

			// HR Analytics (requires HR or ADMIN role)
			analytics := protected.Group("/analytics")
			analytics.Use(middleware.RequireHRAdmin())
			{
				analytics.GET("/dashboard", analyticsHandler.GetDashboard)
				analytics.GET("/recognizers", analyticsHandler.GetMostActiveRecognizers)
				analytics.GET("/teams", analyticsHandler.GetTeamPatterns)
				analytics.GET("/values", analyticsHandler.GetValuesDistribution)
				analytics.GET("/export", analyticsHandler.ExportCards)
			}
		}
	}

	return r
}

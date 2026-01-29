package main

import (
	"log"
	"net/http"

	"github.com/castlery/thank-you-card/internal/cache"
	"github.com/castlery/thank-you-card/internal/events"
	"github.com/castlery/thank-you-card/internal/handlers"
	"github.com/castlery/thank-you-card/internal/repositories"
	"github.com/castlery/thank-you-card/internal/services"
	"github.com/castlery/thank-you-card/internal/shared/config"
	"github.com/castlery/thank-you-card/internal/shared/database"
	"github.com/castlery/thank-you-card/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 0. Load .env file (try .env first, then env for compatibility)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("env"); err != nil {
			log.Println("No .env or env file found, using system environment variables")
		} else {
			log.Println("Loaded config from env file")
		}
	} else {
		log.Println("Loaded config from .env")
	}

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
	if err := database.Initialize(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 2.1 Run Migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 2.2 Seed initial data
	if err := database.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// 3. Initialize repositories
	cardRepo := repositories.NewGormCardRepository(database.DB)
	valueRepo := repositories.NewGormCompanyValueRepository(database.DB)
	milestoneRepo := repositories.NewGormMilestoneRepository(database.DB)
	analyticsRepo := repositories.NewGormAnalyticsRepository(database.DB)

	// 4. Initialize external services
	var employeeService services.EmployeeService
	if cfg.External.AzureClientID != "" && cfg.External.AzureClientSecret != "" {
		employeeService = services.NewGraphEmployeeService(
			cfg.External.AzureTenantID,
			cfg.External.AzureClientID,
			cfg.External.AzureClientSecret,
		)
	} else {
		employeeService = services.NewMockEmployeeService(cfg.External.EmployeeDirectoryURL)
	}

	teamsPublisher := events.NewTeamsWebhookPublisher(cfg.External.TeamsWebhookURL, true)

	// 5. Initialize services
	cardCreationService := services.NewCardCreationService(
		cardRepo,
		valueRepo,
		employeeService,
		teamsPublisher,
	)
	statsService := services.NewPersonalStatisticsService(cardRepo, milestoneRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	exportService := services.NewExportService(analyticsRepo)

	// 6. Initialize cache
	analyticsCache := cache.NewInMemoryCache()

	// 7. Initialize handlers
	cardHandler := handlers.NewCardHandler(cardCreationService, cardRepo, employeeService)
	statsHandler := handlers.NewStatisticsHandler(statsService)
	valueHandler := handlers.NewValueHandler(valueRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, exportService, analyticsCache)
	teamsHandler := handlers.NewTeamsHandler(cardCreationService)

	// 8. Setup router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 9. Apply global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// 10. Register routes
	// Health check (before auth middleware)
	r.GET("/health", func(c *gin.Context) {
		status := "ok"
		if err := database.HealthCheck(); err != nil {
			status = "err"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "thank-you-card-service",
		})
	})

	// Protected routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))
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

	// Teams Bot routes (unprotected by AuthMiddleware, as Teams uses its own auth)
	teams := r.Group("/api/teams")
	{
		teams.POST("/messages", teamsHandler.HandleMessages)
	}

	// 11. Start server
	log.Printf("Thank You Card Service starting on :%s in %s mode", cfg.Server.Port, cfg.Server.Mode)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

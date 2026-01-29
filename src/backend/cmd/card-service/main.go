package main

import (
	"log"
	"net/http"

	"github.com/castlery/thank-you-card/internal/card/domain/services"
	"github.com/castlery/thank-you-card/internal/card/handlers"
	"github.com/castlery/thank-you-card/internal/card/infrastructure/events"
	"github.com/castlery/thank-you-card/internal/card/infrastructure/external"
	"github.com/castlery/thank-you-card/internal/card/infrastructure/persistence"
	"github.com/castlery/thank-you-card/internal/card/routes"
	"github.com/castlery/thank-you-card/internal/shared/config"
	"github.com/castlery/thank-you-card/internal/shared/database"
	"github.com/castlery/thank-you-card/internal/shared/middleware"

	analyticsServices "github.com/castlery/thank-you-card/internal/analytics/domain/services"
	analyticsHandlers "github.com/castlery/thank-you-card/internal/analytics/handlers"
	analyticsCache "github.com/castlery/thank-you-card/internal/analytics/infrastructure/cache"
	analyticsPersistence "github.com/castlery/thank-you-card/internal/analytics/infrastructure/persistence"
	analyticsRoutes "github.com/castlery/thank-you-card/internal/analytics/routes"

	teamsHandlers "github.com/castlery/thank-you-card/internal/teams/handlers"
	teamsRoutes "github.com/castlery/thank-you-card/internal/teams/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 0. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
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
	cardRepo := persistence.NewGormCardRepository(database.DB)
	valueRepo := persistence.NewGormCompanyValueRepository(database.DB)
	milestoneRepo := persistence.NewGormMilestoneRepository(database.DB)

	// 4. Initialize external services
	var employeeService services.EmployeeService
	if cfg.External.AzureClientID != "" && cfg.External.AzureClientSecret != "" {
		employeeService = external.NewGraphEmployeeService(
			cfg.External.AzureTenantID,
			cfg.External.AzureClientID,
			cfg.External.AzureClientSecret,
		)
	} else {
		employeeService = external.NewMockEmployeeService(cfg.External.EmployeeDirectoryURL)
	}

	teamsPublisher := events.NewTeamsWebhookPublisher(cfg.External.TeamsWebhookURL, true)

	// 5. Initialize domain services
	cardCreationService := services.NewCardCreationService(
		cardRepo,
		valueRepo,
		employeeService,
		teamsPublisher,
	)
	statsService := services.NewPersonalStatisticsService(cardRepo, milestoneRepo)

	// 6. Initialize handlers
	cardHandler := handlers.NewCardHandler(cardCreationService, cardRepo, employeeService)
	statsHandler := handlers.NewStatisticsHandler(statsService)
	valueHandler := handlers.NewValueHandler(valueRepo)

	// 7. Initialize analytics components
	analyticsRepo := analyticsPersistence.NewGormAnalyticsRepository(database.DB)
	analyticsCache := analyticsCache.NewInMemoryCache()
	analyticsService := analyticsServices.NewAnalyticsService(analyticsRepo)
	exportService := analyticsServices.NewExportService(analyticsRepo)
	analyticsHandler := analyticsHandlers.NewAnalyticsHandler(analyticsService, exportService, analyticsCache)

	// 7.1 Initialize Teams components
	teamsHandler := teamsHandlers.NewTeamsHandler(cardCreationService)

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
	// Public routes
	r.GET("/health", func(c *gin.Context) {
		status := "ok"
		if err := database.HealthCheck(); err != nil {
			status = "err"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "card-service",
		})
	})

	// Protected routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))
	{
		routes.RegisterCardRoutes(api, cardHandler)
		routes.RegisterStatisticsRoutes(api, statsHandler)
		routes.RegisterValueRoutes(api, valueHandler)
		analyticsRoutes.RegisterAnalyticsRoutes(api, analyticsHandler)
	}

	// Teams Bot routes (unprotected by AuthMiddleware, as Teams uses its own auth)
	teams := r.Group("/api/teams")
	{
		teamsRoutes.RegisterTeamsRoutes(teams, teamsHandler)
	}

	// 11. Start server
	log.Printf("Card Service starting on :%s in %s mode", cfg.Server.Port, cfg.Server.Mode)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

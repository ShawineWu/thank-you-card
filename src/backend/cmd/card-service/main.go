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

	"github.com/gin-gonic/gin"
)

func main() {
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

	// 3. Initialize repositories
	cardRepo := persistence.NewGormCardRepository(database.DB)
	valueRepo := persistence.NewGormCompanyValueRepository(database.DB)
	milestoneRepo := persistence.NewGormMilestoneRepository(database.DB)

	// 4. Initialize external services
	employeeService := external.NewMockEmployeeService(cfg.External.EmployeeDirectoryURL)
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
	cardHandler := handlers.NewCardHandler(cardCreationService, cardRepo)
	statsHandler := handlers.NewStatisticsHandler(statsService)

	// 7. Setup router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 8. Apply global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// 9. Register routes
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
	}

	// 10. Start server
	log.Printf("Card Service starting on :%s in %s mode", cfg.Server.Port, cfg.Server.Mode)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/company/thank-you-card/internal/container"
	"github.com/company/thank-you-card/internal/presentation/routes"
	"github.com/company/thank-you-card/pkg/config"
	"github.com/company/thank-you-card/pkg/database"
	"github.com/company/thank-you-card/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	logger.Init(cfg.LogLevel)
	logger.Info("Starting Thank You Card API server...")
	
	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	
	// Create dependency injection container
	appContainer := container.NewContainer(cfg, database.GetDB())
	
	// Start background processing
	appContainer.StartBackgroundProcessing()
	logger.Info("Background job processing started")
	
	// Initialize Gin router
	router := gin.Default()
	
	// Setup routes with handlers from container
	routes.SetupRoutes(
		router,
		appContainer.CardHandler,
		appContainer.EmployeeHandler,
		appContainer.MilestoneHandler,
	)
	
	// Setup graceful shutdown
	setupGracefulShutdown(appContainer)
	
	// Start server
	logger.Info("Starting server on port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupGracefulShutdown sets up graceful shutdown handling
func setupGracefulShutdown(container *container.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		logger.Info("Shutting down gracefully...")
		
		// Stop background processing
		container.StopBackgroundProcessing()
		logger.Info("Background job processing stopped")
		
		// Additional cleanup can be added here
		
		os.Exit(0)
	}()
}
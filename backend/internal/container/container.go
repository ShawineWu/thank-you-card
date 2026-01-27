package container

import (
	"time"

	"github.com/company/thank-you-card/internal/application/card"
	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/internal/infrastructure/cache"
	"github.com/company/thank-you-card/internal/infrastructure/database"
	"github.com/company/thank-you-card/internal/infrastructure/events"
	"github.com/company/thank-you-card/internal/infrastructure/external"
	"github.com/company/thank-you-card/internal/presentation/handlers"
	"github.com/company/thank-you-card/pkg/config"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB    *gorm.DB
	Cache cache.Cache

	// Repositories
	CardRepository         card.CardRepository
	CompanyValueRepository card.CompanyValueRepository
	MilestoneRepository    card.MilestoneRepository

	// External Services
	EmployeeService card.EmployeeDirectoryService

	// Event Publishing
	EventPublisher events.EventPublisher

	// Domain Services
	CardCreationService       *card.CardCreationService
	MilestoneTrackingService  *card.MilestoneTrackingService
	PersonalStatisticsService *card.PersonalStatisticsService

	// Factories
	CardFactory      *card.CardFactory
	MilestoneFactory *card.MilestoneFactory

	// Application Services - Commands
	CreateCardCommand cardApp.CreateCardCommand
	ShareCardCommand  cardApp.ShareCardCommand
	DeleteCardCommand cardApp.DeleteCardCommand

	// Application Services - Queries
	GetCardQuery          cardApp.GetCardQuery
	GetCardFeedQuery      cardApp.GetCardFeedQuery
	GetPersonalCardsQuery cardApp.GetPersonalCardsQuery
	SearchCardsQuery      cardApp.SearchCardsQuery
	GetStatisticsQuery    cardApp.GetStatisticsQuery
	GetTop10Query         cardApp.GetTop10Query
	GetCompanyValuesQuery cardApp.GetCompanyValuesQuery

	// Event Handlers
	CardCreatedEventHandler      *cardApp.CardCreatedEventHandler
	CardSharedEventHandler       *cardApp.CardSharedEventHandler
	MilestoneAchievedEventHandler *cardApp.MilestoneAchievedEventHandler
	EventDispatcher              *cardApp.EventDispatcher
	BackgroundJobProcessor       *cardApp.BackgroundJobProcessor

	// HTTP Handlers
	CardHandler      *handlers.CardHandler
	EmployeeHandler  *handlers.EmployeeHandler
	MilestoneHandler *handlers.MilestoneHandler
}

// NewContainer creates and wires up all dependencies
func NewContainer(cfg *config.Config, db *gorm.DB) *Container {
	container := &Container{
		Config: cfg,
		DB:     db,
	}

	container.setupInfrastructure()
	container.setupRepositories()
	container.setupExternalServices()
	container.setupEventPublishing()
	container.setupDomainServices()
	container.setupFactories()
	container.setupApplicationServices()
	container.setupEventHandlers()
	container.setupHTTPHandlers()

	return container
}

// setupInfrastructure initializes infrastructure components
func (c *Container) setupInfrastructure() {
	// Initialize cache with 5-minute default TTL
	c.Cache = cache.NewMemoryCache(5 * time.Minute)
}

// setupRepositories initializes repository implementations
func (c *Container) setupRepositories() {
	c.CardRepository = database.NewCardRepository(c.DB)
	c.CompanyValueRepository = database.NewCompanyValueRepository(c.DB)
	c.MilestoneRepository = database.NewMilestoneRepository(c.DB)
}

// setupExternalServices initializes external service clients
func (c *Container) setupExternalServices() {
	// For development, use mock employee service
	// In production, this would use real API endpoints and credentials
	c.EmployeeService = external.NewEmployeeDirectoryService("", "")
}

// setupEventPublishing initializes event publishing
func (c *Container) setupEventPublishing() {
	// For development, use mock event publisher
	// In production, this would use real webhook URLs
	webhookURL := c.Config.GetString("TEAMS_WEBHOOK_URL", "http://localhost:9001/webhook")
	
	if webhookURL == "mock" || webhookURL == "" {
		c.EventPublisher = events.NewMockEventPublisher()
	} else {
		c.EventPublisher = events.NewWebhookPublisher(webhookURL)
	}
}

// setupDomainServices initializes domain services
func (c *Container) setupDomainServices() {
	c.CardCreationService = card.NewCardCreationService(c.EmployeeService)
	c.MilestoneTrackingService = card.NewMilestoneTrackingService(
		c.CardRepository,
		c.MilestoneRepository,
	)
	c.PersonalStatisticsService = card.NewPersonalStatisticsService(c.CardRepository)
}

// setupFactories initializes domain factories
func (c *Container) setupFactories() {
	c.CardFactory = card.NewCardFactory()
	c.MilestoneFactory = card.NewMilestoneFactory()
}

// setupApplicationServices initializes application services
func (c *Container) setupApplicationServices() {
	// Command handlers
	c.CreateCardCommand = cardApp.NewCreateCardCommandHandler(
		c.CardRepository,
		c.CompanyValueRepository,
		c.CardCreationService,
		c.EventPublisher,
	)
	c.ShareCardCommand = cardApp.NewShareCardCommandHandler(
		c.CardRepository,
		c.EventPublisher,
	)
	c.DeleteCardCommand = cardApp.NewDeleteCardCommandHandler(c.CardRepository)

	// Query handlers
	c.GetCardQuery = cardApp.NewGetCardQueryHandler(c.CardRepository, c.Cache)
	c.GetCardFeedQuery = cardApp.NewGetCardFeedQueryHandler(c.CardRepository, c.Cache)
	c.GetPersonalCardsQuery = cardApp.NewGetPersonalCardsQueryHandler(c.CardRepository)
	c.SearchCardsQuery = cardApp.NewSearchCardsQueryHandler(c.CardRepository)
	c.GetStatisticsQuery = cardApp.NewGetStatisticsQueryHandler(c.PersonalStatisticsService, c.Cache)
	c.GetTop10Query = cardApp.NewGetTop10QueryHandler(c.CardRepository, c.Cache)
	c.GetCompanyValuesQuery = cardApp.NewGetCompanyValuesQueryHandler(c.CompanyValueRepository, c.Cache)
}

// setupEventHandlers initializes event handlers and background processing
func (c *Container) setupEventHandlers() {
	c.CardCreatedEventHandler = cardApp.NewCardCreatedEventHandler(
		c.MilestoneTrackingService,
		c.MilestoneRepository,
		c.EventPublisher,
	)
	c.CardSharedEventHandler = cardApp.NewCardSharedEventHandler()
	c.MilestoneAchievedEventHandler = cardApp.NewMilestoneAchievedEventHandler()

	c.EventDispatcher = cardApp.NewEventDispatcher(
		c.CardCreatedEventHandler,
		c.CardSharedEventHandler,
		c.MilestoneAchievedEventHandler,
	)

	// Initialize background job processor with 3 workers
	c.BackgroundJobProcessor = cardApp.NewBackgroundJobProcessor(c.EventDispatcher, 3)
}

// setupHTTPHandlers initializes HTTP handlers
func (c *Container) setupHTTPHandlers() {
	c.CardHandler = handlers.NewCardHandler(
		c.CreateCardCommand,
		c.ShareCardCommand,
		c.DeleteCardCommand,
		c.GetCardQuery,
		c.GetCardFeedQuery,
		c.GetPersonalCardsQuery,
		c.SearchCardsQuery,
		c.GetStatisticsQuery,
		c.GetTop10Query,
		c.GetCompanyValuesQuery,
	)

	c.EmployeeHandler = handlers.NewEmployeeHandler(c.EmployeeService)
	c.MilestoneHandler = handlers.NewMilestoneHandler(c.MilestoneRepository)
}

// StartBackgroundProcessing starts background job processing
func (c *Container) StartBackgroundProcessing() {
	c.BackgroundJobProcessor.Start()
}

// StopBackgroundProcessing stops background job processing
func (c *Container) StopBackgroundProcessing() {
	c.BackgroundJobProcessor.Stop()
}
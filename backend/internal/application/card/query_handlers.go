package card

import (
	"fmt"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/internal/infrastructure/cache"
	"github.com/company/thank-you-card/pkg/logger"
	"github.com/google/uuid"
)

// GetCardQueryHandler handles single card retrieval queries
type GetCardQueryHandler struct {
	cardRepository card.CardRepository
	cache          cache.Cache
}

// NewGetCardQueryHandler creates a new get card query handler
func NewGetCardQueryHandler(cardRepo card.CardRepository, cache cache.Cache) *GetCardQueryHandler {
	return &GetCardQueryHandler{
		cardRepository: cardRepo,
		cache:          cache,
	}
}

// Handle processes a get card query
func (h *GetCardQueryHandler) Handle(cardID uuid.UUID, userID string) (*CardResponse, error) {
	logger.Info("Processing get card query for card:", cardID)

	// Try to get from cache first
	cacheKey := fmt.Sprintf("card:%s", cardID.String())
	if cachedCard, found := h.cache.Get(cacheKey); found {
		if cardResponse, ok := cachedCard.(CardResponse); ok {
			logger.Debug("Card found in cache:", cardID)
			return &cardResponse, nil
		}
	}

	// Get from repository
	cardEntity, err := h.cardRepository.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	// Check if user can view this card
	if !cardEntity.CanBeViewedBy(userID) {
		return nil, fmt.Errorf("user %s is not authorized to view card %s", userID, cardID)
	}

	// Convert to response DTO
	response := MapCardToResponse(cardEntity)

	// Cache the response
	h.cache.Set(cacheKey, response)

	logger.Info("Successfully retrieved card:", cardID)
	return &response, nil
}

// GetCardFeedQueryHandler handles company-wide card feed queries
type GetCardFeedQueryHandler struct {
	cardRepository card.CardRepository
	cache          cache.Cache
}

// NewGetCardFeedQueryHandler creates a new get card feed query handler
func NewGetCardFeedQueryHandler(cardRepo card.CardRepository, cache cache.Cache) *GetCardFeedQueryHandler {
	return &GetCardFeedQueryHandler{
		cardRepository: cardRepo,
		cache:          cache,
	}
}

// Handle processes a get card feed query
func (h *GetCardFeedQueryHandler) Handle(page, pageSize int) (*PaginatedResponse[CardResponse], error) {
	logger.Info("Processing get card feed query, page:", page, "pageSize:", pageSize)

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get cards from repository
	cards, total, err := h.cardRepository.GetCompanyFeed(pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get card feed: %w", err)
	}

	// Convert to response DTOs
	cardResponses := MapCardsToResponse(cards)

	// Create paginated response
	response := NewPaginatedResponse(cardResponses, page, pageSize, total)

	logger.Info("Successfully retrieved card feed, total cards:", total)
	return &response, nil
}

// GetPersonalCardsQueryHandler handles personal card queries (sent/received)
type GetPersonalCardsQueryHandler struct {
	cardRepository card.CardRepository
}

// NewGetPersonalCardsQueryHandler creates a new get personal cards query handler
func NewGetPersonalCardsQueryHandler(cardRepo card.CardRepository) *GetPersonalCardsQueryHandler {
	return &GetPersonalCardsQueryHandler{
		cardRepository: cardRepo,
	}
}

// HandleSentCards processes a query for cards sent by a user
func (h *GetPersonalCardsQueryHandler) HandleSentCards(userID string, page, pageSize int) (*PaginatedResponse[CardResponse], error) {
	logger.Info("Processing get sent cards query for user:", userID)

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get cards from repository
	cards, total, err := h.cardRepository.GetCardsBySender(userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent cards: %w", err)
	}

	// Convert to response DTOs
	cardResponses := MapCardsToResponse(cards)

	// Create paginated response
	response := NewPaginatedResponse(cardResponses, page, pageSize, total)

	logger.Info("Successfully retrieved sent cards for user:", userID, "total:", total)
	return &response, nil
}

// HandleReceivedCards processes a query for cards received by a user
func (h *GetPersonalCardsQueryHandler) HandleReceivedCards(userID string, page, pageSize int) (*PaginatedResponse[CardResponse], error) {
	logger.Info("Processing get received cards query for user:", userID)

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get cards from repository
	cards, total, err := h.cardRepository.GetCardsByRecipient(userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get received cards: %w", err)
	}

	// Convert to response DTOs
	cardResponses := MapCardsToResponse(cards)

	// Create paginated response
	response := NewPaginatedResponse(cardResponses, page, pageSize, total)

	logger.Info("Successfully retrieved received cards for user:", userID, "total:", total)
	return &response, nil
}

// SearchCardsQueryHandler handles card search queries
type SearchCardsQueryHandler struct {
	cardRepository card.CardRepository
}

// NewSearchCardsQueryHandler creates a new search cards query handler
func NewSearchCardsQueryHandler(cardRepo card.CardRepository) *SearchCardsQueryHandler {
	return &SearchCardsQueryHandler{
		cardRepository: cardRepo,
	}
}

// Handle processes a search cards query
func (h *SearchCardsQueryHandler) Handle(request GetCardsRequest) (*PaginatedResponse[CardResponse], error) {
	logger.Info("Processing search cards query")

	// Validate request
	if validationErrors := ValidateGetCardsRequest(request); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Convert to domain filters
	filters := MapGetCardsRequestToSearchFilters(request)

	// Calculate offset
	offset := (request.Page - 1) * request.PageSize

	// Search cards in repository
	cards, total, err := h.cardRepository.SearchCards(filters, request.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search cards: %w", err)
	}

	// Convert to response DTOs
	cardResponses := MapCardsToResponse(cards)

	// Create paginated response
	response := NewPaginatedResponse(cardResponses, request.Page, request.PageSize, total)

	logger.Info("Successfully searched cards, total found:", total)
	return &response, nil
}

// GetStatisticsQueryHandler handles personal statistics queries
type GetStatisticsQueryHandler struct {
	statisticsService *card.PersonalStatisticsService
	cache             cache.Cache
}

// NewGetStatisticsQueryHandler creates a new get statistics query handler
func NewGetStatisticsQueryHandler(statsService *card.PersonalStatisticsService, cache cache.Cache) *GetStatisticsQueryHandler {
	return &GetStatisticsQueryHandler{
		statisticsService: statsService,
		cache:             cache,
	}
}

// Handle processes a get statistics query
func (h *GetStatisticsQueryHandler) Handle(userID string) (*PersonalStatisticsResponse, error) {
	logger.Info("Processing get statistics query for user:", userID)

	// Try to get from cache first
	cacheKey := cache.GetStatisticsKey(userID)
	if cachedStats, found := h.cache.Get(cacheKey); found {
		if statsResponse, ok := cachedStats.(PersonalStatisticsResponse); ok {
			logger.Debug("Statistics found in cache for user:", userID)
			return &statsResponse, nil
		}
	}

	// Calculate statistics using domain service
	stats, err := h.statisticsService.CalculatePersonalStatistics(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate statistics: %w", err)
	}

	// Convert to response DTO
	response := MapPersonalStatisticsToResponse(stats)

	// Cache the response for 5 minutes
	h.cache.SetWithTTL(cacheKey, response, 5*60*1000) // 5 minutes in milliseconds

	logger.Info("Successfully calculated statistics for user:", userID)
	return &response, nil
}

// GetTop10QueryHandler handles top 10 recognized employees queries
type GetTop10QueryHandler struct {
	cardRepository card.CardRepository
	cache          cache.Cache
}

// NewGetTop10QueryHandler creates a new get top 10 query handler
func NewGetTop10QueryHandler(cardRepo card.CardRepository, cache cache.Cache) *GetTop10QueryHandler {
	return &GetTop10QueryHandler{
		cardRepository: cardRepo,
		cache:          cache,
	}
}

// Handle processes a get top 10 query
func (h *GetTop10QueryHandler) Handle() ([]Top10EmployeeResponse, error) {
	logger.Info("Processing get top 10 query")

	// Try to get from cache first
	cacheKey := cache.DefaultCacheKeys().Top10EmployeesKey
	if cachedTop10, found := h.cache.Get(cacheKey); found {
		if top10Response, ok := cachedTop10.([]Top10EmployeeResponse); ok {
			logger.Debug("Top 10 found in cache")
			return top10Response, nil
		}
	}

	// Get from repository
	top10, err := h.cardRepository.GetTop10RecognizedEmployees()
	if err != nil {
		return nil, fmt.Errorf("failed to get top 10 employees: %w", err)
	}

	// Convert to response DTOs
	response := MapTop10EmployeesToResponse(top10)

	// Cache the response for 10 minutes
	h.cache.SetWithTTL(cacheKey, response, 10*60*1000) // 10 minutes in milliseconds

	logger.Info("Successfully retrieved top 10 employees")
	return response, nil
}

// GetCompanyValuesQueryHandler handles company values queries
type GetCompanyValuesQueryHandler struct {
	companyValueRepository card.CompanyValueRepository
	cache                  cache.Cache
}

// NewGetCompanyValuesQueryHandler creates a new get company values query handler
func NewGetCompanyValuesQueryHandler(valueRepo card.CompanyValueRepository, cache cache.Cache) *GetCompanyValuesQueryHandler {
	return &GetCompanyValuesQueryHandler{
		companyValueRepository: valueRepo,
		cache:                  cache,
	}
}

// Handle processes a get company values query
func (h *GetCompanyValuesQueryHandler) Handle() ([]CompanyValueResponse, error) {
	logger.Info("Processing get company values query")

	// Try to get from cache first
	cacheKey := cache.DefaultCacheKeys().CompanyValuesKey
	if cachedValues, found := h.cache.Get(cacheKey); found {
		if valuesResponse, ok := cachedValues.([]CompanyValueResponse); ok {
			logger.Debug("Company values found in cache")
			return valuesResponse, nil
		}
	}

	// Get from repository
	values, err := h.companyValueRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get company values: %w", err)
	}

	// Convert to response DTOs
	response := MapCompanyValuesToResponse(values)

	// Cache the response for 1 hour (company values rarely change)
	h.cache.SetWithTTL(cacheKey, response, 60*60*1000) // 1 hour in milliseconds

	logger.Info("Successfully retrieved company values, count:", len(response))
	return response, nil
}

// Query interfaces for dependency injection

// GetCardQuery represents the interface for getting a single card
type GetCardQuery interface {
	Handle(cardID uuid.UUID, userID string) (*CardResponse, error)
}

// GetCardFeedQuery represents the interface for getting the card feed
type GetCardFeedQuery interface {
	Handle(page, pageSize int) (*PaginatedResponse[CardResponse], error)
}

// GetPersonalCardsQuery represents the interface for getting personal cards
type GetPersonalCardsQuery interface {
	HandleSentCards(userID string, page, pageSize int) (*PaginatedResponse[CardResponse], error)
	HandleReceivedCards(userID string, page, pageSize int) (*PaginatedResponse[CardResponse], error)
}

// SearchCardsQuery represents the interface for searching cards
type SearchCardsQuery interface {
	Handle(request GetCardsRequest) (*PaginatedResponse[CardResponse], error)
}

// GetStatisticsQuery represents the interface for getting statistics
type GetStatisticsQuery interface {
	Handle(userID string) (*PersonalStatisticsResponse, error)
}

// GetTop10Query represents the interface for getting top 10 employees
type GetTop10Query interface {
	Handle() ([]Top10EmployeeResponse, error)
}

// GetCompanyValuesQuery represents the interface for getting company values
type GetCompanyValuesQuery interface {
	Handle() ([]CompanyValueResponse, error)
}

// Ensure handlers implement interfaces
var _ GetCardQuery = (*GetCardQueryHandler)(nil)
var _ GetCardFeedQuery = (*GetCardFeedQueryHandler)(nil)
var _ GetPersonalCardsQuery = (*GetPersonalCardsQueryHandler)(nil)
var _ SearchCardsQuery = (*SearchCardsQueryHandler)(nil)
var _ GetStatisticsQuery = (*GetStatisticsQueryHandler)(nil)
var _ GetTop10Query = (*GetTop10QueryHandler)(nil)
var _ GetCompanyValuesQuery = (*GetCompanyValuesQueryHandler)(nil)
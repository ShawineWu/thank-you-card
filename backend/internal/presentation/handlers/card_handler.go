package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	cardApp "github.com/company/thank-you-card/internal/application/card"
	"github.com/company/thank-you-card/pkg/logger"
)

// CardHandler handles HTTP requests for card operations
type CardHandler struct {
	createCardCommand     cardApp.CreateCardCommand
	shareCardCommand      cardApp.ShareCardCommand
	deleteCardCommand     cardApp.DeleteCardCommand
	getCardQuery          cardApp.GetCardQuery
	getCardFeedQuery      cardApp.GetCardFeedQuery
	getPersonalCardsQuery cardApp.GetPersonalCardsQuery
	searchCardsQuery      cardApp.SearchCardsQuery
	getStatisticsQuery    cardApp.GetStatisticsQuery
	getTop10Query         cardApp.GetTop10Query
	getCompanyValuesQuery cardApp.GetCompanyValuesQuery
}

// NewCardHandler creates a new card handler
func NewCardHandler(
	createCardCmd cardApp.CreateCardCommand,
	shareCardCmd cardApp.ShareCardCommand,
	deleteCardCmd cardApp.DeleteCardCommand,
	getCardQuery cardApp.GetCardQuery,
	getCardFeedQuery cardApp.GetCardFeedQuery,
	getPersonalCardsQuery cardApp.GetPersonalCardsQuery,
	searchCardsQuery cardApp.SearchCardsQuery,
	getStatisticsQuery cardApp.GetStatisticsQuery,
	getTop10Query cardApp.GetTop10Query,
	getCompanyValuesQuery cardApp.GetCompanyValuesQuery,
) *CardHandler {
	return &CardHandler{
		createCardCommand:     createCardCmd,
		shareCardCommand:      shareCardCmd,
		deleteCardCommand:     deleteCardCmd,
		getCardQuery:          getCardQuery,
		getCardFeedQuery:      getCardFeedQuery,
		getPersonalCardsQuery: getPersonalCardsQuery,
		searchCardsQuery:      searchCardsQuery,
		getStatisticsQuery:    getStatisticsQuery,
		getTop10Query:         getTop10Query,
		getCompanyValuesQuery: getCompanyValuesQuery,
	}
}

// CreateCard handles POST /api/v1/cards
func (h *CardHandler) CreateCard(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("CreateCard request received, requestID:", requestID)

	var request cardApp.CreateCardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request body: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserID(c)
	request.SenderID = userID

	// Execute command
	response, err := h.createCardCommand.Handle(request)
	if err != nil {
		logger.Error("Failed to create card:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"CREATE_CARD_FAILED",
			"Failed to create card: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Card created successfully, ID:", response.ID)
	c.JSON(http.StatusCreated, cardApp.NewApiResponse(
		response,
		"Card created successfully",
		requestID,
	))
}

// GetCard handles GET /api/v1/cards/:id
func (h *CardHandler) GetCard(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetCard request received, requestID:", requestID)

	// Parse card ID
	cardIDStr := c.Param("id")
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		logger.Error("Invalid card ID:", cardIDStr)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_CARD_ID",
			"Invalid card ID format",
			requestID,
			nil,
		))
		return
	}

	// Get user ID from context
	userID := getUserID(c)

	// Execute query
	response, err := h.getCardQuery.Handle(cardID, userID)
	if err != nil {
		logger.Error("Failed to get card:", err)
		c.JSON(http.StatusNotFound, cardApp.NewErrorResponse(
			"CARD_NOT_FOUND",
			"Card not found: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Card retrieved successfully, ID:", cardID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Card retrieved successfully",
		requestID,
	))
}

// GetCardFeed handles GET /api/v1/cards
func (h *CardHandler) GetCardFeed(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetCardFeed request received, requestID:", requestID)

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Execute query
	response, err := h.getCardFeedQuery.Handle(page, pageSize)
	if err != nil {
		logger.Error("Failed to get card feed:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_FEED_FAILED",
			"Failed to get card feed: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Card feed retrieved successfully, total:", response.Pagination.TotalItems)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Card feed retrieved successfully",
		requestID,
	))
}

// GetSentCards handles GET /api/v1/cards/sent
func (h *CardHandler) GetSentCards(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetSentCards request received, requestID:", requestID)

	// Get user ID from context
	userID := getUserID(c)

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Execute query
	response, err := h.getPersonalCardsQuery.HandleSentCards(userID, page, pageSize)
	if err != nil {
		logger.Error("Failed to get sent cards:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_SENT_CARDS_FAILED",
			"Failed to get sent cards: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Sent cards retrieved successfully for user:", userID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Sent cards retrieved successfully",
		requestID,
	))
}

// GetReceivedCards handles GET /api/v1/cards/received
func (h *CardHandler) GetReceivedCards(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetReceivedCards request received, requestID:", requestID)

	// Get user ID from context
	userID := getUserID(c)

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Execute query
	response, err := h.getPersonalCardsQuery.HandleReceivedCards(userID, page, pageSize)
	if err != nil {
		logger.Error("Failed to get received cards:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_RECEIVED_CARDS_FAILED",
			"Failed to get received cards: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Received cards retrieved successfully for user:", userID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Received cards retrieved successfully",
		requestID,
	))
}

// SearchCards handles GET /api/v1/cards/search
func (h *CardHandler) SearchCards(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("SearchCards request received, requestID:", requestID)

	// Parse search parameters
	var request cardApp.GetCardsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Error("Invalid search parameters:", err)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_SEARCH_PARAMS",
			"Invalid search parameters: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Set default pagination if not provided
	if request.Page == 0 {
		request.Page = 1
	}
	if request.PageSize == 0 {
		request.PageSize = 20
	}

	// Execute query
	response, err := h.searchCardsQuery.Handle(request)
	if err != nil {
		logger.Error("Failed to search cards:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"SEARCH_CARDS_FAILED",
			"Failed to search cards: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Cards searched successfully, total found:", response.Pagination.TotalItems)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Cards searched successfully",
		requestID,
	))
}

// ShareCard handles POST /api/v1/cards/:id/share
func (h *CardHandler) ShareCard(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("ShareCard request received, requestID:", requestID)

	// Parse card ID
	cardIDStr := c.Param("id")
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		logger.Error("Invalid card ID:", cardIDStr)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_CARD_ID",
			"Invalid card ID format",
			requestID,
			nil,
		))
		return
	}

	var request cardApp.ShareCardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request body: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Set card ID and user ID
	request.CardID = cardID
	request.SharedBy = getUserID(c)

	// Execute command
	err = h.shareCardCommand.Handle(request)
	if err != nil {
		logger.Error("Failed to share card:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"SHARE_CARD_FAILED",
			"Failed to share card: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Card shared successfully, ID:", cardID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		nil,
		"Card shared successfully",
		requestID,
	))
}

// DeleteCard handles DELETE /api/v1/cards/:id
func (h *CardHandler) DeleteCard(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("DeleteCard request received, requestID:", requestID)

	// Parse card ID
	cardIDStr := c.Param("id")
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		logger.Error("Invalid card ID:", cardIDStr)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_CARD_ID",
			"Invalid card ID format",
			requestID,
			nil,
		))
		return
	}

	// Get user ID from context
	userID := getUserID(c)

	// Execute command
	err = h.deleteCardCommand.Handle(cardID, userID)
	if err != nil {
		logger.Error("Failed to delete card:", err)
		c.JSON(http.StatusForbidden, cardApp.NewErrorResponse(
			"DELETE_CARD_FAILED",
			"Failed to delete card: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Card deleted successfully, ID:", cardID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		nil,
		"Card deleted successfully",
		requestID,
	))
}

// GetPersonalStatistics handles GET /api/v1/statistics/personal
func (h *CardHandler) GetPersonalStatistics(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetPersonalStatistics request received, requestID:", requestID)

	// Get user ID from context
	userID := getUserID(c)

	// Execute query
	response, err := h.getStatisticsQuery.Handle(userID)
	if err != nil {
		logger.Error("Failed to get personal statistics:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_STATISTICS_FAILED",
			"Failed to get personal statistics: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Personal statistics retrieved successfully for user:", userID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Personal statistics retrieved successfully",
		requestID,
	))
}

// GetTop10 handles GET /api/v1/statistics/top10
func (h *CardHandler) GetTop10(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetTop10 request received, requestID:", requestID)

	// Execute query
	response, err := h.getTop10Query.Handle()
	if err != nil {
		logger.Error("Failed to get top 10:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_TOP10_FAILED",
			"Failed to get top 10: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Top 10 retrieved successfully")
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Top 10 retrieved successfully",
		requestID,
	))
}

// GetCompanyValues handles GET /api/v1/values
func (h *CardHandler) GetCompanyValues(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetCompanyValues request received, requestID:", requestID)

	// Execute query
	response, err := h.getCompanyValuesQuery.Handle()
	if err != nil {
		logger.Error("Failed to get company values:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_VALUES_FAILED",
			"Failed to get company values: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	logger.Info("Company values retrieved successfully, count:", len(response))
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Company values retrieved successfully",
		requestID,
	))
}

// Helper functions

// getRequestID extracts or generates a request ID
func getRequestID(c *gin.Context) string {
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		return requestID
	}
	return uuid.New().String()
}

// getUserID extracts user ID from context (set by auth middleware)
func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	// For development/testing, return a default user ID
	return "emp001"
}

// getPaginationParams extracts pagination parameters from query string
func getPaginationParams(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
}
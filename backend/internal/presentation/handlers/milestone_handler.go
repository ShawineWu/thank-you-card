package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
	cardApp "github.com/company/thank-you-card/internal/application/card"
	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/pkg/logger"
)

// MilestoneHandler handles HTTP requests for milestone operations
type MilestoneHandler struct {
	milestoneRepository card.MilestoneRepository
}

// NewMilestoneHandler creates a new milestone handler
func NewMilestoneHandler(milestoneRepo card.MilestoneRepository) *MilestoneHandler {
	return &MilestoneHandler{
		milestoneRepository: milestoneRepo,
	}
}

// GetUserMilestones handles GET /api/v1/milestones/:userId
func (h *MilestoneHandler) GetUserMilestones(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetUserMilestones request received, requestID:", requestID)

	// Get user ID from path parameter
	userID := c.Param("userId")
	if userID == "" {
		logger.Error("User ID is required")
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"MISSING_USER_ID",
			"User ID is required",
			requestID,
			nil,
		))
		return
	}

	// Get current user ID for authorization
	currentUserID := getUserID(c)
	
	// For now, allow users to view their own milestones or all milestones (public)
	// In a more restrictive system, you might want to limit this
	_ = currentUserID // Placeholder for future authorization logic

	// Get milestones from repository
	milestones, err := h.milestoneRepository.GetByEmployeeID(userID)
	if err != nil {
		logger.Error("Failed to get user milestones:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_MILESTONES_FAILED",
			"Failed to get user milestones: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTOs
	response := cardApp.MapMilestonesToResponse(milestones)

	logger.Info("User milestones retrieved successfully for user:", userID, "count:", len(response))
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"User milestones retrieved successfully",
		requestID,
	))
}

// GetMyMilestones handles GET /api/v1/milestones/me
func (h *MilestoneHandler) GetMyMilestones(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetMyMilestones request received, requestID:", requestID)

	// Get current user ID
	userID := getUserID(c)

	// Get milestones from repository
	milestones, err := h.milestoneRepository.GetByEmployeeID(userID)
	if err != nil {
		logger.Error("Failed to get my milestones:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_MILESTONES_FAILED",
			"Failed to get my milestones: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTOs
	response := cardApp.MapMilestonesToResponse(milestones)

	logger.Info("My milestones retrieved successfully for user:", userID, "count:", len(response))
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"My milestones retrieved successfully",
		requestID,
	))
}

// GetRecentMilestones handles GET /api/v1/milestones/recent
func (h *MilestoneHandler) GetRecentMilestones(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetRecentMilestones request received, requestID:", requestID)

	// Get limit parameter (default: 10, max: 50)
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// Get recent milestones from repository
	milestones, err := h.milestoneRepository.GetRecentMilestones(limit)
	if err != nil {
		logger.Error("Failed to get recent milestones:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_RECENT_MILESTONES_FAILED",
			"Failed to get recent milestones: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTOs
	response := cardApp.MapMilestonesToResponse(milestones)

	logger.Info("Recent milestones retrieved successfully, count:", len(response))
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Recent milestones retrieved successfully",
		requestID,
	))
}

// GetMilestonesByType handles GET /api/v1/milestones/type/:type
func (h *MilestoneHandler) GetMilestonesByType(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetMilestonesByType request received, requestID:", requestID)

	// Get milestone type from path parameter
	milestoneType := c.Param("type")
	if milestoneType == "" {
		logger.Error("Milestone type is required")
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"MISSING_MILESTONE_TYPE",
			"Milestone type is required",
			requestID,
			nil,
		))
		return
	}

	// Validate milestone type
	if milestoneType != "CardsSent" && milestoneType != "CardsReceived" {
		logger.Error("Invalid milestone type:", milestoneType)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_MILESTONE_TYPE",
			"Milestone type must be 'CardsSent' or 'CardsReceived'",
			requestID,
			nil,
		))
		return
	}

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get milestones from repository
	milestones, total, err := h.milestoneRepository.GetMilestonesByType(milestoneType, pageSize, offset)
	if err != nil {
		logger.Error("Failed to get milestones by type:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_MILESTONES_BY_TYPE_FAILED",
			"Failed to get milestones by type: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTOs
	milestoneResponses := cardApp.MapMilestonesToResponse(milestones)

	// Create paginated response
	response := cardApp.NewPaginatedResponse(milestoneResponses, page, pageSize, total)

	logger.Info("Milestones by type retrieved successfully, type:", milestoneType, "total:", total)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Milestones by type retrieved successfully",
		requestID,
	))
}
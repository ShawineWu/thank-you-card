package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/castlery/thank-you-card/internal/analytics/domain/services"
	"github.com/castlery/thank-you-card/internal/analytics/dtos"
	"github.com/castlery/thank-you-card/internal/analytics/infrastructure/cache"
	sharedDtos "github.com/castlery/thank-you-card/internal/card/dtos"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles HTTP requests for analytics
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	exportService    *services.ExportService
	cache            *cache.InMemoryCache
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(
	analyticsService *services.AnalyticsService,
	exportService *services.ExportService,
	cache *cache.InMemoryCache,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		exportService:    exportService,
		cache:            cache,
	}
}

// GetDashboard handles GET /api/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	// Check cache
	const cacheKey = "analytics:dashboard"
	if cached, found := h.cache.Get(cacheKey); found {
		h.successResponse(c, http.StatusOK, cached, "")
		return
	}

	// Fetch data
	data, err := h.analyticsService.GetDashboardAnalytics()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch dashboard analytics", err.Error())
		return
	}

	// Cache for 5 minutes
	h.cache.Set(cacheKey, data, 5*time.Minute)

	h.successResponse(c, http.StatusOK, data, "")
}

// GetTopRecognizers handles GET /api/analytics/recognizers/top
func (h *AnalyticsHandler) GetTopRecognizers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Check cache
	cacheKey := fmt.Sprintf("analytics:top_recognizers:%d", limit)
	if cached, found := h.cache.Get(cacheKey); found {
		h.successResponse(c, http.StatusOK, gin.H{"recognizers": cached}, "")
		return
	}

	// Fetch data
	data, err := h.analyticsService.GetTopRecognizers(limit)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch top recognizers", err.Error())
		return
	}

	// Cache for 10 minutes
	h.cache.Set(cacheKey, data, 10*time.Minute)

	h.successResponse(c, http.StatusOK, gin.H{"recognizers": data}, "")
}

// GetTeamAnalytics handles GET /api/analytics/teams
func (h *AnalyticsHandler) GetTeamAnalytics(c *gin.Context) {
	// Check cache
	const cacheKey = "analytics:teams"
	if cached, found := h.cache.Get(cacheKey); found {
		h.successResponse(c, http.StatusOK, gin.H{"teams": cached}, "")
		return
	}

	// Fetch data
	data, err := h.analyticsService.GetTeamAnalytics()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch team analytics", err.Error())
		return
	}

	// Cache for 15 minutes
	h.cache.Set(cacheKey, data, 15*time.Minute)

	h.successResponse(c, http.StatusOK, gin.H{"teams": data}, "")
}

// GetValueDistribution handles GET /api/analytics/values/distribution
func (h *AnalyticsHandler) GetValueDistribution(c *gin.Context) {
	// Check cache
	const cacheKey = "analytics:value_distribution"
	if cached, found := h.cache.Get(cacheKey); found {
		h.successResponse(c, http.StatusOK, cached, "")
		return
	}

	// Fetch data
	data, err := h.analyticsService.GetValueDistribution()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch value distribution", err.Error())
		return
	}

	// Cache for 15 minutes
	h.cache.Set(cacheKey, data, 15*time.Minute)

	h.successResponse(c, http.StatusOK, data, "")
}

// ExportData handles POST /api/analytics/export
func (h *AnalyticsHandler) ExportData(c *gin.Context) {
	var req dtos.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_DATE", "Invalid start date format", nil)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_DATE", "Invalid end date format", nil)
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Generate CSV
	filters := dtos.ExportFilters{
		StartDate: startDate,
		EndDate:   endDate,
	}

	csvBuffer, err := h.exportService.ExportToCSV(filters)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "EXPORT_FAILED", "Failed to generate export", err.Error())
		return
	}

	// Send CSV file
	filename := fmt.Sprintf("recognition_cards_%s_%s.csv", req.StartDate, req.EndDate)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv", csvBuffer.Bytes())
}

// Helper methods

func (h *AnalyticsHandler) successResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, sharedDtos.StandardResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

func (h *AnalyticsHandler) errorResponse(c *gin.Context, status int, code, message string, details interface{}) {
	c.JSON(status, sharedDtos.ErrorResponse{
		Success: false,
		Error: sharedDtos.ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

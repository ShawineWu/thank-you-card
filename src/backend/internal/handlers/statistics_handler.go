package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/castlery/thank-you-card/internal/dtos"
	"github.com/castlery/thank-you-card/internal/services"

	"github.com/gin-gonic/gin"
)

// StatisticsHandler handles HTTP requests for statistics
type StatisticsHandler struct {
	statsService *services.PersonalStatisticsService
}

// NewStatisticsHandler creates a new statistics handler
func NewStatisticsHandler(statsService *services.PersonalStatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statsService: statsService,
	}
}

// GetUserStats handles GET /api/v1/statistics/user/:id
func (h *StatisticsHandler) GetUserStats(c *gin.Context) {
	employeeID := c.Param("id")
	if employeeID == "" {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_ID", "Employee ID is required", nil)
		return
	}

	sent, received, valueStats, milestones, err := h.statsService.GetUserStats(employeeID)
	if err != nil {
		log.Printf("Failed to get user stats: %v", err)
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch user statistics", err.Error())
		return
	}

	// Map to DTO
	valDTOs := make([]dtos.ValueUsageStat, len(valueStats))
	for i, v := range valueStats {
		valDTOs[i] = dtos.ValueUsageStat{
			ValueID:   v.ValueID,
			ValueName: v.ValueName,
			Count:     v.Count,
		}
	}

	msDTOs := make([]dtos.MilestoneResponse, len(milestones))
	for i, m := range milestones {
		msDTOs[i] = dtos.MilestoneResponse{
			Type:        string(m.MilestoneType),
			Threshold:   m.Threshold,
			Title:       m.Title,
			Description: m.Description,
			AchievedAt:  m.AchievedAt,
		}
	}

	h.successResponse(c, http.StatusOK, dtos.UserStatsResponse{
		EmployeeID:    employeeID,
		CardsSent:     sent,
		CardsReceived: received,
		ValueStats:    valDTOs,
		Milestones:    msDTOs,
	}, "")
}

// GetTopRecipients handles GET /api/v1/statistics/top10
func (h *StatisticsHandler) GetTopRecipients(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	topRecipients, err := h.statsService.GetTopRecipients(limit)
	if err != nil {
		log.Printf("Failed to get top recipients: %v", err)
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch top recipients", err.Error())
		return
	}

	// Map to DTO
	topDTOs := make([]dtos.TopEmployeeStat, len(topRecipients))
	for i, r := range topRecipients {
		topDTOs[i] = dtos.TopEmployeeStat{
			EmployeeID:    r.EmployeeID,
			EmployeeName:  r.EmployeeName,
			CardsReceived: r.Count,
			Rank:          i + 1,
		}
	}

	h.successResponse(c, http.StatusOK, topDTOs, "")
}

// Helper methods (extracted or replicated from card_handler)

func (h *StatisticsHandler) successResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, dtos.StandardResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

func (h *StatisticsHandler) errorResponse(c *gin.Context, status int, code, message string, details interface{}) {
	c.JSON(status, dtos.ErrorResponse{
		Success: false,
		Error: dtos.ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

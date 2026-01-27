package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/services"
)

type StatsHandler struct {
	statsService services.StatsService
}

func NewStatsHandler(statsService services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetPersonalStats(c *gin.Context) {
	emp := getCurrentEmployee(c)
	if emp == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	stats, err := h.statsService.GetPersonalStats(ctx, emp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load personal stats"})
		return
	}

	resp := dto.PersonalStatsResponse{
		TotalSent:         stats.TotalSent,
		TotalReceived:     stats.TotalReceived,
		TopSentValues:     make([]dto.ValueCountResponse, len(stats.TopSentValues)),
		TopReceivedValues: make([]dto.ValueCountResponse, len(stats.TopReceivedValues)),
	}

	for i, vc := range stats.TopSentValues {
		resp.TopSentValues[i] = dto.ValueCountResponse{
			ValueID: vc.ValueID,
			Code:    vc.Code,
			Name:    vc.Name,
			Type:    vc.Type,
			Count:   vc.Count,
		}
	}

	for i, vc := range stats.TopReceivedValues {
		resp.TopReceivedValues[i] = dto.ValueCountResponse{
			ValueID: vc.ValueID,
			Code:    vc.Code,
			Name:    vc.Name,
			Type:    vc.Type,
			Count:   vc.Count,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *StatsHandler) GetTopRecipients(c *gin.Context) {
	// This endpoint is now handled by analytics service
	// Redirect to analytics endpoint or return message
	c.JSON(http.StatusOK, gin.H{
		"message": "Use /api/analytics/top-recipients for HR admin access, or implement public endpoint here",
	})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/middleware"
	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/services"
)

type MilestoneHandler struct {
	milestoneService services.MilestoneService
}

func NewMilestoneHandler(milestoneService services.MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{
		milestoneService: milestoneService,
	}
}

// GetAllMilestones returns all milestone definitions
// GET /api/milestones
func (h *MilestoneHandler) GetAllMilestones(c *gin.Context) {
	ctx := c.Request.Context()
	milestones, err := h.milestoneService.GetAllMilestones(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load milestones"})
		return
	}

	resp := make([]dto.MilestoneResponse, len(milestones))
	for i, m := range milestones {
		resp[i] = dto.MilestoneResponse{
			ID:          m.ID,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			Type:        m.Type,
			Threshold:   m.Threshold,
			IconURL:     m.IconURL,
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resp})
}

// GetMyAchievements returns the current user's achievements
// GET /api/milestones/me
func (h *MilestoneHandler) GetMyAchievements(c *gin.Context) {
	// Get current user from context
	userVal, ok := c.Get(middleware.ContextCurrentUser)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	user, ok := userVal.(*models.User)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	ctx := c.Request.Context()
	achievements, err := h.milestoneService.GetUserAchievements(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load achievements"})
		return
	}

	resp := make([]dto.UserAchievementResponse, len(achievements))
	for i, a := range achievements {
		resp[i] = dto.UserAchievementResponse{
			Milestone: dto.MilestoneResponse{
				ID:          a.Milestone.ID,
				Code:        a.Milestone.Code,
				Name:        a.Milestone.Name,
				Description: a.Milestone.Description,
				Type:        a.Milestone.Type,
				Threshold:   a.Milestone.Threshold,
				IconURL:     a.Milestone.IconURL,
			},
			AchievedAt: a.AchievedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resp})
}

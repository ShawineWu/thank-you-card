package routes

import (
	"github.com/castlery/thank-you-card/internal/teams/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterTeamsRoutes registers the Teams Bot routes
func RegisterTeamsRoutes(r *gin.RouterGroup, h *handlers.TeamsHandler) {
	r.POST("/messages", h.HandleMessages)
}

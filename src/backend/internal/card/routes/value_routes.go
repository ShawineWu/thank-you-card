package routes

import (
	"github.com/castlery/thank-you-card/internal/card/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterValueRoutes registers all company value-related routes
func RegisterValueRoutes(rg *gin.RouterGroup, h *handlers.ValueHandler) {
	values := rg.Group("/values")
	{
		values.GET("", h.GetValues)
	}
}

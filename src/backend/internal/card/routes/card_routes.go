package routes

import (
	"github.com/castlery/thank-you-card/internal/card/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterCardRoutes registers all card-related routes
func RegisterCardRoutes(rg *gin.RouterGroup, h *handlers.CardHandler) {
	cards := rg.Group("/cards")
	{
		cards.POST("", h.CreateCard)
		cards.GET("", h.GetCards)
		cards.GET("/received", h.GetReceivedCards)
		cards.GET("/sent", h.GetSentCards)
		cards.GET("/:id", h.GetCardByID)
	}

	employees := rg.Group("/employees")
	{
		employees.GET("/search", h.SearchEmployees)
	}
}

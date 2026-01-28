package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/repositories"
	"thank-you-card-backend/internal/services"
)

// TeamsHandler provides simplified endpoints for Teams Bot integration
type TeamsHandler struct {
	cardService services.CardService
}

func NewTeamsHandler(cardService services.CardService) *TeamsHandler {
	return &TeamsHandler{cardService: cardService}
}

// GetFeed returns a simplified feed for Teams Bot
// This endpoint is optimized for Teams channel display
func (h *TeamsHandler) GetFeed(c *gin.Context) {
	ctx := c.Request.Context()
	filters := parseCardFilters(c)
	// Ensure we use proper filters for Teams feed
	filters.Page = parseIntQuery(c, "page", 1)
	filters.PageSize = parseIntQuery(c, "pageSize", 20)

	cards, total, err := h.cardService.GetCompanyFeed(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load feed"})
		return
	}

	// Simplified response for Teams
	resp := make([]dto.TeamsCardResponse, 0, len(cards))
	for _, card := range cards {
		// Collect recipient names
		recipientNames := make([]string, len(card.Recipients))
		for i, r := range card.Recipients {
			recipientNames[i] = r.Recipient.Name
		}

		// Collect value names
		valueNames := make([]string, len(card.Values))
		for i, v := range card.Values {
			valueNames[i] = v.CompanyValue.Name
		}

		// Count reactions
		reactionCount := len(card.Reactions)

		resp = append(resp, dto.TeamsCardResponse{
			ID:             card.ID,
			SenderName:     card.Sender.Name,
			RecipientNames: recipientNames,
			Reason:         card.Reason,
			Values:         valueNames,
			CreatedAt:      card.CreatedAt.Format("2006-01-02 15:04:05"),
			ReactionCount:  reactionCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": total,
	})
}

// GetCardDetail returns a single card detail for Teams Bot
func (h *TeamsHandler) GetCardDetail(c *gin.Context) {
	cardID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	ctx := c.Request.Context()
	// Get feed with large page size to find the card
	filters := repositories.CardFilters{
		Page:     1,
		PageSize: 1000, // Large enough to find the card
	}

	cards, _, err := h.cardService.GetCompanyFeed(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load card"})
		return
	}

	// Find the card by ID
	var foundCard *dto.TeamsCardResponse
	for _, card := range cards {
		if card.ID == cardID {
			recipientNames := make([]string, len(card.Recipients))
			for i, r := range card.Recipients {
				recipientNames[i] = r.Recipient.Name
			}

			valueNames := make([]string, len(card.Values))
			for i, v := range card.Values {
				valueNames[i] = v.CompanyValue.Name
			}

			foundCard = &dto.TeamsCardResponse{
				ID:             card.ID,
				SenderName:     card.Sender.Name,
				RecipientNames: recipientNames,
				Reason:         card.Reason,
				Values:         valueNames,
				CreatedAt:      card.CreatedAt.Format("2006-01-02 15:04:05"),
				ReactionCount:  len(card.Reactions),
			}
			break
		}
	}

	if foundCard == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	c.JSON(http.StatusOK, foundCard)
}

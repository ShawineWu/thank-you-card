package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
	"github.com/castlery/thank-you-card/internal/card/domain/repositories"
	"github.com/castlery/thank-you-card/internal/card/domain/services"
	"github.com/castlery/thank-you-card/internal/card/dtos"
	"github.com/castlery/thank-you-card/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CardHandler handles HTTP requests for cards
type CardHandler struct {
	cardCreationService *services.CardCreationService
	cardRepo            repositories.CardRepository
}

// NewCardHandler creates a new card handler
func NewCardHandler(
	cardCreationService *services.CardCreationService,
	cardRepo repositories.CardRepository,
) *CardHandler {
	return &CardHandler{
		cardCreationService: cardCreationService,
		cardRepo:            cardRepo,
	}
}

// CreateCard handles POST /api/cards
func (h *CardHandler) CreateCard(c *gin.Context) {
	senderID, err := middleware.GetUserID(c)
	if err != nil {
		h.errorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req dtos.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	card, err := h.cardCreationService.CreateCard(
		senderID,
		req.RecipientIDs,
		req.RecognitionReason,
		req.ValueIDs,
	)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "CREATION_FAILED", "Failed to create card", err.Error())
		return
	}

	h.successResponse(c, http.StatusCreated, h.mapCardToResponse(card), "Card created successfully")
}

// GetCards handles GET /api/cards (company-wide feed)
func (h *CardHandler) GetCards(c *gin.Context) {
	var params dtos.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_PARAMS", "Invalid pagination parameters", err.Error())
		return
	}

	cards, total, err := h.cardRepo.FindAll(params.Page, params.PageSize)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch cards", err.Error())
		return
	}

	h.paginatedResponse(c, cards, total, params.Page, params.PageSize)
}

// GetCardByID handles GET /api/cards/:id
func (h *CardHandler) GetCardByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid card ID format", nil)
		return
	}

	card, err := h.cardRepo.FindByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch card", err.Error())
		return
	}

	if card == nil {
		h.errorResponse(c, http.StatusNotFound, "NOT_FOUND", "Card not found", nil)
		return
	}

	h.successResponse(c, http.StatusOK, h.mapCardToResponse(card), "")
}

// GetReceivedCards handles GET /api/cards/received
func (h *CardHandler) GetReceivedCards(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.errorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var params dtos.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_PARAMS", "Invalid pagination parameters", err.Error())
		return
	}

	cards, total, err := h.cardRepo.FindByRecipient(userID, params.Page, params.PageSize)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch received cards", err.Error())
		return
	}

	h.paginatedResponse(c, cards, total, params.Page, params.PageSize)
}

// GetSentCards handles GET /api/cards/sent
func (h *CardHandler) GetSentCards(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.errorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var params dtos.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_PARAMS", "Invalid pagination parameters", err.Error())
		return
	}

	cards, total, err := h.cardRepo.FindBySender(userID, params.Page, params.PageSize)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch sent cards", err.Error())
		return
	}

	h.paginatedResponse(c, cards, total, params.Page, params.PageSize)
}

// Helper methods

func (h *CardHandler) mapCardToResponse(card *models.Card) dtos.CardResponse {
	recipients := make([]string, len(card.Recipients))
	for i, r := range card.Recipients {
		recipients[i] = r.RecipientID
	}

	selectedValues := make([]dtos.ValueResponse, len(card.Values))
	for i, v := range card.Values {
		selectedValues[i] = dtos.ValueResponse{
			ID:          v.ValueID,
			Name:        v.CompanyValue.Name,
			Description: v.CompanyValue.Description,
			Type:        string(v.CompanyValue.Type),
		}
	}

	return dtos.CardResponse{
		ID:                card.ID,
		SenderID:          card.SenderID,
		Recipients:        recipients,
		RecognitionReason: card.RecognitionReason,
		SelectedValues:    selectedValues,
		CreatedAt:         card.CreatedAt,
	}
}

func (h *CardHandler) successResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, dtos.StandardResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

func (h *CardHandler) errorResponse(c *gin.Context, status int, code, message string, details interface{}) {
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

func (h *CardHandler) paginatedResponse(c *gin.Context, cards []*models.Card, total int64, page, pageSize int) {
	cardResponses := make([]dtos.CardResponse, len(cards))
	for i, card := range cards {
		cardResponses[i] = h.mapCardToResponse(card)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, dtos.StandardResponse{
		Success: true,
		Data: dtos.CardListResponse{
			Data: cardResponses,
			Pagination: dtos.PaginationResponse{
				Page:        page,
				PageSize:    pageSize,
				TotalPages:  totalPages,
				TotalItems:  total,
				HasNext:     page < totalPages,
				HasPrevious: page > 1,
			},
		},
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

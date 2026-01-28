package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/middleware"
	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/repositories"
	"thank-you-card-backend/internal/services"
)

type CardHandler struct {
	cardService services.CardService
}

func NewCardHandler(cardService services.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

func getCurrentEmployee(c *gin.Context) *models.Employee {
	// Try to get from context
	v, ok := c.Get(middleware.ContextCurrentEmployee)
	if ok {
		if emp, ok := v.(*models.Employee); ok && emp != nil {
			return emp
		}
	}

	// Fallback: try to get from user
	userVal, ok := c.Get(middleware.ContextCurrentUser)
	if ok {
		if user, ok := userVal.(*models.User); ok && user != nil && user.Employee != nil {
			return user.Employee
		}
	}

	return nil
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp := getCurrentEmployee(c)
	if emp == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "employee not found: user must be linked to an employee"})
		return
	}

	ctx := c.Request.Context()
	card, err := h.cardService.CreateCard(ctx, emp, req.RecipientIDs, req.ValueIDs, req.Reason)
	if err != nil {
		if services.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, toCardResponse(card))
}

func (h *CardHandler) GetFeed(c *gin.Context) {
	ctx := c.Request.Context()
	filters := parseCardFilters(c)

	cards, total, err := h.cardService.GetCompanyFeed(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load feed"})
		return
	}

	resp := make([]dto.CardResponse, 0, len(cards))
	for _, card := range cards {
		cp := card
		resp = append(resp, toCardResponse(&cp))
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": total,
	})
}

func (h *CardHandler) GetMyReceived(c *gin.Context) {
	emp := getCurrentEmployee(c)
	if emp == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "employee not found: user must be linked to an employee"})
		return
	}
	ctx := c.Request.Context()
	filters := parseCardFilters(c)

	cards, total, err := h.cardService.GetPersonalReceived(ctx, emp.ID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load received cards"})
		return
	}

	resp := make([]dto.CardResponse, 0, len(cards))
	for _, card := range cards {
		cp := card
		resp = append(resp, toCardResponse(&cp))
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": total,
	})
}

func (h *CardHandler) GetMySent(c *gin.Context) {
	emp := getCurrentEmployee(c)
	if emp == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "employee not found: user must be linked to an employee"})
		return
	}
	ctx := c.Request.Context()
	filters := parseCardFilters(c)

	cards, total, err := h.cardService.GetPersonalSent(ctx, emp.ID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load sent cards"})
		return
	}

	resp := make([]dto.CardResponse, 0, len(cards))
	for _, card := range cards {
		cp := card
		resp = append(resp, toCardResponse(&cp))
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": total,
	})
}

func (h *CardHandler) React(c *gin.Context) {
	emp := getCurrentEmployee(c)
	if emp == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "employee not found: user must be linked to an employee"})
		return
	}

	cardID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	var body struct {
		EmojiCode string `json:"emojiCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.EmojiCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "emojiCode required"})
		return
	}

	ctx := c.Request.Context()
	if err := h.cardService.ReactToCard(ctx, cardID, emp.ID, body.EmojiCode); err != nil {
		if services.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to react to card"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CardHandler) RemoveReaction(c *gin.Context) {
	emp := getCurrentEmployee(c)
	if emp == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "employee not found: user must be linked to an employee"})
		return
	}

	cardID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	ctx := c.Request.Context()
	if err := h.cardService.RemoveReaction(ctx, cardID, emp.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove reaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseCardFilters(c *gin.Context) repositories.CardFilters {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var fromPtr, toPtr *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			fromPtr = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			toPtr = &t
		}
	}

	// Parse values[] query parameter
	var valueIDs []uint
	if values := c.QueryArray("values[]"); len(values) > 0 {
		valueIDs = make([]uint, 0, len(values))
		for _, v := range values {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				valueIDs = append(valueIDs, uint(id))
			}
		}
	}

	// Parse senderIds[] query parameter
	var senderIDs []uint
	if senders := c.QueryArray("senderIds[]"); len(senders) > 0 {
		senderIDs = make([]uint, 0, len(senders))
		for _, v := range senders {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				senderIDs = append(senderIDs, uint(id))
			}
		}
	}

	// Parse recipientIds[] query parameter
	var recipientIDs []uint
	if recipients := c.QueryArray("recipientIds[]"); len(recipients) > 0 {
		recipientIDs = make([]uint, 0, len(recipients))
		for _, v := range recipients {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				recipientIDs = append(recipientIDs, uint(id))
			}
		}
	}

	return repositories.CardFilters{
		Page:         page,
		PageSize:     size,
		From:         fromPtr,
		To:           toPtr,
		Query:        c.Query("q"),
		ValueIDs:     valueIDs,
		SenderIDs:    senderIDs,
		RecipientIDs: recipientIDs,
	}
}

func parseUintParam(c *gin.Context, name string) (uint, error) {
	raw := c.Param(name)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func toCardResponse(card *models.Card) dto.CardResponse {
	sender := dto.EmployeeSummary{
		ID:         card.Sender.ID,
		Name:       card.Sender.Name,
		Department: card.Sender.Department,
	}

	// Initialize as empty slices to ensure JSON serializes as [] instead of null
	recipients := make([]dto.EmployeeSummary, 0, len(card.Recipients))
	for _, cr := range card.Recipients {
		recipients = append(recipients, dto.EmployeeSummary{
			ID:         cr.Recipient.ID,
			Name:       cr.Recipient.Name,
			Department: cr.Recipient.Department,
		})
	}

	values := make([]dto.CompanyValueSummary, 0, len(card.Values))
	for _, cv := range card.Values {
		values = append(values, dto.CompanyValueSummary{
			ID:   cv.CompanyValue.ID,
			Code: cv.CompanyValue.Code,
			Name: cv.CompanyValue.Name,
			Type: cv.CompanyValue.Type,
		})
	}

	emojiMap := map[string]*dto.EmojiSummary{}
	for _, r := range card.Reactions {
		es, exists := emojiMap[r.EmojiCode]
		if !exists {
			es = &dto.EmojiSummary{
				EmojiCode: r.EmojiCode,
			}
			emojiMap[r.EmojiCode] = es
		}
		es.Count++
		es.UserIDs = append(es.UserIDs, r.UserID)
	}

	// Initialize as empty slice to ensure JSON serializes as [] instead of null
	reactions := make([]dto.EmojiSummary, 0, len(emojiMap))
	for _, v := range emojiMap {
		reactions = append(reactions, *v)
	}

	return dto.CardResponse{
		ID:         card.ID,
		Sender:     sender,
		Recipients: recipients,
		Reason:     card.Reason,
		Values:     values,
		CreatedAt:  card.CreatedAt.Format(time.RFC3339),
		Reactions:  reactions,
	}
}

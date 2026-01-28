package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/services"
)

type AIHandler struct {
	deepseekService services.DeepSeekService
}

func NewAIHandler(deepseekService services.DeepSeekService) *AIHandler {
	return &AIHandler{
		deepseekService: deepseekService,
	}
}

type GenerateTextRequest struct {
	RecipientName string   `json:"recipientName" binding:"required"`
	UserInput     string   `json:"userInput" binding:"required"`
	ValueNames    []string `json:"valueNames"`
}

type GenerateTextResponse struct {
	Text string `json:"text"`
}

// GenerateRecognitionText generates recognition text using AI
func (h *AIHandler) GenerateRecognitionText(c *gin.Context) {
	var req GenerateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if req.RecipientName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipientName is required"})
		return
	}

	if req.UserInput == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userInput is required"})
		return
	}

	// Generate text
	text, err := h.deepseekService.GenerateRecognitionText(
		c.Request.Context(),
		req.RecipientName,
		req.UserInput,
		req.ValueNames,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate text: %v", err)})
		return
	}

	c.JSON(http.StatusOK, GenerateTextResponse{Text: text})
}

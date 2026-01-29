package handlers

import (
	"fmt"
	"log"
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
	RecipientName    string   `json:"recipientName" binding:"required"`
	SenderName       string   `json:"senderName"`
	SenderDepartment string   `json:"senderDepartment"`
	UserInput        string   `json:"userInput" binding:"required"`
	ValueNames       []string `json:"valueNames"`
	Language         string   `json:"language"` // "zh" for Chinese, "en" for English, defaults to "en"
}

type GenerateTextResponse struct {
	Text string `json:"text"`
}

// GenerateRecognitionText generates recognition text using AI
func (h *AIHandler) GenerateRecognitionText(c *gin.Context) {
	log.Printf("INFO: Received AI generation request\n")

	var req GenerateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: Failed to bind request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to English if language is not provided
	if req.Language == "" {
		req.Language = "en"
	}

	log.Printf("INFO: Request data - recipient: %s, sender: %s (%s), input: %s, values: %v, language: %s\n",
		req.RecipientName, req.SenderName, req.SenderDepartment, req.UserInput, req.ValueNames, req.Language)

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
	log.Printf("INFO: Calling DeepSeek service to generate text...\n")
	text, err := h.deepseekService.GenerateRecognitionText(
		c.Request.Context(),
		req.RecipientName,
		req.SenderName,
		req.SenderDepartment,
		req.UserInput,
		req.ValueNames,
		req.Language,
	)
	if err != nil {
		log.Printf("ERROR: Failed to generate text: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate text: %v", err)})
		return
	}

	log.Printf("INFO: Successfully generated text, length: %d\n", len(text))
	c.JSON(http.StatusOK, GenerateTextResponse{Text: text})
}

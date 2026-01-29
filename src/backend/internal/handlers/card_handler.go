package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/castlery/thank-you-card/internal/dtos"
	"github.com/castlery/thank-you-card/internal/models"
	"github.com/castlery/thank-you-card/internal/repositories"
	"github.com/castlery/thank-you-card/internal/services"
	"github.com/castlery/thank-you-card/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CardHandler handles HTTP requests for cards
type CardHandler struct {
	cardCreationService *services.CardCreationService
	cardRepo            repositories.CardRepository
	employeeService     services.EmployeeService
	deepSeekAPIKey      string
	deepSeekBaseURL     string
}

// NewCardHandler creates a new card handler
func NewCardHandler(
	cardCreationService *services.CardCreationService,
	cardRepo repositories.CardRepository,
	employeeService services.EmployeeService,
	deepSeekAPIKey, deepSeekBaseURL string,
) *CardHandler {
	if deepSeekBaseURL == "" {
		deepSeekBaseURL = "https://api.deepseek.com"
	}
	return &CardHandler{
		cardCreationService: cardCreationService,
		cardRepo:            cardRepo,
		employeeService:     employeeService,
		deepSeekAPIKey:      deepSeekAPIKey,
		deepSeekBaseURL:     strings.TrimSuffix(deepSeekBaseURL, "/"),
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

	senderName := middleware.GetUserName(c)

	// Convert string valueIds to UUID
	valueUUIDs := make([]uuid.UUID, len(req.ValueIDs))
	for i, valueIDStr := range req.ValueIDs {
		valueUUID, err := uuid.Parse(valueIDStr)
		if err != nil {
			h.errorResponse(c, http.StatusBadRequest, "INVALID_VALUE_ID", fmt.Sprintf("Invalid value ID format: %s", valueIDStr), err.Error())
			return
		}
		valueUUIDs[i] = valueUUID
	}

	// Map recipients for domain service
	recipients := make([]map[string]string, len(req.Recipients))
	for i, r := range req.Recipients {
		recipients[i] = map[string]string{
			"id":   r.ID,
			"name": r.Name,
		}
	}

	card, err := h.cardCreationService.CreateCard(
		senderID,
		senderName,
		recipients,
		req.RecognitionReason,
		valueUUIDs,
	)
	if err != nil {
		// Business validation errors should be 400, not 500
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "not found") {
			status = http.StatusBadRequest
		}

		h.errorResponse(c, status, "CREATION_FAILED", "Failed to create card", err.Error())
		return
	}

	h.successResponse(c, http.StatusCreated, h.mapCardToResponse(card), "Card created successfully")
}

// deepSeekChatRequest is OpenAI-compatible chat request
type deepSeekChatRequest struct {
	Model     string                `json:"model"`
	Messages  []deepSeekChatMessage `json:"messages"`
	MaxTokens int                   `json:"max_tokens,omitempty"`
}

type deepSeekChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// deepSeekChatResponse is OpenAI-compatible chat response
type deepSeekChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateRecognitionReason handles POST /api/v1/cards/generate-reason
func (h *CardHandler) GenerateRecognitionReason(c *gin.Context) {
	if h.deepSeekAPIKey == "" {
		h.errorResponse(c, http.StatusServiceUnavailable, "DEEPSEEK_DISABLED", "Recognition reason generation is not configured", nil)
		return
	}

	var req dtos.GenerateReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	recipientList := strings.Join(req.RecipientNames, ", ")
	valueList := strings.Join(req.ValueNames, ", ")

	systemPrompt := `你的角色：你是一位资深的人力资源专家，正在推行 thank you card 项目让员工之间表达日常工作支持的感激，你懂得如何将公司文化价值观与员工事例做紧密贴合。`
	userPrompt := fmt.Sprintf(`你的任务：根据以下信息，生成一段用于感谢同事的 thank you card 内容。

发件人：%s
接收人：%s
用户输入的关键信息：%s
用户输入的价值观及信条：%s

生成结果的要求：
- 语言要温暖、真诚、充满感激之情
- 要体现对接收人的认可和赞美（内容必须是发件人感谢接收人，不要搞错发件人和接收人）
- 将用户已经选择的价值观及信条自然地融入文本中
- 长度控制在300-500字之间
- 使用中文或英文（根据用户输入的语言选择）
- 语气要自然、亲切，不要太正式
- 使用 STAR 原则

请直接输出感谢文本，不要包含其他说明文字。`, req.SenderName, recipientList, req.KeyInfo, valueList)

	body := deepSeekChatRequest{
		Model: "deepseek-chat",
		Messages: []deepSeekChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens: 800,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "GENERATE_FAILED", "Failed to build request", err.Error())
		return
	}

	url := h.deepSeekBaseURL + "/v1/chat/completions"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "GENERATE_FAILED", "Failed to create request", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+h.deepSeekAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("DeepSeek request error: %v", err)
		h.errorResponse(c, http.StatusBadGateway, "GENERATE_FAILED", "Failed to call DeepSeek", err.Error())
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "GENERATE_FAILED", "Failed to read response", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("DeepSeek API error: status=%d body=%s", resp.StatusCode, string(respBytes))
		h.errorResponse(c, http.StatusBadGateway, "GENERATE_FAILED", "DeepSeek API returned an error", string(respBytes))
		return
	}

	var chatResp deepSeekChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "GENERATE_FAILED", "Failed to parse DeepSeek response", err.Error())
		return
	}

	if len(chatResp.Choices) == 0 {
		h.errorResponse(c, http.StatusInternalServerError, "GENERATE_FAILED", "No content in DeepSeek response", nil)
		return
	}

	reason := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if len(reason) > 1000 {
		reason = reason[:1000]
	}

	h.successResponse(c, http.StatusOK, dtos.GenerateReasonResponse{RecognitionReason: reason}, "")
}

// SearchEmployees handles GET /api/employees/search
func (h *CardHandler) SearchEmployees(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_QUERY", "Search query is required", nil)
		return
	}

	employees, err := h.employeeService.SearchEmployees(query)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "SEARCH_FAILED", "Failed to search employees", err.Error())
		return
	}

	// Map to DTO
	resp := make([]dtos.EmployeeDTO, len(employees))
	for i, emp := range employees {
		resp[i] = dtos.EmployeeDTO{
			ID:         emp.ID,
			AADID:      emp.AADID,
			Name:       emp.Name,
			Email:      emp.Email,
			Department: emp.Department,
		}
	}

	h.successResponse(c, http.StatusOK, resp, "")
}

// GetCards handles GET /api/cards (company-wide feed)
func (h *CardHandler) GetCards(c *gin.Context) {
	var params dtos.CardFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_PARAMS", "Invalid filtering parameters", err.Error())
		return
	}

	var startDate, endDate *string
	if params.StartDate != "" {
		startDate = &params.StartDate
	}
	if params.EndDate != "" {
		endDate = &params.EndDate
	}

	cards, total, err := h.cardRepo.FindWithFilters(
		params.SenderID,
		params.RecipientID,
		params.ValueIDs,
		startDate,
		endDate,
		params.Search,
		params.Page,
		params.PageSize,
	)
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

	var params dtos.CardFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_PARAMS", "Invalid filtering parameters", err.Error())
		return
	}

	var startDate, endDate *string
	if params.StartDate != "" {
		startDate = &params.StartDate
	}
	if params.EndDate != "" {
		endDate = &params.EndDate
	}

	cards, total, err := h.cardRepo.FindWithFilters(
		params.SenderID,
		userID, // Force current user as recipient
		params.ValueIDs,
		startDate,
		endDate,
		params.Search,
		params.Page,
		params.PageSize,
	)
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

	var params dtos.CardFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "INVALID_PARAMS", "Invalid filtering parameters", err.Error())
		return
	}

	var startDate, endDate *string
	if params.StartDate != "" {
		startDate = &params.StartDate
	}
	if params.EndDate != "" {
		endDate = &params.EndDate
	}

	cards, total, err := h.cardRepo.FindWithFilters(
		userID, // Force current user as sender
		params.RecipientID,
		params.ValueIDs,
		startDate,
		endDate,
		params.Search,
		params.Page,
		params.PageSize,
	)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch sent cards", err.Error())
		return
	}

	h.paginatedResponse(c, cards, total, params.Page, params.PageSize)
}

// Helper methods

func (h *CardHandler) mapCardToResponse(card *models.Card) dtos.CardResponse {
	recipients := make([]dtos.RecipientResponse, len(card.Recipients))
	for i, r := range card.Recipients {
		recipients[i] = dtos.RecipientResponse{
			ID:   r.RecipientID,
			Name: r.RecipientName,
		}
	}

	selectedValues := make([]dtos.ValueResponse, len(card.Values))
	for i, v := range card.Values {
		selectedValues[i] = dtos.ValueResponse{
			ID:          v.ValueID,
			Code:        v.CompanyValue.Code,
			Name:        v.CompanyValue.Name,
			Description: v.CompanyValue.Description,
			Type:        string(v.CompanyValue.Type),
			Examples:    []string(v.CompanyValue.Examples),
		}
	}

	return dtos.CardResponse{
		ID:                card.ID,
		SenderID:          card.SenderID,
		SenderName:        card.SenderName,
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

package handlers

import (
	"net/http"
	"time"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
	"github.com/castlery/thank-you-card/internal/card/domain/repositories"
	"github.com/castlery/thank-you-card/internal/card/dtos"

	"github.com/gin-gonic/gin"
)

// ValueHandler handles HTTP requests for company values
type ValueHandler struct {
	valueRepo repositories.CompanyValueRepository
}

// NewValueHandler creates a new value handler
func NewValueHandler(valueRepo repositories.CompanyValueRepository) *ValueHandler {
	return &ValueHandler{valueRepo: valueRepo}
}

// GetValues handles GET /api/values
func (h *ValueHandler) GetValues(c *gin.Context) {
	values, err := h.valueRepo.FindAll()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", "Failed to fetch company values", err.Error())
		return
	}

	response := make([]dtos.ValueResponse, len(values))
	for i, v := range values {
		response[i] = h.mapValueToResponse(v)
	}

	h.successResponse(c, http.StatusOK, response, "")
}

func (h *ValueHandler) mapValueToResponse(v *models.CompanyValue) dtos.ValueResponse {
	return dtos.ValueResponse{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
		Type:        string(v.Type),
	}
}

func (h *ValueHandler) successResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, dtos.StandardResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC(),
		RequestID: c.GetString("requestId"),
	})
}

func (h *ValueHandler) errorResponse(c *gin.Context, status int, code, message string, details interface{}) {
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

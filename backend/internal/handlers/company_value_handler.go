package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/services"
)

type CompanyValueHandler struct {
	service services.CompanyValueService
}

func NewCompanyValueHandler(service services.CompanyValueService) *CompanyValueHandler {
	return &CompanyValueHandler{service: service}
}

// GetAll returns all company values and credos
// GET /api/company-values
func (h *CompanyValueHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	values, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load company values"})
		return
	}

	resp := make([]dto.CompanyValueResponse, len(values))
	for i, v := range values {
		resp[i] = dto.CompanyValueResponse{
			ID:          v.ID,
			Code:        v.Code,
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
		}
	}

	c.JSON(http.StatusOK, dto.CompanyValueListResponse{
		Items: resp,
		Total: len(resp),
	})
}

// GetByID returns a single company value by ID
// GET /api/company-values/:id
func (h *CompanyValueHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	ctx := c.Request.Context()
	value, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "company value not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "company value not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load company value"})
		return
	}

	resp := dto.CompanyValueResponse{
		ID:          value.ID,
		Code:        value.Code,
		Name:        value.Name,
		Type:        value.Type,
		Description: value.Description,
	}

	c.JSON(http.StatusOK, resp)
}

// GetByType returns company values filtered by type (VALUE or CREDO)
// GET /api/company-values?type=VALUE or ?type=CREDO
func (h *CompanyValueHandler) GetByType(c *gin.Context) {
	valueType := c.Query("type")
	if valueType == "" {
		h.GetAll(c)
		return
	}

	ctx := c.Request.Context()
	values, err := h.service.GetByType(ctx, valueType)
	if err != nil {
		fmt.Printf("ERROR: Failed to get company values by type %s: %v\n", valueType, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("INFO: Found %d company values of type %s\n", len(values), valueType)

	resp := make([]dto.CompanyValueResponse, len(values))
	for i, v := range values {
		resp[i] = dto.CompanyValueResponse{
			ID:          v.ID,
			Code:        v.Code,
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
		}
	}

	c.JSON(http.StatusOK, dto.CompanyValueListResponse{
		Items: resp,
		Total: len(resp),
	})
}

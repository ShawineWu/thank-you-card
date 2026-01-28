package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var employeeSummary *dto.EmployeeDetailSummary
	if result.Employee != nil {
		employeeSummary = &dto.EmployeeDetailSummary{
			ID:         result.Employee.ID,
			Name:       result.Employee.Name,
			Email:      result.Employee.Email,
			Department: result.Employee.Department,
		}
	}

	resp := dto.LoginResponse{
		Token: result.Token,
		User: dto.UserResponse{
			ID:       result.User.ID,
			Username: result.User.Username,
			Role:     result.User.Role,
			Employee: employeeSummary,
		},
	}

	c.JSON(http.StatusOK, resp)
}

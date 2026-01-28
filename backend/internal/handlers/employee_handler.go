package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/repositories"
)

type EmployeeHandler struct {
	employeeRepo repositories.EmployeeRepository
}

func NewEmployeeHandler(employeeRepo repositories.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{
		employeeRepo: employeeRepo,
	}
}

// GetAll returns all employees
// GET /api/employees
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	employees, err := h.employeeRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load employees"})
		return
	}

	resp := make([]dto.EmployeeSummary, len(employees))
	for i, emp := range employees {
		resp[i] = dto.EmployeeSummary{
			ID:         emp.ID,
			Name:       emp.Name,
			Department: emp.Department,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": len(resp),
	})
}

// Search searches employees by name, department, or email
// GET /api/employees/search?q=query
func (h *EmployeeHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Query("q")

	employees, err := h.employeeRepo.Search(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search employees"})
		return
	}

	resp := make([]dto.EmployeeSummary, len(employees))
	for i, emp := range employees {
		resp[i] = dto.EmployeeSummary{
			ID:         emp.ID,
			Name:       emp.Name,
			Department: emp.Department,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": len(resp),
	})
}

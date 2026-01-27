package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
	cardApp "github.com/company/thank-you-card/internal/application/card"
	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/pkg/logger"
)

// EmployeeHandler handles HTTP requests for employee operations
type EmployeeHandler struct {
	employeeService card.EmployeeDirectoryService
}

// NewEmployeeHandler creates a new employee handler
func NewEmployeeHandler(employeeService card.EmployeeDirectoryService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

// SearchEmployees handles GET /api/v1/employees/search
func (h *EmployeeHandler) SearchEmployees(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("SearchEmployees request received, requestID:", requestID)

	// Get search query parameter
	query := c.Query("q")
	if query == "" {
		logger.Error("Search query is required")
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"MISSING_QUERY",
			"Search query parameter 'q' is required",
			requestID,
			nil,
		))
		return
	}

	// Get limit parameter (default: 10, max: 50)
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// Execute search
	employees, err := h.employeeService.SearchEmployees(query, limit)
	if err != nil {
		logger.Error("Failed to search employees:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"SEARCH_EMPLOYEES_FAILED",
			"Failed to search employees: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTOs
	response := cardApp.MapEmployeesToResponse(employees)

	logger.Info("Employees searched successfully, found:", len(response))
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Employees searched successfully",
		requestID,
	))
}

// GetEmployee handles GET /api/v1/employees/:id
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("GetEmployee request received, requestID:", requestID)

	// Get employee ID from path parameter
	employeeID := c.Param("id")
	if employeeID == "" {
		logger.Error("Employee ID is required")
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"MISSING_EMPLOYEE_ID",
			"Employee ID is required",
			requestID,
			nil,
		))
		return
	}

	// Execute query
	employee, err := h.employeeService.GetEmployee(employeeID)
	if err != nil {
		logger.Error("Failed to get employee:", err)
		c.JSON(http.StatusNotFound, cardApp.NewErrorResponse(
			"EMPLOYEE_NOT_FOUND",
			"Employee not found: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTO
	response := cardApp.MapEmployeeToResponse(employee)

	logger.Info("Employee retrieved successfully, ID:", employeeID)
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Employee retrieved successfully",
		requestID,
	))
}

// ValidateEmployees handles POST /api/v1/employees/validate
func (h *EmployeeHandler) ValidateEmployees(c *gin.Context) {
	requestID := getRequestID(c)
	logger.Info("ValidateEmployees request received, requestID:", requestID)

	var request struct {
		EmployeeIDs []string `json:"employee_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request body: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Validate employee IDs
	err := h.employeeService.ValidateEmployeeIDs(request.EmployeeIDs)
	if err != nil {
		logger.Error("Employee validation failed:", err)
		c.JSON(http.StatusBadRequest, cardApp.NewErrorResponse(
			"INVALID_EMPLOYEES",
			"Employee validation failed: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Get employee details
	employees, err := h.employeeService.GetEmployeesByIDs(request.EmployeeIDs)
	if err != nil {
		logger.Error("Failed to get employees:", err)
		c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
			"GET_EMPLOYEES_FAILED",
			"Failed to get employees: "+err.Error(),
			requestID,
			nil,
		))
		return
	}

	// Convert to response DTOs
	response := cardApp.MapEmployeesToResponse(employees)

	logger.Info("Employees validated successfully, count:", len(response))
	c.JSON(http.StatusOK, cardApp.NewApiResponse(
		response,
		"Employees validated successfully",
		requestID,
	))
}
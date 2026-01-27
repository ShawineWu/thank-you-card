package external

import (
	"log"

	"github.com/castlery/thank-you-card/internal/card/domain/services"
)

// MockEmployeeService implements services.EmployeeValidator for development
type MockEmployeeService struct {
	baseURL string
}

// NewMockEmployeeService creates a new mock employee service
func NewMockEmployeeService(baseURL string) services.EmployeeValidator {
	return &MockEmployeeService{
		baseURL: baseURL,
	}
}

// ValidateEmployees validates that all given employee IDs correspond to valid, active employees
func (s *MockEmployeeService) ValidateEmployees(employeeIDs []string) (bool, error) {
	log.Printf("Validating employees against Employee Directory Service (%s): %v", s.baseURL, employeeIDs)

	// In a mock implementation, we'll assume all IDs are valid for now
	// unless they are empty or "invalid_id"
	for _, id := range employeeIDs {
		if id == "" || id == "invalid_id" {
			return false, nil
		}
	}

	return true, nil
}

package external

import (
	"log"
	"strings"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
	"github.com/castlery/thank-you-card/internal/card/domain/services"
)

// MockEmployeeService implements services.EmployeeValidator for development
type MockEmployeeService struct {
	baseURL string
}

// NewMockEmployeeService creates a new mock employee service
func NewMockEmployeeService(baseURL string) services.EmployeeService {
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

// SearchEmployees searches for employees by query
func (s *MockEmployeeService) SearchEmployees(query string) ([]models.Employee, error) {
	log.Printf("Searching employees in Mock Directory (%s): %s", s.baseURL, query)

	allMocks := []models.Employee{
		{ID: "d3dc438c-1de2-49dc-ae8b-8dffb7fcb2cc", AADID: "d3dc438c-1de2-49dc-ae8b-8dffb7fcb2cc", Name: "Youshan Li (SZX)", Email: "youshan.li@castlery.com", Department: "Eng"},
		{ID: "cd7e9fb0-d844-441e-aac4-b8c21d36ef7c", AADID: "cd7e9fb0-d844-441e-aac4-b8c21d36ef7c", Name: "Christopher Li (SZX)", Email: "christopher.li@castlery.com", Department: "Eng"},
		{ID: "02129f90-2224-486b-859f-b60c2da6b826", AADID: "02129f90-2224-486b-859f-b60c2da6b826", Name: "Zehui Lin (SZX)", Email: "zehui.lin@castlery.com", Department: "Eng"},
		{ID: "3ef2d96d-35b0-42ee-8ace-e4e7880f12c0", AADID: "3ef2d96d-35b0-42ee-8ace-e4e7880f12c0", Name: "Shawine Wu (SZX)", Email: "shawine.wu@castlery.com", Department: "Eng"},
		{ID: "f1f7264d-5a6c-4fef-82e9-a122f9a332d0", AADID: "f1f7264d-5a6c-4fef-82e9-a122f9a332d0", Name: "Alex Zu (SZX)", Email: "alex.zu@castlery.com", Department: "Eng"},
	}

	query = strings.ToLower(query)
	results := make([]models.Employee, 0)
	for _, emp := range allMocks {
		if strings.Contains(strings.ToLower(emp.Name), query) || strings.Contains(strings.ToLower(emp.Email), query) {
			results = append(results, emp)
		}
	}

	return results, nil
}

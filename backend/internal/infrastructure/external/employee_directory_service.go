package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/pkg/logger"
)

// EmployeeDirectoryServiceImpl implements the EmployeeDirectoryService interface
type EmployeeDirectoryServiceImpl struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewEmployeeDirectoryService creates a new employee directory service
func NewEmployeeDirectoryService(baseURL, apiKey string) card.EmployeeDirectoryService {
	return &EmployeeDirectoryServiceImpl{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// EmployeeExists checks if an employee exists in the directory
func (s *EmployeeDirectoryServiceImpl) EmployeeExists(employeeID string) bool {
	_, err := s.GetEmployee(employeeID)
	return err == nil
}

// GetEmployee retrieves an employee by ID
func (s *EmployeeDirectoryServiceImpl) GetEmployee(employeeID string) (*card.Employee, error) {
	// For now, return a mock implementation
	// In a real implementation, this would make an HTTP call to the employee directory service
	
	// Mock implementation - simulate some employees
	mockEmployees := map[string]*card.Employee{
		"emp001": {
			ID:         "emp001",
			Name:       "John Doe",
			Email:      "john.doe@company.com",
			Department: "Engineering",
			Team:       "Backend",
		},
		"emp002": {
			ID:         "emp002",
			Name:       "Jane Smith",
			Email:      "jane.smith@company.com",
			Department: "Product",
			Team:       "Design",
		},
		"emp003": {
			ID:         "emp003",
			Name:       "Bob Johnson",
			Email:      "bob.johnson@company.com",
			Department: "Engineering",
			Team:       "Frontend",
		},
		"emp004": {
			ID:         "emp004",
			Name:       "Alice Brown",
			Email:      "alice.brown@company.com",
			Department: "HR",
			Team:       "People Operations",
		},
		"emp005": {
			ID:         "emp005",
			Name:       "Charlie Wilson",
			Email:      "charlie.wilson@company.com",
			Department: "Sales",
			Team:       "Enterprise",
		},
	}

	employee, exists := mockEmployees[employeeID]
	if !exists {
		return nil, fmt.Errorf("employee not found: %s", employeeID)
	}

	logger.Debug("Retrieved employee:", employee.Name)
	return employee, nil
}

// SearchEmployees searches for employees by name or email
func (s *EmployeeDirectoryServiceImpl) SearchEmployees(query string, limit int) ([]card.Employee, error) {
	// Mock implementation - in reality, this would call the external service
	allEmployees := []card.Employee{
		{
			ID:         "emp001",
			Name:       "John Doe",
			Email:      "john.doe@company.com",
			Department: "Engineering",
			Team:       "Backend",
		},
		{
			ID:         "emp002",
			Name:       "Jane Smith",
			Email:      "jane.smith@company.com",
			Department: "Product",
			Team:       "Design",
		},
		{
			ID:         "emp003",
			Name:       "Bob Johnson",
			Email:      "bob.johnson@company.com",
			Department: "Engineering",
			Team:       "Frontend",
		},
		{
			ID:         "emp004",
			Name:       "Alice Brown",
			Email:      "alice.brown@company.com",
			Department: "HR",
			Team:       "People Operations",
		},
		{
			ID:         "emp005",
			Name:       "Charlie Wilson",
			Email:      "charlie.wilson@company.com",
			Department: "Sales",
			Team:       "Enterprise",
		},
	}

	// Simple search implementation
	var results []card.Employee
	queryLower := strings.ToLower(query)
	
	for _, emp := range allEmployees {
		if strings.Contains(strings.ToLower(emp.Name), queryLower) ||
		   strings.Contains(strings.ToLower(emp.Email), queryLower) {
			results = append(results, emp)
			if len(results) >= limit {
				break
			}
		}
	}

	logger.Debug("Search results for query:", query, "found:", len(results))
	return results, nil
}

// GetEmployeesByIDs retrieves multiple employees by their IDs
func (s *EmployeeDirectoryServiceImpl) GetEmployeesByIDs(employeeIDs []string) ([]card.Employee, error) {
	var employees []card.Employee
	
	for _, id := range employeeIDs {
		employee, err := s.GetEmployee(id)
		if err != nil {
			logger.Error("Failed to get employee:", id, err)
			continue // Skip invalid employees rather than failing the entire request
		}
		employees = append(employees, *employee)
	}

	return employees, nil
}

// ValidateEmployeeIDs validates that all provided employee IDs exist
func (s *EmployeeDirectoryServiceImpl) ValidateEmployeeIDs(employeeIDs []string) error {
	var invalidIDs []string
	
	for _, id := range employeeIDs {
		if !s.EmployeeExists(id) {
			invalidIDs = append(invalidIDs, id)
		}
	}

	if len(invalidIDs) > 0 {
		return fmt.Errorf("invalid employee IDs: %s", strings.Join(invalidIDs, ", "))
	}

	return nil
}

// Real HTTP implementation methods (commented out for mock)

// makeHTTPRequest makes an HTTP request to the employee directory service
func (s *EmployeeDirectoryServiceImpl) makeHTTPRequest(endpoint string, params url.Values) ([]byte, error) {
	fullURL := s.baseURL + endpoint
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var response []byte
	_, err = resp.Body.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return response, nil
}

// parseEmployeeResponse parses the employee directory API response
func (s *EmployeeDirectoryServiceImpl) parseEmployeeResponse(data []byte) (*card.Employee, error) {
	var employee card.Employee
	err := json.Unmarshal(data, &employee)
	if err != nil {
		return nil, fmt.Errorf("failed to parse employee response: %w", err)
	}
	return &employee, nil
}

// parseEmployeesResponse parses the employee directory API response for multiple employees
func (s *EmployeeDirectoryServiceImpl) parseEmployeesResponse(data []byte) ([]card.Employee, error) {
	var employees []card.Employee
	err := json.Unmarshal(data, &employees)
	if err != nil {
		return nil, fmt.Errorf("failed to parse employees response: %w", err)
	}
	return employees, nil
}
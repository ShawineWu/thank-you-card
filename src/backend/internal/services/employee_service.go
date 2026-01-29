package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/castlery/thank-you-card/internal/models"
)

// MockEmployeeService implements EmployeeService for development
type MockEmployeeService struct {
	baseURL string
}

// NewMockEmployeeService creates a new mock employee service
func NewMockEmployeeService(baseURL string) EmployeeService {
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

// GraphEmployeeService implements EmployeeService using Microsoft Graph API
type GraphEmployeeService struct {
	tenantID     string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	token        string
	tokenExpiry  time.Time
}

// NewGraphEmployeeService creates a new Graph API based employee service
func NewGraphEmployeeService(tenantID, clientID, clientSecret string) EmployeeService {
	return &GraphEmployeeService{
		tenantID:     tenantID,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *GraphEmployeeService) getAccessToken() (string, error) {
	if s.token != "" && time.Now().Before(s.tokenExpiry) {
		return s.token, nil
	}

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", s.tenantID)
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", s.clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	s.token = result.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)

	return s.token, nil
}

func (s *GraphEmployeeService) ValidateEmployees(employeeIDs []string) (bool, error) {
	// For production, we can verify each ID exists in Graph
	// For now, satisfy the interface
	token, err := s.getAccessToken()
	if err != nil {
		return false, err
	}

	for _, id := range employeeIDs {
		// Use AAD Object ID to fetch user
		reqURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", id)
		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("User %s validation failed: %d", id, resp.StatusCode)
			return false, nil
		}
	}

	return true, nil
}

func (s *GraphEmployeeService) SearchEmployees(query string) ([]models.Employee, error) {
	token, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}

	// Search by display name or mail
	filter := fmt.Sprintf("startswith(displayName,'%s') or startswith(mail,'%s')", query, query)
	reqURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/users?$filter=%s&$select=id,displayName,mail,department,jobTitle", url.QueryEscape(filter))

	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Value []struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			Mail        string `json:"mail"`
			Department  string `json:"department"`
		} `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	employees := make([]models.Employee, 0, len(result.Value))
	for _, v := range result.Value {
		employees = append(employees, models.Employee{
			ID:         v.ID, // Using AAD ID as primary ID for now
			AADID:      v.ID,
			Name:       v.DisplayName,
			Email:      v.Mail,
			Department: v.Department,
		})
	}

	return employees, nil
}

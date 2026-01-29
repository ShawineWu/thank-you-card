package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	mu           sync.RWMutex // Protects token and tokenExpiry
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
	// First, try to read cached token with read lock
	s.mu.RLock()
	if s.token != "" && time.Now().Add(60*time.Second).Before(s.tokenExpiry) {
		token := s.token
		s.mu.RUnlock()
		return token, nil
	}
	s.mu.RUnlock()

	// Need to refresh token, acquire write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have refreshed)
	if s.token != "" && time.Now().Add(60*time.Second).Before(s.tokenExpiry) {
		return s.token, nil
	}

	log.Printf("Requesting new access token for tenant: %s, client: %s", s.tenantID, s.clientID)

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", s.tenantID)
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", s.clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		log.Printf("Token request failed for tenant %s: %v", s.tenantID, err)
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response body for detailed error message
		bodyBytes, _ := io.ReadAll(resp.Body)
		var errorResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.Unmarshal(bodyBytes, &errorResp)

		errMsg := fmt.Sprintf("token request failed with status %d", resp.StatusCode)
		if errorResp.ErrorDescription != "" {
			errMsg = fmt.Sprintf("%s: %s", errMsg, errorResp.ErrorDescription)
		} else if errorResp.Error != "" {
			errMsg = fmt.Sprintf("%s: %s", errMsg, errorResp.Error)
		}

		log.Printf("Token request failed for tenant %s: %s", s.tenantID, errMsg)
		return "", fmt.Errorf(errMsg)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	s.token = result.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	log.Printf("Successfully obtained access token, expires in %d seconds", result.ExpiresIn)

	return s.token, nil
}

func (s *GraphEmployeeService) ValidateEmployees(employeeIDs []string) (bool, error) {
	if len(employeeIDs) == 0 {
		return true, nil
	}

	log.Printf("Validating %d employees against Graph API", len(employeeIDs))

	token, err := s.getAccessToken()
	if err != nil {
		return false, fmt.Errorf("failed to get access token: %w", err)
	}

	for _, id := range employeeIDs {
		if id == "" {
			log.Printf("Empty employee ID provided, validation failed")
			return false, nil
		}

		reqURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s?$select=id", url.PathEscape(id))
		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			return false, fmt.Errorf("failed to create validation request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := s.httpClient.Do(req)
		if err != nil {
			log.Printf("Validation request failed for user %s: %v", id, err)
			return false, fmt.Errorf("validation request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			log.Printf("User %s not found in Azure AD", id)
			return false, nil
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			var errorResp struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			json.Unmarshal(bodyBytes, &errorResp)

			errMsg := fmt.Sprintf("validation failed for user %s with status %d", id, resp.StatusCode)
			if errorResp.Error.Message != "" {
				errMsg = fmt.Sprintf("%s: %s", errMsg, errorResp.Error.Message)
			}

			log.Printf("Graph API validation error: %s", errMsg)
			return false, fmt.Errorf(errMsg)
		}
	}

	log.Printf("All %d employees validated successfully", len(employeeIDs))
	return true, nil
}

func (s *GraphEmployeeService) SearchEmployees(query string) ([]models.Employee, error) {
	log.Printf("Searching employees in Graph API with query: %s", query)

	token, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Build the filter with proper escaping for OData
	// Escape single quotes in query by doubling them
	escapedQuery := strings.ReplaceAll(query, "'", "''")
	filter := fmt.Sprintf("startswith(displayName,'%s') or startswith(mail,'%s')", escapedQuery, escapedQuery)

	// Build URL with proper encoding
	baseURL := "https://graph.microsoft.com/v1.0/users"
	params := url.Values{}
	params.Set("$filter", filter)
	params.Set("$select", "id,displayName,mail,department,jobTitle")
	params.Set("$top", "50") // Limit results

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("ConsistencyLevel", "eventual") // Required for some filter operations

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("Graph API search request failed: %v", err)
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		var errorResp struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		json.Unmarshal(bodyBytes, &errorResp)

		errMsg := fmt.Sprintf("search failed with status %d", resp.StatusCode)
		if errorResp.Error.Message != "" {
			errMsg = fmt.Sprintf("%s: %s", errMsg, errorResp.Error.Message)
		}

		// Handle rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			log.Printf("Graph API rate limited, retry-after: %s", retryAfter)
		}

		log.Printf("Graph API search failed: %s", errMsg)
		return nil, fmt.Errorf(errMsg)
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
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	employees := make([]models.Employee, 0, len(result.Value))
	for _, v := range result.Value {
		employees = append(employees, models.Employee{
			ID:         v.ID,
			AADID:      v.ID,
			Name:       v.DisplayName,
			Email:      v.Mail,
			Department: v.Department,
		})
	}

	log.Printf("Graph API search completed: query=%s, results=%d", query, len(employees))

	return employees, nil
}

// EmployeeServiceConfig holds configuration for creating an EmployeeService
type EmployeeServiceConfig struct {
	AzureTenantID     string
	AzureClientID     string
	AzureClientSecret string
	MockBaseURL       string
}

// NewEmployeeService creates an EmployeeService based on configuration
// If Azure credentials are provided, it returns a GraphEmployeeService
// Otherwise, it falls back to MockEmployeeService
func NewEmployeeService(config EmployeeServiceConfig) EmployeeService {
	if config.AzureTenantID != "" && config.AzureClientID != "" && config.AzureClientSecret != "" {
		log.Printf("Using GraphEmployeeService with Microsoft Graph API (tenant: %s)", config.AzureTenantID)
		return NewGraphEmployeeService(
			config.AzureTenantID,
			config.AzureClientID,
			config.AzureClientSecret,
		)
	}

	log.Println("Using MockEmployeeService (Azure credentials not configured)")
	return NewMockEmployeeService(config.MockBaseURL)
}

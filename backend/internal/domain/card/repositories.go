package card

import (
	"github.com/google/uuid"
)

// CardRepository defines the interface for card data access
type CardRepository interface {
	// Basic CRUD operations
	Create(card *Card, recipients []CardRecipient, values []CardValue) error
	GetByID(id uuid.UUID) (*Card, error)
	Update(card *Card) error
	Delete(id uuid.UUID) error

	// Query operations
	GetCompanyFeed(limit, offset int) ([]Card, int, error)
	GetCardsBySender(senderID string, limit, offset int) ([]Card, int, error)
	GetCardsByRecipient(recipientID string, limit, offset int) ([]Card, int, error)
	SearchCards(filters SearchFilters, limit, offset int) ([]Card, int, error)

	// Statistics operations
	CountCardsBySender(senderID string) (int, error)
	CountCardsByRecipient(recipientID string) (int, error)
	GetValueDistributionBySender(senderID string) ([]ValueDistribution, error)
	GetRecentActivityByEmployee(employeeID string, limit int) ([]RecentActivity, error)
	GetTop10RecognizedEmployees() ([]Top10Employee, error)

	// Card loading with relationships
	LoadCardWithRelationships(cardID uuid.UUID) (*Card, error)
}

// CompanyValueRepository defines the interface for company values data access
type CompanyValueRepository interface {
	// Basic CRUD operations
	GetAll() ([]CompanyValue, error)
	GetByID(id uuid.UUID) (*CompanyValue, error)
	GetByType(valueType string) ([]CompanyValue, error)
	Create(value *CompanyValue) error
	Update(value *CompanyValue) error
	Delete(id uuid.UUID) error

	// Validation operations
	ValidateValueIDs(valueIDs []uuid.UUID) error
	GetValuesByIDs(valueIDs []uuid.UUID) ([]CompanyValue, error)
}

// MilestoneRepository defines the interface for milestone data access
type MilestoneRepository interface {
	// Basic CRUD operations
	Create(milestone *EmployeeMilestone) error
	GetByEmployeeID(employeeID string) ([]EmployeeMilestone, error)
	GetByID(id uuid.UUID) (*EmployeeMilestone, error)

	// Query operations
	MilestoneExists(employeeID, milestoneType string, threshold int) (bool, error)
	GetMilestonesByType(milestoneType string, limit, offset int) ([]EmployeeMilestone, int, error)
	GetRecentMilestones(limit int) ([]EmployeeMilestone, error)

	// Batch operations
	CreateBatch(milestones []EmployeeMilestone) error
}

// EmployeeDirectoryService defines the interface for employee directory integration
type EmployeeDirectoryService interface {
	// Employee validation
	EmployeeExists(employeeID string) bool
	GetEmployee(employeeID string) (*Employee, error)
	SearchEmployees(query string, limit int) ([]Employee, error)
	
	// Bulk operations
	GetEmployeesByIDs(employeeIDs []string) ([]Employee, error)
	ValidateEmployeeIDs(employeeIDs []string) error
}

// SearchFilters defines the structure for card search filters
type SearchFilters struct {
	SenderID     *string     `json:"sender_id,omitempty"`
	RecipientID  *string     `json:"recipient_id,omitempty"`
	ValueIDs     []uuid.UUID `json:"value_ids,omitempty"`
	DateFrom     *string     `json:"date_from,omitempty"`
	DateTo       *string     `json:"date_to,omitempty"`
	SearchText   *string     `json:"search_text,omitempty"`
}

// Employee represents an employee from the directory service
type Employee struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department,omitempty"`
	Team       string `json:"team,omitempty"`
}

// Top10Employee represents an employee in the top 10 recognized list
type Top10Employee struct {
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	CardCount    int    `json:"card_count"`
	Rank         int    `json:"rank"`
}
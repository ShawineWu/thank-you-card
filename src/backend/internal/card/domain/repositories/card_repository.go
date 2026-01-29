package repositories

import (
	"github.com/castlery/thank-you-card/internal/card/domain/models"

	"github.com/google/uuid"
)

// CardRepository defines the interface for card persistence operations
type CardRepository interface {
	// Save persists a new card
	Save(card *models.Card) error

	// FindByID retrieves a card by its ID
	FindByID(id uuid.UUID) (*models.Card, error)

	// FindBySender retrieves cards sent by a specific employee with pagination
	FindBySender(senderID string, page, pageSize int) ([]*models.Card, int64, error)

	// FindByRecipient retrieves cards received by a specific employee with pagination
	FindByRecipient(recipientID string, page, pageSize int) ([]*models.Card, int64, error)

	// FindAll retrieves all cards (company-wide feed) with pagination
	FindAll(page, pageSize int) ([]*models.Card, int64, error)

	// FindWithFilters retrieves cards matching multiple filter criteria
	FindWithFilters(senderID, recipientID string, valueIDs []uuid.UUID, startDate, endDate *string, search string, page, pageSize int) ([]*models.Card, int64, error)

	// CountBySender counts cards sent by an employee
	CountBySender(senderID string) (int64, error)

	// CountByRecipient counts cards received by an employee
	CountByRecipient(recipientID string) (int64, error)

	// GetTopRecipients retrieves the employees who received the most cards
	GetTopRecipients(limit int) ([]TopRecipient, error)

	// GetValueStatsForUser retrieves statistics on values associated with cards received by a user
	GetValueStatsForUser(employeeID string) ([]ValueStat, error)
}

// TopRecipient represents an employee and their received card count
type TopRecipient struct {
	EmployeeID string
	Count      int64
}

// ValueStat represents a company value and its usage count for a user
type ValueStat struct {
	ValueID   uuid.UUID
	ValueName string
	Count     int64
}

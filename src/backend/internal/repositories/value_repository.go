package repositories

import (
	"github.com/castlery/thank-you-card/internal/models"

	"github.com/google/uuid"
)

// CompanyValueRepository defines the interface for company value persistence operations
type CompanyValueRepository interface {
	// FindAll retrieves all company values and credos
	FindAll() ([]*models.CompanyValue, error)

	// FindByID retrieves a specific value by ID
	FindByID(id uuid.UUID) (*models.CompanyValue, error)

	// FindByIDs retrieves multiple values by their IDs
	FindByIDs(ids []uuid.UUID) ([]*models.CompanyValue, error)

	// FindByType retrieves values by type (Value or Credo)
	FindByType(valueType models.ValueType) ([]*models.CompanyValue, error)
}

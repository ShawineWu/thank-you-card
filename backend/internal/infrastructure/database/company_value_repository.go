package database

import (
	"fmt"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CompanyValueRepositoryImpl implements the CompanyValueRepository interface
type CompanyValueRepositoryImpl struct {
	db *gorm.DB
}

// NewCompanyValueRepository creates a new company value repository
func NewCompanyValueRepository(db *gorm.DB) card.CompanyValueRepository {
	return &CompanyValueRepositoryImpl{db: db}
}

// GetAll retrieves all company values
func (r *CompanyValueRepositoryImpl) GetAll() ([]card.CompanyValue, error) {
	var values []card.CompanyValue
	err := r.db.Order("type ASC, name ASC").Find(&values).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all company values: %w", err)
	}
	return values, nil
}

// GetByID retrieves a company value by ID
func (r *CompanyValueRepositoryImpl) GetByID(id uuid.UUID) (*card.CompanyValue, error) {
	var value card.CompanyValue
	err := r.db.First(&value, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company value not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get company value: %w", err)
	}
	return &value, nil
}

// GetByType retrieves company values by type (Value or Credo)
func (r *CompanyValueRepositoryImpl) GetByType(valueType string) ([]card.CompanyValue, error) {
	var values []card.CompanyValue
	err := r.db.Where("type = ?", valueType).Order("name ASC").Find(&values).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get company values by type: %w", err)
	}
	return values, nil
}

// Create creates a new company value
func (r *CompanyValueRepositoryImpl) Create(value *card.CompanyValue) error {
	err := r.db.Create(value).Error
	if err != nil {
		return fmt.Errorf("failed to create company value: %w", err)
	}
	return nil
}

// Update updates a company value
func (r *CompanyValueRepositoryImpl) Update(value *card.CompanyValue) error {
	err := r.db.Save(value).Error
	if err != nil {
		return fmt.Errorf("failed to update company value: %w", err)
	}
	return nil
}

// Delete deletes a company value by ID
func (r *CompanyValueRepositoryImpl) Delete(id uuid.UUID) error {
	result := r.db.Delete(&card.CompanyValue{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete company value: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("company value not found: %s", id)
	}
	return nil
}

// ValidateValueIDs validates that all provided value IDs exist
func (r *CompanyValueRepositoryImpl) ValidateValueIDs(valueIDs []uuid.UUID) error {
	if len(valueIDs) == 0 {
		return fmt.Errorf("no value IDs provided")
	}

	var count int64
	err := r.db.Model(&card.CompanyValue{}).Where("id IN ?", valueIDs).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate value IDs: %w", err)
	}

	if int(count) != len(valueIDs) {
		return fmt.Errorf("some value IDs are invalid: expected %d, found %d", len(valueIDs), count)
	}

	return nil
}

// GetValuesByIDs retrieves company values by their IDs
func (r *CompanyValueRepositoryImpl) GetValuesByIDs(valueIDs []uuid.UUID) ([]card.CompanyValue, error) {
	if len(valueIDs) == 0 {
		return []card.CompanyValue{}, nil
	}

	var values []card.CompanyValue
	err := r.db.Where("id IN ?", valueIDs).Order("type ASC, name ASC").Find(&values).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get values by IDs: %w", err)
	}

	return values, nil
}
package repositories

import (
	"fmt"

	"github.com/castlery/thank-you-card/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormCompanyValueRepository implements CompanyValueRepository using GORM
type GormCompanyValueRepository struct {
	db *gorm.DB
}

// NewGormCompanyValueRepository creates a new GORM-based company value repository
func NewGormCompanyValueRepository(db *gorm.DB) CompanyValueRepository {
	return &GormCompanyValueRepository{db: db}
}

// FindAll retrieves all company values and credos
func (r *GormCompanyValueRepository) FindAll() ([]*models.CompanyValue, error) {
	var values []*models.CompanyValue
	err := r.db.Order("type, name").Find(&values).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find all values: %w", err)
	}

	return values, nil
}

// FindByID retrieves a specific value by ID
func (r *GormCompanyValueRepository) FindByID(id uuid.UUID) (*models.CompanyValue, error) {
	var value models.CompanyValue
	err := r.db.First(&value, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find value: %w", err)
	}

	return &value, nil
}

// FindByIDs retrieves multiple values by their IDs
func (r *GormCompanyValueRepository) FindByIDs(ids []uuid.UUID) ([]*models.CompanyValue, error) {
	var values []*models.CompanyValue
	err := r.db.Where("id IN ?", ids).Find(&values).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find values by IDs: %w", err)
	}

	return values, nil
}

// FindByType retrieves values by type (Value or Credo)
func (r *GormCompanyValueRepository) FindByType(valueType models.ValueType) ([]*models.CompanyValue, error) {
	var values []*models.CompanyValue
	err := r.db.Where("type = ?", valueType).Order("name").Find(&values).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find values by type: %w", err)
	}

	return values, nil
}

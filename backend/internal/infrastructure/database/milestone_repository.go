package database

import (
	"fmt"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MilestoneRepositoryImpl implements the MilestoneRepository interface
type MilestoneRepositoryImpl struct {
	db *gorm.DB
}

// NewMilestoneRepository creates a new milestone repository
func NewMilestoneRepository(db *gorm.DB) card.MilestoneRepository {
	return &MilestoneRepositoryImpl{db: db}
}

// Create creates a new milestone
func (r *MilestoneRepositoryImpl) Create(milestone *card.EmployeeMilestone) error {
	err := r.db.Create(milestone).Error
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}
	return nil
}

// GetByEmployeeID retrieves all milestones for an employee
func (r *MilestoneRepositoryImpl) GetByEmployeeID(employeeID string) ([]card.EmployeeMilestone, error) {
	var milestones []card.EmployeeMilestone
	err := r.db.Where("employee_id = ?", employeeID).
		Order("achieved_at DESC").
		Find(&milestones).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones by employee ID: %w", err)
	}
	return milestones, nil
}

// GetByID retrieves a milestone by ID
func (r *MilestoneRepositoryImpl) GetByID(id uuid.UUID) (*card.EmployeeMilestone, error) {
	var milestone card.EmployeeMilestone
	err := r.db.First(&milestone, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("milestone not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get milestone: %w", err)
	}
	return &milestone, nil
}

// MilestoneExists checks if a specific milestone already exists for an employee
func (r *MilestoneRepositoryImpl) MilestoneExists(employeeID, milestoneType string, threshold int) (bool, error) {
	var count int64
	err := r.db.Model(&card.EmployeeMilestone{}).
		Where("employee_id = ? AND milestone_type = ? AND threshold = ?", employeeID, milestoneType, threshold).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check milestone existence: %w", err)
	}
	return count > 0, nil
}

// GetMilestonesByType retrieves milestones by type with pagination
func (r *MilestoneRepositoryImpl) GetMilestonesByType(milestoneType string, limit, offset int) ([]card.EmployeeMilestone, int, error) {
	var milestones []card.EmployeeMilestone
	var total int64

	// Get total count
	if err := r.db.Model(&card.EmployeeMilestone{}).
		Where("milestone_type = ?", milestoneType).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count milestones by type: %w", err)
	}

	// Get milestones with pagination
	err := r.db.Where("milestone_type = ?", milestoneType).
		Order("achieved_at DESC").
		Limit(limit).Offset(offset).
		Find(&milestones).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get milestones by type: %w", err)
	}

	return milestones, int(total), nil
}

// GetRecentMilestones retrieves the most recent milestones across all employees
func (r *MilestoneRepositoryImpl) GetRecentMilestones(limit int) ([]card.EmployeeMilestone, error) {
	var milestones []card.EmployeeMilestone
	err := r.db.Order("achieved_at DESC").
		Limit(limit).
		Find(&milestones).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get recent milestones: %w", err)
	}
	return milestones, nil
}

// CreateBatch creates multiple milestones in a single transaction
func (r *MilestoneRepositoryImpl) CreateBatch(milestones []card.EmployeeMilestone) error {
	if len(milestones) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, milestone := range milestones {
			if err := tx.Create(&milestone).Error; err != nil {
				return fmt.Errorf("failed to create milestone in batch: %w", err)
			}
		}
		return nil
	})
}
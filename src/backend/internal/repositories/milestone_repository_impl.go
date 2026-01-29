package repositories

import (
	"fmt"

	"github.com/castlery/thank-you-card/internal/models"

	"gorm.io/gorm"
)

// GormMilestoneRepository implements MilestoneRepository using GORM
type GormMilestoneRepository struct {
	db *gorm.DB
}

// NewGormMilestoneRepository creates a new GORM-based milestone repository
func NewGormMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &GormMilestoneRepository{db: db}
}

// Save persists a new milestone achievement
func (r *GormMilestoneRepository) Save(milestone *models.EmployeeMilestone) error {
	if err := r.db.Create(milestone).Error; err != nil {
		return fmt.Errorf("failed to save milestone: %w", err)
	}
	return nil
}

// FindByEmployee retrieves all milestones for an employee
func (r *GormMilestoneRepository) FindByEmployee(employeeID string) ([]*models.EmployeeMilestone, error) {
	var milestones []*models.EmployeeMilestone
	err := r.db.
		Where("employee_id = ?", employeeID).
		Order("achieved_at DESC").
		Find(&milestones).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find milestones: %w", err)
	}

	return milestones, nil
}

// FindByEmployeeAndType retrieves milestones of a specific type for an employee
func (r *GormMilestoneRepository) FindByEmployeeAndType(employeeID string, milestoneType models.MilestoneType) ([]*models.EmployeeMilestone, error) {
	var milestones []*models.EmployeeMilestone
	err := r.db.
		Where("employee_id = ? AND milestone_type = ?", employeeID, milestoneType).
		Order("threshold DESC").
		Find(&milestones).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find milestones by type: %w", err)
	}

	return milestones, nil
}

// Exists checks if a milestone already exists
func (r *GormMilestoneRepository) Exists(employeeID string, milestoneType models.MilestoneType, threshold int) (bool, error) {
	var count int64
	err := r.db.Model(&models.EmployeeMilestone{}).
		Where("employee_id = ? AND milestone_type = ? AND threshold = ?", employeeID, milestoneType, threshold).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check milestone existence: %w", err)
	}

	return count > 0, nil
}

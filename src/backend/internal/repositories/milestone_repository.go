package repositories

import (
	"github.com/castlery/thank-you-card/internal/models"
)

// MilestoneRepository defines the interface for milestone persistence operations
type MilestoneRepository interface {
	// Save persists a new milestone achievement
	Save(milestone *models.EmployeeMilestone) error

	// FindByEmployee retrieves all milestones for an employee
	FindByEmployee(employeeID string) ([]*models.EmployeeMilestone, error)

	// FindByEmployeeAndType retrieves milestones of a specific type for an employee
	FindByEmployeeAndType(employeeID string, milestoneType models.MilestoneType) ([]*models.EmployeeMilestone, error)

	// Exists checks if a milestone already exists
	Exists(employeeID string, milestoneType models.MilestoneType, threshold int) (bool, error)
}

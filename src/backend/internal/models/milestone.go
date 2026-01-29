package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MilestoneType represents the type of milestone
type MilestoneType string

const (
	MilestoneTypeCardsSent     MilestoneType = "CardsSent"
	MilestoneTypeCardsReceived MilestoneType = "CardsReceived"
)

// Standard milestone thresholds
var MilestoneThresholds = []int{1, 5, 10, 25, 50, 100}

// EmployeeMilestone represents a milestone achievement aggregate root
type EmployeeMilestone struct {
	ID            uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EmployeeID    string        `gorm:"type:varchar(255);not null;index;uniqueIndex:idx_employee_milestones_unique,priority:1" json:"employeeId"`
	MilestoneType MilestoneType `gorm:"type:varchar(50);not null;index;uniqueIndex:idx_employee_milestones_unique,priority:2" json:"milestoneType"`
	Threshold     int           `gorm:"not null;uniqueIndex:idx_employee_milestones_unique,priority:3" json:"threshold"`
	AchievedAt    time.Time     `gorm:"not null" json:"achievedAt"`
	Title         string        `gorm:"type:varchar(255);not null" json:"title"`
	Description   string        `gorm:"type:text;not null" json:"description"`
}

// TableName specifies the table name for EmployeeMilestone
func (EmployeeMilestone) TableName() string {
	return "employee_milestones"
}

// NewEmployeeMilestone creates a new milestone with generated title and description
func NewEmployeeMilestone(employeeID string, milestoneType MilestoneType, threshold int) *EmployeeMilestone {
	title, description := generateMilestoneContent(milestoneType, threshold)

	return &EmployeeMilestone{
		EmployeeID:    employeeID,
		MilestoneType: milestoneType,
		Threshold:     threshold,
		AchievedAt:    time.Now().UTC(),
		Title:         title,
		Description:   description,
	}
}

// generateMilestoneContent generates title and description based on type and threshold
func generateMilestoneContent(milestoneType MilestoneType, threshold int) (string, string) {
	var title, description string

	switch milestoneType {
	case MilestoneTypeCardsSent:
		title = fmt.Sprintf("%d Cards Sent", threshold)
		if threshold == 1 {
			description = "Congratulations on sending your first thank you card! You're spreading recognition and appreciation."
		} else {
			description = fmt.Sprintf("Amazing! You've sent %d thank you cards, recognizing and appreciating your colleagues.", threshold)
		}
	case MilestoneTypeCardsReceived:
		title = fmt.Sprintf("%d Cards Received", threshold)
		if threshold == 1 {
			description = "You've received your first recognition card! Your great work is being noticed."
		} else {
			description = fmt.Sprintf("Impressive! You've received %d recognition cards. Keep up the excellent work!", threshold)
		}
	default:
		title = "Milestone Achieved"
		description = "You've reached a new milestone!"
	}

	return title, description
}

// IsValidThreshold checks if the threshold is a standard milestone threshold
func IsValidThreshold(threshold int) bool {
	for _, t := range MilestoneThresholds {
		if t == threshold {
			return true
		}
	}
	return false
}

package models

import (
	"time"

	"github.com/google/uuid"
)

// ValueType represents the type of company value
type ValueType string

const (
	ValueTypeValue ValueType = "Value"
	ValueTypeCredo ValueType = "Credo"
)

// CompanyValue represents a company value or credo aggregate root
type CompanyValue struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique;index" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Type        ValueType `gorm:"type:varchar(50);not null;index" json:"type"`
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName specifies the table name for CompanyValue
func (CompanyValue) TableName() string {
	return "company_values"
}

// IsValid checks if the value type is valid
func (v *CompanyValue) IsValid() bool {
	return v.Type == ValueTypeValue || v.Type == ValueTypeCredo
}

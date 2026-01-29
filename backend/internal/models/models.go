package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account with login credentials.
type User struct {
	ID         uint           `gorm:"primaryKey"`
	Username   string         `gorm:"size:255;uniqueIndex;not null"`
	Password   string         `gorm:"size:255;not null"`                   // bcrypt hashed
	Role       string         `gorm:"size:32;not null;default:'EMPLOYEE'"` // EMPLOYEE, HR, ADMIN
	EmployeeID *uint          `gorm:"index"`                               // Optional link to Employee
	Employee   *Employee      `gorm:"foreignKey:EmployeeID"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Employee represents an employee that can send/receive cards and may be an HR admin.
type Employee struct {
	ID         uint           `gorm:"primaryKey"`
	Email      string         `gorm:"size:255;uniqueIndex"`
	Name       string         `gorm:"size:255;not null"`
	Department string         `gorm:"size:255;not null"`
	IsHRAdmin  bool           `gorm:"default:false;not null"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// CompanyValue represents a company value or credo that a card can be tagged with.
type CompanyValue struct {
	ID          uint           `gorm:"primaryKey"`
	Code        string         `gorm:"size:64;uniqueIndex;not null"`
	Name        string         `gorm:"size:255;not null"`
	Type        string         `gorm:"size:32;not null"` // e.g. "VALUE" or "CREDO"
	Description string         `gorm:"type:text;not null"`
	Examples    string         `gorm:"type:jsonb"` // JSON array of example strings
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Card is the thank-you card entity.
type Card struct {
	ID           uint           `gorm:"primaryKey"`
	SenderID     uint           `gorm:"index;not null"`
	Sender       Employee       `gorm:"foreignKey:SenderID"`
	Reason       string         `gorm:"type:text;not null"`
	TemplateCode *string        `gorm:"size:64"` // reserved for future template features
	CreatedAt    time.Time      `gorm:"not null;index"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Recipients []CardRecipient
	Values     []CardValue
	Reactions  []EmojiReaction
}

// CardRecipient links a card to one recipient (many-to-many).
type CardRecipient struct {
	ID          uint           `gorm:"primaryKey"`
	CardID      uint           `gorm:"index;not null"`
	Card        Card           `gorm:"foreignKey:CardID"`
	RecipientID uint           `gorm:"index;not null"`
	Recipient   Employee       `gorm:"foreignKey:RecipientID"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// CardValue links a card to a company value/credo (many-to-many).
type CardValue struct {
	ID             uint           `gorm:"primaryKey"`
	CardID         uint           `gorm:"index;not null"`
	Card           Card           `gorm:"foreignKey:CardID"`
	CompanyValueID uint           `gorm:"index;not null"`
	CompanyValue   CompanyValue   `gorm:"foreignKey:CompanyValueID"`
	CreatedAt      time.Time      `gorm:"not null"`
	UpdatedAt      time.Time      `gorm:"not null"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// EmojiReaction represents a single user's emoji reaction on a card.
type EmojiReaction struct {
	ID        uint           `gorm:"primaryKey"`
	CardID    uint           `gorm:"index;not null;uniqueIndex:idx_card_user"`
	Card      Card           `gorm:"foreignKey:CardID"`
	UserID    uint           `gorm:"index;not null;uniqueIndex:idx_card_user"`
	User      Employee       `gorm:"foreignKey:UserID"`
	EmojiCode string         `gorm:"size:64;not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Milestone represents a milestone achievement definition.
type Milestone struct {
	ID          uint           `gorm:"primaryKey"`
	Code        string         `gorm:"size:64;uniqueIndex;not null"`
	Name        string         `gorm:"size:255;not null"`
	Description string         `gorm:"type:text;not null"`
	Type        string         `gorm:"size:32;not null"` // SENT, RECEIVED, TOTAL
	Threshold   int            `gorm:"not null"`
	IconURL     *string        `gorm:"size:512"` // Optional icon URL
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// UserMilestone tracks user achievements.
type UserMilestone struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      uint           `gorm:"index;not null;uniqueIndex:idx_user_milestone"`
	User        User           `gorm:"foreignKey:UserID"`
	MilestoneID uint           `gorm:"index;not null;uniqueIndex:idx_user_milestone"`
	Milestone   Milestone      `gorm:"foreignKey:MilestoneID"`
	AchievedAt  time.Time      `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

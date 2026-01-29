package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Card represents a recognition/thank you card aggregate root
type Card struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SenderID          string          `gorm:"type:varchar(255);not null;index" json:"senderId"`
	SenderName        string          `gorm:"type:varchar(255)" json:"senderName"`
	RecognitionReason string          `gorm:"type:text;not null" json:"recognitionReason"`
	Recipients        []CardRecipient `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"recipients"`
	Values            []CardValue     `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE" json:"values"`
	CreatedAt         time.Time       `gorm:"not null;index:idx_cards_created_at,sort:desc" json:"createdAt"`
	UpdatedAt         time.Time       `gorm:"not null" json:"updatedAt"`
}

// TableName specifies the table name for Card
func (Card) TableName() string {
	return "cards"
}

// CardRecipient represents a recipient of a card (junction table)
type CardRecipient struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CardID        uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_card_recipients_unique,priority:1" json:"cardId"`
	RecipientID   string    `gorm:"type:varchar(255);not null;index;uniqueIndex:idx_card_recipients_unique,priority:2" json:"recipientId"`
	RecipientName string    `gorm:"type:varchar(255)" json:"recipientName"`
	CreatedAt     time.Time `gorm:"not null" json:"createdAt"`
}

// TableName specifies the table name for CardRecipient
func (CardRecipient) TableName() string {
	return "card_recipients"
}

// CardValue represents a company value associated with a card (junction table)
type CardValue struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CardID       uuid.UUID    `gorm:"type:uuid;not null;index;uniqueIndex:idx_card_values_unique,priority:1" json:"cardId"`
	ValueID      uuid.UUID    `gorm:"type:uuid;not null;index;uniqueIndex:idx_card_values_unique,priority:2" json:"valueId"`
	CompanyValue CompanyValue `gorm:"foreignKey:ValueID" json:"companyValue"`
	CreatedAt    time.Time    `gorm:"not null" json:"createdAt"`
}

// TableName specifies the table name for CardValue
func (CardValue) TableName() string {
	return "card_values"
}

// Validation errors
var (
	ErrInvalidSender      = errors.New("sender ID cannot be empty")
	ErrNoRecipients       = errors.New("at least one recipient is required")
	ErrSenderAsRecipient  = errors.New("sender cannot be a recipient")
	ErrDuplicateRecipient = errors.New("duplicate recipients found")
	ErrInvalidReason      = errors.New("recognition reason must be between 10 and 1000 characters")
	ErrInvalidValueCount  = errors.New("must select 1-3 company values")
	ErrDuplicateValue     = errors.New("duplicate values found")
)

// Validate validates the Card aggregate
func (c *Card) Validate() error {
	// Validate sender
	if c.SenderID == "" {
		return ErrInvalidSender
	}

	// Validate recipients
	if len(c.Recipients) == 0 {
		return ErrNoRecipients
	}

	// Check for sender as recipient
	recipientIDs := make(map[string]bool)
	for _, recipient := range c.Recipients {
		if recipient.RecipientID == c.SenderID {
			return ErrSenderAsRecipient
		}

		// Check for duplicates
		if recipientIDs[recipient.RecipientID] {
			return ErrDuplicateRecipient
		}
		recipientIDs[recipient.RecipientID] = true
	}

	// Validate recognition reason
	reasonLen := len(c.RecognitionReason)
	if reasonLen < 10 || reasonLen > 1000 {
		return ErrInvalidReason
	}

	// Validate values count
	if len(c.Values) < 1 || len(c.Values) > 3 {
		return ErrInvalidValueCount
	}

	// Check for duplicate values
	valueIDs := make(map[uuid.UUID]bool)
	for _, value := range c.Values {
		if valueIDs[value.ValueID] {
			return ErrDuplicateValue
		}
		valueIDs[value.ValueID] = true
	}

	return nil
}

// NewCard creates a new Card instance with validation
func NewCard(senderID, senderName string, recipients []map[string]string, reason string, valueIDs []uuid.UUID) (*Card, error) {
	card := &Card{
		SenderID:          senderID,
		SenderName:        senderName,
		RecognitionReason: reason,
		Recipients:        make([]CardRecipient, 0, len(recipients)),
		Values:            make([]CardValue, 0, len(valueIDs)),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// Create recipients
	for _, r := range recipients {
		card.Recipients = append(card.Recipients, CardRecipient{
			RecipientID:   r["id"],
			RecipientName: r["name"],
			CreatedAt:     time.Now().UTC(),
		})
	}

	// Create values
	for _, valueID := range valueIDs {
		card.Values = append(card.Values, CardValue{
			ValueID:   valueID,
			CreatedAt: time.Now().UTC(),
		})
	}

	// Validate
	if err := card.Validate(); err != nil {
		return nil, err
	}

	return card, nil
}

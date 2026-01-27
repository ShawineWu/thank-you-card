package card

import (
	"time"
	"github.com/google/uuid"
)

// Card represents the main card aggregate root
type Card struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SenderID  string    `json:"sender_id" gorm:"not null"`
	Reason    string    `json:"reason" gorm:"not null;check:length(reason) >= 10 AND length(reason) <= 1000"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:now()"`
	
	// Relationships
	Recipients []CardRecipient `json:"recipients" gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE"`
	Values     []CardValue     `json:"values" gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE"`
}

// CardRecipient represents a recipient of a card
type CardRecipient struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CardID      uuid.UUID `json:"card_id" gorm:"type:uuid;not null"`
	RecipientID string    `json:"recipient_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// CardValue represents a company value associated with a card
type CardValue struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CardID    uuid.UUID `json:"card_id" gorm:"type:uuid;not null"`
	ValueID   uuid.UUID `json:"value_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// CompanyValue represents company values and credos
type CompanyValue struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null;check:type IN ('Value', 'Credo')"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// EmployeeMilestone represents milestone achievements
type EmployeeMilestone struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EmployeeID    string    `json:"employee_id" gorm:"not null"`
	MilestoneType string    `json:"milestone_type" gorm:"not null;check:milestone_type IN ('CardsSent', 'CardsReceived')"`
	Threshold     int       `json:"threshold" gorm:"not null"`
	AchievedAt    time.Time `json:"achieved_at" gorm:"not null;default:now()"`
	Title         string    `json:"title" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
}

// Value Objects

// Sender represents the sender of a card
type Sender struct {
	EmployeeID string `json:"employee_id"`
}

// NewSender creates a new sender value object
func NewSender(employeeID string) (Sender, error) {
	if err := ValidateEmployeeID(employeeID); err != nil {
		return Sender{}, err
	}
	return Sender{EmployeeID: employeeID}, nil
}

// Recipient represents a recipient of a card
type Recipient struct {
	EmployeeID string `json:"employee_id"`
}

// NewRecipient creates a new recipient value object
func NewRecipient(employeeID string) (Recipient, error) {
	if err := ValidateEmployeeID(employeeID); err != nil {
		return Recipient{}, err
	}
	return Recipient{EmployeeID: employeeID}, nil
}

// RecognitionReason represents the reason for recognition
type RecognitionReason struct {
	Text string `json:"text"`
}

// NewRecognitionReason creates a new recognition reason value object
func NewRecognitionReason(text string) (RecognitionReason, error) {
	if err := ValidateRecognitionReason(text); err != nil {
		return RecognitionReason{}, err
	}
	return RecognitionReason{Text: text}, nil
}

// Domain Events

// DomainEvent represents a base domain event
type DomainEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// CardCreated represents a card creation event
type CardCreated struct {
	DomainEvent
	CardID            uuid.UUID   `json:"card_id"`
	SenderID          string      `json:"sender_id"`
	RecipientIDs      []string    `json:"recipient_ids"`
	RecognitionReason string      `json:"recognition_reason"`
	SelectedValueIDs  []uuid.UUID `json:"selected_value_ids"`
	CreatedAt         time.Time   `json:"created_at"`
}

// NewCardCreated creates a new CardCreated event
func NewCardCreated(
	cardID uuid.UUID,
	senderID string,
	recipientIDs []string,
	reason string,
	valueIDs []uuid.UUID,
) CardCreated {
	return CardCreated{
		DomainEvent: DomainEvent{
			EventID:   uuid.New(),
			EventType: "CardCreated",
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		CardID:            cardID,
		SenderID:          senderID,
		RecipientIDs:      recipientIDs,
		RecognitionReason: reason,
		SelectedValueIDs:  valueIDs,
		CreatedAt:         time.Now(),
	}
}

// CardShared represents a card sharing event
type CardShared struct {
	DomainEvent
	CardID        uuid.UUID `json:"card_id"`
	SharedBy      string    `json:"shared_by"`
	TargetChannel string    `json:"target_channel"`
	SharedAt      time.Time `json:"shared_at"`
}

// NewCardShared creates a new CardShared event
func NewCardShared(cardID uuid.UUID, sharedBy, targetChannel string) CardShared {
	return CardShared{
		DomainEvent: DomainEvent{
			EventID:   uuid.New(),
			EventType: "CardShared",
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		CardID:        cardID,
		SharedBy:      sharedBy,
		TargetChannel: targetChannel,
		SharedAt:      time.Now(),
	}
}

// MilestoneAchieved represents a milestone achievement event
type MilestoneAchieved struct {
	DomainEvent
	EmployeeID    string    `json:"employee_id"`
	MilestoneType string    `json:"milestone_type"`
	Threshold     int       `json:"threshold"`
	AchievedAt    time.Time `json:"achieved_at"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
}

// NewMilestoneAchieved creates a new MilestoneAchieved event
func NewMilestoneAchieved(milestone EmployeeMilestone) MilestoneAchieved {
	return MilestoneAchieved{
		DomainEvent: DomainEvent{
			EventID:   uuid.New(),
			EventType: "MilestoneAchieved",
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		EmployeeID:    milestone.EmployeeID,
		MilestoneType: milestone.MilestoneType,
		Threshold:     milestone.Threshold,
		AchievedAt:    milestone.AchievedAt,
		Title:         milestone.Title,
		Description:   milestone.Description,
	}
}

// Domain Methods

// GetRecipientIDs returns all recipient IDs for the card
func (c *Card) GetRecipientIDs() []string {
	var recipientIDs []string
	for _, recipient := range c.Recipients {
		recipientIDs = append(recipientIDs, recipient.RecipientID)
	}
	return recipientIDs
}

// GetValueIDs returns all value IDs associated with the card
func (c *Card) GetValueIDs() []uuid.UUID {
	var valueIDs []uuid.UUID
	for _, value := range c.Values {
		valueIDs = append(valueIDs, value.ValueID)
	}
	return valueIDs
}

// HasRecipient checks if the card has a specific recipient
func (c *Card) HasRecipient(employeeID string) bool {
	for _, recipient := range c.Recipients {
		if recipient.RecipientID == employeeID {
			return true
		}
	}
	return false
}

// HasValue checks if the card is associated with a specific value
func (c *Card) HasValue(valueID uuid.UUID) bool {
	for _, value := range c.Values {
		if value.ValueID == valueID {
			return true
		}
	}
	return false
}

// IsOwnedBy checks if the card is owned by a specific employee
func (c *Card) IsOwnedBy(employeeID string) bool {
	return c.SenderID == employeeID
}

// CanBeEditedBy checks if the card can be edited by a specific employee
func (c *Card) CanBeEditedBy(employeeID string) bool {
	return CanEditCard(employeeID, c)
}

// CanBeDeletedBy checks if the card can be deleted by a specific employee
func (c *Card) CanBeDeletedBy(employeeID string) bool {
	return CanDeleteCard(employeeID, c)
}

// CanBeViewedBy checks if the card can be viewed by a specific employee
func (c *Card) CanBeViewedBy(employeeID string) bool {
	return CanViewCard(employeeID, c)
}
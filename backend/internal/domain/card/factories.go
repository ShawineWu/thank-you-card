package card

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CardFactory creates Card aggregates with proper validation
type CardFactory struct{}

// NewCardFactory creates a new card factory
func NewCardFactory() *CardFactory {
	return &CardFactory{}
}

// CreateCard creates a new card with validation
func (f *CardFactory) CreateCard(
	senderID string,
	recipientIDs []string,
	reason string,
	valueIDs []uuid.UUID,
) (*Card, error) {
	// Validate inputs
	if senderID == "" {
		return nil, fmt.Errorf("sender ID cannot be empty")
	}

	if len(recipientIDs) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}

	if err := ValidateRecognitionReason(reason); err != nil {
		return nil, fmt.Errorf("invalid recognition reason: %w", err)
	}

	if len(valueIDs) == 0 {
		return nil, fmt.Errorf("at least one company value must be selected")
	}

	// Create the card
	card := &Card{
		ID        : uuid.New(),
		SenderID  : senderID,
		Reason    : reason,
		CreatedAt : time.Now(),
		UpdatedAt : time.Now(),
	}

	return card, nil
}

// MilestoneFactory creates milestone instances
type MilestoneFactory struct{}

// NewMilestoneFactory creates a new milestone factory
func NewMilestoneFactory() *MilestoneFactory {
	return &MilestoneFactory{}
}

// CreateMilestone creates a new milestone with proper title and description
func (f *MilestoneFactory) CreateMilestone(
	employeeID string,
	milestoneType string,
	threshold int,
) EmployeeMilestone {
	title, description := f.generateMilestoneContent(milestoneType, threshold)
	
	return EmployeeMilestone{
		ID            : uuid.New(),
		EmployeeID    : employeeID,
		MilestoneType : milestoneType,
		Threshold     : threshold,
		AchievedAt    : time.Now(),
		Title         : title,
		Description   : description,
	}
}

// generateMilestoneContent generates appropriate title and description for milestones
func (f *MilestoneFactory) generateMilestoneContent(milestoneType string, threshold int) (string, string) {
	switch milestoneType {
	case "CardsSent":
		return f.generateSentCardsContent(threshold)
	case "CardsReceived":
		return f.generateReceivedCardsContent(threshold)
	default:
		return "Milestone Achieved", "You've reached a milestone!"
	}
}

// generateSentCardsContent generates content for cards sent milestones
func (f *MilestoneFactory) generateSentCardsContent(threshold int) (string, string) {
	switch threshold {
	case 1:
		return "First Thank You Card Sent! 🎉", 
			   "Congratulations on sending your first thank you card! You're spreading positivity and recognition in our company."
	case 5:
		return "Recognition Enthusiast 🌟", 
			   "You've sent 5 thank you cards! Your appreciation for others is making a real difference."
	case 10:
		return "Gratitude Ambassador 🏆", 
			   "Amazing! You've sent 10 thank you cards. You're truly an ambassador of gratitude and recognition."
	case 25:
		return "Recognition Champion 🥇", 
			   "Incredible! 25 thank you cards sent. You're a champion at recognizing and appreciating your colleagues."
	case 50:
		return "Appreciation Master 👑", 
			   "Outstanding! You've sent 50 thank you cards. You're a master at spreading appreciation throughout the company."
	case 100:
		return "Recognition Legend 🌟👑", 
			   "Legendary! 100 thank you cards sent. You're a true legend in fostering a culture of recognition and gratitude."
	default:
		return fmt.Sprintf("%d Cards Sent", threshold), 
			   fmt.Sprintf("Congratulations! You've sent %d thank you cards and continue to spread positivity.", threshold)
	}
}

// generateReceivedCardsContent generates content for cards received milestones
func (f *MilestoneFactory) generateReceivedCardsContent(threshold int) (string, string) {
	switch threshold {
	case 1:
		return "First Recognition Received! 🎊", 
			   "Congratulations! You've received your first thank you card. Your contributions are being noticed and appreciated."
	case 5:
		return "Valued Team Member 💎", 
			   "You've received 5 thank you cards! You're clearly a valued member of our team."
	case 10:
		return "Highly Appreciated 🌟", 
			   "Wonderful! 10 thank you cards received. Your work and attitude are highly appreciated by your colleagues."
	case 25:
		return "Team Favorite 🏅", 
			   "Fantastic! 25 thank you cards received. You're clearly a favorite among your teammates."
	case 50:
		return "Company Star ⭐", 
			   "Exceptional! You've received 50 thank you cards. You're truly a star in our company."
	case 100:
		return "Recognition Superstar 🌟⭐", 
			   "Incredible! 100 thank you cards received. You're a superstar who consistently delivers exceptional value."
	default:
		return fmt.Sprintf("%d Cards Received", threshold), 
			   fmt.Sprintf("Amazing! You've received %d thank you cards. Your contributions continue to be recognized and valued.", threshold)
	}
}

// Helper functions for creating specific milestone types

// CreateSentCardsMilestone creates a milestone for cards sent
func CreateSentCardsMilestone(employeeID string, threshold int) EmployeeMilestone {
	factory := NewMilestoneFactory()
	return factory.CreateMilestone(employeeID, "CardsSent", threshold)
}

// CreateReceivedCardsMilestone creates a milestone for cards received
func CreateReceivedCardsMilestone(employeeID string, threshold int) EmployeeMilestone {
	factory := NewMilestoneFactory()
	return factory.CreateMilestone(employeeID, "CardsReceived", threshold)
}
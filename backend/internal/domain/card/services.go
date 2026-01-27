package card

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CardCreationService handles the business logic for creating cards
type CardCreationService struct {
	employeeService EmployeeDirectoryService
}

// NewCardCreationService creates a new card creation service
func NewCardCreationService(employeeService EmployeeDirectoryService) *CardCreationService {
	return &CardCreationService{
		employeeService: employeeService,
	}
}

// CreateCard creates a new card with validation and business rules
func (s *CardCreationService) CreateCard(
	senderID string,
	recipientIDs []string,
	reason string,
	valueIDs []uuid.UUID,
) (*Card, []CardRecipient, []CardValue, error) {
	// Validate sender exists
	if !s.employeeService.EmployeeExists(senderID) {
		return nil, nil, nil, fmt.Errorf("sender employee %s does not exist", senderID)
	}

	// Validate recipients exist and are not the sender
	for _, recipientID := range recipientIDs {
		if recipientID == senderID {
			return nil, nil, nil, fmt.Errorf("cannot send card to yourself")
		}
		if !s.employeeService.EmployeeExists(recipientID) {
			return nil, nil, nil, fmt.Errorf("recipient employee %s does not exist", recipientID)
		}
	}

	// Validate recognition reason
	if err := ValidateRecognitionReason(reason); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid recognition reason: %w", err)
	}

	// Validate at least one recipient
	if len(recipientIDs) == 0 {
		return nil, nil, nil, fmt.Errorf("at least one recipient is required")
	}

	// Validate at least one value is selected
	if len(valueIDs) == 0 {
		return nil, nil, nil, fmt.Errorf("at least one company value must be selected")
	}

	// Create the card
	card := &Card{
		ID        : uuid.New(),
		SenderID  : senderID,
		Reason    : reason,
		CreatedAt : time.Now(),
		UpdatedAt : time.Now(),
	}

	// Create recipients
	var recipients []CardRecipient
	for _, recipientID := range recipientIDs {
		recipient := CardRecipient{
			ID          : uuid.New(),
			CardID      : card.ID,
			RecipientID : recipientID,
			CreatedAt   : time.Now(),
		}
		recipients = append(recipients, recipient)
	}

	// Create card values
	var cardValues []CardValue
	for _, valueID := range valueIDs {
		cardValue := CardValue{
			ID        : uuid.New(),
			CardID    : card.ID,
			ValueID   : valueID,
			CreatedAt : time.Now(),
		}
		cardValues = append(cardValues, cardValue)
	}

	return card, recipients, cardValues, nil
}

// MilestoneTrackingService handles milestone detection and tracking
type MilestoneTrackingService struct {
	cardRepository      CardRepository
	milestoneRepository MilestoneRepository
}

// NewMilestoneTrackingService creates a new milestone tracking service
func NewMilestoneTrackingService(
	cardRepo CardRepository,
	milestoneRepo MilestoneRepository,
) *MilestoneTrackingService {
	return &MilestoneTrackingService{
		cardRepository:      cardRepo,
		milestoneRepository: milestoneRepo,
	}
}

// CheckAndRecordMilestones checks for milestone achievements after a card is created
func (s *MilestoneTrackingService) CheckAndRecordMilestones(cardCreatedEvent CardCreated) ([]EmployeeMilestone, error) {
	var newMilestones []EmployeeMilestone

	// Check milestones for sender (cards sent)
	senderMilestones, err := s.checkSenderMilestones(cardCreatedEvent.SenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check sender milestones: %w", err)
	}
	newMilestones = append(newMilestones, senderMilestones...)

	// Check milestones for each recipient (cards received)
	for _, recipientID := range cardCreatedEvent.RecipientIDs {
		recipientMilestones, err := s.checkRecipientMilestones(recipientID)
		if err != nil {
			return nil, fmt.Errorf("failed to check recipient milestones for %s: %w", recipientID, err)
		}
		newMilestones = append(newMilestones, recipientMilestones...)
	}

	return newMilestones, nil
}

// checkSenderMilestones checks for cards sent milestones
func (s *MilestoneTrackingService) checkSenderMilestones(employeeID string) ([]EmployeeMilestone, error) {
	// Get total cards sent by this employee
	totalSent, err := s.cardRepository.CountCardsBySender(employeeID)
	if err != nil {
		return nil, err
	}

	// Define milestone thresholds for cards sent
	thresholds := []int{1, 5, 10, 25, 50, 100}
	
	var newMilestones []EmployeeMilestone
	for _, threshold := range thresholds {
		if totalSent >= threshold {
			// Check if this milestone already exists
			exists, err := s.milestoneRepository.MilestoneExists(employeeID, "CardsSent", threshold)
			if err != nil {
				return nil, err
			}
			
			if !exists {
				milestone := CreateSentCardsMilestone(employeeID, threshold)
				newMilestones = append(newMilestones, milestone)
			}
		}
	}

	return newMilestones, nil
}

// checkRecipientMilestones checks for cards received milestones
func (s *MilestoneTrackingService) checkRecipientMilestones(employeeID string) ([]EmployeeMilestone, error) {
	// Get total cards received by this employee
	totalReceived, err := s.cardRepository.CountCardsByRecipient(employeeID)
	if err != nil {
		return nil, err
	}

	// Define milestone thresholds for cards received
	thresholds := []int{1, 5, 10, 25, 50, 100}
	
	var newMilestones []EmployeeMilestone
	for _, threshold := range thresholds {
		if totalReceived >= threshold {
			// Check if this milestone already exists
			exists, err := s.milestoneRepository.MilestoneExists(employeeID, "CardsReceived", threshold)
			if err != nil {
				return nil, err
			}
			
			if !exists {
				milestone := CreateReceivedCardsMilestone(employeeID, threshold)
				newMilestones = append(newMilestones, milestone)
			}
		}
	}

	return newMilestones, nil
}

// PersonalStatisticsService calculates personal statistics for employees
type PersonalStatisticsService struct {
	cardRepository CardRepository
}

// NewPersonalStatisticsService creates a new personal statistics service
func NewPersonalStatisticsService(cardRepo CardRepository) *PersonalStatisticsService {
	return &PersonalStatisticsService{
		cardRepository: cardRepo,
	}
}

// CalculatePersonalStatistics calculates comprehensive personal statistics
func (s *PersonalStatisticsService) CalculatePersonalStatistics(employeeID string) (*PersonalStatistics, error) {
	// Get cards sent count
	cardsSent, err := s.cardRepository.CountCardsBySender(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to count cards sent: %w", err)
	}

	// Get cards received count
	cardsReceived, err := s.cardRepository.CountCardsByRecipient(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to count cards received: %w", err)
	}

	// Get value distribution for sent cards
	valueDistribution, err := s.cardRepository.GetValueDistributionBySender(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get value distribution: %w", err)
	}

	// Get recent activity
	recentActivity, err := s.cardRepository.GetRecentActivityByEmployee(employeeID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}

	return &PersonalStatistics{
		CardsSent:         cardsSent,
		CardsReceived:     cardsReceived,
		ValueDistribution: valueDistribution,
		RecentActivity:    recentActivity,
	}, nil
}

// PersonalStatistics represents personal recognition statistics
type PersonalStatistics struct {
	CardsSent         int                   `json:"cards_sent"`
	CardsReceived     int                   `json:"cards_received"`
	ValueDistribution []ValueDistribution   `json:"value_distribution"`
	RecentActivity    []RecentActivity      `json:"recent_activity"`
}

// ValueDistribution represents the distribution of values used
type ValueDistribution struct {
	ValueID    uuid.UUID `json:"value_id"`
	ValueName  string    `json:"value_name"`
	Count      int       `json:"count"`
	Percentage float64   `json:"percentage"`
}

// RecentActivity represents recent card activity
type RecentActivity struct {
	Type       string    `json:"type"` // "sent" or "received"
	CardID     uuid.UUID `json:"card_id"`
	Date       time.Time `json:"date"`
	OtherParty string    `json:"other_party"` // employee ID of the other party
}
package services

import (
	"fmt"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
	"github.com/castlery/thank-you-card/internal/card/domain/repositories"

	"github.com/google/uuid"
)

// EmployeeValidator defines the interface for validating employees
type EmployeeValidator interface {
	ValidateEmployees(employeeIDs []string) (bool, error)
}

// EventPublisher defines the interface for publishing card events
type EventPublisher interface {
	PublishCardCreated(card *models.Card) error
}

// CardCreationService handles the complex card creation process
type CardCreationService struct {
	cardRepo          repositories.CardRepository
	valueRepo         repositories.CompanyValueRepository
	employeeValidator EmployeeValidator
	eventPublisher    EventPublisher
}

// NewCardCreationService creates a new card creation service
func NewCardCreationService(
	cardRepo repositories.CardRepository,
	valueRepo repositories.CompanyValueRepository,
	employeeValidator EmployeeValidator,
	eventPublisher EventPublisher,
) *CardCreationService {
	return &CardCreationService{
		cardRepo:          cardRepo,
		valueRepo:         valueRepo,
		employeeValidator: employeeValidator,
		eventPublisher:    eventPublisher,
	}
}

// CreateCard creates a new card with full validation
func (s *CardCreationService) CreateCard(
	senderID string,
	recipientIDs []string,
	reason string,
	valueIDs []uuid.UUID,
) (*models.Card, error) {
	// Validate employees (sender + recipients)
	allEmployeeIDs := append([]string{senderID}, recipientIDs...)
	valid, err := s.employeeValidator.ValidateEmployees(allEmployeeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to validate employees: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("one or more employees are invalid or inactive")
	}

	// Validate values exist
	values, err := s.valueRepo.FindByIDs(valueIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to validate values: %w", err)
	}
	if len(values) != len(valueIDs) {
		return nil, fmt.Errorf("one or more values not found")
	}

	// Create card (will validate business rules)
	card, err := models.NewCard(senderID, recipientIDs, reason, valueIDs)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save card
	if err := s.cardRepo.Save(card); err != nil {
		return nil, fmt.Errorf("failed to save card: %w", err)
	}

	// Publish event (synchronously as requested)
	if err := s.eventPublisher.PublishCardCreated(card); err != nil {
		// We log the error but don't fail the creation as the card is already saved
		fmt.Printf("Warning: failed to publish card created event: %v\n", err)
	}

	return card, nil
}

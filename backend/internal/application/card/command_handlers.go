package card

import (
	"fmt"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/internal/infrastructure/events"
	"github.com/company/thank-you-card/pkg/logger"
	"github.com/google/uuid"
)

// CreateCardCommandHandler handles card creation commands
type CreateCardCommandHandler struct {
	cardRepository         card.CardRepository
	companyValueRepository card.CompanyValueRepository
	cardCreationService    *card.CardCreationService
	eventPublisher         events.EventPublisher
}

// NewCreateCardCommandHandler creates a new create card command handler
func NewCreateCardCommandHandler(
	cardRepo card.CardRepository,
	valueRepo card.CompanyValueRepository,
	creationService *card.CardCreationService,
	eventPublisher events.EventPublisher,
) *CreateCardCommandHandler {
	return &CreateCardCommandHandler{
		cardRepository:         cardRepo,
		companyValueRepository: valueRepo,
		cardCreationService:    creationService,
		eventPublisher:         eventPublisher,
	}
}

// Handle processes a create card command
func (h *CreateCardCommandHandler) Handle(request CreateCardRequest) (*CardResponse, error) {
	logger.Info("Processing create card command for sender:", request.SenderID)

	// Validate request
	if validationErrors := ValidateCreateCardRequest(request); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Validate company values exist
	if err := h.companyValueRepository.ValidateValueIDs(request.ValueIDs); err != nil {
		return nil, fmt.Errorf("invalid company values: %w", err)
	}

	// Create card using domain service
	cardEntity, recipients, values, err := h.cardCreationService.CreateCard(
		request.SenderID,
		request.RecipientIDs,
		request.Reason,
		request.ValueIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	// Save to repository
	if err := h.cardRepository.Create(cardEntity, recipients, values); err != nil {
		return nil, fmt.Errorf("failed to save card: %w", err)
	}

	// Load the complete card with relationships
	savedCard, err := h.cardRepository.LoadCardWithRelationships(cardEntity.ID)
	if err != nil {
		logger.Error("Failed to load saved card:", err)
		// Continue with the original card entity
		savedCard = cardEntity
		savedCard.Recipients = recipients
		savedCard.Values = values
	}

	// Publish domain event
	cardCreatedEvent := card.NewCardCreated(
		savedCard.ID,
		savedCard.SenderID,
		savedCard.GetRecipientIDs(),
		savedCard.Reason,
		savedCard.GetValueIDs(),
	)

	if err := h.eventPublisher.PublishCardCreated(cardCreatedEvent); err != nil {
		logger.Error("Failed to publish CardCreated event:", err)
		// Don't fail the command, just log the error
	}

	// Convert to response DTO
	response := MapCardToResponse(savedCard)
	
	logger.Info("Successfully created card:", savedCard.ID)
	return &response, nil
}

// ShareCardCommandHandler handles card sharing commands
type ShareCardCommandHandler struct {
	cardRepository card.CardRepository
	eventPublisher events.EventPublisher
}

// NewShareCardCommandHandler creates a new share card command handler
func NewShareCardCommandHandler(
	cardRepo card.CardRepository,
	eventPublisher events.EventPublisher,
) *ShareCardCommandHandler {
	return &ShareCardCommandHandler{
		cardRepository: cardRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle processes a share card command
func (h *ShareCardCommandHandler) Handle(request ShareCardRequest) error {
	logger.Info("Processing share card command for card:", request.CardID)

	// Validate that the card exists
	cardEntity, err := h.cardRepository.GetByID(request.CardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	// Check if the user can share this card (business rule)
	if !cardEntity.CanBeViewedBy(request.SharedBy) {
		return fmt.Errorf("user %s is not authorized to share card %s", request.SharedBy, request.CardID)
	}

	// Publish domain event
	cardSharedEvent := card.NewCardShared(
		request.CardID,
		request.SharedBy,
		request.TargetChannel,
	)

	if err := h.eventPublisher.PublishCardShared(cardSharedEvent); err != nil {
		return fmt.Errorf("failed to publish CardShared event: %w", err)
	}

	logger.Info("Successfully shared card:", request.CardID)
	return nil
}

// DeleteCardCommandHandler handles card deletion commands
type DeleteCardCommandHandler struct {
	cardRepository card.CardRepository
}

// NewDeleteCardCommandHandler creates a new delete card command handler
func NewDeleteCardCommandHandler(cardRepo card.CardRepository) *DeleteCardCommandHandler {
	return &DeleteCardCommandHandler{
		cardRepository: cardRepo,
	}
}

// Handle processes a delete card command
func (h *DeleteCardCommandHandler) Handle(cardID uuid.UUID, userID string) error {
	logger.Info("Processing delete card command for card:", cardID)

	// Get the card to check permissions
	cardEntity, err := h.cardRepository.GetByID(cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	// Check if the user can delete this card
	if !cardEntity.CanBeDeletedBy(userID) {
		return fmt.Errorf("user %s is not authorized to delete card %s", userID, cardID)
	}

	// Delete the card
	if err := h.cardRepository.Delete(cardID); err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	logger.Info("Successfully deleted card:", cardID)
	return nil
}

// Command interfaces for dependency injection

// CreateCardCommand represents the interface for creating cards
type CreateCardCommand interface {
	Handle(request CreateCardRequest) (*CardResponse, error)
}

// ShareCardCommand represents the interface for sharing cards
type ShareCardCommand interface {
	Handle(request ShareCardRequest) error
}

// DeleteCardCommand represents the interface for deleting cards
type DeleteCardCommand interface {
	Handle(cardID uuid.UUID, userID string) error
}

// Ensure handlers implement interfaces
var _ CreateCardCommand = (*CreateCardCommandHandler)(nil)
var _ ShareCardCommand = (*ShareCardCommandHandler)(nil)
var _ DeleteCardCommand = (*DeleteCardCommandHandler)(nil)
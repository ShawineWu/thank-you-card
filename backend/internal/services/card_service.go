package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/repositories"
)

var (
	ErrValidation = errors.New("validation error")
)

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

type CardService interface {
	CreateCard(ctx context.Context, sender *models.Employee, recipients []uint, valueIDs []uint, reason string) (*models.Card, error)
	GetCompanyFeed(ctx context.Context, filters repositories.CardFilters) ([]models.Card, int64, error)
	GetPersonalReceived(ctx context.Context, employeeID uint, filters repositories.CardFilters) ([]models.Card, int64, error)
	GetPersonalSent(ctx context.Context, employeeID uint, filters repositories.CardFilters) ([]models.Card, int64, error)
	ReactToCard(ctx context.Context, cardID uint, employeeID uint, emojiCode string) error
	RemoveReaction(ctx context.Context, cardID uint, employeeID uint) error
}

type cardService struct {
	cardRepo          repositories.CardRepository
	emojiRepo         repositories.EmojiReactionRepository
	teamsNotification TeamsNotificationService
}

func NewCardService(cardRepo repositories.CardRepository, emojiRepo repositories.EmojiReactionRepository, teamsNotification TeamsNotificationService) CardService {
	return &cardService{
		cardRepo:          cardRepo,
		emojiRepo:         emojiRepo,
		teamsNotification: teamsNotification,
	}
}

func (s *cardService) CreateCard(ctx context.Context, sender *models.Employee, recipients []uint, valueIDs []uint, reason string) (*models.Card, error) {
	// Trim whitespace and validate length
	reason = strings.TrimSpace(reason)
	if len(reason) == 0 || len(reason) > 2000 {
		return nil, fmt.Errorf("%w: reason must be 1-2000 chars (current: %d)", ErrValidation, len(reason))
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("%w: at least one recipient required", ErrValidation)
	}
	if len(valueIDs) == 0 || len(valueIDs) > 3 {
		return nil, fmt.Errorf("%w: valueIDs must be between 1 and 3", ErrValidation)
	}

	card := &models.Card{
		SenderID:  sender.ID,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	if err := s.cardRepo.CreateCardWithRelations(ctx, card, recipients, valueIDs); err != nil {
		return nil, err
	}

	// Load the complete card with all relations for Teams notification
	completeCard, err := s.cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		// Log error but don't fail card creation if notification fails
		// Card was successfully created, notification is optional
		return card, nil
	}

	// Send Teams notification asynchronously (non-blocking)
	// We don't wait for the notification to complete to avoid blocking the response
	go func() {
		// Use background context for async notification
		notificationCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.teamsNotification.SendCardNotification(notificationCtx, completeCard); err != nil {
			// Log error but don't fail the request
			log.Printf("ERROR: Failed to send Teams notification for card ID %d: %v", completeCard.ID, err)
		} else {
			log.Printf("SUCCESS: Sent Teams notification for card ID: %d", completeCard.ID)
		}
	}()

	return completeCard, nil
}

func (s *cardService) GetCompanyFeed(ctx context.Context, filters repositories.CardFilters) ([]models.Card, int64, error) {
	return s.cardRepo.GetFeed(ctx, filters)
}

func (s *cardService) GetPersonalReceived(ctx context.Context, employeeID uint, filters repositories.CardFilters) ([]models.Card, int64, error) {
	return s.cardRepo.GetByRecipient(ctx, employeeID, filters)
}

func (s *cardService) GetPersonalSent(ctx context.Context, employeeID uint, filters repositories.CardFilters) ([]models.Card, int64, error) {
	return s.cardRepo.GetBySender(ctx, employeeID, filters)
}

func (s *cardService) ReactToCard(ctx context.Context, cardID uint, employeeID uint, emojiCode string) error {
	if emojiCode == "" {
		return fmt.Errorf("%w: emojiCode required", ErrValidation)
	}
	return s.emojiRepo.SetReaction(ctx, cardID, employeeID, emojiCode)
}

func (s *cardService) RemoveReaction(ctx context.Context, cardID uint, employeeID uint) error {
	return s.emojiRepo.RemoveReaction(ctx, cardID, employeeID)
}

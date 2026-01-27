package services

import (
	"context"
	"errors"
	"fmt"
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
	cardRepo  repositories.CardRepository
	emojiRepo repositories.EmojiReactionRepository
}

func NewCardService(cardRepo repositories.CardRepository, emojiRepo repositories.EmojiReactionRepository) CardService {
	return &cardService{
		cardRepo:  cardRepo,
		emojiRepo: emojiRepo,
	}
}

func (s *cardService) CreateCard(ctx context.Context, sender *models.Employee, recipients []uint, valueIDs []uint, reason string) (*models.Card, error) {
	if len(reason) == 0 || len(reason) > 2000 {
		return nil, fmt.Errorf("%w: reason must be 1-2000 chars", ErrValidation)
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

	// For now, we skip actual Teams HTTP push.
	// Future: inject NotificationService to send card to Teams channel.

	return card, nil
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

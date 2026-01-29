package services

import (
	"context"
	"time"

	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/repositories"

	"gorm.io/gorm"
)

type MilestoneService interface {
	GetAllMilestones(ctx context.Context) ([]Milestone, error)
	GetUserAchievements(ctx context.Context, userID uint) ([]UserAchievement, error)
	CheckAndAwardMilestones(ctx context.Context, userID uint, sentCount int64, receivedCount int64) error
}

type Milestone struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Type        string
	Threshold   int
	IconURL     *string
}

type UserAchievement struct {
	Milestone  Milestone
	AchievedAt time.Time
}

type milestoneService struct {
	milestoneRepo repositories.MilestoneRepository
	cardRepo      repositories.CardRepository
	userRepo      repositories.UserRepository
}

func NewMilestoneService(milestoneRepo repositories.MilestoneRepository, cardRepo repositories.CardRepository, userRepo repositories.UserRepository) MilestoneService {
	return &milestoneService{
		milestoneRepo: milestoneRepo,
		cardRepo:      cardRepo,
		userRepo:      userRepo,
	}
}

func (s *milestoneService) GetAllMilestones(ctx context.Context) ([]Milestone, error) {
	milestones, err := s.milestoneRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Milestone, len(milestones))
	for i, m := range milestones {
		result[i] = Milestone{
			ID:          m.ID,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			Type:        m.Type,
			Threshold:   m.Threshold,
			IconURL:     m.IconURL,
		}
	}

	return result, nil
}

func (s *milestoneService) GetUserAchievements(ctx context.Context, userID uint) ([]UserAchievement, error) {
	// Get user to find employeeID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.EmployeeID == nil {
		// User has no linked employee, return empty achievements
		return []UserAchievement{}, nil
	}

	employeeID := *user.EmployeeID

	// First, check and award any new milestones
	// Get user's current stats
	_, sentTotal, err := s.cardRepo.GetBySender(ctx, employeeID, repositories.CardFilters{Page: 1, PageSize: 1})
	if err != nil {
		// If user has no sent cards, sentTotal will be 0, which is fine
		sentTotal = 0
	}

	_, receivedTotal, err := s.cardRepo.GetByRecipient(ctx, employeeID, repositories.CardFilters{Page: 1, PageSize: 1})
	if err != nil {
		// If user has no received cards, receivedTotal will be 0, which is fine
		receivedTotal = 0
	}

	// Check and award milestones
	if err := s.CheckAndAwardMilestones(ctx, userID, sentTotal, receivedTotal); err != nil {
		return nil, err
	}

	// Get all user milestones
	userMilestones, err := s.milestoneRepo.GetUserMilestones(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]UserAchievement, len(userMilestones))
	for i, um := range userMilestones {
		result[i] = UserAchievement{
			Milestone: Milestone{
				ID:          um.Milestone.ID,
				Code:        um.Milestone.Code,
				Name:        um.Milestone.Name,
				Description: um.Milestone.Description,
				Type:        um.Milestone.Type,
				Threshold:   um.Milestone.Threshold,
				IconURL:     um.Milestone.IconURL,
			},
			AchievedAt: um.AchievedAt,
		}
	}

	return result, nil
}

func (s *milestoneService) CheckAndAwardMilestones(ctx context.Context, userID uint, sentCount int64, receivedCount int64) error {
	// Note: This method is called from GetUserAchievements which already validated user has employee
	// Get all milestones
	milestones, err := s.milestoneRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	totalCount := sentCount + receivedCount

	// Check each milestone
	for _, milestone := range milestones {
		var shouldAward bool

		switch milestone.Type {
		case "SENT":
			shouldAward = sentCount >= int64(milestone.Threshold)
		case "RECEIVED":
			shouldAward = receivedCount >= int64(milestone.Threshold)
		case "TOTAL":
			shouldAward = totalCount >= int64(milestone.Threshold)
		default:
			continue
		}

		if shouldAward {
			// Check if user already has this milestone
			existing, err := s.milestoneRepo.GetUserMilestoneByMilestoneID(ctx, userID, milestone.ID)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			// If not awarded yet, create it
			if existing == nil || existing.ID == 0 {
				userMilestone := &models.UserMilestone{
					UserID:      userID,
					MilestoneID: milestone.ID,
					AchievedAt:  time.Now(),
				}
				if err := s.milestoneRepo.CreateUserMilestone(ctx, userMilestone); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

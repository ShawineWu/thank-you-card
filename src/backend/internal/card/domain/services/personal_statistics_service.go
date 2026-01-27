package services

import (
	"fmt"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
	"github.com/castlery/thank-you-card/internal/card/domain/repositories"
)

// PersonalStatisticsService handles the calculation of recognition statistics
type PersonalStatisticsService struct {
	cardRepo      repositories.CardRepository
	milestoneRepo repositories.MilestoneRepository
}

// NewPersonalStatisticsService creates a new personal statistics service
func NewPersonalStatisticsService(
	cardRepo repositories.CardRepository,
	milestoneRepo repositories.MilestoneRepository,
) *PersonalStatisticsService {
	return &PersonalStatisticsService{
		cardRepo:      cardRepo,
		milestoneRepo: milestoneRepo,
	}
}

// GetUserStats calculates statistics for a specific employee
func (s *PersonalStatisticsService) GetUserStats(employeeID string) (int64, int64, []repositories.ValueStat, []*models.EmployeeMilestone, error) {
	// Count cards sent
	sentCount, err := s.cardRepo.CountBySender(employeeID)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("failed to get sent count: %w", err)
	}

	// Count cards received
	receivedCount, err := s.cardRepo.CountByRecipient(employeeID)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("failed to get received count: %w", err)
	}

	// Get value usage stats
	valueStats, err := s.cardRepo.GetValueStatsForUser(employeeID)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("failed to get value stats: %w", err)
	}

	// Get milestones
	milestones, err := s.milestoneRepo.FindByEmployee(employeeID)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("failed to get milestones: %w", err)
	}

	return sentCount, receivedCount, valueStats, milestones, nil
}

// GetTopRecipients retrieves the top N recognized employees
func (s *PersonalStatisticsService) GetTopRecipients(limit int) ([]repositories.TopRecipient, error) {
	return s.cardRepo.GetTopRecipients(limit)
}

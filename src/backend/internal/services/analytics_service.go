package services

import (
	"fmt"

	"github.com/castlery/thank-you-card/internal/dtos"
	"github.com/castlery/thank-you-card/internal/repositories"
)

// AnalyticsService provides business logic for analytics operations
type AnalyticsService struct {
	analyticsRepo repositories.AnalyticsRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(analyticsRepo repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
	}
}

// GetDashboardAnalytics returns aggregated data for HR dashboard
func (s *AnalyticsService) GetDashboardAnalytics() (*dtos.DashboardResponse, error) {
	// Get total cards
	totalCards, err := s.analyticsRepo.GetTotalCards()
	if err != nil {
		return nil, fmt.Errorf("failed to get total cards: %w", err)
	}

	// Get active users
	activeUsers, err := s.analyticsRepo.GetActiveUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	// Get top values
	topValuesData, err := s.analyticsRepo.GetTopValues(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top values: %w", err)
	}

	topValues := make([]dtos.ValueUsageResponse, len(topValuesData))
	for i, v := range topValuesData {
		topValues[i] = dtos.ValueUsageResponse{
			ValueID:   v.ValueID.String(),
			ValueName: v.ValueName,
			Count:     v.Count,
		}
	}

	// Get cards trend
	thisMonth, lastMonth, err := s.analyticsRepo.GetCardsTrend()
	if err != nil {
		return nil, fmt.Errorf("failed to get cards trend: %w", err)
	}

	// Calculate engagement rate (simplified: active users / total cards ratio)
	engagementRate := 0.0
	if totalCards > 0 {
		engagementRate = float64(activeUsers) / float64(totalCards) * 100
		if engagementRate > 100 {
			engagementRate = 100
		}
	}

	return &dtos.DashboardResponse{
		TotalCards:  totalCards,
		ActiveUsers: activeUsers,
		TopValues:   topValues,
		CardsTrend: dtos.TrendResponse{
			ThisMonth: thisMonth,
			LastMonth: lastMonth,
		},
		EngagementRate: engagementRate,
	}, nil
}

// GetTopRecognizers returns the most active recognizers
func (s *AnalyticsService) GetTopRecognizers(limit int) ([]dtos.TopRecognizerResponse, error) {
	data, err := s.analyticsRepo.GetTopRecognizers(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top recognizers: %w", err)
	}

	results := make([]dtos.TopRecognizerResponse, len(data))
	for i, r := range data {
		results[i] = dtos.TopRecognizerResponse{
			EmployeeID:   r.EmployeeID,
			EmployeeName: r.EmployeeName,
			CardsSent:    r.CardsSent,
			Rank:         i + 1,
		}
	}

	return results, nil
}

// GetTeamAnalytics returns team-level recognition statistics
func (s *AnalyticsService) GetTeamAnalytics() ([]dtos.TeamAnalyticsResponse, error) {
	data, err := s.analyticsRepo.GetTeamAnalytics()
	if err != nil {
		return nil, fmt.Errorf("failed to get team analytics: %w", err)
	}

	results := make([]dtos.TeamAnalyticsResponse, len(data))
	for i, t := range data {
		results[i] = dtos.TeamAnalyticsResponse{
			TeamID:        t.TeamID,
			TeamName:      t.TeamName,
			CardsSent:     t.CardsSent,
			CardsReceived: t.CardsReceived,
			TopValue:      t.TopValue,
		}
	}

	return results, nil
}

// GetValueDistribution returns the distribution of company values
func (s *AnalyticsService) GetValueDistribution() (*dtos.ValueDistributionResponse, error) {
	data, err := s.analyticsRepo.GetValueDistribution()
	if err != nil {
		return nil, fmt.Errorf("failed to get value distribution: %w", err)
	}

	// Calculate total for percentages
	var total int64
	for _, v := range data {
		total += v.Count
	}

	distribution := make([]dtos.ValueUsageResponse, len(data))
	for i, v := range data {
		percentage := 0.0
		if total > 0 {
			percentage = float64(v.Count) / float64(total) * 100
		}

		distribution[i] = dtos.ValueUsageResponse{
			ValueID:    v.ValueID.String(),
			ValueName:  v.ValueName,
			Count:      v.Count,
			Percentage: percentage,
		}
	}

	return &dtos.ValueDistributionResponse{
		Distribution: distribution,
		TotalCards:   total,
	}, nil
}

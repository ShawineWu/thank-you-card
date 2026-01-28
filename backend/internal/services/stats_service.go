package services

import (
	"context"

	"thank-you-card-backend/internal/repositories"
)

type StatsService interface {
	GetPersonalStats(ctx context.Context, employeeID uint) (*PersonalStats, error)
}

type PersonalStats struct {
	TotalSent         int64
	TotalReceived     int64
	TopSentValues     []ValueCount
	TopReceivedValues []ValueCount
}

type ValueCount struct {
	ValueID uint
	Code    string
	Name    string
	Type    string
	Count   int64
}

type statsService struct {
	cardRepo repositories.CardRepository
}

func NewStatsService(cardRepo repositories.CardRepository) StatsService {
	return &statsService{cardRepo: cardRepo}
}

func (s *statsService) GetPersonalStats(ctx context.Context, employeeID uint) (*PersonalStats, error) {
	// Get sent count
	_, sentTotal, err := s.cardRepo.GetBySender(ctx, employeeID, repositories.CardFilters{Page: 1, PageSize: 1})
	if err != nil {
		return nil, err
	}

	// Get received count
	_, receivedTotal, err := s.cardRepo.GetByRecipient(ctx, employeeID, repositories.CardFilters{Page: 1, PageSize: 1})
	if err != nil {
		return nil, err
	}

	// Get all sent cards to calculate value distribution
	allSentCards, _, err := s.cardRepo.GetBySender(ctx, employeeID, repositories.CardFilters{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}

	// Get all received cards to calculate value distribution
	allReceivedCards, _, err := s.cardRepo.GetByRecipient(ctx, employeeID, repositories.CardFilters{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}

	// Calculate top sent values
	sentValueMap := make(map[uint]*ValueCount)
	for _, card := range allSentCards {
		for _, cv := range card.Values {
			vc, exists := sentValueMap[cv.CompanyValueID]
			if !exists {
				vc = &ValueCount{
					ValueID: cv.CompanyValueID,
					Code:    cv.CompanyValue.Code,
					Name:    cv.CompanyValue.Name,
					Type:    cv.CompanyValue.Type,
				}
				sentValueMap[cv.CompanyValueID] = vc
			}
			vc.Count++
		}
	}

	// Calculate top received values
	receivedValueMap := make(map[uint]*ValueCount)
	for _, card := range allReceivedCards {
		for _, cv := range card.Values {
			vc, exists := receivedValueMap[cv.CompanyValueID]
			if !exists {
				vc = &ValueCount{
					ValueID: cv.CompanyValueID,
					Code:    cv.CompanyValue.Code,
					Name:    cv.CompanyValue.Name,
					Type:    cv.CompanyValue.Type,
				}
				receivedValueMap[cv.CompanyValueID] = vc
			}
			vc.Count++
		}
	}

	// Convert maps to slices and sort
	topSentValues := make([]ValueCount, 0, len(sentValueMap))
	for _, vc := range sentValueMap {
		topSentValues = append(topSentValues, *vc)
	}
	// Sort by count descending (simple bubble sort for small datasets)
	for i := 0; i < len(topSentValues)-1; i++ {
		for j := i + 1; j < len(topSentValues); j++ {
			if topSentValues[i].Count < topSentValues[j].Count {
				topSentValues[i], topSentValues[j] = topSentValues[j], topSentValues[i]
			}
		}
	}
	if len(topSentValues) > 5 {
		topSentValues = topSentValues[:5]
	}

	topReceivedValues := make([]ValueCount, 0, len(receivedValueMap))
	for _, vc := range receivedValueMap {
		topReceivedValues = append(topReceivedValues, *vc)
	}
	// Sort by count descending
	for i := 0; i < len(topReceivedValues)-1; i++ {
		for j := i + 1; j < len(topReceivedValues); j++ {
			if topReceivedValues[i].Count < topReceivedValues[j].Count {
				topReceivedValues[i], topReceivedValues[j] = topReceivedValues[j], topReceivedValues[i]
			}
		}
	}
	if len(topReceivedValues) > 5 {
		topReceivedValues = topReceivedValues[:5]
	}

	return &PersonalStats{
		TotalSent:         sentTotal,
		TotalReceived:     receivedTotal,
		TopSentValues:     topSentValues,
		TopReceivedValues: topReceivedValues,
	}, nil
}

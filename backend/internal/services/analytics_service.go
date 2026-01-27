package services

import (
	"context"
	"fmt"
	"time"

	"thank-you-card-backend/internal/repositories"
)

type AnalyticsService interface {
	GetDashboardOverview(ctx context.Context, from, to time.Time) (*DashboardOverview, error)
	GetTopRecognizedEmployees(ctx context.Context, from, to time.Time, limit int) ([]repositories.TopEmployee, error)
	GetMostActiveRecognizers(ctx context.Context, from, to time.Time, limit int) ([]repositories.TopEmployee, error)
	GetTeamRecognitionPatterns(ctx context.Context, from, to time.Time) ([]repositories.TeamPattern, error)
	GetValuesDistribution(ctx context.Context, from, to time.Time) ([]repositories.ValueDistributionItem, error)
	ExportCards(ctx context.Context, filters repositories.CardFilters) ([]CardExportRow, error)
}

type DashboardOverview struct {
	TotalCards        int64
	TotalEmployees    int64
	TotalDepartments  int64
	ActiveRecognizers int64
	Period            string
}

type CardExportRow struct {
	CardID         uint
	SenderName     string
	SenderEmail    string
	RecipientName  string
	RecipientEmail string
	Reason         string
	Values         string
	Credos         string
	CreatedAt      time.Time
	ReactionCounts string
}

type analyticsService struct {
	analyticsRepo repositories.AnalyticsRepository
	cardRepo      repositories.CardRepository
	employeeRepo  repositories.EmployeeRepository
}

func NewAnalyticsService(
	analyticsRepo repositories.AnalyticsRepository,
	cardRepo repositories.CardRepository,
	employeeRepo repositories.EmployeeRepository,
) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
		cardRepo:      cardRepo,
		employeeRepo:  employeeRepo,
	}
}

func (s *analyticsService) GetDashboardOverview(ctx context.Context, from, to time.Time) (*DashboardOverview, error) {
	// Get total cards in period
	_, totalCards, err := s.cardRepo.GetFeed(ctx, repositories.CardFilters{
		From:     &from,
		To:       &to,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, err
	}

	// Get active recognizers (distinct senders)
	recognizers, err := s.analyticsRepo.MostActiveRecognizers(ctx, from, to, 1000)
	if err != nil {
		return nil, err
	}

	// Get unique departments from team patterns
	patterns, err := s.analyticsRepo.TeamRecognitionPatterns(ctx, from, to)
	if err != nil {
		return nil, err
	}

	period := from.Format("2006-01-02") + " to " + to.Format("2006-01-02")

	return &DashboardOverview{
		TotalCards:        totalCards,
		ActiveRecognizers: int64(len(recognizers)),
		TotalDepartments:  int64(len(patterns)),
		Period:            period,
	}, nil
}

func (s *analyticsService) GetTopRecognizedEmployees(ctx context.Context, from, to time.Time, limit int) ([]repositories.TopEmployee, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.analyticsRepo.TopRecognizedEmployees(ctx, from, to, limit)
}

func (s *analyticsService) GetMostActiveRecognizers(ctx context.Context, from, to time.Time, limit int) ([]repositories.TopEmployee, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.analyticsRepo.MostActiveRecognizers(ctx, from, to, limit)
}

func (s *analyticsService) GetTeamRecognitionPatterns(ctx context.Context, from, to time.Time) ([]repositories.TeamPattern, error) {
	return s.analyticsRepo.TeamRecognitionPatterns(ctx, from, to)
}

func (s *analyticsService) GetValuesDistribution(ctx context.Context, from, to time.Time) ([]repositories.ValueDistributionItem, error) {
	return s.analyticsRepo.ValuesDistribution(ctx, from, to)
}

func (s *analyticsService) ExportCards(ctx context.Context, filters repositories.CardFilters) ([]CardExportRow, error) {
	// Get all cards matching filters (no pagination for export)
	filters.Page = 1
	filters.PageSize = 10000 // Large limit for export
	cards, _, err := s.cardRepo.FilterAndSearch(ctx, filters)
	if err != nil {
		return nil, err
	}

	rows := make([]CardExportRow, 0, len(cards))
	for _, card := range cards {
		for _, recipient := range card.Recipients {
			// Collect values and credos
			var values, credos []string
			for _, cv := range card.Values {
				if cv.CompanyValue.Type == "VALUE" {
					values = append(values, cv.CompanyValue.Name)
				} else {
					credos = append(credos, cv.CompanyValue.Name)
				}
			}

			// Collect reaction counts
			reactionMap := make(map[string]int)
			for _, r := range card.Reactions {
				reactionMap[r.EmojiCode]++
			}
			var reactionStrs []string
			for emoji, count := range reactionMap {
				reactionStrs = append(reactionStrs, emoji+":"+fmt.Sprintf("%d", count))
			}

			rows = append(rows, CardExportRow{
				CardID:         card.ID,
				SenderName:     card.Sender.Name,
				SenderEmail:    card.Sender.Email,
				RecipientName:  recipient.Recipient.Name,
				RecipientEmail: recipient.Recipient.Email,
				Reason:         card.Reason,
				Values:         joinStrings(values, "; "),
				Credos:         joinStrings(credos, "; "),
				CreatedAt:      card.CreatedAt,
				ReactionCounts: joinStrings(reactionStrs, ", "),
			})
		}
	}

	return rows, nil
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

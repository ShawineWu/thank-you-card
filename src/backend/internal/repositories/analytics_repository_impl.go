package repositories

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GormAnalyticsRepository implements AnalyticsRepository using GORM
type GormAnalyticsRepository struct {
	db *gorm.DB
}

// NewGormAnalyticsRepository creates a new analytics repository
func NewGormAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &GormAnalyticsRepository{db: db}
}

// GetTotalCards returns the total number of cards
func (r *GormAnalyticsRepository) GetTotalCards() (int64, error) {
	var count int64
	err := r.db.Table("cards").Count(&count).Error
	return count, err
}

// GetActiveUsers returns count of unique users who participated
func (r *GormAnalyticsRepository) GetActiveUsers() (int64, error) {
	var count int64
	query := `
		SELECT COUNT(DISTINCT employee_id) FROM (
			SELECT sender_id as employee_id FROM cards
			UNION
			SELECT recipient_id as employee_id FROM card_recipients
		) AS active_users
	`
	err := r.db.Raw(query).Scan(&count).Error
	return count, err
}

// GetTopValues returns the most used company values
func (r *GormAnalyticsRepository) GetTopValues(limit int) ([]ValueUsage, error) {
	var results []ValueUsage
	query := `
		SELECT 
			cv.id as value_id,
			cv.name as value_name,
			COUNT(*) as count
		FROM card_values cvs
		JOIN company_values cv ON cv.id = cvs.value_id
		GROUP BY cv.id, cv.name
		ORDER BY count DESC
		LIMIT ?
	`
	err := r.db.Raw(query, limit).Scan(&results).Error
	return results, err
}

// GetCardsTrend returns card counts for current and previous month
func (r *GormAnalyticsRepository) GetCardsTrend() (int64, int64, error) {
	now := time.Now()
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)

	var thisMonth, lastMonth int64

	// This month
	err := r.db.Table("cards").
		Where("created_at >= ?", thisMonthStart).
		Count(&thisMonth).Error
	if err != nil {
		return 0, 0, err
	}

	// Last month
	err = r.db.Table("cards").
		Where("created_at >= ? AND created_at < ?", lastMonthStart, thisMonthStart).
		Count(&lastMonth).Error
	if err != nil {
		return 0, 0, err
	}

	return thisMonth, lastMonth, nil
}

// GetTopRecognizers returns employees who sent the most cards
func (r *GormAnalyticsRepository) GetTopRecognizers(limit int) ([]TopRecognizer, error) {
	var results []TopRecognizer
	query := `
		SELECT 
			sender_id as employee_id,
			COUNT(*) as cards_sent
		FROM cards
		GROUP BY sender_id
		ORDER BY cards_sent DESC
		LIMIT ?
	`
	err := r.db.Raw(query, limit).Scan(&results).Error
	return results, err
}

// GetTeamAnalytics returns recognition statistics by team
func (r *GormAnalyticsRepository) GetTeamAnalytics() ([]TeamStats, error) {
	// Note: This is a placeholder implementation as we don't have team data yet
	// In a real implementation, you'd join with employee service data
	var results []TeamStats

	// Mock implementation - returns empty for now
	// TODO: Integrate with employee service to get team information
	query := `
		SELECT 
			'placeholder' as team_id,
			'Placeholder Team' as team_name,
			0 as cards_sent,
			0 as cards_received,
			'' as top_value
		WHERE 1=0
	`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err
}

// GetValueDistribution returns distribution of all values
func (r *GormAnalyticsRepository) GetValueDistribution() ([]ValueUsage, error) {
	var results []ValueUsage
	query := `
		SELECT 
			cv.id as value_id,
			cv.name as value_name,
			COUNT(*) as count
		FROM card_values cvs
		JOIN company_values cv ON cv.id = cvs.value_id
		GROUP BY cv.id, cv.name
		ORDER BY count DESC
	`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err
}

// GetCardsForExport returns cards within date range for CSV export
func (r *GormAnalyticsRepository) GetCardsForExport(startDate, endDate time.Time) ([]ExportCard, error) {
	var cards []struct {
		ID                string
		SenderID          string
		RecognitionReason string
		CreatedAt         time.Time
	}

	err := r.db.Table("cards").
		Select("id, sender_id, recognition_reason, created_at").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&cards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch cards for export: %w", err)
	}

	// Enrich with recipients and values
	exportCards := make([]ExportCard, len(cards))
	for i, card := range cards {
		// Get recipients
		var recipients []string
		err := r.db.Table("card_recipients").
			Select("recipient_id").
			Where("card_id = ?", card.ID).
			Pluck("recipient_id", &recipients).Error
		if err != nil {
			return nil, fmt.Errorf("failed to fetch recipients: %w", err)
		}

		// Get values
		var values []string
		query := `
			SELECT cv.name
			FROM card_values cvs
			JOIN company_values cv ON cv.id = cvs.value_id
			WHERE cvs.card_id = ?
		`
		err = r.db.Raw(query, card.ID).Pluck("name", &values).Error
		if err != nil {
			return nil, fmt.Errorf("failed to fetch values: %w", err)
		}

		exportCards[i] = ExportCard{
			ID:                card.ID,
			SenderID:          card.SenderID,
			Recipients:        strings.Join(recipients, "; "),
			RecognitionReason: card.RecognitionReason,
			Values:            strings.Join(values, "; "),
			CreatedAt:         card.CreatedAt,
		}
	}

	return exportCards, nil
}

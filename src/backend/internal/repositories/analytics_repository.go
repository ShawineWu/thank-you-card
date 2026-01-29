package repositories

import (
	"time"

	"github.com/google/uuid"
)

// AnalyticsRepository defines methods for analytics data aggregation
type AnalyticsRepository interface {
	// GetTotalCards returns the total number of cards
	GetTotalCards() (int64, error)

	// GetActiveUsers returns count of users who sent or received cards
	GetActiveUsers() (int64, error)

	// GetTopValues returns the most used company values
	GetTopValues(limit int) ([]ValueUsage, error)

	// GetCardsTrend returns card counts for current and previous month
	GetCardsTrend() (thisMonth int64, lastMonth int64, err error)

	// GetTopRecognizers returns employees who sent the most cards
	GetTopRecognizers(limit int) ([]TopRecognizer, error)

	// GetTeamAnalytics returns recognition statistics by team
	GetTeamAnalytics() ([]TeamStats, error)

	// GetValueDistribution returns distribution of all values
	GetValueDistribution() ([]ValueUsage, error)

	// GetCardsForExport returns cards within date range for CSV export
	GetCardsForExport(startDate, endDate time.Time) ([]ExportCard, error)
}

// ValueUsage represents a company value and its usage count
type ValueUsage struct {
	ValueID   uuid.UUID
	ValueName string
	Count     int64
}

// TopRecognizer represents an employee and their sent card count
type TopRecognizer struct {
	EmployeeID   string
	EmployeeName string
	CardsSent    int64
}

// TeamStats represents team-level recognition statistics
type TeamStats struct {
	TeamID        string
	TeamName      string
	CardsSent     int64
	CardsReceived int64
	TopValue      string
}

// ExportCard represents a card record for CSV export
type ExportCard struct {
	ID                string
	SenderID          string
	SenderName        string
	Recipients        string // Comma-separated IDs
	RecipientNames    string // Comma-separated names
	RecognitionReason string
	Values            string // Comma-separated value names
	CreatedAt         time.Time
}

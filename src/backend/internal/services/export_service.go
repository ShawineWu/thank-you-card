package services

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"time"

	"github.com/castlery/thank-you-card/internal/dtos"
	"github.com/castlery/thank-you-card/internal/repositories"
)

// ExportService handles data export operations
type ExportService struct {
	analyticsRepo repositories.AnalyticsRepository
}

// NewExportService creates a new export service
func NewExportService(analyticsRepo repositories.AnalyticsRepository) *ExportService {
	return &ExportService{
		analyticsRepo: analyticsRepo,
	}
}

// ExportToCSV generates a CSV file with card data
func (s *ExportService) ExportToCSV(filters dtos.ExportFilters) (*bytes.Buffer, error) {
	// Fetch cards within date range
	cards, err := s.analyticsRepo.GetCardsForExport(filters.StartDate, filters.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cards for export: %w", err)
	}

	// Create CSV buffer
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Write header
	header := []string{"Card ID", "Sender ID", "Sender Name", "Recipient IDs", "Recipient Names", "Recognition Reason", "Values", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, card := range cards {
		row := []string{
			card.ID,
			card.SenderID,
			card.SenderName,
			card.Recipients,
			card.RecipientNames,
			card.RecognitionReason,
			card.Values,
			card.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf, nil
}

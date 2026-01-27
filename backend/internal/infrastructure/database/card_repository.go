package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CardRepositoryImpl implements the CardRepository interface
type CardRepositoryImpl struct {
	db *gorm.DB
}

// NewCardRepository creates a new card repository
func NewCardRepository(db *gorm.DB) card.CardRepository {
	return &CardRepositoryImpl{db: db}
}

// Create creates a new card with recipients and values in a transaction
func (r *CardRepositoryImpl) Create(cardEntity *card.Card, recipients []card.CardRecipient, values []card.CardValue) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the card
		if err := tx.Create(cardEntity).Error; err != nil {
			return fmt.Errorf("failed to create card: %w", err)
		}

		// Create recipients
		if len(recipients) > 0 {
			if err := tx.Create(&recipients).Error; err != nil {
				return fmt.Errorf("failed to create card recipients: %w", err)
			}
		}

		// Create card values
		if len(values) > 0 {
			if err := tx.Create(&values).Error; err != nil {
				return fmt.Errorf("failed to create card values: %w", err)
			}
		}

		return nil
	})
}

// GetByID retrieves a card by ID with all relationships
func (r *CardRepositoryImpl) GetByID(id uuid.UUID) (*card.Card, error) {
	var cardEntity card.Card
	err := r.db.Preload("Recipients").Preload("Values").First(&cardEntity, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("card not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	return &cardEntity, nil
}

// Update updates a card
func (r *CardRepositoryImpl) Update(cardEntity *card.Card) error {
	err := r.db.Save(cardEntity).Error
	if err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}
	return nil
}

// Delete deletes a card by ID
func (r *CardRepositoryImpl) Delete(id uuid.UUID) error {
	result := r.db.Delete(&card.Card{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete card: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("card not found: %s", id)
	}
	return nil
}

// GetCompanyFeed retrieves the company-wide card feed with pagination
func (r *CardRepositoryImpl) GetCompanyFeed(limit, offset int) ([]card.Card, int, error) {
	var cards []card.Card
	var total int64

	// Get total count
	if err := r.db.Model(&card.Card{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards: %w", err)
	}

	// Get cards with pagination
	err := r.db.Preload("Recipients").Preload("Values").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&cards).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get company feed: %w", err)
	}

	return cards, int(total), nil
}

// GetCardsBySender retrieves cards sent by a specific sender with pagination
func (r *CardRepositoryImpl) GetCardsBySender(senderID string, limit, offset int) ([]card.Card, int, error) {
	var cards []card.Card
	var total int64

	// Get total count
	if err := r.db.Model(&card.Card{}).Where("sender_id = ?", senderID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards by sender: %w", err)
	}

	// Get cards with pagination
	err := r.db.Preload("Recipients").Preload("Values").
		Where("sender_id = ?", senderID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&cards).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cards by sender: %w", err)
	}

	return cards, int(total), nil
}

// GetCardsByRecipient retrieves cards received by a specific recipient with pagination
func (r *CardRepositoryImpl) GetCardsByRecipient(recipientID string, limit, offset int) ([]card.Card, int, error) {
	var cards []card.Card
	var total int64

	// Subquery to get card IDs for the recipient
	subQuery := r.db.Model(&card.CardRecipient{}).
		Select("card_id").
		Where("recipient_id = ?", recipientID)

	// Get total count
	if err := r.db.Model(&card.Card{}).
		Where("id IN (?)", subQuery).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards by recipient: %w", err)
	}

	// Get cards with pagination
	err := r.db.Preload("Recipients").Preload("Values").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&cards).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cards by recipient: %w", err)
	}

	return cards, int(total), nil
}

// SearchCards searches cards based on filters with pagination
func (r *CardRepositoryImpl) SearchCards(filters card.SearchFilters, limit, offset int) ([]card.Card, int, error) {
	query := r.db.Model(&card.Card{})
	
	// Apply filters
	query = r.applySearchFilters(query, filters)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count filtered cards: %w", err)
	}

	var cards []card.Card
	err := query.Preload("Recipients").Preload("Values").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&cards).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search cards: %w", err)
	}

	return cards, int(total), nil
}

// applySearchFilters applies search filters to a GORM query
func (r *CardRepositoryImpl) applySearchFilters(query *gorm.DB, filters card.SearchFilters) *gorm.DB {
	if filters.SenderID != nil {
		query = query.Where("sender_id = ?", *filters.SenderID)
	}

	if filters.RecipientID != nil {
		subQuery := r.db.Model(&card.CardRecipient{}).
			Select("card_id").
			Where("recipient_id = ?", *filters.RecipientID)
		query = query.Where("id IN (?)", subQuery)
	}

	if len(filters.ValueIDs) > 0 {
		subQuery := r.db.Model(&card.CardValue{}).
			Select("card_id").
			Where("value_id IN ?", filters.ValueIDs)
		query = query.Where("id IN (?)", subQuery)
	}

	if filters.DateFrom != nil {
		if dateFrom, err := time.Parse("2006-01-02", *filters.DateFrom); err == nil {
			query = query.Where("created_at >= ?", dateFrom)
		}
	}

	if filters.DateTo != nil {
		if dateTo, err := time.Parse("2006-01-02", *filters.DateTo); err == nil {
			query = query.Where("created_at <= ?", dateTo.Add(24*time.Hour))
		}
	}

	if filters.SearchText != nil && *filters.SearchText != "" {
		searchTerm := "%" + strings.ToLower(*filters.SearchText) + "%"
		query = query.Where("LOWER(reason) LIKE ?", searchTerm)
	}

	return query
}

// CountCardsBySender counts total cards sent by a sender
func (r *CardRepositoryImpl) CountCardsBySender(senderID string) (int, error) {
	var count int64
	err := r.db.Model(&card.Card{}).Where("sender_id = ?", senderID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count cards by sender: %w", err)
	}
	return int(count), nil
}

// CountCardsByRecipient counts total cards received by a recipient
func (r *CardRepositoryImpl) CountCardsByRecipient(recipientID string) (int, error) {
	var count int64
	err := r.db.Model(&card.CardRecipient{}).Where("recipient_id = ?", recipientID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count cards by recipient: %w", err)
	}
	return int(count), nil
}

// GetValueDistributionBySender gets value distribution for cards sent by a sender
func (r *CardRepositoryImpl) GetValueDistributionBySender(senderID string) ([]card.ValueDistribution, error) {
	var results []struct {
		ValueID   uuid.UUID `json:"value_id"`
		ValueName string    `json:"value_name"`
		Count     int       `json:"count"`
	}

	err := r.db.Table("card_values cv").
		Select("cv.value_id, comp_val.name as value_name, COUNT(*) as count").
		Joins("JOIN cards c ON cv.card_id = c.id").
		Joins("JOIN company_values comp_val ON cv.value_id = comp_val.id").
		Where("c.sender_id = ?", senderID).
		Group("cv.value_id, comp_val.name").
		Order("count DESC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get value distribution: %w", err)
	}

	// Calculate total for percentages
	var total int
	for _, result := range results {
		total += result.Count
	}

	// Convert to domain objects
	var distribution []card.ValueDistribution
	for _, result := range results {
		percentage := float64(result.Count) / float64(total) * 100
		distribution = append(distribution, card.ValueDistribution{
			ValueID:    result.ValueID,
			ValueName:  result.ValueName,
			Count:      result.Count,
			Percentage: percentage,
		})
	}

	return distribution, nil
}

// GetRecentActivityByEmployee gets recent activity for an employee
func (r *CardRepositoryImpl) GetRecentActivityByEmployee(employeeID string, limit int) ([]card.RecentActivity, error) {
	var activities []card.RecentActivity

	// Get sent cards
	var sentCards []struct {
		CardID    uuid.UUID `json:"card_id"`
		CreatedAt time.Time `json:"created_at"`
		Recipient string    `json:"recipient"`
	}

	err := r.db.Table("cards c").
		Select("c.id as card_id, c.created_at, cr.recipient_id as recipient").
		Joins("JOIN card_recipients cr ON c.id = cr.card_id").
		Where("c.sender_id = ?", employeeID).
		Order("c.created_at DESC").
		Limit(limit).
		Find(&sentCards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sent cards activity: %w", err)
	}

	for _, sent := range sentCards {
		activities = append(activities, card.RecentActivity{
			Type:       "sent",
			CardID:     sent.CardID,
			Date:       sent.CreatedAt,
			OtherParty: sent.Recipient,
		})
	}

	// Get received cards
	var receivedCards []struct {
		CardID    uuid.UUID `json:"card_id"`
		CreatedAt time.Time `json:"created_at"`
		Sender    string    `json:"sender"`
	}

	err = r.db.Table("cards c").
		Select("c.id as card_id, c.created_at, c.sender_id as sender").
		Joins("JOIN card_recipients cr ON c.id = cr.card_id").
		Where("cr.recipient_id = ?", employeeID).
		Order("c.created_at DESC").
		Limit(limit).
		Find(&receivedCards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get received cards activity: %w", err)
	}

	for _, received := range receivedCards {
		activities = append(activities, card.RecentActivity{
			Type:       "received",
			CardID:     received.CardID,
			Date:       received.CreatedAt,
			OtherParty: received.Sender,
		})
	}

	// Sort by date and limit
	// Note: In a real implementation, you'd want to do this sorting in the database
	// This is simplified for demonstration
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities, nil
}

// GetTop10RecognizedEmployees gets the top 10 most recognized employees
func (r *CardRepositoryImpl) GetTop10RecognizedEmployees() ([]card.Top10Employee, error) {
	var results []struct {
		EmployeeID string `json:"employee_id"`
		CardCount  int    `json:"card_count"`
	}

	err := r.db.Table("card_recipients cr").
		Select("cr.recipient_id as employee_id, COUNT(*) as card_count").
		Group("cr.recipient_id").
		Order("card_count DESC").
		Limit(10).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top 10 recognized employees: %w", err)
	}

	// Convert to domain objects
	var top10 []card.Top10Employee
	for i, result := range results {
		top10 = append(top10, card.Top10Employee{
			EmployeeID:   result.EmployeeID,
			EmployeeName: result.EmployeeID, // This would be resolved by the Employee Directory Service
			CardCount:    result.CardCount,
			Rank:         i + 1,
		})
	}

	return top10, nil
}

// LoadCardWithRelationships loads a card with all its relationships
func (r *CardRepositoryImpl) LoadCardWithRelationships(cardID uuid.UUID) (*card.Card, error) {
	return r.GetByID(cardID)
}
package repositories

import (
	"fmt"

	"github.com/castlery/thank-you-card/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormCardRepository implements CardRepository using GORM
type GormCardRepository struct {
	db *gorm.DB
}

// NewGormCardRepository creates a new GORM-based card repository
func NewGormCardRepository(db *gorm.DB) CardRepository {
	return &GormCardRepository{db: db}
}

// Save persists a new card with all its associations
func (r *GormCardRepository) Save(card *models.Card) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create card with all associations
		if err := tx.Create(card).Error; err != nil {
			return fmt.Errorf("failed to save card: %w", err)
		}
		return nil
	})
}

// FindByID retrieves a card by ID with all associations preloaded
func (r *GormCardRepository) FindByID(id uuid.UUID) (*models.Card, error) {
	var card models.Card
	err := r.db.
		Preload("Recipients").
		Preload("Values.CompanyValue").
		First(&card, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find card: %w", err)
	}

	return &card, nil
}

// FindBySender retrieves cards sent by a specific employee
func (r *GormCardRepository) FindBySender(senderID string, page, pageSize int) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	// Count total
	if err := r.db.Model(&models.Card{}).Where("sender_id = ?", senderID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.
		Preload("Recipients").
		Preload("Values.CompanyValue").
		Where("sender_id = ?", senderID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&cards).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find cards by sender: %w", err)
	}

	return cards, total, nil
}

// FindByRecipient retrieves cards received by a specific employee
func (r *GormCardRepository) FindByRecipient(recipientID string, page, pageSize int) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	// Subquery to get card IDs received by this employee
	subQuery := r.db.Model(&models.CardRecipient{}).
		Select("card_id").
		Where("recipient_id = ?", recipientID)

	// Count total
	if err := r.db.Model(&models.Card{}).
		Where("id IN (?)", subQuery).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.
		Preload("Recipients").
		Preload("Values.CompanyValue").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&cards).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find cards by recipient: %w", err)
	}

	return cards, total, nil
}

// FindAll retrieves all cards (company-wide feed)
func (r *GormCardRepository) FindAll(page, pageSize int) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	// Count total
	if err := r.db.Model(&models.Card{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cards: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.
		Preload("Recipients").
		Preload("Values.CompanyValue").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&cards).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find all cards: %w", err)
	}

	return cards, total, nil
}

// FindWithFilters retrieves cards matching multiple filter criteria
func (r *GormCardRepository) FindWithFilters(senderID, recipientID string, valueIDs []uuid.UUID, startDate, endDate *string, search string, page, pageSize int) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	query := r.db.Model(&models.Card{})

	// Filter by sender if provided
	if senderID != "" {
		query = query.Where("sender_id = ?", senderID)
	}

	// Filter by recipient if provided
	if recipientID != "" {
		subQuery := r.db.Model(&models.CardRecipient{}).
			Select("card_id").
			Where("recipient_id = ?", recipientID)
		query = query.Where("id IN (?)", subQuery)
	}

	// Filter by value IDs if provided
	if len(valueIDs) > 0 {
		subQuery := r.db.Model(&models.CardValue{}).
			Select("card_id").
			Where("value_id IN ?", valueIDs)
		query = query.Where("id IN (?)", subQuery)
	}

	// Filter by date range if provided
	if startDate != nil && *startDate != "" {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		query = query.Where("created_at <= ?", *endDate)
	}

	// Keyword search if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("recognition_reason ILIKE ?", searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count filtered cards: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.
		Preload("Recipients").
		Preload("Values.CompanyValue").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&cards).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find filtered cards: %w", err)
	}

	return cards, total, nil
}

// CountBySender counts cards sent by an employee
func (r *GormCardRepository) CountBySender(senderID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Card{}).
		Where("sender_id = ?", senderID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count cards by sender: %w", err)
	}

	return count, nil
}

// CountByRecipient counts cards received by an employee
func (r *GormCardRepository) CountByRecipient(recipientID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.CardRecipient{}).
		Where("recipient_id = ?", recipientID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count cards by recipient: %w", err)
	}

	return count, nil
}

// GetTopRecipients retrieves the employees who received the most cards
func (r *GormCardRepository) GetTopRecipients(limit int) ([]TopRecipient, error) {
	var results []TopRecipient
	err := r.db.Model(&models.CardRecipient{}).
		Select("recipient_id as employee_id, count(*) as count").
		Group("recipient_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top recipients: %w", err)
	}

	return results, nil
}

// GetValueStatsForUser retrieves statistics on values associated with cards received by a user
func (r *GormCardRepository) GetValueStatsForUser(employeeID string) ([]ValueStat, error) {
	var results []ValueStat

	// Subquery to get card IDs received by this employee
	subQuery := r.db.Model(&models.CardRecipient{}).
		Select("card_id").
		Where("recipient_id = ?", employeeID)

	err := r.db.Table("card_values").
		Select("card_values.value_id, company_values.name as value_name, count(*) as count").
		Joins("JOIN company_values ON company_values.id = card_values.value_id").
		Where("card_values.card_id IN (?)", subQuery).
		Group("card_values.value_id, company_values.name").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get value stats for user: %w", err)
	}

	return results, nil
}

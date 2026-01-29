package repositories

import (
	"context"
	"time"

	"thank-you-card-backend/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Employee, error)
	GetAll(ctx context.Context) ([]models.Employee, error)
	Search(ctx context.Context, query string) ([]models.Employee, error)
}

type CardFilters struct {
	ValueIDs     []uint
	SenderIDs    []uint
	RecipientIDs []uint
	From         *time.Time
	To           *time.Time
	Query        string
	Page         int
	PageSize     int
}

type CardRepository interface {
	CreateCardWithRelations(ctx context.Context, card *models.Card, recipients []uint, valueIDs []uint) error
	GetByID(ctx context.Context, id uint) (*models.Card, error)
	GetFeed(ctx context.Context, filters CardFilters) ([]models.Card, int64, error)
	GetByRecipient(ctx context.Context, recipientID uint, filters CardFilters) ([]models.Card, int64, error)
	GetBySender(ctx context.Context, senderID uint, filters CardFilters) ([]models.Card, int64, error)
	FilterAndSearch(ctx context.Context, filters CardFilters) ([]models.Card, int64, error)
}

type EmojiReactionRepository interface {
	SetReaction(ctx context.Context, cardID uint, userID uint, emojiCode string) error
	RemoveReaction(ctx context.Context, cardID uint, userID uint) error
}

type AnalyticsRepository interface {
	TopRecognizedEmployees(ctx context.Context, from, to time.Time, limit int) ([]TopEmployee, error)
	MostActiveRecognizers(ctx context.Context, from, to time.Time, limit int) ([]TopEmployee, error)
	TeamRecognitionPatterns(ctx context.Context, from, to time.Time) ([]TeamPattern, error)
	ValuesDistribution(ctx context.Context, from, to time.Time) ([]ValueDistributionItem, error)
}

type CompanyValueRepository interface {
	GetAll(ctx context.Context) ([]models.CompanyValue, error)
	GetByID(ctx context.Context, id uint) (*models.CompanyValue, error)
	GetByType(ctx context.Context, valueType string) ([]models.CompanyValue, error)
}

type MilestoneRepository interface {
	GetAll(ctx context.Context) ([]models.Milestone, error)
	GetByID(ctx context.Context, id uint) (*models.Milestone, error)
	GetUserMilestones(ctx context.Context, userID uint) ([]models.UserMilestone, error)
	CreateUserMilestone(ctx context.Context, userMilestone *models.UserMilestone) error
	GetUserMilestoneByMilestoneID(ctx context.Context, userID uint, milestoneID uint) (*models.UserMilestone, error)
}

type TopEmployee struct {
	EmployeeID uint
	Name       string
	Department string
	Count      int64
}

type TeamPattern struct {
	Department              string
	TotalSent               int64
	TotalReceived           int64
	AverageCardsPerEmployee float64
}

type ValueDistributionItem struct {
	CompanyValueID uint
	Code           string
	Name           string
	Type           string
	Count          int64
}

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) CreateCardWithRelations(ctx context.Context, card *models.Card, recipients []uint, valueIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(card).Error; err != nil {
			return err
		}

		for _, rid := range recipients {
			cr := models.CardRecipient{
				CardID:      card.ID,
				RecipientID: rid,
			}
			if err := tx.Create(&cr).Error; err != nil {
				return err
			}
		}

		for _, vid := range valueIDs {
			cv := models.CardValue{
				CardID:         card.ID,
				CompanyValueID: vid,
			}
			if err := tx.Create(&cv).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByID retrieves a card by ID with all relations loaded
func (r *cardRepository) GetByID(ctx context.Context, id uint) (*models.Card, error) {
	var card models.Card
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Recipients.Recipient").
		Preload("Values.CompanyValue").
		Preload("Reactions").
		First(&card, id).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

// For now, feed and filter queries will be simple joins with pagination.
// More complex analytics are handled in AnalyticsRepository.

func (r *cardRepository) GetFeed(ctx context.Context, filters CardFilters) ([]models.Card, int64, error) {
	return r.queryCards(ctx, filters, nil)
}

func (r *cardRepository) GetByRecipient(ctx context.Context, recipientID uint, filters CardFilters) ([]models.Card, int64, error) {
	return r.queryCards(ctx, filters, func(tx *gorm.DB) *gorm.DB {
		return tx.Joins("JOIN card_recipients ON card_recipients.card_id = cards.id AND card_recipients.recipient_id = ?", recipientID)
	})
}

func (r *cardRepository) GetBySender(ctx context.Context, senderID uint, filters CardFilters) ([]models.Card, int64, error) {
	return r.queryCards(ctx, filters, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("sender_id = ?", senderID)
	})
}

func (r *cardRepository) FilterAndSearch(ctx context.Context, filters CardFilters) ([]models.Card, int64, error) {
	return r.queryCards(ctx, filters, nil)
}

func (r *cardRepository) queryCards(ctx context.Context, filters CardFilters, base func(*gorm.DB) *gorm.DB) ([]models.Card, int64, error) {
	db := r.db.WithContext(ctx).Model(&models.Card{}).
		Preload("Sender").
		Preload("Recipients").
		Preload("Recipients.Recipient").
		Preload("Values").
		Preload("Values.CompanyValue").
		Preload("Reactions").
		Preload("Reactions.User")

	if base != nil {
		db = base(db)
	}

	if len(filters.ValueIDs) > 0 {
		db = db.Joins("JOIN card_values ON card_values.card_id = cards.id").
			Where("card_values.company_value_id IN ?", filters.ValueIDs)
	}

	if len(filters.SenderIDs) > 0 {
		db = db.Where("sender_id IN ?", filters.SenderIDs)
	}

	if len(filters.RecipientIDs) > 0 {
		db = db.Joins("JOIN card_recipients ON card_recipients.card_id = cards.id").
			Where("card_recipients.recipient_id IN ?", filters.RecipientIDs)
	}

	if filters.From != nil {
		db = db.Where("cards.created_at >= ?", *filters.From)
	}
	if filters.To != nil {
		db = db.Where("cards.created_at <= ?", *filters.To)
	}
	if filters.Query != "" {
		db = db.Where("cards.reason ILIKE ?", "%"+filters.Query+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var cards []models.Card
	if err := db.Order("cards.created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&cards).Error; err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}

// Gorm implementations

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) GetByID(ctx context.Context, id uint) (*models.Employee, error) {
	var e models.Employee
	if err := r.db.WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *employeeRepository) GetAll(ctx context.Context) ([]models.Employee, error) {
	var employees []models.Employee
	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepository) Search(ctx context.Context, query string) ([]models.Employee, error) {
	var employees []models.Employee
	db := r.db.WithContext(ctx)
	if query != "" {
		searchPattern := "%" + query + "%"
		db = db.Where("name ILIKE ? OR department ILIKE ? OR email ILIKE ?", searchPattern, searchPattern, searchPattern)
	}
	if err := db.Order("name ASC").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

type emojiReactionRepository struct {
	db *gorm.DB
}

func NewEmojiReactionRepository(db *gorm.DB) EmojiReactionRepository {
	return &emojiReactionRepository{db: db}
}

func (r *emojiReactionRepository) SetReaction(ctx context.Context, cardID uint, userID uint, emojiCode string) error {
	// upsert by (card_id, user_id)
	var existing models.EmojiReaction
	err := r.db.WithContext(ctx).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			reaction := models.EmojiReaction{
				CardID:    cardID,
				UserID:    userID,
				EmojiCode: emojiCode,
			}
			return r.db.WithContext(ctx).Create(&reaction).Error
		}
		return err
	}

	existing.EmojiCode = emojiCode
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *emojiReactionRepository) RemoveReaction(ctx context.Context, cardID uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Delete(&models.EmojiReaction{}).Error
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) TopRecognizedEmployees(ctx context.Context, from, to time.Time, limit int) ([]TopEmployee, error) {
	var results []TopEmployee

	query := r.db.WithContext(ctx).
		Model(&models.CardRecipient{}).
		Select(`
			card_recipients.recipient_id as employee_id,
			employees.name,
			employees.department,
			COUNT(*) as count
		`).
		Joins("JOIN cards ON cards.id = card_recipients.card_id").
		Joins("JOIN employees ON employees.id = card_recipients.recipient_id").
		Where("cards.created_at >= ? AND cards.created_at <= ?", from, to).
		Group("card_recipients.recipient_id, employees.name, employees.department").
		Order("count DESC").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *analyticsRepository) MostActiveRecognizers(ctx context.Context, from, to time.Time, limit int) ([]TopEmployee, error) {
	var results []TopEmployee

	query := r.db.WithContext(ctx).
		Model(&models.Card{}).
		Select(`
			cards.sender_id as employee_id,
			employees.name,
			employees.department,
			COUNT(*) as count
		`).
		Joins("JOIN employees ON employees.id = cards.sender_id").
		Where("cards.created_at >= ? AND cards.created_at <= ?", from, to).
		Group("cards.sender_id, employees.name, employees.department").
		Order("count DESC").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *analyticsRepository) TeamRecognitionPatterns(ctx context.Context, from, to time.Time) ([]TeamPattern, error) {
	var results []TeamPattern

	// Get sent counts by department
	sentQuery := r.db.WithContext(ctx).
		Model(&models.Card{}).
		Select(`
			employees.department,
			COUNT(*) as total_sent
		`).
		Joins("JOIN employees ON employees.id = cards.sender_id").
		Where("cards.created_at >= ? AND cards.created_at <= ?", from, to).
		Group("employees.department")

	var sentCounts []struct {
		Department string
		TotalSent  int64
	}
	if err := sentQuery.Scan(&sentCounts).Error; err != nil {
		return nil, err
	}

	// Get received counts by department
	receivedQuery := r.db.WithContext(ctx).
		Model(&models.CardRecipient{}).
		Select(`
			employees.department,
			COUNT(*) as total_received
		`).
		Joins("JOIN cards ON cards.id = card_recipients.card_id").
		Joins("JOIN employees ON employees.id = card_recipients.recipient_id").
		Where("cards.created_at >= ? AND cards.created_at <= ?", from, to).
		Group("employees.department")

	var receivedCounts []struct {
		Department    string
		TotalReceived int64
	}
	if err := receivedQuery.Scan(&receivedCounts).Error; err != nil {
		return nil, err
	}

	// Get employee counts by department
	var deptEmployees []struct {
		Department string
		Count      int64
	}
	if err := r.db.WithContext(ctx).
		Model(&models.Employee{}).
		Select("department, COUNT(*) as count").
		Group("department").
		Scan(&deptEmployees).Error; err != nil {
		return nil, err
	}

	// Combine results
	deptMap := make(map[string]*TeamPattern)
	for _, sc := range sentCounts {
		if deptMap[sc.Department] == nil {
			deptMap[sc.Department] = &TeamPattern{Department: sc.Department}
		}
		deptMap[sc.Department].TotalSent = sc.TotalSent
	}
	for _, rc := range receivedCounts {
		if deptMap[rc.Department] == nil {
			deptMap[rc.Department] = &TeamPattern{Department: rc.Department}
		}
		deptMap[rc.Department].TotalReceived = rc.TotalReceived
	}
	for _, de := range deptEmployees {
		if deptMap[de.Department] == nil {
			deptMap[de.Department] = &TeamPattern{Department: de.Department}
		}
		if de.Count > 0 {
			deptMap[de.Department].AverageCardsPerEmployee = float64(deptMap[de.Department].TotalReceived) / float64(de.Count)
		}
	}

	for _, tp := range deptMap {
		results = append(results, *tp)
	}

	return results, nil
}

func (r *analyticsRepository) ValuesDistribution(ctx context.Context, from, to time.Time) ([]ValueDistributionItem, error) {
	var results []ValueDistributionItem

	query := r.db.WithContext(ctx).
		Model(&models.CardValue{}).
		Select(`
			card_values.company_value_id,
			company_values.code,
			company_values.name,
			company_values.type,
			COUNT(*) as count
		`).
		Joins("JOIN cards ON cards.id = card_values.card_id").
		Joins("JOIN company_values ON company_values.id = card_values.company_value_id").
		Where("cards.created_at >= ? AND cards.created_at <= ?", from, to).
		Group("card_values.company_value_id, company_values.code, company_values.name, company_values.type").
		Order("count DESC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// CompanyValueRepository implementation
type companyValueRepository struct {
	db *gorm.DB
}

func NewCompanyValueRepository(db *gorm.DB) CompanyValueRepository {
	return &companyValueRepository{db: db}
}

func (r *companyValueRepository) GetAll(ctx context.Context) ([]models.CompanyValue, error) {
	var values []models.CompanyValue
	if err := r.db.WithContext(ctx).
		Order("type ASC, name ASC").
		Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

func (r *companyValueRepository) GetByID(ctx context.Context, id uint) (*models.CompanyValue, error) {
	var value models.CompanyValue
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&value).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (r *companyValueRepository) GetByType(ctx context.Context, valueType string) ([]models.CompanyValue, error) {
	var values []models.CompanyValue
	if err := r.db.WithContext(ctx).
		Where("type = ?", valueType).
		Order("name ASC").
		Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

// MilestoneRepository implementation
type milestoneRepository struct {
	db *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &milestoneRepository{db: db}
}

func (r *milestoneRepository) GetAll(ctx context.Context) ([]models.Milestone, error) {
	var milestones []models.Milestone
	if err := r.db.WithContext(ctx).
		Order("type ASC, threshold ASC").
		Find(&milestones).Error; err != nil {
		return nil, err
	}
	return milestones, nil
}

func (r *milestoneRepository) GetByID(ctx context.Context, id uint) (*models.Milestone, error) {
	var milestone models.Milestone
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&milestone).Error; err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *milestoneRepository) GetUserMilestones(ctx context.Context, userID uint) ([]models.UserMilestone, error) {
	var userMilestones []models.UserMilestone
	if err := r.db.WithContext(ctx).
		Preload("Milestone").
		Where("user_id = ?", userID).
		Order("achieved_at DESC").
		Find(&userMilestones).Error; err != nil {
		return nil, err
	}
	return userMilestones, nil
}

func (r *milestoneRepository) CreateUserMilestone(ctx context.Context, userMilestone *models.UserMilestone) error {
	return r.db.WithContext(ctx).Create(userMilestone).Error
}

func (r *milestoneRepository) GetUserMilestoneByMilestoneID(ctx context.Context, userID uint, milestoneID uint) (*models.UserMilestone, error) {
	var userMilestone models.UserMilestone
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND milestone_id = ?", userID, milestoneID).
		First(&userMilestone).Error; err != nil {
		return nil, err
	}
	return &userMilestone, nil
}

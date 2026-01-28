package repositories

import (
	"context"

	"thank-you-card-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

package services

import (
	"context"
	"errors"

	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/repositories"

	"gorm.io/gorm"
)

type CompanyValueService interface {
	GetAll(ctx context.Context) ([]models.CompanyValue, error)
	GetByID(ctx context.Context, id uint) (*models.CompanyValue, error)
	GetByType(ctx context.Context, valueType string) ([]models.CompanyValue, error)
}

type companyValueService struct {
	repo repositories.CompanyValueRepository
}

func NewCompanyValueService(repo repositories.CompanyValueRepository) CompanyValueService {
	return &companyValueService{repo: repo}
}

func (s *companyValueService) GetAll(ctx context.Context) ([]models.CompanyValue, error) {
	return s.repo.GetAll(ctx)
}

func (s *companyValueService) GetByID(ctx context.Context, id uint) (*models.CompanyValue, error) {
	value, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("company value not found")
		}
		return nil, err
	}
	return value, nil
}

func (s *companyValueService) GetByType(ctx context.Context, valueType string) ([]models.CompanyValue, error) {
	if valueType != "VALUE" && valueType != "CREDO" {
		return nil, errors.New("invalid value type, must be VALUE or CREDO")
	}
	return s.repo.GetByType(ctx, valueType)
}

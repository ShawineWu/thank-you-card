package database

import (
	"log"
	"time"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
	"github.com/google/uuid"
)

// Seed populates the database with initial data
func Seed() error {
	if DB == nil {
		return log.New(nil, "", 0).Output(0, "database not initialized")
	}

	log.Println("Seeding initial data...")

	// Seed Company Values
	values := []models.CompanyValue{
		{
			ID:          uuid.New(),
			Name:        "Customer Obsessed",
			Description: "We put our customers at the heart of everything we do.",
			Type:        models.ValueTypeValue,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Name:        "Own It",
			Description: "We take responsibility for our actions and results.",
			Type:        models.ValueTypeValue,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Name:        "Do More with Less",
			Description: "We are resourceful and find creative ways to achieve more.",
			Type:        models.ValueTypeValue,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Name:        "Be Bold",
			Description: "We take calculated risks and learn from failure.",
			Type:        models.ValueTypeValue,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Name:        "Better Together",
			Description: "We collaborate and support each other to achieve our goals.",
			Type:        models.ValueTypeValue,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Name:        "Stay Humble",
			Description: "We are open to feedback and always look for ways to improve.",
			Type:        models.ValueTypeCredo,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, v := range values {
		var count int64
		DB.Model(&models.CompanyValue{}).Where("name = ?", v.Name).Count(&count)
		if count == 0 {
			if err := DB.Create(&v).Error; err != nil {
				return err
			}
			log.Printf("Seeded company value: %s", v.Name)
		}
	}

	log.Println("Initial data seeding completed successfully")
	return nil
}

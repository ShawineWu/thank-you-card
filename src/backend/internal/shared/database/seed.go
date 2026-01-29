package database

import (
	"log"
	"time"

	"github.com/castlery/thank-you-card/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Seed populates the database with initial data
func Seed() error {
	if DB == nil {
		return log.New(nil, "", 0).Output(0, "database not initialized")
	}

	log.Println("Seeding initial data...")

	// Seed Company Values (Code must be unique). Exists (by code or name) => skip.
	values := []models.CompanyValue{
		{
			ID:          uuid.New(),
			Code:        "CUSTOMER_OBSESSED",
			Name:        "Customer Obsessed",
			Description: "We put our customers at the heart of everything we do.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Led a project that improved customer satisfaction by 30%",
				"Advocated for a customer-requested feature that wasn't initially prioritized",
				"Fixed a customer-reported issue within 24 hours of receiving the feedback",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "OWN_IT",
			Name:        "Own It",
			Description: "We take responsibility for our actions and results.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Took responsibility for a project outcome, both successes and failures",
				"Admitted a mistake early and worked proactively to fix it",
				"Fixed an issue that wasn't technically my responsibility but affected the team",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "DO_MORE_WITH_LESS",
			Name:        "Do More with Less",
			Description: "We are resourceful and find creative ways to achieve more.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Optimized a process that reduced costs by 30% without sacrificing quality",
				"Simplified a complex process that saved the team 5 hours per week",
				"Found a creative solution that reduced implementation time by 50%",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "BE_BOLD",
			Name:        "Be Bold",
			Description: "We take calculated risks and learn from failure.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Quickly prototyped a solution to test a hypothesis within 2 days",
				"Made a decision and moved forward when others were stuck in analysis paralysis",
				"Proposed a new approach that improved team velocity",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "BETTER_TOGETHER",
			Name:        "Better Together",
			Description: "We collaborate and support each other to achieve our goals.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Helped a teammate meet a tight deadline by working late together",
				"Organized team building activities that improved team cohesion",
				"Supported the team during a difficult project launch, staying positive and encouraging",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "STAY_HUMBLE",
			Name:        "Stay Humble",
			Description: "We are open to feedback and always look for ways to improve.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Maintained humility after a major project success",
				"Welcomed feedback and incorporated suggestions from diverse team members",
				"Asked questions to understand the 'why' behind existing processes",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	for _, v := range values {
		var count int64
		DB.Model(&models.CompanyValue{}).Where("code = ? OR name = ?", v.Code, v.Name).Count(&count)
		if count > 0 {
			continue // 存在即跳过
		}
		if err := DB.Create(&v).Error; err != nil {
			return err
		}
		log.Printf("Seeded company value: %s (%s)", v.Name, v.Code)
	}

	log.Println("Initial data seeding completed successfully")
	return nil
}

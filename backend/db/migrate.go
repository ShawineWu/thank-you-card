package db

import (
	"context"
	"fmt"
	"time"

	"thank-you-card-backend/internal/models"

	"gorm.io/gorm"
)

// AutoMigrateAndSeed runs Gorm migrations and seeds static reference data.
func AutoMigrateAndSeed(gormDB *gorm.DB) error {
	if err := gormDB.AutoMigrate(
		&models.Employee{},
		&models.CompanyValue{},
		&models.Card{},
		&models.CardRecipient{},
		&models.CardValue{},
		&models.EmojiReaction{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	if err := seedCompanyValues(gormDB); err != nil {
		return fmt.Errorf("seed company values: %w", err)
	}

	if err := seedMockEmployees(gormDB); err != nil {
		return fmt.Errorf("seed mock employees: %w", err)
	}

	return nil
}

func seedCompanyValues(gormDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type seedValue struct {
		Code string
		Name string
		Type string
		Desc string
	}

	var seeds = []seedValue{
		// Values
		{"MAKE_IMPACT", "Make an Impact", "VALUE", "Be driven by the desire to build something that can touch millions of lives."},
		{"STRIVE_EXCELLENCE", "Strive for Excellence", "VALUE", "Today's great is not good enough for tomorrow."},
		{"STAND_TOGETHER", "Stand Together", "VALUE", "Embrace the spirit of solidarity & commitment through thick and thin."},
		{"BE_OPEN_MINDED", "Be Open-Minded", "VALUE", "Diversity is strength. Magic often happens at the intersection of two different worlds."},
		{"STAY_GROUNDED", "Stay Grounded", "VALUE", "Always remember our roots and higher purpose. Plus, life is too short to be around jerks."},

		// Credos
		{"BIAS_FOR_ACTION", "Bias for Action", "CREDO", "Speed matters, take action to deliver a high quality result with calculated risk taking. Be time conscious and set deadlines on initiatives. Promptly ask for help when faced with roadblocks, be the roadblock-remover where we can. It is better to be moving than not, most consequences are manageable."},
		{"CUSTOMER_CENTRIC", "Customer Centric", "CREDO", "We start with the customer experience, centring our processes and decisions around it. Never stick to what has been done for the sake of tradition. Advocate for our customers, add value wherever possible and build trust with sincere interactions."},
		{"THINK_STRATEGICALLY", "Think Strategically", "CREDO", "We know our business inside out, acutely aware of the market dynamics and crystal clear of our positioning and what we offer. We stay ahead of the curve by being flexible enough to pivot quickly, being nimble in problem solving and being obsessed with value creation for our customers, employees and partners. We connect the dots where others don't."},
		{"DEEP_DIVE", "Deep Dive", "CREDO", "No detail is too small, no task too unimportant. Always ask why, validate with data, slicing for insights and looking for opportunities. Only when we know the nuts and bolts of everything we work on, can we troubleshoot effectively and innovate creatively. Identify root causes. When in doubt, keep questioning constructively."},
		{"INVENT_SIMPLIFY", "Invent and Simplify", "CREDO", "There is no need to over-complicate. Get creative with solutions, streamline where we can and tap on the expertise of different teams. Innovation and invention are expected and achieved by leveraging all available resources."},
		{"EARN_TRUST", "Earn Trust", "CREDO", "Always deliver on promises. We listen attentively, speak candidly and maintain the highest standards of ethics - towards our customers and our team."},
		{"TAKE_OWNERSHIP", "Take Ownership", "CREDO", "Each of us represents the company, beyond just ourselves or our team. When times are good, we celebrate together. When times are bad, we stay and face the adversity as one. Act with a view of the long term, today's great piece of work will have a lasting impact on the bigger picture. The work we produce speaks for the company, be proud of it, own it."},
		{"CHALLENGE_DISAGREE_COMMIT", "Challenge Disagree and Commit", "CREDO", "When in doubt, challenge respectfully and objectively, even if it is uncomfortable. Unemotionally review the objective of the things we are doing and whether how we are doing it is the best way to do it. Present facts instead of succumbing to feelings. Ideas evolve and improve when scrutinised. Once a decision is determined, commit wholeheartedly and own the consequences together."},
		{"LEARN_BE_CURIOUS", "Learn and Be Curious", "CREDO", "Learn from team mates, customers and competitors. Always be curious about the whys and how things are done, be eager to go deeper and bring our learnings back to our work."},
		{"DO_MORE_WITH_LESS", "Do More with Less", "CREDO", "Know where we should be investing, and invest it better. Everyone is empowered to implement efficient processes and reduce spending. Do so wisely because we cut costs but we don't cut corners."},
	}

	for _, sv := range seeds {
		var existing models.CompanyValue
		if err := gormDB.WithContext(ctx).
			Where("code = ?", sv.Code).
			First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		if existing.ID != 0 {
			// Ensure name/description/type are up to date.
			existing.Name = sv.Name
			existing.Type = sv.Type
			existing.Description = sv.Desc
			if err := gormDB.WithContext(ctx).Save(&existing).Error; err != nil {
				return err
			}
			continue
		}

		cv := models.CompanyValue{
			Code:        sv.Code,
			Name:        sv.Name,
			Type:        sv.Type,
			Description: sv.Desc,
		}
		if err := gormDB.WithContext(ctx).Create(&cv).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedMockEmployees(gormDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	employees := []models.Employee{
		{
			Email:      "employee@example.com",
			Name:       "Mock Employee",
			Department: "Engineering",
			IsHRAdmin:  false,
		},
		{
			Email:      "hradmin@example.com",
			Name:       "Mock HR Admin",
			Department: "HR",
			IsHRAdmin:  true,
		},
	}

	for _, e := range employees {
		var existing models.Employee
		if err := gormDB.WithContext(ctx).
			Where("email = ?", e.Email).
			First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if existing.ID != 0 {
			continue
		}
		if err := gormDB.WithContext(ctx).Create(&e).Error; err != nil {
			return err
		}
	}

	return nil
}

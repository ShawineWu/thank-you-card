package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"thank-you-card-backend/internal/models"

	"gorm.io/gorm"
)

// AutoMigrateAndSeed runs Gorm migrations and seeds static reference data.
func AutoMigrateAndSeed(gormDB *gorm.DB) error {
	if err := gormDB.AutoMigrate(
		&models.User{},
		&models.Employee{},
		&models.CompanyValue{},
		&models.Card{},
		&models.CardRecipient{},
		&models.CardValue{},
		&models.EmojiReaction{},
		&models.Milestone{},
		&models.UserMilestone{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	if err := seedCompanyValues(gormDB); err != nil {
		return fmt.Errorf("seed company values: %w", err)
	}

	if err := seedMockEmployees(gormDB); err != nil {
		return fmt.Errorf("seed mock employees: %w", err)
	}

	if err := seedUsers(gormDB); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	if err := seedMilestones(gormDB); err != nil {
		return fmt.Errorf("seed milestones: %w", err)
	}

	return nil
}

func seedCompanyValues(gormDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type seedValue struct {
		Code     string
		Name     string
		Type     string
		Desc     string
		Examples []string
	}

	var seeds = []seedValue{
		// Values
		{"MAKE_IMPACT", "Make an Impact", "VALUE", "Be driven by the desire to build something that can touch millions of lives.", []string{
			"Led a project that improved customer satisfaction by 30%",
			"Volunteered to mentor new team members and helped 5 people onboard successfully",
			"Identified and fixed a critical bug that affected thousands of users",
		}},
		{"STRIVE_EXCELLENCE", "Strive for Excellence", "VALUE", "Today's great is not good enough for tomorrow.", []string{
			"Continuously improved code quality, reducing technical debt by 40%",
			"Set up automated testing that caught 90% of bugs before production",
			"Redesigned a feature based on user feedback, increasing engagement by 25%",
		}},
		{"STAND_TOGETHER", "Stand Together", "VALUE", "Embrace the spirit of solidarity & commitment through thick and thin.", []string{
			"Helped a teammate meet a tight deadline by working late together",
			"Organized team building activities that improved team cohesion",
			"Supported the team during a difficult project launch, staying positive and encouraging",
		}},
		{"BE_OPEN_MINDED", "Be Open-Minded", "VALUE", "Diversity is strength. Magic often happens at the intersection of two different worlds.", []string{
			"Actively listened to different perspectives in team discussions",
			"Collaborated with cross-functional teams to create innovative solutions",
			"Welcomed feedback and incorporated suggestions from diverse team members",
		}},
		{"STAY_GROUNDED", "Stay Grounded", "VALUE", "Always remember our roots and higher purpose. Plus, life is too short to be around jerks.", []string{
			"Maintained humility after a major project success",
			"Helped create a positive and inclusive team environment",
			"Remembered to celebrate small wins and appreciate team efforts",
		}},

		// Credos
		{"BIAS_FOR_ACTION", "Bias for Action", "CREDO", "Speed matters, take action to deliver a high quality result with calculated risk taking. Be time conscious and set deadlines on initiatives. Promptly ask for help when faced with roadblocks, be the roadblock-remover where we can. It is better to be moving than not, most consequences are manageable.", []string{
			"Quickly prototyped a solution to test a hypothesis within 2 days",
			"Identified a blocker and immediately reached out to the right people to resolve it",
			"Made a decision and moved forward when others were stuck in analysis paralysis",
		}},
		{"CUSTOMER_CENTRIC", "Customer Centric", "CREDO", "We start with the customer experience, centring our processes and decisions around it. Never stick to what has been done for the sake of tradition. Advocate for our customers, add value wherever possible and build trust with sincere interactions.", []string{
			"Conducted user interviews to understand pain points before building a feature",
			"Advocated for a customer-requested feature that wasn't initially prioritized",
			"Fixed a customer-reported issue within 24 hours of receiving the feedback",
		}},
		{"THINK_STRATEGICALLY", "Think Strategically", "CREDO", "We know our business inside out, acutely aware of the market dynamics and crystal clear of our positioning and what we offer. We stay ahead of the curve by being flexible enough to pivot quickly, being nimble in problem solving and being obsessed with value creation for our customers, employees and partners. We connect the dots where others don't.", []string{
			"Identified a market opportunity and proposed a strategic initiative",
			"Connected insights from different departments to create a comprehensive solution",
			"Pivoted project direction based on market research and competitive analysis",
		}},
		{"DEEP_DIVE", "Deep Dive", "CREDO", "No detail is too small, no task too unimportant. Always ask why, validate with data, slicing for insights and looking for opportunities. Only when we know the nuts and bolts of everything we work on, can we troubleshoot effectively and innovate creatively. Identify root causes. When in doubt, keep questioning constructively.", []string{
			"Investigated a performance issue and found the root cause after thorough analysis",
			"Asked probing questions during code review that led to a better solution",
			"Validated assumptions with data before making important decisions",
		}},
		{"INVENT_SIMPLIFY", "Invent and Simplify", "CREDO", "There is no need to over-complicate. Get creative with solutions, streamline where we can and tap on the expertise of different teams. Innovation and invention are expected and achieved by leveraging all available resources.", []string{
			"Simplified a complex process that saved the team 5 hours per week",
			"Created a reusable component that eliminated code duplication across 3 projects",
			"Found a creative solution that reduced implementation time by 50%",
		}},
		{"EARN_TRUST", "Earn Trust", "CREDO", "Always deliver on promises. We listen attentively, speak candidly and maintain the highest standards of ethics - towards our customers and our team.", []string{
			"Delivered a project on time despite unexpected challenges",
			"Admitted a mistake early and worked proactively to fix it",
			"Provided honest feedback in a constructive and respectful manner",
		}},
		{"TAKE_OWNERSHIP", "Take Ownership", "CREDO", "Each of us represents the company, beyond just ourselves or our team. When times are good, we celebrate together. When times are bad, we stay and face the adversity as one. Act with a view of the long term, today's great piece of work will have a lasting impact on the bigger picture. The work we produce speaks for the company, be proud of it, own it.", []string{
			"Took responsibility for a project outcome, both successes and failures",
			"Went beyond assigned tasks to ensure overall project success",
			"Fixed an issue that wasn't technically my responsibility but affected the team",
		}},
		{"CHALLENGE_DISAGREE_COMMIT", "Challenge Disagree and Commit", "CREDO", "When in doubt, challenge respectfully and objectively, even if it is uncomfortable. Unemotionally review the objective of the things we are doing and whether how we are doing it is the best way to do it. Present facts instead of succumbing to feelings. Ideas evolve and improve when scrutinised. Once a decision is determined, commit wholeheartedly and own the consequences together.", []string{
			"Respectfully challenged a technical approach with data and alternative solutions",
			"Disagreed with a decision but fully committed once the team made a choice",
			"Facilitated a healthy debate that led to a better final decision",
		}},
		{"LEARN_BE_CURIOUS", "Learn and Be Curious", "CREDO", "Learn from team mates, customers and competitors. Always be curious about the whys and how things are done, be eager to go deeper and bring our learnings back to our work.", []string{
			"Learned a new technology and applied it to improve our codebase",
			"Asked questions to understand the 'why' behind existing processes",
			"Shared learnings from a conference with the team and implemented best practices",
		}},
		{"DO_MORE_WITH_LESS", "Do More with Less", "CREDO", "Know where we should be investing, and invest it better. Everyone is empowered to implement efficient processes and reduce spending. Do so wisely because we cut costs but we don't cut corners.", []string{
			"Optimized a process that reduced costs by 30% without sacrificing quality",
			"Found a free tool that replaced a paid service, saving the company money",
			"Streamlined workflows to achieve the same results with fewer resources",
		}},
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

		// Serialize examples to JSON
		examplesJSON, _ := json.Marshal(sv.Examples)

		if existing.ID != 0 {
			// Ensure name/description/type/examples are up to date.
			existing.Name = sv.Name
			existing.Type = sv.Type
			existing.Description = sv.Desc
			existing.Examples = string(examplesJSON)
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
			Examples:    string(examplesJSON),
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

func seedUsers(gormDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get employees for linking
	var employee1, hrAdmin models.Employee
	gormDB.WithContext(ctx).Where("email = ?", "employee@example.com").First(&employee1)
	gormDB.WithContext(ctx).Where("email = ?", "hradmin@example.com").First(&hrAdmin)

	// Hash passwords (default password: "password123")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPasswordHR, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPasswordAdmin, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		{
			Username:   "employee",
			Password:   string(hashedPassword),
			Role:       "EMPLOYEE",
			EmployeeID: &employee1.ID,
		},
		{
			Username:   "hr",
			Password:   string(hashedPasswordHR),
			Role:       "HR",
			EmployeeID: &hrAdmin.ID,
		},
		{
			Username:   "admin",
			Password:   string(hashedPasswordAdmin),
			Role:       "ADMIN",
			EmployeeID: &hrAdmin.ID, // Admin also needs an employee for reactions
		},
	}

	for _, u := range users {
		var existing models.User
		if err := gormDB.WithContext(ctx).
			Where("username = ?", u.Username).
			First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if existing.ID != 0 {
			// Update existing user if EmployeeID is missing (e.g., admin user)
			if existing.EmployeeID == nil && u.EmployeeID != nil {
				existing.EmployeeID = u.EmployeeID
				if err := gormDB.WithContext(ctx).Save(&existing).Error; err != nil {
					return err
				}
			}
			continue
		}
		if err := gormDB.WithContext(ctx).Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedMilestones(gormDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type seedMilestone struct {
		Code        string
		Name        string
		Description string
		Type        string
		Threshold   int
	}

	seeds := []seedMilestone{
		// SENT milestones
		{"SENT_5", "Sent 5 Cards", "Congratulations! You've sent 5 thank you cards. Keep spreading appreciation!", "SENT", 5},
		{"SENT_10", "Sent 10 Cards", "Amazing! You've sent 10 thank you cards. Your recognition makes a difference!", "SENT", 10},
		{"SENT_50", "Sent 50 Cards", "Outstanding! You've sent 50 thank you cards. You're a recognition champion!", "SENT", 50},
		{"SENT_100", "Sent 100 Cards", "Incredible! You've sent 100 thank you cards. You're a true appreciation leader!", "SENT", 100},
		// RECEIVED milestones
		{"RECEIVED_5", "Received 5 Cards", "Well done! You've received 5 thank you cards. Your impact is being recognized!", "RECEIVED", 5},
		{"RECEIVED_10", "Received 10 Cards", "Excellent! You've received 10 thank you cards. Your contributions are valued!", "RECEIVED", 10},
		{"RECEIVED_50", "Received 50 Cards", "Remarkable! You've received 50 thank you cards. You're making a real difference!", "RECEIVED", 50},
		{"RECEIVED_100", "Received 100 Cards", "Extraordinary! You've received 100 thank you cards. You're an inspiration!", "RECEIVED", 100},
		// TOTAL milestones
		{"TOTAL_5", "5 Cards Total", "Great start! You've sent or received 5 thank you cards combined.", "TOTAL", 5},
		{"TOTAL_10", "10 Cards Total", "Nice progress! You've sent or received 10 thank you cards combined.", "TOTAL", 10},
		{"TOTAL_50", "50 Cards Total", "Impressive! You've sent or received 50 thank you cards combined.", "TOTAL", 50},
		{"TOTAL_100", "100 Cards Total", "Phenomenal! You've sent or received 100 thank you cards combined.", "TOTAL", 100},
	}

	for _, sm := range seeds {
		var existing models.Milestone
		if err := gormDB.WithContext(ctx).
			Where("code = ?", sm.Code).
			First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		if existing.ID != 0 {
			// Update existing milestone if needed
			existing.Name = sm.Name
			existing.Description = sm.Description
			existing.Type = sm.Type
			existing.Threshold = sm.Threshold
			if err := gormDB.WithContext(ctx).Save(&existing).Error; err != nil {
				return err
			}
			continue
		}

		milestone := models.Milestone{
			Code:        sm.Code,
			Name:        sm.Name,
			Description: sm.Description,
			Type:        sm.Type,
			Threshold:   sm.Threshold,
		}
		if err := gormDB.WithContext(ctx).Create(&milestone).Error; err != nil {
			return err
		}
	}

	return nil
}

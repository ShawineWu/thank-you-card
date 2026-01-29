package database

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/castlery/thank-you-card/internal/models"
)

// Seed populates the database with initial data
func Seed() error {
	if DB == nil {
		return log.New(nil, "", 0).Output(0, "database not initialized")
	}

	log.Println("Seeding initial data...")

	// Seed Company Values (Code must be unique). Exists (by code or name) => skip.
	// Based on the 15 company values from the database
	values := []models.CompanyValue{
		{
			ID:          uuid.New(),
			Code:        "MAKE_IMPACT",
			Name:        "Make an Impact",
			Description: "Be driven by the desire to build something that can touch millions of lives.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Led a project that improved customer satisfaction by 30%",
				"Volunteered to mentor new team members and helped 5 people onboard successfully",
				"Identified and fixed a critical bug that affected thousands of users",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "STRIVE_EXCELLENCE",
			Name:        "Strive for Excellence",
			Description: "Today's great is not good enough for tomorrow.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Continuously improved code quality, reducing technical debt by 40%",
				"Set up automated testing that caught 90% of bugs before production",
				"Redesigned a feature based on user feedback, increasing engagement by 25%",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "STAND_TOGETHER",
			Name:        "Stand Together",
			Description: "Embrace the spirit of solidarity & commitment through thick and thin.",
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
			Code:        "BE_OPEN_MINDED",
			Name:        "Be Open-Minded",
			Description: "Diversity is strength. Magic often happens at the intersection of two different worlds.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Actively listened to different perspectives in team discussions",
				"Collaborated with cross-functional teams to create innovative solutions",
				"Welcomed feedback and incorporated suggestions from diverse team members",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "STAY_GROUNDED",
			Name:        "Stay Grounded",
			Description: "Always remember our roots and higher purpose. Plus, life is too short to be around jerks.",
			Type:        models.ValueTypeValue,
			Examples: datatypes.NewJSONSlice([]string{
				"Maintained humility after a major project success",
				"Helped create a positive and inclusive team environment",
				"Remembered to celebrate small wins and appreciate team efforts",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "BIAS_FOR_ACTION",
			Name:        "Bias for Action",
			Description: "Speed matters, take action to deliver a high quality result with calculated risk taking. Be time conscious and set deadlines on initiatives. Promptly ask for help when faced with roadblocks, be the roadblock-remover where we can. It is better to be moving than not, most consequences are manageable.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Quickly prototyped a solution to test a hypothesis within 2 days",
				"Identified a blocker and immediately reached out to the right people to resolve it",
				"Made a decision and moved forward when others were stuck in analysis paralysis",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "CUSTOMER_CENTRIC",
			Name:        "Customer Centric",
			Description: "We start with the customer experience, centring our processes and decisions around it. Never stick to what has been done for the sake of tradition. Advocate for our customers, add value wherever possible and build trust with sincere interactions.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Conducted user interviews to understand pain points before building a feature",
				"Advocated for a customer-requested feature that wasn't initially prioritized",
				"Fixed a customer-reported issue within 24 hours of receiving the feedback",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "THINK_STRATEGICALLY",
			Name:        "Think Strategically",
			Description: "We know our business inside out, acutely aware of the market dynamics and crystal clear of our positioning and what we offer. We stay ahead of the curve by being flexible enough to pivot quickly, being nimble in problem solving and being obsessed with value creation for our customers, employees and partners. We connect the dots where others don't.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Identified a market opportunity and proposed a strategic initiative",
				"Connected insights from different departments to create a comprehensive solution",
				"Pivoted project direction based on market research and competitive analysis",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "DEEP_DIVE",
			Name:        "Deep Dive",
			Description: "No detail is too small, no task too unimportant. Always ask why, validate with data, slicing for insights and looking for opportunities. Only when we know the nuts and bolts of everything we work on, can we troubleshoot effectively and innovate creatively. Identify root causes. When in doubt, keep questioning constructively.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Investigated a performance issue and found the root cause after thorough analysis",
				"Asked probing questions during code review that led to a better solution",
				"Validated assumptions with data before making important decisions",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "INVENT_SIMPLIFY",
			Name:        "Invent and Simplify",
			Description: "There is no need to over-complicate. Get creative with solutions, streamline where we can and tap on the expertise of different teams. Innovation and invention are expected and achieved by leveraging all available resources.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Simplified a complex process that saved the team 5 hours per week",
				"Created a reusable component that eliminated code duplication across 3 projects",
				"Found a creative solution that reduced implementation time by 50%",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "EARN_TRUST",
			Name:        "Earn Trust",
			Description: "Always deliver on promises. We listen attentively, speak candidly and maintain the highest standards of ethics - towards our customers and our team.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Delivered a project on time despite unexpected challenges",
				"Admitted a mistake early and worked proactively to fix it",
				"Provided honest feedback in a constructive and respectful manner",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "TAKE_OWNERSHIP",
			Name:        "Take Ownership",
			Description: "Each of us represents the company, beyond just ourselves or our team. When times are good, we celebrate together. When times are bad, we stay and face the adversity as one. Act with a view of the long term, today's great piece of work will have a lasting impact on the bigger picture. The work we produce speaks for the company, be proud of it, own it.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Took responsibility for a project outcome, both successes and failures",
				"Went beyond assigned tasks to ensure overall project success",
				"Fixed an issue that wasn't technically my responsibility but affected the team",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "CHALLENGE_DISAGREE_COMMIT",
			Name:        "Challenge Disagree and Commit",
			Description: "When in doubt, challenge respectfully and objectively, even if it is uncomfortable. Unemotionally review the objective of the things we are doing and whether how we are doing it is the best way to do it. Present facts instead of succumbing to feelings. Ideas evolve and improve when scrutinised. Once a decision is determined, commit wholeheartedly and own the consequences together.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Respectfully challenged a technical approach with data and alternative solutions",
				"Disagreed with a decision but fully committed once the team made a choice",
				"Facilitated a healthy debate that led to a better final decision",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "LEARN_BE_CURIOUS",
			Name:        "Learn and Be Curious",
			Description: "Learn from team mates, customers and competitors. Always be curious about the whys and how things are done, be eager to go deeper and bring our learnings back to our work.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Learned a new technology and applied it to improve our codebase",
				"Asked questions to understand the 'why' behind existing processes",
				"Shared learnings from a conference with the team and implemented best practices",
			}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			Code:        "DO_MORE_WITH_LESS",
			Name:        "Do More with Less",
			Description: "Know where we should be investing, and invest it better. Everyone is empowered to implement efficient processes and reduce spending. Do so wisely because we cut costs but we don't cut corners.",
			Type:        models.ValueTypeCredo,
			Examples: datatypes.NewJSONSlice([]string{
				"Optimized a process that reduced costs by 30% without sacrificing quality",
				"Found a free tool that replaced a paid service, saving the company money",
				"Streamlined workflows to achieve the same results with fewer resources",
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

	// Seed sample milestones for demo users
	if err := seedMilestones(); err != nil {
		log.Printf("Warning: Failed to seed milestones: %v", err)
	}

	log.Println("Initial data seeding completed successfully")
	return nil
}

// seedMilestones creates sample milestone achievements for demo purposes
func seedMilestones() error {
	// Sample employee IDs - these would typically come from your auth system
	// Using common test user IDs
	sampleEmployees := []string{
		"demo-user-001",
		"demo-user-002",
	}

	milestones := []models.EmployeeMilestone{
		// Demo user 1 - has achieved several milestones
		{
			EmployeeID:    sampleEmployees[0],
			MilestoneType: models.MilestoneTypeCardsReceived,
			Threshold:     5,
			AchievedAt:    time.Now().AddDate(0, -3, 0).UTC(),
			Title:         "5 Cards Received",
			Description:   "Congratulations! You've received 5 recognition cards. Your great work is being noticed!",
		},
		{
			EmployeeID:    sampleEmployees[0],
			MilestoneType: models.MilestoneTypeCardsReceived,
			Threshold:     10,
			AchievedAt:    time.Now().AddDate(0, -1, 0).UTC(),
			Title:         "10 Cards Received",
			Description:   "Impressive! You've received 10 recognition cards. Keep up the excellent work!",
		},
		{
			EmployeeID:    sampleEmployees[0],
			MilestoneType: models.MilestoneTypeCardsSent,
			Threshold:     5,
			AchievedAt:    time.Now().AddDate(0, -2, 0).UTC(),
			Title:         "5 Cards Sent",
			Description:   "Amazing! You've sent 5 thank you cards, recognizing and appreciating your colleagues.",
		},
		// Demo user 2 - fewer milestones
		{
			EmployeeID:    sampleEmployees[1],
			MilestoneType: models.MilestoneTypeCardsReceived,
			Threshold:     5,
			AchievedAt:    time.Now().AddDate(0, -1, -15).UTC(),
			Title:         "5 Cards Received",
			Description:   "Congratulations! You've received 5 recognition cards. Your great work is being noticed!",
		},
	}

	for _, m := range milestones {
		// Check if milestone already exists
		var count int64
		DB.Model(&models.EmployeeMilestone{}).
			Where("employee_id = ? AND milestone_type = ? AND threshold = ?",
				m.EmployeeID, m.MilestoneType, m.Threshold).
			Count(&count)
		if count > 0 {
			continue // Already exists, skip
		}
		if err := DB.Create(&m).Error; err != nil {
			log.Printf("Warning: Failed to seed milestone for %s: %v", m.EmployeeID, err)
			continue
		}
		log.Printf("Seeded milestone: %s - %s (%d)", m.EmployeeID, m.MilestoneType, m.Threshold)
	}

	return nil
}

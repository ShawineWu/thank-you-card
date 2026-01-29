package dto

import "time"

type MilestoneResponse struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Threshold   int     `json:"threshold"`
	IconURL     *string `json:"iconUrl,omitempty"`
}

type UserAchievementResponse struct {
	Milestone  MilestoneResponse `json:"milestone"`
	AchievedAt time.Time         `json:"achievedAt"`
}

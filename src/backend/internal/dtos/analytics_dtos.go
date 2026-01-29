package dtos

import "time"

// DashboardResponse contains overall analytics for HR dashboard
type DashboardResponse struct {
	TotalCards     int64                `json:"totalCards"`
	ActiveUsers    int64                `json:"activeUsers"`
	TopValues      []ValueUsageResponse `json:"topValues"`
	CardsTrend     TrendResponse        `json:"cardsTrend"`
	EngagementRate float64              `json:"engagementRate"`
}

// ValueUsageResponse represents value usage statistics
type ValueUsageResponse struct {
	ValueID    string  `json:"valueId"`
	ValueName  string  `json:"valueName"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage,omitempty"`
}

// TrendResponse shows time-based trends
type TrendResponse struct {
	ThisMonth int64 `json:"thisMonth"`
	LastMonth int64 `json:"lastMonth"`
}

// TopRecognizerResponse represents a top recognizer
type TopRecognizerResponse struct {
	EmployeeID   string `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	CardsSent    int64  `json:"cardsSent"`
	Rank         int    `json:"rank"`
}

// TeamAnalyticsResponse contains team-level recognition data
type TeamAnalyticsResponse struct {
	TeamID        string `json:"teamId"`
	TeamName      string `json:"teamName"`
	CardsSent     int64  `json:"cardsSent"`
	CardsReceived int64  `json:"cardsReceived"`
	TopValue      string `json:"topValue"`
}

// ValueDistributionResponse shows distribution of company values
type ValueDistributionResponse struct {
	Distribution []ValueUsageResponse `json:"distribution"`
	TotalCards   int64                `json:"totalCards"`
}

// ExportRequest defines parameters for CSV export
type ExportRequest struct {
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Format    string `json:"format" binding:"required,oneof=csv"`
}

// ExportFilters contains filter parameters for export
type ExportFilters struct {
	StartDate time.Time
	EndDate   time.Time
}

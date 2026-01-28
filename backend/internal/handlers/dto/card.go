package dto

type CreateCardRequest struct {
	RecipientIDs []uint `json:"recipientIds" binding:"required,min=1"` // selected recipients
	ValueIDs     []uint `json:"valueIds" binding:"required,min=1,max=3"`
	Reason       string `json:"reason" binding:"required"`
}

type EmployeeSummary struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

type CompanyValueSummary struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type EmojiSummary struct {
	EmojiCode string `json:"emojiCode"`
	Count     int    `json:"count"`
	UserIDs   []uint `json:"userIds"`
}

type CardResponse struct {
	ID         uint                  `json:"id"`
	Sender     EmployeeSummary       `json:"sender"`
	Recipients []EmployeeSummary     `json:"recipients"`
	Reason     string                `json:"reason"`
	Values     []CompanyValueSummary `json:"values"`
	CreatedAt  string                `json:"createdAt"`
	Reactions  []EmojiSummary        `json:"reactions"`
}

// Teams-specific simplified response
type TeamsCardResponse struct {
	ID             uint     `json:"id"`
	SenderName     string   `json:"senderName"`
	RecipientNames []string `json:"recipientNames"`
	Reason         string   `json:"reason"`
	Values         []string `json:"values"`
	CreatedAt      string   `json:"createdAt"`
	ReactionCount  int      `json:"reactionCount"`
}

// Stats DTOs
type PersonalStatsResponse struct {
	TotalSent         int64                `json:"totalSent"`
	TotalReceived     int64                `json:"totalReceived"`
	TopSentValues     []ValueCountResponse `json:"topSentValues"`
	TopReceivedValues []ValueCountResponse `json:"topReceivedValues"`
}

type ValueCountResponse struct {
	ValueID uint   `json:"valueId"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Count   int64  `json:"count"`
}

// Analytics DTOs
type DashboardResponse struct {
	TotalCards        int64  `json:"totalCards"`
	TotalEmployees    int64  `json:"totalEmployees"`
	TotalDepartments  int64  `json:"totalDepartments"`
	ActiveRecognizers int64  `json:"activeRecognizers"`
	Period            string `json:"period"`
}

type TopEmployeeResponse struct {
	EmployeeID uint   `json:"employeeId"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Count      int64  `json:"count"`
}

type TeamPatternResponse struct {
	Department              string  `json:"department"`
	TotalSent               int64   `json:"totalSent"`
	TotalReceived           int64   `json:"totalReceived"`
	AverageCardsPerEmployee float64 `json:"averageCardsPerEmployee"`
}

type ValueDistributionResponse struct {
	ValueID    uint    `json:"valueId"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

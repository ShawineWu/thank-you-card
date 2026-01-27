package card

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs

// CreateCardRequest represents a request to create a new card
type CreateCardRequest struct {
	SenderID     string      `json:"sender_id" binding:"required"`
	RecipientIDs []string    `json:"recipient_ids" binding:"required,min=1,max=10"`
	Reason       string      `json:"reason" binding:"required,min=10,max=1000"`
	ValueIDs     []uuid.UUID `json:"value_ids" binding:"required,min=1,max=5"`
}

// ShareCardRequest represents a request to share a card
type ShareCardRequest struct {
	CardID        uuid.UUID `json:"card_id" binding:"required"`
	SharedBy      string    `json:"shared_by" binding:"required"`
	TargetChannel string    `json:"target_channel" binding:"required"`
}

// GetCardsRequest represents a request to get cards with filters
type GetCardsRequest struct {
	SenderID    *string     `json:"sender_id,omitempty"`
	RecipientID *string     `json:"recipient_id,omitempty"`
	ValueIDs    []uuid.UUID `json:"value_ids,omitempty"`
	DateFrom    *string     `json:"date_from,omitempty"`
	DateTo      *string     `json:"date_to,omitempty"`
	SearchText  *string     `json:"search_text,omitempty"`
	Page        int         `json:"page" binding:"min=1"`
	PageSize    int         `json:"page_size" binding:"min=1,max=100"`
}

// Response DTOs

// CardResponse represents a card in API responses
type CardResponse struct {
	ID         uuid.UUID              `json:"id"`
	SenderID   string                 `json:"sender_id"`
	Reason     string                 `json:"reason"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	Recipients []CardRecipientResponse `json:"recipients"`
	Values     []CardValueResponse     `json:"values"`
}

// CardRecipientResponse represents a card recipient in API responses
type CardRecipientResponse struct {
	ID          uuid.UUID `json:"id"`
	RecipientID string    `json:"recipient_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// CardValueResponse represents a card value in API responses
type CardValueResponse struct {
	ID        uuid.UUID `json:"id"`
	ValueID   uuid.UUID `json:"value_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CompanyValueResponse represents a company value in API responses
type CompanyValueResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EmployeeResponse represents an employee in API responses
type EmployeeResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department,omitempty"`
	Team       string `json:"team,omitempty"`
}

// PersonalStatisticsResponse represents personal statistics in API responses
type PersonalStatisticsResponse struct {
	CardsSent         int                         `json:"cards_sent"`
	CardsReceived     int                         `json:"cards_received"`
	ValueDistribution []ValueDistributionResponse `json:"value_distribution"`
	RecentActivity    []RecentActivityResponse    `json:"recent_activity"`
}

// ValueDistributionResponse represents value distribution in API responses
type ValueDistributionResponse struct {
	ValueID    uuid.UUID `json:"value_id"`
	ValueName  string    `json:"value_name"`
	Count      int       `json:"count"`
	Percentage float64   `json:"percentage"`
}

// RecentActivityResponse represents recent activity in API responses
type RecentActivityResponse struct {
	Type       string    `json:"type"`
	CardID     uuid.UUID `json:"card_id"`
	Date       time.Time `json:"date"`
	OtherParty string    `json:"other_party"`
}

// Top10EmployeeResponse represents a top 10 employee in API responses
type Top10EmployeeResponse struct {
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	CardCount    int    `json:"card_count"`
	Rank         int    `json:"rank"`
}

// MilestoneResponse represents a milestone in API responses
type MilestoneResponse struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    string    `json:"employee_id"`
	MilestoneType string    `json:"milestone_type"`
	Threshold     int       `json:"threshold"`
	AchievedAt    time.Time `json:"achieved_at"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page         int  `json:"page"`
	PageSize     int  `json:"page_size"`
	TotalPages   int  `json:"total_pages"`
	TotalItems   int  `json:"total_items"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

// ApiResponse represents a standard API response
type ApiResponse[T any] struct {
	Success   bool      `json:"success"`
	Data      T         `json:"data,omitempty"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Success   bool                 `json:"success"`
	Error     ErrorDetail          `json:"error"`
	Timestamp time.Time            `json:"timestamp"`
	RequestID string               `json:"request_id"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Helper functions for creating responses

// NewApiResponse creates a new successful API response
func NewApiResponse[T any](data T, message string, requestID string) ApiResponse[T] {
	return ApiResponse[T]{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

// NewErrorResponse creates a new error API response
func NewErrorResponse(code, message, requestID string, details []ValidationError) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, page, pageSize, totalItems int) PaginatedResponse[T] {
	totalPages := (totalItems + pageSize - 1) / pageSize
	
	return PaginatedResponse[T]{
		Data: data,
		Pagination: PaginationInfo{
			Page:        page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalItems:  totalItems,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}
}
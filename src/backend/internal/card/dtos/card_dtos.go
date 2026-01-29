package dtos

import (
	"time"

	"github.com/google/uuid"
)

// CreateCardRequest defines the payload for creating a new card
type CreateCardRequest struct {
	Recipients        []RecipientRequest `json:"recipients" binding:"required,min=1"`
	RecognitionReason string             `json:"recognitionReason" binding:"required,min=10,max=1000"`
	ValueIDs          []string           `json:"valueIds" binding:"required,min=1,max=3"`
}

type RecipientRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// EmployeeDTO represents an employee from the directory
type EmployeeDTO struct {
	ID         string `json:"id"`
	AADID      string `json:"aadId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

// CardResponse represents the standardized card data in responses
type CardResponse struct {
	ID                uuid.UUID           `json:"id"`
	SenderID          string              `json:"senderId"`
	SenderName        string              `json:"senderName"`
	Recipients        []RecipientResponse `json:"recipients"`
	RecognitionReason string              `json:"recognitionReason"`
	SelectedValues    []ValueResponse     `json:"selectedValues"`
	CreatedAt         time.Time           `json:"createdAt"`
}

type RecipientResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ValueResponse represents a company value or credo in responses
type ValueResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "Value" or "Credo"
}

// PaginationParams represents the query parameters for pagination
type PaginationParams struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"pageSize,default=20" binding:"min=1,max=100"`
}

// CardFilterParams represents the query parameters for filtering cards
type CardFilterParams struct {
	PaginationParams
	SenderID    string      `form:"senderId"`
	RecipientID string      `form:"recipientId"`
	ValueIDs    []uuid.UUID `form:"valueIds[]"`
	StartDate   string      `form:"startDate"`
	EndDate     string      `form:"endDate"`
	Search      string      `form:"search"`
}

// PaginationResponse metadata for paginated results
type PaginationResponse struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int   `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

// CardListResponse represents a paginated list of cards
type CardListResponse struct {
	Data       []CardResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// StandardResponse follows the standard response format
type StandardResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"requestId,omitempty"`
}

// ErrorResponse follows the error response format
type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     ErrorInfo `json:"error"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"requestId,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// UserStatsResponse represents statistics for a specific user
type UserStatsResponse struct {
	EmployeeID    string              `json:"employeeId"`
	CardsSent     int64               `json:"cardsSent"`
	CardsReceived int64               `json:"cardsReceived"`
	ValueStats    []ValueUsageStat    `json:"valueStats"`
	Milestones    []MilestoneResponse `json:"milestones"`
}

// ValueUsageStat represents how many times a specific value has been associated with the user's cards
type ValueUsageStat struct {
	ValueID   uuid.UUID `json:"valueId"`
	ValueName string    `json:"valueName"`
	Count     int64     `json:"count"`
}

// MilestoneResponse represents a milestone achievement
type MilestoneResponse struct {
	Type        string    `json:"type"`
	Threshold   int       `json:"threshold"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AchievedAt  time.Time `json:"achievedAt"`
}

// TopEmployeeStat represents an employee in the Top 10 list
type TopEmployeeStat struct {
	EmployeeID    string `json:"employeeId"`
	CardsReceived int64  `json:"cardsReceived"`
	Rank          int    `json:"rank"`
}

package card

import (
	"fmt"
	"time"

	"github.com/company/thank-you-card/internal/domain/card"
)

// Domain to DTO mappers

// MapCardToResponse converts a domain Card to CardResponse
func MapCardToResponse(domainCard *card.Card) CardResponse {
	response := CardResponse{
		ID:        domainCard.ID,
		SenderID:  domainCard.SenderID,
		Reason:    domainCard.Reason,
		CreatedAt: domainCard.CreatedAt,
		UpdatedAt: domainCard.UpdatedAt,
	}

	// Map recipients
	for _, recipient := range domainCard.Recipients {
		response.Recipients = append(response.Recipients, CardRecipientResponse{
			ID:          recipient.ID,
			RecipientID: recipient.RecipientID,
			CreatedAt:   recipient.CreatedAt,
		})
	}

	// Map values
	for _, value := range domainCard.Values {
		response.Values = append(response.Values, CardValueResponse{
			ID:        value.ID,
			ValueID:   value.ValueID,
			CreatedAt: value.CreatedAt,
		})
	}

	return response
}

// MapCardsToResponse converts a slice of domain Cards to CardResponse slice
func MapCardsToResponse(domainCards []card.Card) []CardResponse {
	responses := make([]CardResponse, len(domainCards))
	for i, domainCard := range domainCards {
		responses[i] = MapCardToResponse(&domainCard)
	}
	return responses
}

// MapCompanyValueToResponse converts a domain CompanyValue to CompanyValueResponse
func MapCompanyValueToResponse(domainValue *card.CompanyValue) CompanyValueResponse {
	return CompanyValueResponse{
		ID:          domainValue.ID,
		Name:        domainValue.Name,
		Description: domainValue.Description,
		Type:        domainValue.Type,
		CreatedAt:   domainValue.CreatedAt,
		UpdatedAt:   domainValue.UpdatedAt,
	}
}

// MapCompanyValuesToResponse converts a slice of domain CompanyValues to CompanyValueResponse slice
func MapCompanyValuesToResponse(domainValues []card.CompanyValue) []CompanyValueResponse {
	responses := make([]CompanyValueResponse, len(domainValues))
	for i, domainValue := range domainValues {
		responses[i] = MapCompanyValueToResponse(&domainValue)
	}
	return responses
}

// MapEmployeeToResponse converts a domain Employee to EmployeeResponse
func MapEmployeeToResponse(domainEmployee *card.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:         domainEmployee.ID,
		Name:       domainEmployee.Name,
		Email:      domainEmployee.Email,
		Department: domainEmployee.Department,
		Team:       domainEmployee.Team,
	}
}

// MapEmployeesToResponse converts a slice of domain Employees to EmployeeResponse slice
func MapEmployeesToResponse(domainEmployees []card.Employee) []EmployeeResponse {
	responses := make([]EmployeeResponse, len(domainEmployees))
	for i, domainEmployee := range domainEmployees {
		responses[i] = MapEmployeeToResponse(&domainEmployee)
	}
	return responses
}

// MapPersonalStatisticsToResponse converts domain PersonalStatistics to PersonalStatisticsResponse
func MapPersonalStatisticsToResponse(domainStats *card.PersonalStatistics) PersonalStatisticsResponse {
	response := PersonalStatisticsResponse{
		CardsSent:     domainStats.CardsSent,
		CardsReceived: domainStats.CardsReceived,
	}

	// Map value distribution
	for _, dist := range domainStats.ValueDistribution {
		response.ValueDistribution = append(response.ValueDistribution, ValueDistributionResponse{
			ValueID:    dist.ValueID,
			ValueName:  dist.ValueName,
			Count:      dist.Count,
			Percentage: dist.Percentage,
		})
	}

	// Map recent activity
	for _, activity := range domainStats.RecentActivity {
		response.RecentActivity = append(response.RecentActivity, RecentActivityResponse{
			Type:       activity.Type,
			CardID:     activity.CardID,
			Date:       activity.Date,
			OtherParty: activity.OtherParty,
		})
	}

	return response
}

// MapTop10EmployeesToResponse converts domain Top10Employee slice to Top10EmployeeResponse slice
func MapTop10EmployeesToResponse(domainTop10 []card.Top10Employee) []Top10EmployeeResponse {
	responses := make([]Top10EmployeeResponse, len(domainTop10))
	for i, domainEmployee := range domainTop10 {
		responses[i] = Top10EmployeeResponse{
			EmployeeID:   domainEmployee.EmployeeID,
			EmployeeName: domainEmployee.EmployeeName,
			CardCount:    domainEmployee.CardCount,
			Rank:         domainEmployee.Rank,
		}
	}
	return responses
}

// MapMilestoneToResponse converts a domain EmployeeMilestone to MilestoneResponse
func MapMilestoneToResponse(domainMilestone *card.EmployeeMilestone) MilestoneResponse {
	return MilestoneResponse{
		ID:            domainMilestone.ID,
		EmployeeID:    domainMilestone.EmployeeID,
		MilestoneType: domainMilestone.MilestoneType,
		Threshold:     domainMilestone.Threshold,
		AchievedAt:    domainMilestone.AchievedAt,
		Title:         domainMilestone.Title,
		Description:   domainMilestone.Description,
	}
}

// MapMilestonesToResponse converts a slice of domain EmployeeMilestones to MilestoneResponse slice
func MapMilestonesToResponse(domainMilestones []card.EmployeeMilestone) []MilestoneResponse {
	responses := make([]MilestoneResponse, len(domainMilestones))
	for i, domainMilestone := range domainMilestones {
		responses[i] = MapMilestoneToResponse(&domainMilestone)
	}
	return responses
}

// DTO to Domain mappers

// MapCreateCardRequestToSearchFilters converts GetCardsRequest to domain SearchFilters
func MapGetCardsRequestToSearchFilters(request GetCardsRequest) card.SearchFilters {
	return card.SearchFilters{
		SenderID:    request.SenderID,
		RecipientID: request.RecipientID,
		ValueIDs:    request.ValueIDs,
		DateFrom:    request.DateFrom,
		DateTo:      request.DateTo,
		SearchText:  request.SearchText,
	}
}

// Validation helpers

// ValidateCreateCardRequest validates a CreateCardRequest
func ValidateCreateCardRequest(request CreateCardRequest) []ValidationError {
	var errors []ValidationError

	// Validate sender ID
	if request.SenderID == "" {
		errors = append(errors, ValidationError{
			Field:   "sender_id",
			Message: "sender ID is required",
		})
	}

	// Validate recipients
	if len(request.RecipientIDs) == 0 {
		errors = append(errors, ValidationError{
			Field:   "recipient_ids",
			Message: "at least one recipient is required",
		})
	} else if len(request.RecipientIDs) > 10 {
		errors = append(errors, ValidationError{
			Field:   "recipient_ids",
			Message: "cannot have more than 10 recipients",
		})
	}

	// Check for duplicate recipients
	recipientMap := make(map[string]bool)
	for i, recipientID := range request.RecipientIDs {
		if recipientID == "" {
			errors = append(errors, ValidationError{
				Field:   "recipient_ids",
				Message: fmt.Sprintf("recipient at index %d cannot be empty", i),
			})
		} else if recipientMap[recipientID] {
			errors = append(errors, ValidationError{
				Field:   "recipient_ids",
				Message: fmt.Sprintf("duplicate recipient: %s", recipientID),
			})
		} else {
			recipientMap[recipientID] = true
		}

		// Check if recipient is the sender
		if recipientID == request.SenderID {
			errors = append(errors, ValidationError{
				Field:   "recipient_ids",
				Message: "cannot send a card to yourself",
			})
		}
	}

	// Validate reason
	if len(request.Reason) < 10 {
		errors = append(errors, ValidationError{
			Field:   "reason",
			Message: "reason must be at least 10 characters long",
		})
	} else if len(request.Reason) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "reason",
			Message: "reason cannot exceed 1000 characters",
		})
	}

	// Validate values
	if len(request.ValueIDs) == 0 {
		errors = append(errors, ValidationError{
			Field:   "value_ids",
			Message: "at least one company value must be selected",
		})
	} else if len(request.ValueIDs) > 5 {
		errors = append(errors, ValidationError{
			Field:   "value_ids",
			Message: "cannot select more than 5 company values",
		})
	}

	// Check for duplicate values
	valueMap := make(map[string]bool)
	for _, valueID := range request.ValueIDs {
		valueStr := valueID.String()
		if valueMap[valueStr] {
			errors = append(errors, ValidationError{
				Field:   "value_ids",
				Message: fmt.Sprintf("duplicate value selected: %s", valueStr),
			})
		} else {
			valueMap[valueStr] = true
		}
	}

	return errors
}

// ValidateGetCardsRequest validates a GetCardsRequest
func ValidateGetCardsRequest(request GetCardsRequest) []ValidationError {
	var errors []ValidationError

	// Validate pagination
	if request.Page < 1 {
		errors = append(errors, ValidationError{
			Field:   "page",
			Message: "page must be at least 1",
		})
	}

	if request.PageSize < 1 {
		errors = append(errors, ValidationError{
			Field:   "page_size",
			Message: "page size must be at least 1",
		})
	} else if request.PageSize > 100 {
		errors = append(errors, ValidationError{
			Field:   "page_size",
			Message: "page size cannot exceed 100",
		})
	}

	// Validate date formats if provided
	if request.DateFrom != nil {
		if _, err := time.Parse("2006-01-02", *request.DateFrom); err != nil {
			errors = append(errors, ValidationError{
				Field:   "date_from",
				Message: "date_from must be in YYYY-MM-DD format",
			})
		}
	}

	if request.DateTo != nil {
		if _, err := time.Parse("2006-01-02", *request.DateTo); err != nil {
			errors = append(errors, ValidationError{
				Field:   "date_to",
				Message: "date_to must be in YYYY-MM-DD format",
			})
		}
	}

	return errors
}
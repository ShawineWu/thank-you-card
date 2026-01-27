package card

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ValidationError represents a domain validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidateRecognitionReason validates the recognition reason text
func ValidateRecognitionReason(reason string) error {
	// Trim whitespace
	reason = strings.TrimSpace(reason)
	
	// Check if empty
	if reason == "" {
		return ValidationError{
			Field:   "reason",
			Message: "recognition reason cannot be empty",
		}
	}
	
	// Check minimum length (10 characters)
	if utf8.RuneCountInString(reason) < 10 {
		return ValidationError{
			Field:   "reason",
			Message: "recognition reason must be at least 10 characters long",
		}
	}
	
	// Check maximum length (1000 characters)
	if utf8.RuneCountInString(reason) > 1000 {
		return ValidationError{
			Field:   "reason",
			Message: "recognition reason cannot exceed 1000 characters",
		}
	}
	
	// Check for inappropriate content (basic check)
	if containsInappropriateContent(reason) {
		return ValidationError{
			Field:   "reason",
			Message: "recognition reason contains inappropriate content",
		}
	}
	
	return nil
}

// ValidateEmployeeID validates an employee ID format
func ValidateEmployeeID(employeeID string) error {
	// Trim whitespace
	employeeID = strings.TrimSpace(employeeID)
	
	// Check if empty
	if employeeID == "" {
		return ValidationError{
			Field:   "employee_id",
			Message: "employee ID cannot be empty",
		}
	}
	
	// Check length (assuming employee IDs are between 3 and 50 characters)
	if len(employeeID) < 3 || len(employeeID) > 50 {
		return ValidationError{
			Field:   "employee_id",
			Message: "employee ID must be between 3 and 50 characters",
		}
	}
	
	return nil
}

// ValidateRecipients validates a list of recipient IDs
func ValidateRecipients(recipientIDs []string, senderID string) error {
	// Check if empty
	if len(recipientIDs) == 0 {
		return ValidationError{
			Field:   "recipients",
			Message: "at least one recipient is required",
		}
	}
	
	// Check maximum number of recipients (business rule: max 10 recipients per card)
	if len(recipientIDs) > 10 {
		return ValidationError{
			Field:   "recipients",
			Message: "cannot have more than 10 recipients per card",
		}
	}
	
	// Validate each recipient ID and check for duplicates
	seen := make(map[string]bool)
	for i, recipientID := range recipientIDs {
		// Validate format
		if err := ValidateEmployeeID(recipientID); err != nil {
			return ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: err.Error(),
			}
		}
		
		// Check if recipient is the sender
		if recipientID == senderID {
			return ValidationError{
				Field:   "recipients",
				Message: "cannot send a card to yourself",
			}
		}
		
		// Check for duplicates
		if seen[recipientID] {
			return ValidationError{
				Field:   "recipients",
				Message: fmt.Sprintf("duplicate recipient: %s", recipientID),
			}
		}
		seen[recipientID] = true
	}
	
	return nil
}

// ValidateCompanyValues validates selected company values
func ValidateCompanyValues(valueIDs []string) error {
	// Check if empty
	if len(valueIDs) == 0 {
		return ValidationError{
			Field:   "values",
			Message: "at least one company value must be selected",
		}
	}
	
	// Check maximum number of values (business rule: max 5 values per card)
	if len(valueIDs) > 5 {
		return ValidationError{
			Field:   "values",
			Message: "cannot select more than 5 company values per card",
		}
	}
	
	// Check for duplicates
	seen := make(map[string]bool)
	for i, valueID := range valueIDs {
		if valueID == "" {
			return ValidationError{
				Field:   fmt.Sprintf("values[%d]", i),
				Message: "value ID cannot be empty",
			}
		}
		
		if seen[valueID] {
			return ValidationError{
				Field:   "values",
				Message: fmt.Sprintf("duplicate value selected: %s", valueID),
			}
		}
		seen[valueID] = true
	}
	
	return nil
}

// containsInappropriateContent performs basic inappropriate content detection
func containsInappropriateContent(text string) bool {
	// Convert to lowercase for case-insensitive matching
	lowerText := strings.ToLower(text)
	
	// Basic inappropriate words list (this would be more comprehensive in production)
	inappropriateWords := []string{
		"hate", "stupid", "idiot", "dumb", "loser", "failure",
		// Add more inappropriate words as needed
	}
	
	for _, word := range inappropriateWords {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	
	return false
}

// Business Rules and Policies

// CanSendCard checks if a sender can send a card based on business rules
func CanSendCard(senderID string, recipientIDs []string) error {
	// Validate sender ID
	if err := ValidateEmployeeID(senderID); err != nil {
		return fmt.Errorf("invalid sender: %w", err)
	}
	
	// Validate recipients
	if err := ValidateRecipients(recipientIDs, senderID); err != nil {
		return fmt.Errorf("invalid recipients: %w", err)
	}
	
	// Additional business rules can be added here:
	// - Rate limiting (e.g., max cards per day)
	// - Department restrictions
	// - Approval workflows
	
	return nil
}

// CanViewCard checks if a user can view a specific card
func CanViewCard(userID string, card *Card) bool {
	// Business rule: All cards are publicly visible to everyone in the company
	// This could be extended to include privacy settings in the future
	return true
}

// CanEditCard checks if a user can edit a specific card
func CanEditCard(userID string, card *Card) bool {
	// Business rule: Only the sender can edit their own cards
	// And only within a certain time window (e.g., 24 hours)
	return userID == card.SenderID
}

// CanDeleteCard checks if a user can delete a specific card
func CanDeleteCard(userID string, card *Card) bool {
	// Business rule: Only the sender can delete their own cards
	// And only within a certain time window (e.g., 24 hours)
	return userID == card.SenderID
}
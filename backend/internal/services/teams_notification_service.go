package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"thank-you-card-backend/internal/models"
)

// TeamsNotificationService handles sending notifications to Teams channel via webhook
type TeamsNotificationService interface {
	SendCardNotification(ctx context.Context, card *models.Card) error
}

type teamsNotificationService struct {
	webhookURL string
	httpClient *http.Client
}

// NewTeamsNotificationService creates a new Teams notification service
func NewTeamsNotificationService() TeamsNotificationService {
	webhookURL := os.Getenv("THANK_YOU_CARD_WEB_HOOK")
	if webhookURL == "" {
		// Return a no-op service if webhook URL is not configured
		return &noOpTeamsNotificationService{}
	}

	return &teamsNotificationService{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendCardNotification sends a card notification to Teams channel
func (s *teamsNotificationService) SendCardNotification(ctx context.Context, card *models.Card) error {
	// Build recipient names
	recipientNames := make([]string, len(card.Recipients))
	for i, r := range card.Recipients {
		recipientNames[i] = r.Recipient.Name
	}

	// Build value names
	valueNames := make([]string, len(card.Values))
	for i, v := range card.Values {
		valueNames[i] = v.CompanyValue.Name
	}

	// Format recipients
	recipientsText := ""
	if len(recipientNames) == 1 {
		recipientsText = recipientNames[0]
	} else if len(recipientNames) == 2 {
		recipientsText = fmt.Sprintf("%s and %s", recipientNames[0], recipientNames[1])
	} else if len(recipientNames) > 2 {
		recipientsText = fmt.Sprintf("%s and %d others", recipientNames[0], len(recipientNames)-1)
	}

	// Build clean and simple Adaptive Card for Teams
	// Simple, elegant design without company logo, matching the reference style
	adaptiveCardBody := []map[string]interface{}{
		// Title section - centered with emoji
		{
			"type": "Container",
			"items": []map[string]interface{}{
				{
					"type":                "TextBlock",
					"text":                "✨ Thank You Card",
					"weight":              "Bolder",
					"size":                "ExtraLarge",
					"wrap":                true,
					"horizontalAlignment": "Center",
					"spacing":             "None",
				},
				{
					"type":                "TextBlock",
					"text":                "Some help may be routine, but it's worth being sincerely thanked ❤️",
					"wrap":                true,
					"size":                "Medium",
					"horizontalAlignment": "Center",
					"spacing":             "Small",
					"isSubtle":            true,
				},
			},
			"spacing": "Medium",
		},
		// To section - small font at the beginning
		{
			"type": "Container",
			"items": []map[string]interface{}{
				{
					"type":    "TextBlock",
					"text":    recipientsText,
					"wrap":    true,
					"spacing": "Small",
					"size":    "Small",
				},
			},
			"spacing": "Small",
		},
	}

	// Main message/reason section - before values
	adaptiveCardBody = append(adaptiveCardBody, map[string]interface{}{
		"type":  "Container",
		"style": "default",
		"items": []map[string]interface{}{
			{
				"type":    "TextBlock",
				"text":    card.Reason,
				"wrap":    true,
				"spacing": "None",
				"size":    "Medium",
			},
		},
		"spacing": "Medium",
	})

	// Add values section after reason - as tags without "Values:" prefix
	// Use FactSet to create a tag-like appearance
	if len(valueNames) > 0 {
		// Create value tags as a comma-separated list with bold formatting
		valueTagsText := ""
		for i, v := range valueNames {
			if i > 0 {
				valueTagsText += ", "
			}
			valueTagsText += fmt.Sprintf("**%s**", v)
		}
		adaptiveCardBody = append(adaptiveCardBody, map[string]interface{}{
			"type": "Container",
			"items": []map[string]interface{}{
				{
					"type":    "TextBlock",
					"text":    valueTagsText,
					"wrap":    true,
					"spacing": "None",
					"size":    "Small",
				},
			},
			"spacing": "Small",
		})
	}

	// Add closing sentiment
	adaptiveCardBody = append(adaptiveCardBody, map[string]interface{}{
		"type":    "Container",
		"spacing": "Medium",
		"items": []map[string]interface{}{
			{
				"type":                "TextBlock",
				"text":                "🌈 Working with you is a lucky thing!",
				"wrap":                true,
				"horizontalAlignment": "Center",
				"spacing":             "None",
				"size":                "Medium",
				"weight":              "Bolder",
			},
		},
	})

	// Add sender signature at the end
	adaptiveCardBody = append(adaptiveCardBody, map[string]interface{}{
		"type":    "Container",
		"spacing": "Small",
		"items": []map[string]interface{}{
			{
				"type":     "TextBlock",
				"text":     fmt.Sprintf("%s (%s)", card.Sender.Name, card.Sender.Department),
				"wrap":     true,
				"spacing":  "None",
				"size":     "Small",
				"isSubtle": true,
			},
		},
	})

	// Build Teams message with clean and simple Adaptive Card
	messageCard := map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]interface{}{
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"type":    "AdaptiveCard",
					"version": "1.4",
					"body":    adaptiveCardBody,
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(messageCard)
	if err != nil {
		return fmt.Errorf("failed to marshal message card: %w", err)
	}

	// Log the request payload for debugging
	fmt.Printf("Sending Teams webhook to: %s\n", s.webhookURL)
	fmt.Printf("Request payload: %s\n", string(jsonData))

	// Send HTTP POST request to webhook
	req, err := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	bodyBytes, readErr := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Always log response for debugging
	if bodyStr != "" {
		fmt.Printf("Teams webhook response (status %d): %s\n", resp.StatusCode, bodyStr)
	} else {
		fmt.Printf("Teams webhook response (status %d): [empty body]\n", resp.StatusCode)
	}
	if readErr != nil {
		fmt.Printf("Warning: failed to read response body: %v\n", readErr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d, response: %s", resp.StatusCode, bodyStr)
	}

	// Check if response indicates failure even with 200 status
	// Power Automate may return 200 with error message in body
	if bodyStr != "" && (resp.StatusCode == 200 || resp.StatusCode == 201) {
		// Check for common error indicators in response (case-insensitive)
		lowerBody := fmt.Sprintf("%s", bodyStr)
		if contains(lowerBody, "error") || contains(lowerBody, "failed") || contains(lowerBody, "invalid") {
			fmt.Printf("WARNING: Response may indicate failure despite %d status: %s\n", resp.StatusCode, bodyStr)
		}
	}

	return nil
}

// noOpTeamsNotificationService is a no-op implementation when webhook is not configured
type noOpTeamsNotificationService struct{}

func (s *noOpTeamsNotificationService) SendCardNotification(ctx context.Context, card *models.Card) error {
	// No-op: webhook URL not configured
	return nil
}

// Helper function to join strings with separator
func joinStringsWithSep(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + " " + strs[i]
	}
	return result
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

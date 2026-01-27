package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
)

// EventPublisher defines the interface for publishing card events
type EventPublisher interface {
	PublishCardCreated(card *models.Card) error
}

// TeamsWebhookPublisher implements EventPublisher using HTTP webhooks
type TeamsWebhookPublisher struct {
	webhookURL string
	enabled    bool
	httpClient *http.Client
}

// NewTeamsWebhookPublisher creates a new Teams webhook publisher
func NewTeamsWebhookPublisher(webhookURL string, enabled bool) EventPublisher {
	return &TeamsWebhookPublisher{
		webhookURL: webhookURL,
		enabled:    enabled,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// PublishCardCreated sends a notification to Teams via webhook
func (p *TeamsWebhookPublisher) PublishCardCreated(card *models.Card) error {
	if !p.enabled || p.webhookURL == "" {
		log.Println("Teams notification is disabled or webhook URL is empty, skipping.")
		return nil
	}

	// For MVP, we'll send a simple JSON payload.
	// In the future, this can be expanded to Adaptive Cards.
	payload := map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]interface{}{
					"type": "AdaptiveCard",
					"body": []map[string]interface{}{
						{
							"type":   "TextBlock",
							"size":   "Medium",
							"weight": "Bolder",
							"text":   "🎉 New Thank You Card!",
						},
						{
							"type": "TextBlock",
							"text": fmt.Sprintf("From: **%s**", card.SenderID),
							"wrap": true,
						},
						{
							"type": "TextBlock",
							"text": fmt.Sprintf("To: **%s**", hmapRecipients(card.Recipients)),
							"wrap": true,
						},
						{
							"type":     "TextBlock",
							"text":     card.RecognitionReason,
							"wrap":     true,
							"fontType": "Monospace",
						},
					},
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"version": "1.2",
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal teams payload: %w", err)
	}

	// Simple retry logic (3 attempts)
	var lastErr error
	for i := 0; i < 3; i++ {
		resp, err := p.httpClient.Post(p.webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				log.Printf("Successfully published card creation to Teams: %s", card.ID)
				return nil
			}
			lastErr = fmt.Errorf("teams webhook returned status: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		log.Printf("Attempt %d failed to publish to Teams: %v. Retrying...", i+1, lastErr)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return fmt.Errorf("failed to publish to Teams after 3 attempts: %w", lastErr)
}

func hmapRecipients(recipients []models.CardRecipient) string {
	res := ""
	for i, r := range recipients {
		if i > 0 {
			res += ", "
		}
		res += r.RecipientID
	}
	return res
}

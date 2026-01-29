package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/castlery/thank-you-card/internal/models"
	"github.com/castlery/thank-you-card/internal/services"
)

// TeamsWebhookPublisher implements EventPublisher using HTTP webhooks
type TeamsWebhookPublisher struct {
	webhookURL string
	enabled    bool
	httpClient *http.Client
}

// NewTeamsWebhookPublisher creates a new Teams webhook publisher
func NewTeamsWebhookPublisher(webhookURL string, enabled bool) services.EventPublisher {
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

	// Generate mentions and recipient text
	recipientNames := make([]string, len(card.Recipients))
	entities := make([]map[string]interface{}, len(card.Recipients))

	for i, r := range card.Recipients {
		mentionTag := fmt.Sprintf("<at>%s</at>", r.RecipientName)
		recipientNames[i] = mentionTag

		entities[i] = map[string]interface{}{
			"type": "mention",
			"text": mentionTag,
			"mentioned": map[string]interface{}{
				"id":   r.RecipientID, // Expected to be AAD Object ID
				"name": r.RecipientName,
			},
		}
	}

	recipientText := strings.Join(recipientNames, ", ")

	payload := map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]interface{}{
					"type": "AdaptiveCard",
					"body": []map[string]interface{}{
						{
							"type":  "Container",
							"style": "emphasis", // Light grey background
							"items": []map[string]interface{}{
								{
									"type": "ColumnSet",
									"columns": []map[string]interface{}{
										{
											"type":  "Column",
											"width": "auto",
											"items": []map[string]interface{}{
												{
													"type":  "Image",
													"url":   "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Unicorn/3D/unicorn_3d.png",
													"size":  "Medium",
													"style": "Person",
												},
											},
										},
										{
											"type":                     "Column",
											"width":                    "stretch",
											"verticalContentAlignment": "Center",
											"items": []map[string]interface{}{
												{
													"type":   "TextBlock",
													"text":   "Congratulations",
													"size":   "Large",
													"weight": "Bolder",
													"color":  "Accent",
												},
												{
													"type":   "TextBlock",
													"text":   recipientText,
													"size":   "Medium",
													"weight": "Bolder",
													"wrap":   true,
												},
											},
										},
									},
								},
								{
									"type":    "TextBlock",
									"text":    card.RecognitionReason,
									"wrap":    true,
									"size":    "Medium",
									"spacing": "Medium",
								},
								{
									"type":     "TextBlock",
									"text":     fmt.Sprintf("From **%s**", card.SenderName),
									"size":     "Small",
									"isSubtle": true,
									"spacing":  "Large",
								},
							},
						},
					},
					"actions": []map[string]interface{}{
						{
							"type":  "Action.OpenUrl",
							"title": "Review your praise history",
							"url":   "https://teams.microsoft.com/l/entity/com.castlery.thankyou/thankyou-tab?context={\"subEntityId\":\"history\"}",
						},
						{
							"type":  "Action.OpenUrl",
							"title": "Send praise",
							"url":   "https://teams.microsoft.com/l/entity/com.castlery.thankyou/thankyou-tab?context={\"subEntityId\":\"create\"}",
						},
					},
					"msteams": map[string]interface{}{
						"entities": entities,
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

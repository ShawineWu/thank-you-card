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

// ValueColorConfig holds color configuration for a company value
type ValueColorConfig struct {
	BgColor   string // Background color hex
	TextColor string // Text color hex
}

// valueColorMap maps value names to their colors (matching frontend valueColors.ts)
var valueColorMap = map[string]ValueColorConfig{
	// Values (5)
	"Make an Impact":        {BgColor: "#ede9fe", TextColor: "#6d28d9"},
	"Strive for Excellence": {BgColor: "#fef3c7", TextColor: "#b45309"},
	"Stand Together":        {BgColor: "#d1fae5", TextColor: "#047857"},
	"Be Open-Minded":        {BgColor: "#e0f2fe", TextColor: "#0369a1"},
	"Stay Grounded":         {BgColor: "#ffe4e6", TextColor: "#be123c"},
	// Credos (10)
	"Bias for Action":               {BgColor: "#ffedd5", TextColor: "#c2410c"},
	"Customer Centric":              {BgColor: "#fce7f3", TextColor: "#be185d"},
	"Think Strategically":           {BgColor: "#e0e7ff", TextColor: "#4338ca"},
	"Deep Dive":                     {BgColor: "#cffafe", TextColor: "#0e7490"},
	"Invent and Simplify":           {BgColor: "#ccfbf1", TextColor: "#0f766e"},
	"Earn Trust":                    {BgColor: "#dbeafe", TextColor: "#1d4ed8"},
	"Take Ownership":                {BgColor: "#f3e8ff", TextColor: "#7e22ce"},
	"Challenge Disagree and Commit": {BgColor: "#fee2e2", TextColor: "#b91c1c"},
	"Learn and Be Curious":          {BgColor: "#ecfccb", TextColor: "#4d7c0f"},
	"Do More with Less":             {BgColor: "#fae8ff", TextColor: "#a21caf"},
}

// Default color for unknown values
var defaultValueColor = ValueColorConfig{BgColor: "#f3f4f6", TextColor: "#374151"}

// getValueColor returns the color config for a value name
func getValueColor(valueName string) ValueColorConfig {
	if color, ok := valueColorMap[valueName]; ok {
		return color
	}
	return defaultValueColor
}

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

// PublishCardCreated sends a notification to Teams via webhook asynchronously
func (p *TeamsWebhookPublisher) PublishCardCreated(card *models.Card) error {
	if !p.enabled || p.webhookURL == "" {
		log.Println("Teams notification is disabled or webhook URL is empty, skipping.")
		return nil
	}

	// Run webhook notification asynchronously to avoid blocking the API response
	go p.sendTeamsNotification(card)
	return nil
}

// sendTeamsNotification sends the actual notification to Teams (runs in background)
func (p *TeamsWebhookPublisher) sendTeamsNotification(card *models.Card) {
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

	// Build value badges as clickable Action buttons with links
	valueActions := make([]map[string]interface{}, 0)
	for _, v := range card.Values {
		// Use value ID in the URL query parameter
		valueURL := fmt.Sprintf("http://localhost:3000/values?id=%s", v.CompanyValue.ID.String())

		valueActions = append(valueActions, map[string]interface{}{
			"type":  "Action.OpenUrl",
			"title": fmt.Sprintf("🏷️ %s", v.CompanyValue.Name),
			"url":   valueURL,
		})
	}

	payload := map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]interface{}{
					"type": "AdaptiveCard",
					"body": []map[string]interface{}{
						// Recipient section
						{
							"type":    "Container",
							"spacing": "Medium",
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
													"url":   "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Trophy/3D/trophy_3d.png",
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
													"text":   "Recognition Card",
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
													"color":  "Good",
												},
											},
										},
									},
								},
							},
						},
						// Recognition reason
						{
							"type":    "Container",
							"style":   "emphasis",
							"spacing": "Medium",
							"items": []map[string]interface{}{
								{
									"type":    "TextBlock",
									"text":    "💬 Recognition Message",
									"size":    "Small",
									"weight":  "Bolder",
									"color":   "Accent",
									"spacing": "Small",
								},
								{
									"type":    "TextBlock",
									"text":    card.RecognitionReason,
									"wrap":    true,
									"size":    "Medium",
									"spacing": "Small",
								},
							},
						},
						// Company values section
						{
							"type":    "Container",
							"spacing": "Medium",
							"items": []map[string]interface{}{
								{
									"type":   "TextBlock",
									"text":   "🎯 Company Values Demonstrated",
									"size":   "Small",
									"weight": "Bolder",
									"color":  "Accent",
								},
								{
									"type":    "ActionSet",
									"actions": valueActions,
									"spacing": "Small",
								},
							},
						},
						// Sender info with decorative element
						{
							"type":    "Container",
							"spacing": "Large",
							"items": []map[string]interface{}{
								{
									"type": "ColumnSet",
									"columns": []map[string]interface{}{
										{
											"type":  "Column",
											"width": "stretch",
											"items": []map[string]interface{}{
												{
													"type":     "TextBlock",
													"text":     fmt.Sprintf("💝 Sent with appreciation by **%s**", card.SenderName),
													"size":     "Small",
													"isSubtle": true,
													"wrap":     true,
												},
											},
										},
										{
											"type":  "Column",
											"width": "auto",
											"items": []map[string]interface{}{
												{
													"type": "Image",
													"url":  "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Sparkling%20heart/3D/sparkling_heart_3d.png",
													"size": "Small",
												},
											},
										},
									},
								},
							},
						},
					},
					"actions": []map[string]interface{}{
						{
							"type":  "Action.OpenUrl",
							"title": "🏆 View Praise History",
							"url":   "http://localhost:3000/profile",
						},
						{
							"type":  "Action.OpenUrl",
							"title": "✨ Send Praise",
							"url":   "http://localhost:3000/dashboard?action=send-card",
							"style": "positive",
						},
					},
					"msteams": map[string]interface{}{
						"entities": entities,
					},
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"version": "1.4",
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal teams payload: %v", err)
		return
	}

	// Simple retry logic (3 attempts)
	var lastErr error
	for i := 0; i < 3; i++ {
		resp, err := p.httpClient.Post(p.webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				log.Printf("Successfully published card creation to Teams: %s", card.ID)
				return
			}
			lastErr = fmt.Errorf("teams webhook returned status: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		log.Printf("Attempt %d failed to publish to Teams: %v. Retrying...", i+1, lastErr)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	log.Printf("Failed to publish to Teams after 3 attempts: %v", lastErr)
}

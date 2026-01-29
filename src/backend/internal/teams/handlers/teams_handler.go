package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/castlery/thank-you-card/internal/card/domain/services"
	"github.com/gin-gonic/gin"
)

// TeamsHandler handles requests from Microsoft Teams
type TeamsHandler struct {
	cardCreationService *services.CardCreationService
}

// NewTeamsHandler creates a new Teams handler
func NewTeamsHandler(cardCreationService *services.CardCreationService) *TeamsHandler {
	return &TeamsHandler{
		cardCreationService: cardCreationService,
	}
}

// HandleMessages handles POST requests from Teams
func (h *TeamsHandler) HandleMessages(c *gin.Context) {
	var activity Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity"})
		return
	}

	// Handle Message Extension Submit Action
	if activity.Name == "composeExtension/submitAction" {
		h.handleSubmitAction(c, &activity)
		return
	}

	// Handle Adaptive Card Actions (Action.Execute)
	if activity.Type == "invoke" && activity.Name == "adaptiveCard/action" {
		h.handleAdaptiveCardAction(c, &activity)
		return
	}

	// Default response for other activities
	c.Status(http.StatusAccepted)
}

func (h *TeamsHandler) handleSubmitAction(c *gin.Context, activity *Activity) {
	// Extract data from the activity
	var data struct {
		Data CardData `json:"data"`
	}

	if err := json.Unmarshal(activity.Value, &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"composeExtension": gin.H{
				"type": "message",
				"text": "Failed to parse submitted data",
			},
		})
		return
	}

	// Generate the full content card
	// This will be the default view for non-recipients, and the refresh result for opening
	card := h.generateContentCard(data.Data, true)

	c.JSON(http.StatusOK, gin.H{
		"composeExtension": gin.H{
			"type":             "result",
			"attachmentLayout": "list",
			"attachments": []interface{}{
				map[string]interface{}{
					"contentType": "application/vnd.microsoft.card.adaptive",
					"content":     card,
				},
			},
		},
	})
}

func (h *TeamsHandler) handleAdaptiveCardAction(c *gin.Context, activity *Activity) {
	var actionData struct {
		Action struct {
			Verb string          `json:"verb"`
			Data json.RawMessage `json:"data"`
		} `json:"action"`
	}

	if err := json.Unmarshal(activity.Value, &actionData); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var responseCard map[string]interface{}

	switch actionData.Action.Verb {
	case "checkView":
		// This is triggered automatically when the card is viewed
		// We check if the viewer is the recipient
		var payload struct {
			CardData   CardData          `json:"cardData"`
			Recipients []RecipientStruct `json:"recipients"`
		}
		if err := json.Unmarshal(actionData.Action.Data, &payload); err == nil {
			if h.isRecipient(activity.From.AadObjectId, payload.Recipients) {
				// Viewer is a recipient, show the Envelope
				responseCard = h.generateEnvelopeCard(payload.CardData)
			} else {
				// Viewer is NOT a recipient
				// We return nil/empty to indicate "no change" or we can return the content card without refresh
				// Returning nil usually means "keep matching what you have" but in refresh context, we must return a card.
				// We return the content card (without refresh to stop looping if needed, or with it).
				// Safest: Return the content card WITHOUT refresh property to avoid constant refreshing.
				responseCard = h.generateContentCard(payload.CardData, false)
			}
		}

	case "openCard":
		// This is technically legacy fallback now, as we moved to ToggleVisibility.
		// But in case of errors or old clients, we can keep it handling server-side open.
		var payload struct {
			CardData CardData `json:"cardData"`
		}
		if err := json.Unmarshal(actionData.Action.Data, &payload); err == nil {
			// Return the full content card, no refresh needed anymore
			responseCard = h.generateContentCard(payload.CardData, false)
		}
	}

	if responseCard != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 200,
			"type":       "application/vnd.microsoft.card.adaptive",
			"value":      responseCard,
		})
	} else {
		// If something failed or no match, return 200 OK with no body to do nothing
		c.Status(http.StatusOK)
	}
}

// Helpers

func (h *TeamsHandler) isRecipient(userAadID string, recipients []RecipientStruct) bool {
	// If AAD ID is missing, we can't match safely.
	if userAadID == "" {
		return false
	}
	for _, r := range recipients {
		// We assume r.ID might be the AAD ID.
		// In the frontend, we saw it comes from api.searchEmployees.
		if strings.EqualFold(r.ID, userAadID) {
			return true
		}
	}
	return false
}

func (h *TeamsHandler) generateEnvelopeCard(data CardData) map[string]interface{} {
	// Generate the content body first
	contentBody := h.generateContentBody(data)

	// Wrap content body in a Container that is initially HIDDEN
	contentContainer := map[string]interface{}{
		"type":      "Container",
		"id":        "contentContainer",
		"isVisible": false, // Hidden by default
		"items":     contentBody,
	}

	// Create Envelope Container (Visible by default)
	envelopeContainer := map[string]interface{}{
		"type":      "Container",
		"id":        "envelopeContainer",
		"style":     "accent",
		"isVisible": true,
		"items": []interface{}{
			map[string]interface{}{
				"type":                "ColumnSet",
				"horizontalAlignment": "Center",
				"columns": []interface{}{
					map[string]interface{}{
						"type":  "Column",
						"width": "auto",
						"items": []interface{}{
							map[string]interface{}{
								"type":   "TextBlock",
								"text":   "🎁", // Gift box emoji for a more "surprise" feel
								"size":   "ExtraLarge",
								"weight": "Bolder",
							},
						},
					},
					map[string]interface{}{
						"type":                     "Column",
						"width":                    "stretch",
						"verticalContentAlignment": "Center",
						"items": []interface{}{
							map[string]interface{}{
								"type":   "TextBlock",
								"text":   "You've received a Thank You Card!",
								"weight": "Bolder",
								"size":   "Medium",
								"wrap":   true,
								"color":  "Light", // Better contrast on Accent background
							},
							map[string]interface{}{
								"type":     "TextBlock",
								"text":     fmt.Sprintf("From %s", data.SenderName),
								"isSubtle": true,
								"wrap":     true,
								"color":    "Light",
							},
						},
					},
				},
			},
			map[string]interface{}{
				"type":    "ActionSet",
				"spacing": "Medium",
				"actions": []interface{}{
					map[string]interface{}{
						"type":  "Action.ToggleVisibility", // Client-side animation!
						"title": "Open Card",
						"targetElements": []string{
							"envelopeContainer", // Hide this
							"contentContainer",  // Show this
						},
					},
				},
			},
		},
	}

	return map[string]interface{}{
		"type":    "AdaptiveCard",
		"version": "1.4",
		"body": []interface{}{
			envelopeContainer,
			contentContainer,
		},
	}
}

func (h *TeamsHandler) generateContentCard(data CardData, includeRefresh bool) map[string]interface{} {
	// Generate the body
	body := h.generateContentBody(data)

	// In the standard Content Card view (for others), it's just the body at the root
	card := map[string]interface{}{
		"type":    "AdaptiveCard",
		"version": "1.4",
		"body":    body, // Reuse the same body structure
	}

	if includeRefresh {
		// Parse recipients for the refresh payload
		var recipientsStruct []RecipientStruct
		json.Unmarshal(data.Recipients, &recipientsStruct)

		card["refresh"] = map[string]interface{}{
			"action": map[string]interface{}{
				"type": "Action.Execute",
				"verb": "checkView",
				"data": map[string]interface{}{
					"cardData":   data,
					"recipients": recipientsStruct,
				},
			},
			// We define no userIds so EVERYONE triggers this.
			// The bot decides what to show.
		}
	}

	return card
}

// generateContentBody returns the shared internal body structure of the thank you card
func (h *TeamsHandler) generateContentBody(data CardData) []interface{} {
	// Parse recipients to string
	var recipientNames []string
	var recipientsStruct []RecipientStruct
	json.Unmarshal(data.Recipients, &recipientsStruct)

	for _, r := range recipientsStruct {
		recipientNames = append(recipientNames, r.Name)
	}
	recipientText := strings.Join(recipientNames, ", ")

	return []interface{}{
		map[string]interface{}{
			"type":  "Container",
			"style": "emphasis",
			"items": []interface{}{
				map[string]interface{}{
					"type":   "TextBlock",
					"text":   "🎉 Recognition Card",
					"weight": "Bolder",
					"size":   "Large",
					"color":  "Accent",
				},
			},
		},
		map[string]interface{}{
			"type":    "Container",
			"spacing": "Medium",
			"items": []interface{}{
				map[string]interface{}{
					"type": "FactSet",
					"facts": []interface{}{
						map[string]interface{}{
							"title": "To:",
							"value": recipientText,
						},
						map[string]interface{}{
							"title": "From:",
							"value": data.SenderName,
						},
					},
				},
				map[string]interface{}{
					"type":    "TextBlock",
					"text":    data.RecognitionReason,
					"wrap":    true,
					"size":    "Medium",
					"spacing": "Medium",
				},
			},
		},
	}
}

// Teams Activity structures
type Activity struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
	From  struct {
		Id          string `json:"id"`
		AadObjectId string `json:"aadObjectId"`
	} `json:"from"`
}

type CardData struct {
	Recipients        json.RawMessage `json:"recipients"`
	RecognitionReason string          `json:"recognitionReason"`
	ValueIds          []string        `json:"valueIds"`
	SenderName        string          `json:"senderName"`
}

type RecipientStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

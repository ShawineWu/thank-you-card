package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/pkg/logger"
)

// WebhookPublisher publishes domain events via HTTP webhooks
type WebhookPublisher struct {
	httpClient   *http.Client
	webhookURL   string
	retryAttempts int
	retryDelay   time.Duration
}

// NewWebhookPublisher creates a new webhook publisher
func NewWebhookPublisher(webhookURL string) *WebhookPublisher {
	return &WebhookPublisher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		webhookURL:   webhookURL,
		retryAttempts: 3,
		retryDelay:   time.Second * 2,
	}
}

// EventPublisher interface for publishing domain events
type EventPublisher interface {
	PublishCardCreated(event card.CardCreated) error
	PublishCardShared(event card.CardShared) error
	PublishMilestoneAchieved(event card.MilestoneAchieved) error
}

// PublishCardCreated publishes a CardCreated event
func (p *WebhookPublisher) PublishCardCreated(event card.CardCreated) error {
	logger.Info("Publishing CardCreated event for card:", event.CardID)
	
	payload := WebhookPayload{
		EventType: "CardCreated",
		EventID:   event.EventID.String(),
		Timestamp: event.Timestamp,
		Version:   event.Version,
		Data:      event,
	}

	return p.publishWithRetry(payload)
}

// PublishCardShared publishes a CardShared event
func (p *WebhookPublisher) PublishCardShared(event card.CardShared) error {
	logger.Info("Publishing CardShared event for card:", event.CardID)
	
	payload := WebhookPayload{
		EventType: "CardShared",
		EventID:   event.EventID.String(),
		Timestamp: event.Timestamp,
		Version:   event.Version,
		Data:      event,
	}

	return p.publishWithRetry(payload)
}

// PublishMilestoneAchieved publishes a MilestoneAchieved event
func (p *WebhookPublisher) PublishMilestoneAchieved(event card.MilestoneAchieved) error {
	logger.Info("Publishing MilestoneAchieved event for employee:", event.EmployeeID)
	
	payload := WebhookPayload{
		EventType: "MilestoneAchieved",
		EventID:   event.EventID.String(),
		Timestamp: event.Timestamp,
		Version:   event.Version,
		Data:      event,
	}

	return p.publishWithRetry(payload)
}

// WebhookPayload represents the structure of webhook payloads
type WebhookPayload struct {
	EventType string      `json:"eventType"`
	EventID   string      `json:"eventId"`
	Timestamp time.Time   `json:"timestamp"`
	Version   string      `json:"version"`
	Data      interface{} `json:"data"`
}

// publishWithRetry publishes a webhook with retry logic
func (p *WebhookPublisher) publishWithRetry(payload WebhookPayload) error {
	var lastErr error
	
	for attempt := 0; attempt <= p.retryAttempts; attempt++ {
		if attempt > 0 {
			logger.Info("Retrying webhook publish, attempt:", attempt)
			time.Sleep(p.retryDelay * time.Duration(attempt)) // Exponential backoff
		}

		err := p.publishWebhook(payload)
		if err == nil {
			logger.Info("Successfully published webhook event:", payload.EventType)
			return nil
		}

		lastErr = err
		logger.Error("Failed to publish webhook, attempt", attempt+1, "error:", err)
	}

	return fmt.Errorf("failed to publish webhook after %d attempts: %w", p.retryAttempts+1, lastErr)
}

// publishWebhook makes the actual HTTP request to publish the webhook
func (p *WebhookPublisher) publishWebhook(payload WebhookPayload) error {
	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", p.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ThankYouCard-Webhook/1.0")
	req.Header.Set("X-Event-Type", payload.EventType)
	req.Header.Set("X-Event-ID", payload.EventID)

	// Add webhook signature for security (in production)
	// signature := p.generateSignature(jsonData)
	// req.Header.Set("X-Webhook-Signature", signature)

	// Make the request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// generateSignature generates a webhook signature for security
// This would use HMAC-SHA256 with a secret key in production
func (p *WebhookPublisher) generateSignature(payload []byte) string {
	// Implementation would use crypto/hmac and crypto/sha256
	// For now, return a placeholder
	return "sha256=placeholder_signature"
}

// MockEventPublisher is a mock implementation for testing
type MockEventPublisher struct {
	PublishedEvents []interface{}
}

// NewMockEventPublisher creates a new mock event publisher
func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		PublishedEvents: make([]interface{}, 0),
	}
}

// PublishCardCreated mock implementation
func (m *MockEventPublisher) PublishCardCreated(event card.CardCreated) error {
	m.PublishedEvents = append(m.PublishedEvents, event)
	logger.Info("Mock: Published CardCreated event for card:", event.CardID)
	return nil
}

// PublishCardShared mock implementation
func (m *MockEventPublisher) PublishCardShared(event card.CardShared) error {
	m.PublishedEvents = append(m.PublishedEvents, event)
	logger.Info("Mock: Published CardShared event for card:", event.CardID)
	return nil
}

// PublishMilestoneAchieved mock implementation
func (m *MockEventPublisher) PublishMilestoneAchieved(event card.MilestoneAchieved) error {
	m.PublishedEvents = append(m.PublishedEvents, event)
	logger.Info("Mock: Published MilestoneAchieved event for employee:", event.EmployeeID)
	return nil
}

// GetPublishedEvents returns all published events (for testing)
func (m *MockEventPublisher) GetPublishedEvents() []interface{} {
	return m.PublishedEvents
}

// ClearEvents clears all published events (for testing)
func (m *MockEventPublisher) ClearEvents() {
	m.PublishedEvents = make([]interface{}, 0)
}
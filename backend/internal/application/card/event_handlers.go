package card

import (
	"fmt"

	"github.com/company/thank-you-card/internal/domain/card"
	"github.com/company/thank-you-card/internal/infrastructure/events"
	"github.com/company/thank-you-card/pkg/logger"
)

// CardCreatedEventHandler handles CardCreated domain events
type CardCreatedEventHandler struct {
	milestoneService   *card.MilestoneTrackingService
	milestoneRepository card.MilestoneRepository
	eventPublisher     events.EventPublisher
}

// NewCardCreatedEventHandler creates a new CardCreated event handler
func NewCardCreatedEventHandler(
	milestoneService *card.MilestoneTrackingService,
	milestoneRepo card.MilestoneRepository,
	eventPublisher events.EventPublisher,
) *CardCreatedEventHandler {
	return &CardCreatedEventHandler{
		milestoneService:   milestoneService,
		milestoneRepository: milestoneRepo,
		eventPublisher:     eventPublisher,
	}
}

// Handle processes a CardCreated event
func (h *CardCreatedEventHandler) Handle(event card.CardCreated) error {
	logger.Info("Processing CardCreated event for card:", event.CardID)

	// Check and record milestones asynchronously
	go h.handleMilestoneDetection(event)

	return nil
}

// handleMilestoneDetection handles milestone detection in a separate goroutine
func (h *CardCreatedEventHandler) handleMilestoneDetection(event card.CardCreated) {
	logger.Info("Checking milestones for CardCreated event:", event.CardID)

	// Check for milestone achievements
	newMilestones, err := h.milestoneService.CheckAndRecordMilestones(event)
	if err != nil {
		logger.Error("Failed to check milestones for card:", event.CardID, "error:", err)
		return
	}

	if len(newMilestones) == 0 {
		logger.Debug("No new milestones achieved for card:", event.CardID)
		return
	}

	// Save new milestones
	if err := h.milestoneRepository.CreateBatch(newMilestones); err != nil {
		logger.Error("Failed to save milestones for card:", event.CardID, "error:", err)
		return
	}

	// Publish milestone achieved events
	for _, milestone := range newMilestones {
		milestoneEvent := card.NewMilestoneAchieved(milestone)
		if err := h.eventPublisher.PublishMilestoneAchieved(milestoneEvent); err != nil {
			logger.Error("Failed to publish MilestoneAchieved event:", err)
			// Continue with other milestones
		} else {
			logger.Info("Published MilestoneAchieved event for employee:", milestone.EmployeeID, "type:", milestone.MilestoneType)
		}
	}

	logger.Info("Successfully processed", len(newMilestones), "new milestones for card:", event.CardID)
}

// CardSharedEventHandler handles CardShared domain events
type CardSharedEventHandler struct {
	// Add dependencies as needed for card sharing logic
}

// NewCardSharedEventHandler creates a new CardShared event handler
func NewCardSharedEventHandler() *CardSharedEventHandler {
	return &CardSharedEventHandler{}
}

// Handle processes a CardShared event
func (h *CardSharedEventHandler) Handle(event card.CardShared) error {
	logger.Info("Processing CardShared event for card:", event.CardID, "shared to:", event.TargetChannel)

	// Add any additional logic for card sharing
	// For example:
	// - Send notifications to Teams channel
	// - Update sharing statistics
	// - Log sharing activity

	logger.Info("Successfully processed CardShared event for card:", event.CardID)
	return nil
}

// MilestoneAchievedEventHandler handles MilestoneAchieved domain events
type MilestoneAchievedEventHandler struct {
	// Add dependencies as needed for milestone achievement logic
}

// NewMilestoneAchievedEventHandler creates a new MilestoneAchieved event handler
func NewMilestoneAchievedEventHandler() *MilestoneAchievedEventHandler {
	return &MilestoneAchievedEventHandler{}
}

// Handle processes a MilestoneAchieved event
func (h *MilestoneAchievedEventHandler) Handle(event card.MilestoneAchieved) error {
	logger.Info("Processing MilestoneAchieved event for employee:", event.EmployeeID, "milestone:", event.Title)

	// Add any additional logic for milestone achievements
	// For example:
	// - Send congratulatory notifications
	// - Update employee profiles
	// - Generate milestone certificates
	// - Update leaderboards

	logger.Info("Successfully processed MilestoneAchieved event for employee:", event.EmployeeID)
	return nil
}

// EventDispatcher coordinates event handling
type EventDispatcher struct {
	cardCreatedHandler      *CardCreatedEventHandler
	cardSharedHandler       *CardSharedEventHandler
	milestoneAchievedHandler *MilestoneAchievedEventHandler
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher(
	cardCreatedHandler *CardCreatedEventHandler,
	cardSharedHandler *CardSharedEventHandler,
	milestoneAchievedHandler *MilestoneAchievedEventHandler,
) *EventDispatcher {
	return &EventDispatcher{
		cardCreatedHandler:      cardCreatedHandler,
		cardSharedHandler:       cardSharedHandler,
		milestoneAchievedHandler: milestoneAchievedHandler,
	}
}

// DispatchCardCreated dispatches a CardCreated event
func (d *EventDispatcher) DispatchCardCreated(event card.CardCreated) error {
	return d.cardCreatedHandler.Handle(event)
}

// DispatchCardShared dispatches a CardShared event
func (d *EventDispatcher) DispatchCardShared(event card.CardShared) error {
	return d.cardSharedHandler.Handle(event)
}

// DispatchMilestoneAchieved dispatches a MilestoneAchieved event
func (d *EventDispatcher) DispatchMilestoneAchieved(event card.MilestoneAchieved) error {
	return d.milestoneAchievedHandler.Handle(event)
}

// BackgroundJobProcessor processes background jobs for event handling
type BackgroundJobProcessor struct {
	eventDispatcher *EventDispatcher
	jobQueue        chan Job
	workerCount     int
}

// Job represents a background job
type Job struct {
	Type string
	Data interface{}
}

// NewBackgroundJobProcessor creates a new background job processor
func NewBackgroundJobProcessor(eventDispatcher *EventDispatcher, workerCount int) *BackgroundJobProcessor {
	return &BackgroundJobProcessor{
		eventDispatcher: eventDispatcher,
		jobQueue:        make(chan Job, 1000), // Buffer for 1000 jobs
		workerCount:     workerCount,
	}
}

// Start starts the background job processor
func (p *BackgroundJobProcessor) Start() {
	logger.Info("Starting background job processor with", p.workerCount, "workers")
	
	for i := 0; i < p.workerCount; i++ {
		go p.worker(i)
	}
}

// worker processes jobs from the queue
func (p *BackgroundJobProcessor) worker(workerID int) {
	logger.Info("Starting background worker:", workerID)
	
	for job := range p.jobQueue {
		logger.Debug("Worker", workerID, "processing job:", job.Type)
		
		if err := p.processJob(job); err != nil {
			logger.Error("Worker", workerID, "failed to process job:", job.Type, "error:", err)
		}
	}
}

// processJob processes a single job
func (p *BackgroundJobProcessor) processJob(job Job) error {
	switch job.Type {
	case "CardCreated":
		if event, ok := job.Data.(card.CardCreated); ok {
			return p.eventDispatcher.DispatchCardCreated(event)
		}
		return fmt.Errorf("invalid CardCreated event data")
		
	case "CardShared":
		if event, ok := job.Data.(card.CardShared); ok {
			return p.eventDispatcher.DispatchCardShared(event)
		}
		return fmt.Errorf("invalid CardShared event data")
		
	case "MilestoneAchieved":
		if event, ok := job.Data.(card.MilestoneAchieved); ok {
			return p.eventDispatcher.DispatchMilestoneAchieved(event)
		}
		return fmt.Errorf("invalid MilestoneAchieved event data")
		
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

// EnqueueCardCreated enqueues a CardCreated event for processing
func (p *BackgroundJobProcessor) EnqueueCardCreated(event card.CardCreated) {
	job := Job{
		Type: "CardCreated",
		Data: event,
	}
	
	select {
	case p.jobQueue <- job:
		logger.Debug("Enqueued CardCreated job for card:", event.CardID)
	default:
		logger.Error("Job queue is full, dropping CardCreated job for card:", event.CardID)
	}
}

// EnqueueCardShared enqueues a CardShared event for processing
func (p *BackgroundJobProcessor) EnqueueCardShared(event card.CardShared) {
	job := Job{
		Type: "CardShared",
		Data: event,
	}
	
	select {
	case p.jobQueue <- job:
		logger.Debug("Enqueued CardShared job for card:", event.CardID)
	default:
		logger.Error("Job queue is full, dropping CardShared job for card:", event.CardID)
	}
}

// EnqueueMilestoneAchieved enqueues a MilestoneAchieved event for processing
func (p *BackgroundJobProcessor) EnqueueMilestoneAchieved(event card.MilestoneAchieved) {
	job := Job{
		Type: "MilestoneAchieved",
		Data: event,
	}
	
	select {
	case p.jobQueue <- job:
		logger.Debug("Enqueued MilestoneAchieved job for employee:", event.EmployeeID)
	default:
		logger.Error("Job queue is full, dropping MilestoneAchieved job for employee:", event.EmployeeID)
	}
}

// Stop stops the background job processor
func (p *BackgroundJobProcessor) Stop() {
	logger.Info("Stopping background job processor")
	close(p.jobQueue)
}
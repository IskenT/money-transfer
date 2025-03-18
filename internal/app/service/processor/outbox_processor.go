package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// OutboxEvent
type OutboxEvent struct {
	ID            int64      `db:"id"`
	AggregateType string     `db:"aggregate_type"`
	AggregateID   string     `db:"aggregate_id"`
	EventType     string     `db:"event_type"`
	Payload       []byte     `db:"payload"`
	CreatedAt     time.Time  `db:"created_at"`
	ProcessedAt   *time.Time `db:"processed_at"`
}

// OutboxProcessor
type OutboxProcessor struct {
	db      *sqlx.DB
	running bool
	done    chan struct{}
}

// NewOutboxProcessor
func NewOutboxProcessor(db *sqlx.DB) *OutboxProcessor {
	return &OutboxProcessor{
		db:      db,
		running: false,
		done:    make(chan struct{}),
	}
}

// Start
func (p *OutboxProcessor) Start() {
	if p.running {
		return
	}

	p.running = true
	go p.processEvents()
}

// Stop
func (p *OutboxProcessor) Stop() {
	if !p.running {
		return
	}

	p.running = false
	close(p.done)
}

// processEvents
func (p *OutboxProcessor) processEvents() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.processUnprocessedEvents()
		case <-p.done:
			return
		}
	}
}

// processUnprocessedEvents
func (p *OutboxProcessor) processUnprocessedEvents() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var events []OutboxEvent
	err := p.db.SelectContext(ctx, &events, `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, created_at, processed_at
		FROM money_transfer.outbox_events
		WHERE processed_at IS NULL
		ORDER BY created_at
		LIMIT 100
	`)

	if err != nil {
		log.Printf("Error getting unprocessed events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	for _, event := range events {
		err := p.processEvent(ctx, event)

		if err != nil {
			log.Printf("Error processing event %d: %v", event.ID, err)
			continue
		}

		_, err = p.db.ExecContext(ctx, `
			UPDATE money_transfer.outbox_events
			SET processed_at = NOW()
			WHERE id = $1
		`, event.ID)

		if err != nil {
			log.Printf("Error marking event %d as processed: %v", event.ID, err)
		}
	}
}

// processEvent
func (p *OutboxProcessor) processEvent(ctx context.Context, event OutboxEvent) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("error unmarshaling payload: %w", err)
	}

	// Process
	switch event.EventType {
	case "transfer_completed":
		return p.processTransferCompletedEvent(ctx, event, payload)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

// processTransferCompletedEvent
func (p *OutboxProcessor) processTransferCompletedEvent(ctx context.Context, event OutboxEvent, payload map[string]interface{}) error {

	// For now, just log the event
	log.Printf("Transfer completed: %s from user %v to user %v for amount %v",
		payload["transfer_id"],
		payload["from_user_id"],
		payload["to_user_id"],
		payload["amount"])

	return nil
}

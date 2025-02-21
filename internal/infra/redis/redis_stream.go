package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const ProductStream = "sync:products"

// EventType represents the type of database operation
type EventType string

// ProductEvent represents an event related to product changes
type ProductEvent struct {
	ID        string          `json:"id"`        // Event ID (will be auto-generated if empty)
	Type      EventType       `json:"type"`      // create, update, delete
	ProductID string          `json:"productId"` // Database product ID
	Data      json.RawMessage `json:"data"`      // Product data as JSON
	Timestamp time.Time       `json:"timestamp"` // Event timestamp
}

// ProductEventHandler defines a function type for event handlers
type ProductEventHandler func(context.Context, ProductEvent) error

// PublishProductEvent publishes a product event to Redis stream
func PublishProductEvent(ctx context.Context, event ProductEvent) (string, error) {
	// Set timestamp if not already set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Marshal event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	// Determine stream ID
	streamID := "*" // Auto-generate ID if not provided
	if event.ID != "" {
		streamID = event.ID
	}

	// Publish to Redis stream
	cmd := masterClient.B().Xadd().
		Key(ProductStream).
		Id(streamID).
		FieldValue().
		FieldValue("event", string(eventJSON)).
		Build()

	return masterClient.Do(ctx, cmd).ToString()
}

// ConsumeProductEvents starts consuming product events from Redis stream
func ConsumeProductEvents(ctx context.Context, handler ProductEventHandler) {
	log.Printf("Starting Redis stream consumer for %s", ProductStream)

	// Start the consumer in a goroutine
	go func() {
		lastID := "0-0" // Start from beginning of stream
		heartbeatTicker := time.NewTicker(60 * time.Second)
		defer heartbeatTicker.Stop()

		for {
			select {
			case <-heartbeatTicker.C:
				log.Printf("Redis consumer heartbeat: waiting for events on %s (last ID: %s)", ProductStream, lastID)
			case <-ctx.Done():
				log.Printf("Redis stream consumer shutting down: %v", ctx.Err())
				return // Parent context canceled
			default:
				// Create a timeout context for this specific read operation that will automatically cancel itself after 6 seconds.
				readCtx, cancel := context.WithTimeout(ctx, 6*time.Second)

				// Read new events with blocking call (5 sec timeout)
				cmd := masterClient.B().Xread().
					Count(10).
					Block(5000). // 5 second timeout
					Streams().
					Key(ProductStream).
					Id(lastID).
					Build()

				streamEntries, err := masterClient.Do(readCtx, cmd).AsXRead()
				cancel() // Always cancel the timeout context when done

				if err != nil {
					if err.Error() == "redis: nil" || err.Error() == "redis nil message" || err == context.DeadlineExceeded {
						// Expected when no messages or timeout - continue silently
						continue
					} else if err == context.Canceled {
						// Parent context was canceled
						log.Printf("Redis stream read canceled: %v", err)
						return
					} else {
						// Unexpected error - log and backoff
						log.Printf("Error reading from Redis stream: %v", err)
						time.Sleep(time.Second)
						continue
					}
				}

				// If we received an empty result map but no error, just continue
				if len(streamEntries) == 0 {
					continue
				}

				// Process messages if any
				if entries, ok := streamEntries[ProductStream]; ok && len(entries) > 0 {
					log.Printf("Processing %d new messages from Redis stream", len(entries))

					for _, entry := range entries {
						// Extract event JSON
						var eventJSON string
						for key, value := range entry.FieldValues {
							if key == "event" {
								eventJSON = value
								break
							}
						}

						if eventJSON == "" {
							log.Printf("Received message without event field: %v", entry)
							continue
						}

						// Parse event
						var event ProductEvent
						if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
							log.Printf("Error unmarshaling event: %v", err)
							continue
						}

						// Set the event ID
						event.ID = entry.ID

						// Process event based on its type
						if err := processEvent(ctx, event); err != nil {
							log.Printf("Error processing event %s of type %s: %v",
								entry.ID, event.Type, err)
						}

						// Update last processed ID
						lastID = entry.ID
					}
				}
			}
		}
	}()
}

// Helper function to process events based on type
func processEvent(ctx context.Context, event ProductEvent) error {
	switch event.Type {
	case "create":
		return handleCreateEvent(ctx, event)
	case "update":
		return handleUpdateEvent(ctx, event)
	case "delete":
		return handleDeleteEvent(ctx, event)
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}

// HandleCreateEvent processes product creation events
func handleCreateEvent(ctx context.Context, event ProductEvent) error {
	log.Printf("Processing CREATE event for product ID %s", event.ProductID)

	// Convert event data to JSON string
	eventData, err := json.Marshal(event.Data)
	if err != nil {
		log.Printf("Failed to marshal event data for product ID %s: %v", event.ProductID, err)
		return err
	}

	// Store product in Redis (master)
	SetKey(ctx, fmt.Sprintf("product:%s", event.ProductID), string(eventData))
	return nil
}

// HandleUpdateEvent processes product update events
func handleUpdateEvent(ctx context.Context, event ProductEvent) error {
	log.Printf("Processing UPDATE event for product ID %s", event.ProductID)

	// Convert event data to JSON string
	eventData, err := json.Marshal(event.Data)
	if err != nil {
		log.Printf("Failed to marshal event data for product ID %s: %v", event.ProductID, err)
		return err
	}

	// Update product in Redis (master)
	SetKey(ctx, fmt.Sprintf("product:%s", event.ProductID), string(eventData))
	return nil
}

// HandleDeleteEvent processes product deletion events
func handleDeleteEvent(ctx context.Context, event ProductEvent) error {
	log.Printf("Processing DELETE event for product ID %s", event.ProductID)

	// Delete product from Redis (master)
	DeleteKey(ctx, fmt.Sprintf("product:%s", event.ProductID))
	return nil
}

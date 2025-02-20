package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

const ProductStream = "sync:products"

// EventType represents the type of database operation
type EventType string

// ProductEvent represents an event related to product changes
type ProductEvent struct {
	ID        string          `json:"id"`        // Event ID (will be auto-generated if empty)
	Type      EventType       `json:"type"`      // CREATE, UPDATE, DELETE
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
	// Create a derived context that we can cancel
	processingCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the consumer in a goroutine
	go func() {
		lastID := "0-0" // Start from beginning of stream

		for {
			select {
			case <-ctx.Done():
				return // Parent context canceled
			default:
				// Read new events with blocking call (5 sec timeout)
				cmd := masterClient.B().Xread().
					Count(10).
					Block(5000).
					Streams().
					Key(ProductStream).
					Id(lastID).
					Build()

				streamEntries, err := masterClient.Do(processingCtx, cmd).AsXRead()
				if err != nil {
					if err.Error() == "redis: nil" {
						// No new messages, continue polling
						continue
					}
					log.Printf("Error reading from stream: %v", err)
					time.Sleep(time.Second) // Avoid tight loop on persistent errors
					continue
				}

				// Process messages if any
				if entries, ok := streamEntries[ProductStream]; ok && len(entries) > 0 {
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
						switch event.Type {
						case "CREATE":
							if err := handleCreateEvent(processingCtx, event); err != nil {
								log.Printf("Error processing CREATE event %s: %v", entry.ID, err)
							}

						case "UPDATE":
							if err := handleUpdateEvent(processingCtx, event); err != nil {
								log.Printf("Error processing UPDATE event %s: %v", entry.ID, err)
							}

						case "DELETE":
							if err := handleDeleteEvent(processingCtx, event); err != nil {
								log.Printf("Error processing DELETE event %s: %v", entry.ID, err)
							}

						default:
							log.Printf("Unknown event type %s for event %s", event.Type, entry.ID)
						}

						// Update last processed ID
						lastID = entry.ID
					}
				}
			}
		}
	}()
}

// Handler functions for each event type
func handleCreateEvent(ctx context.Context, event ProductEvent) error {
	// Specialized logic for product creation
	log.Printf("Processing CREATE event for product ID %s", event.ProductID)
	// Implementation details...
	return nil
}

func handleUpdateEvent(ctx context.Context, event ProductEvent) error {
	// Specialized logic for product updates
	log.Printf("Processing UPDATE event for product ID %s", event.ProductID)
	// Implementation details...
	return nil
}

func handleDeleteEvent(ctx context.Context, event ProductEvent) error {
	// Specialized logic for product deletion
	log.Printf("Processing DELETE event for product ID %s", event.ProductID)
	// Implementation details...
	return nil
}

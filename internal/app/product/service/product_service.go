package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/infra/redis"
)

// CreateProduct creates a product and publishes the event asynchronously
func CreateProduct(ctx context.Context, id string, data json.RawMessage) {
	event := redis.ProductEvent{
		Type:      "create",
		ProductID: id,
		Data:      data,
		Timestamp: time.Now(),
	}

	go func() {
		_, _ = redis.PublishProductEvent(ctx, event) // Fire and forget
	}()
}

// UpdateProduct updates a product and publishes the event asynchronously
func UpdateProduct(ctx context.Context, id string, data json.RawMessage) {
	event := redis.ProductEvent{
		Type:      "update",
		ProductID: id,
		Data:      data,
		Timestamp: time.Now(),
	}

	go func() {
		_, _ = redis.PublishProductEvent(ctx, event)
	}()
}

// DeleteProduct deletes a product and publishes the event asynchronously
func DeleteProduct(ctx context.Context, id string, data json.RawMessage) {
	event := redis.ProductEvent{
		Type:      "delete",
		ProductID: id,
		Data:      data,
		Timestamp: time.Now(),
	}

	go func() {
		_, _ = redis.PublishProductEvent(ctx, event)
	}()
}

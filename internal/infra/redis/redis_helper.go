package redis

import (
	"context"
	"log"
)

// SetKey sets a key-value pair in Redis (writes to master).
func SetKey(ctx context.Context, key, value string) {
	cmd := masterClient.B().Set().Key(key).Value(value).Build()
	err := masterClient.Do(ctx, cmd).Error()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
	}
}

// GetKey retrieves a value from Redis (reads from replica).
func GetKey(ctx context.Context, key string) string {
	cmd := replicaClient.B().Get().Key(key).Build()
	res, err := replicaClient.Do(ctx, cmd).ToString()
	if err != nil {
		log.Printf("Failed to get key %s: %v", key, err)
		return ""
	}
	return res
}

// DeleteKey deletes a key in Redis (writes to master).
func DeleteKey(ctx context.Context, key string) {
	cmd := masterClient.B().Del().Key(key).Build()
	err := masterClient.Do(ctx, cmd).Error()
	if err != nil {
		log.Printf("Failed to delete key %s: %v", key, err)
	}
}

package redis

import (
	"context"
	"log"
)

// SetKey sets a key-value pair in Redis.
func SetKey(ctx context.Context, key, value string) {
	err := client.Do(ctx, client.B().Set().Key(key).Value(value).Build()).Error()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
	}
}

// GetKey retrieves a value from Redis using ReadOnly mode.
func GetKey(ctx context.Context, key string) string {
	client.Do(ctx, client.B().Readonly().Build()) // Ensure read commands go to replica
	resp := client.Do(ctx, client.B().Get().Key(key).Build())

	if resp.Error() != nil {
		log.Printf("Failed to get key %s: %v", key, resp.Error())
		return ""
	}

	value, _ := resp.ToString()
	return value
}

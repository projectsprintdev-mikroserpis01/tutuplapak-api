package redis

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/rueidis"
)

var client rueidis.Client

// InitRedis initializes the Redis client with master and replica addresses.
func InitRedis() {
	// Load environment variables
	masterAddr := os.Getenv("REDIS_MASTER_IP") + ":" + os.Getenv("REDIS_PORT")
	replicaAddr := os.Getenv("REDIS_REPLICA_IP") + ":" + os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASS")

	var err error
	client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{masterAddr, replicaAddr}, // Connects to both
		Password:    password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
}

func GetClient() rueidis.Client {
	return client
}

// CloseRedis closes the Redis client.
func CloseRedis() {
	if client != nil {
		client.Close()
	}
}

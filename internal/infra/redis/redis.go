package redis

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/rueidis"
)

var (
	masterClient  rueidis.Client
	replicaClient rueidis.Client
)

// InitRedis initializes separate Redis clients for master and replica.
func InitRedis() {
	// Load environment variables
	masterAddr := os.Getenv("REDIS_MASTER_IP") + ":" + os.Getenv("REDIS_PORT")
	replicaAddr := os.Getenv("REDIS_REPLICA_IP") + ":" + os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASS")

	var err error

	// Connect to Master
	masterClient, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{masterAddr},
		Password:    password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis Master: %v", err)
	}

	// Connect to Replica
	replicaClient, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{replicaAddr},
		Password:    password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis Replica: %v", err)
	}

	fmt.Println("Connected to Redis (Master & Replica)")
}

// GetMasterClient returns the Redis master client.
func GetMasterClient() rueidis.Client {
	return masterClient
}

// GetReplicaClient returns the Redis replica client.
func GetReplicaClient() rueidis.Client {
	return replicaClient
}

// CloseRedis closes both Redis clients.
func CloseRedis() {
	if masterClient != nil {
		masterClient.Close()
	}
	if replicaClient != nil {
		replicaClient.Close()
	}
}

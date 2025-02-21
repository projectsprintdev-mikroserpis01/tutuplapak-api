package redis

import (
	"log"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/infra/env"
	"github.com/redis/rueidis"
)

var (
	masterClient  rueidis.Client
	replicaClient rueidis.Client
)

// InitRedis initializes separate Redis clients for master and replica.
func InitRedis() {
	log.Println("====== Initializing Redis Clients ======")

	// Load environment variables
	masterAddr := env.AppEnv.RedisMasterIp + ":" + env.AppEnv.RedisPort
	replicaAddr := env.AppEnv.RedisReplicaIp + ":" + env.AppEnv.RedisPort
	password := env.AppEnv.RedisPass

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

	log.Println("====== Connected to Redis (Master & Replica) ======")
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

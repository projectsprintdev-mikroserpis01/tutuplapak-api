package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/infra/redis"
)

func main() {
	// Initialize Redis
	log.Println("====== Initializing Redis WORKER ======")
	redis.InitRedis()
	defer redis.CloseRedis()

	// Create a context that is canceled on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Define the event handler
	handler := func(ctx context.Context, event redis.ProductEvent) error {
		// Handle the event (e.g., log it)
		log.Printf("REDIS WORKER => Handling event: %v", event)
		return nil
	}

	// Start consuming product events
	redis.ConsumeProductEvents(ctx, handler)

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Shutting down worker...")
}

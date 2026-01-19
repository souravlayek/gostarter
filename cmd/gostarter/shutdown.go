package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// HandleShutdown sets up signal handling for graceful shutdown
func HandleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		sig := <-sigChan
		log.Printf("[INFO] Received signal: %s", sig.String())
		log.Printf("[INFO] Shutting down gracefully...")

		// Cancel the context to signal shutdown
		cancel()

		// Wait a bit for cleanup
		time.Sleep(100 * time.Millisecond)

		// Exit the program
		os.Exit(0)
	}()
}

// WaitForShutdown waits for shutdown signal
func WaitForShutdown() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handling
	HandleShutdown(cancel)

	// Wait for context to be cancelled
	<-ctx.Done()
}
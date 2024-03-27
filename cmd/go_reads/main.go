package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/kaanserin/go-reads/internal/api"
)

func main() {
	// Initialize an http server
	server, err := api.NewServer(":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Listening inside a go-routine to not block main thread
	go func() {
		log.Printf("Server running on %s\n", ":8080")
		server.ListenAndServe()
	}()

	// Initialize a channel to receive termination signals
	sigChan := make(chan os.Signal, 1)

	// Relay Interrupt and SIGTERM signals to channel
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	// Receive from signal channel
	// This expression blocks
	sig := <-sigChan

	// Only reachable if a termination signal is received from sigChan
	log.Println("Received terminate, gracefully shutting down...", sig)

	// Create a context with 30 second timeout
	tcContext, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Shut down the server
	// If the current connections are not handled in 30 seconds(tcContext), forcefully close them
	server.Shutdown(tcContext)
}

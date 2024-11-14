package main

import (
	"log"
	"txparser/internal"
)

func main() {
	// Initialize Ethereum client to interact with the blockchain
	ethClient, err := internal.NewEthClient()
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %v", err)
	}

	// Initialize storage (in-memory for now)
	storer := internal.NewStorer()

	// Initialize the parser with the Ethereum client and storage
	parser, err := internal.NewManager(ethClient, storer)
	if err != nil {
		log.Fatalf("Failed to initialize parser: %v", err)
	}

	// Initialize the HTTP server with the parser
	s, err := internal.NewServer(parser)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start the HTTP server
	log.Printf("Starting server on port 8080...")
	if err := s.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

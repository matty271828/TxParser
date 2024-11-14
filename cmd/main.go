package main

import (
	"log"
	"txparser/internal"
)

func main() {
	// Initialize Ethereum client to interact with the blockchain
	ethClient, err := internal.NewEthClient()
	if err != nil {
		log.Fatal("Failed to initialize Ethereum client:", err)
	}

	// Initialize storage (in-memory for now)
	storer := internal.NewStorer()

	// Initialize the parser with the Ethereum client and storage
	parser, err := internal.NewManager(ethClient, storer)
	if err != nil {
		log.Fatal("Failed to initialize parser:", err)
	}

	// Initialize the HTTP server with the parser
	s, err := internal.NewServer(parser)
	if err != nil {
		log.Fatal("Failed to initialize server:", err)
	}

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	if err := s.Start(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

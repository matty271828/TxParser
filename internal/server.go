package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	parser Parser
}

func NewServer(parser Parser) (*Server, error) {
	return &Server{parser: parser}, nil
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Note: If I had more time I would add enforcement of HTTP methods.
	mux.HandleFunc("/subscribe", s.HandleSubscribe)
	mux.HandleFunc("/getcurrentblock", s.HandleGetCurrentBlock)
	mux.HandleFunc("/gettransactions", s.HandleGetTransactions)

	return http.ListenAndServe(port, mux)
}

// HandleGetCurrentBlock returns the current block number.
// Returns an error if the current block number cannot be retrieved.
func (s *Server) HandleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Request: GetCurrentBlock - %s", r.URL.String())
	blockNumber := s.parser.GetCurrentBlock()
	if blockNumber == -1 {
		http.Error(w, "Failed to get current block", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"block_number": blockNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleSubscribe subscribes to notifications for a given address.
// Returns an error if the address is not provided or if the subscription fails.
func (s *Server) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Request: Subscribe - %s", r.URL.String())
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	success := s.parser.SubscribeAddress(address)
	if !success {
		http.Error(w, "Address already subscribed", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": success,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetTransactions returns all transactions for a given address.
// Returns an error if the address is not provided or if the transactions cannot be retrieved.
func (s *Server) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Request: GetTransactions - %s", r.URL.String())
	address := r.URL.Query().Get("address")
	if address == "" {
		log.Printf("Error: missing address parameter")
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	transactions := s.parser.GetTransactions(address)
	if transactions == nil {
		log.Printf("Error: failed to get transactions for address %s", address)
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"transactions": transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

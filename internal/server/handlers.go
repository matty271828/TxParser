package server

import (
	"encoding/json"
	"net/http"
)

// HandleGetCurrentBlock returns the current block number.
// Returns an error if the current block number cannot be retrieved.
func (s *Server) HandleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
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
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	success := s.parser.SubscribeAddress(address)
	if !success {
		http.Error(w, "Failed to subscribe to address", http.StatusInternalServerError)
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
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	transactions := s.parser.GetTransactions(address)
	if transactions == nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"transactions": transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

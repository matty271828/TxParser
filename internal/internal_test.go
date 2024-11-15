package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration_GetCurrentBlock(t *testing.T) {
	// Initialize dependencies
	server := setupTestServer(t)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		assert     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "[POSITIVE] successful request",
			method:     http.MethodGet,
			path:       "/getcurrentblock",
			wantStatus: http.StatusOK,
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				blockNumber, ok := response["block_number"].(float64)
				if !ok {
					t.Fatal("Response missing block_number or invalid type")
				}
				if blockNumber <= 0 {
					t.Errorf("Expected block number > 0, got %v", blockNumber)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			server.HandleGetCurrentBlock(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}

			tt.assert(t, w)
		})
	}
}

func TestIntegration_Subscribe(t *testing.T) {
	// Initialize dependencies
	server := setupTestServer(t)

	tests := []struct {
		name       string
		method     string
		address    string
		wantStatus int
		assert     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "[POSITIVE] valid address",
			method:     http.MethodPost,
			address:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			wantStatus: http.StatusOK,
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				success, ok := response["success"].(bool)
				if !ok || !success {
					t.Error("Expected successful subscription")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/subscribe?address="+tt.address, nil)
			w := httptest.NewRecorder()

			server.HandleSubscribe(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}

			tt.assert(t, w)
		})
	}
}

func TestIntegration_GetTransactions(t *testing.T) {
	// Initialize dependencies
	server := setupTestServer(t)

	tests := []struct {
		name       string
		method     string
		address    string
		wantStatus int
		assert     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "[POSITIVE] valid address with transactions",
			method:     http.MethodGet,
			address:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			wantStatus: http.StatusOK,
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				body := w.Body.String()
				t.Logf("Response body: %s", body)

				var response struct {
					Transactions []Transaction `json:"transactions"`
				}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v\nResponse body: %s", err, body)
				}

				if len(response.Transactions) == 0 {
					t.Error("Expected non-empty transaction list")
					return
				}

				tx := response.Transactions[0]
				if tx.TransactionHash == "" {
					t.Error("Transaction hash is empty")
				}
				if tx.Address == "" {
					t.Error("Transaction address is empty")
				}
			},
		},
		{
			name:       "[NEGATIVE] empty address",
			method:     http.MethodGet,
			address:    "",
			wantStatus: http.StatusBadRequest,
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Code != http.StatusBadRequest {
					t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/gettransactions?address="+tt.address, nil)
			w := httptest.NewRecorder()

			server.HandleGetTransactions(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}

			tt.assert(t, w)
		})
	}
}

// Helper function to setup test server
func setupTestServer(t *testing.T) *Server {
	ethClient, err := NewEthClient()
	if err != nil {
		t.Fatalf("Failed to create eth client: %v", err)
	}
	store := NewStorer()
	parser, err := NewManager(ethClient, store)
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	server, err := NewServer(parser)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	return server
}

package ethclient

import (
	"encoding/json"
	"testing"
)

// Using an integration test since this package doesn't do a whole lot
// apart from interacting with the Ethereum node. Additionally, need to
// keep the time to keep this task down to 4hrs so serves as
// an effective way to check my work.
func TestEthClient_Integration_Test(t *testing.T) {
	// Initialize client
	client, err := NewEthClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create eth_blockNumber request
	blockNumberReq, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	})

	// Execute request
	resp, err := client.Execute(blockNumberReq)
	if err != nil {
		t.Fatalf("Failed to get block number: %v", err)
	}

	// Verify we got a hex string response
	result, ok := resp["result"].(string)
	if !ok {
		t.Fatal("Expected block number result to be a string")
	}
	if len(result) < 3 || result[:2] != "0x" {
		t.Fatalf("Expected hex string starting with 0x, got %s", result)
	}
}

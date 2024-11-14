package ethclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Using an integration test since this package doesn't do a whole lot
// apart from interacting with the Ethereum node. Additionally, need to
// keep the time to keep this task down to 4hrs so serves as
// an effective way to check my work.
func TestEthClient_Integration_Test(t *testing.T) {
	type testCase struct {
		name    string
		request map[string]interface{}
		assert  func(t *testing.T, resp map[string]interface{})
	}

	// Initialize client
	client, err := NewEthClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []testCase{
		{
			name: "valid block number request",
			request: map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "eth_blockNumber",
				"params":  []interface{}{},
				"id":      1,
			},
			assert: func(t *testing.T, resp map[string]interface{}) {
				if resp["error"] != nil {
					t.Fatal("unexpected error response")
				}
				if resp["result"] == nil {
					t.Fatal("missing result in response")
				}
				hexStr := fmt.Sprint(resp["result"])
				if !strings.HasPrefix(hexStr, "0x") {
					t.Fatalf("expected hex string starting with 0x, got %s", hexStr)
				}
			},
		},
		{
			name: "invalid method request",
			request: map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "eth_invalidMethod",
				"params":  []interface{}{},
				"id":      2,
			},
			assert: func(t *testing.T, resp map[string]interface{}) {
				if resp["error"] == nil {
					t.Fatal("expected error response")
				}
				errorResp := resp["error"].(map[string]interface{})
				if errorResp["code"] == nil {
					t.Fatal("missing error code")
				}
				if errorResp["message"] == nil {
					t.Fatal("missing error message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			resp, err := client.Execute(reqBody)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}

			tt.assert(t, resp)
		})
	}
}

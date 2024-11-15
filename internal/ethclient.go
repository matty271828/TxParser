package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// EthereumClient defines the interface for interacting with an Ethereum node
type EthereumClient interface {
	// Execute sends a JSON-RPC request to the Ethereum node and returns the decoded response
	Execute(body []byte) (*jsonRPCResponse, error)
}

// EthClient implements the EthereumClient interface.
// Not really going to use the interface much but would allow
// us to mock and also ensure a clear contract for the EthClient
// struct.
var _ EthereumClient = (*EthClient)(nil)

type EthClient struct {
	rpcURL string
	client *http.Client
}

func NewEthClient() (*EthClient, error) {
	return &EthClient{
		rpcURL: "https://ethereum-rpc.publicnode.com",
		client: &http.Client{},
	}, nil
}

// Add this struct at the top level
type jsonRPCResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Modify Execute to return the parsed response
func (c *EthClient) Execute(body []byte) (*jsonRPCResponse, error) {
	resp, err := c.client.Post(c.rpcURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", response.Error.Message)
	}

	return &response, nil
}

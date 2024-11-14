package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// EthereumClient defines the interface for interacting with an Ethereum node
type EthereumClient interface {
	// Execute sends a JSON-RPC request to the Ethereum node and returns the decoded response
	Execute(body []byte) (map[string]interface{}, error)
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

// Execute sends a JSON request and returns the decoded response
func (c *EthClient) Execute(body []byte) (map[string]interface{}, error) {
	resp, err := c.client.Post(c.rpcURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

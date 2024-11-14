package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

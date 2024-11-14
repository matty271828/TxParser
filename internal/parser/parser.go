package parser

import (
	"encoding/json"
	"fmt"
	"txparser/internal/ethclient"
	"txparser/internal/storer"
)

type Manager struct {
	ethClient *ethclient.EthClient
	storer    *storer.Storer
}

// NewManager creates a new parser manager.
func NewManager(
	ethClient *ethclient.EthClient,
	storer *storer.Storer,
) (*Manager, error) {
	return &Manager{
		ethClient: ethClient,
		storer:    storer,
	}, nil
}

// GetCurrentBlock returns the current block number.
func (m *Manager) GetCurrentBlock() int {
	// Create the JSON-RPC request
	body, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	})
	if err != nil {
		return -1
	}

	// Send request using ethclient
	response, err := m.ethClient.Execute(body)
	if err != nil {
		return -1
	}

	// Handle response
	result, ok := response["result"].(string)
	if !ok {
		return -1
	}

	var blockNum int
	fmt.Sscanf(result, "0x%x", &blockNum)
	return blockNum
}

// SubscribeAddress is used to subscribe to notifications
// for a given address.
func (m *Manager) SubscribeAddress(address string) bool {
	return m.storer.Subscribe(address)
}

// GetTransactions returns all transactions for a given address.
func (m *Manager) GetTransactions(address string) []Transaction {
	params := map[string]interface{}{
		"fromBlock": "0x0",
		"toBlock":   "latest",
		"address":   address,
	}

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getLogs",
		"params":  []interface{}{params},
		"id":      1,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil
	}

	response, err := m.ethClient.Execute(body)
	if err != nil {
		return nil
	}

	return parseTransactions(response)
}

// parseTransactions converts a JSON-RPC response into a slice of transactions.
// Returns nil if the response cannot be parsed.
func parseTransactions(response map[string]interface{}) []Transaction {
	result, ok := response["result"].([]interface{})
	if !ok {
		return nil
	}

	var transactions []Transaction
	for _, log := range result {
		tx, ok := parseTransaction(log)
		if !ok {
			continue
		}
		transactions = append(transactions, tx)
	}
	return transactions
}

// parseTransaction converts a single log entry into a Transaction.
// Returns the Transaction and true if successful, or an empty Transaction and false if parsing fails.
func parseTransaction(log interface{}) (Transaction, bool) {
	logMap, ok := log.(map[string]interface{})
	if !ok {
		return Transaction{}, false
	}

	return Transaction{
		Hash:  logMap["transactionHash"].(string),
		From:  logMap["from"].(string),
		To:    logMap["to"].(string),
		Value: logMap["value"].(string),
	}, true
}

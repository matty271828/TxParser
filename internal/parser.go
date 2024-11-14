package internal

import (
	"encoding/json"
	"fmt"
	"log"
)

// Public interface that defines operations that can be performed
// to interact with the ethereum blockchain.
type Parser interface {
	// Get the last parsed block number.
	GetCurrentBlock() int

	// Add an address to the observer.
	SubscribeAddress(address string) bool

	// List of outbound or inbound transactions for a given address.
	GetTransactions(address string) []Transaction
}

// Represents a transaction on the ethereum blockchain.
type Transaction struct {
	// Hash of the transaction.
	Hash string
	// Address of the sender.
	From string
	// Address of the receiver.
	To string
	// Value of the transaction in ETH.
	Value string
}

type Manager struct {
	ethClient *EthClient
	storer    *Storer
}

// Ensure Manager implements Parser interface
var _ Parser = (*Manager)(nil)

// NewManager creates a new parser manager.
func NewManager(
	ethClient *EthClient,
	storer *Storer,
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
		log.Printf("Error marshaling block number request: %v", err)
		return -1
	}

	// Send request using ethclient
	response, err := m.ethClient.Execute(body)
	if err != nil {
		log.Printf("Error getting current block: %v", err)
		return -1
	}

	// Handle response
	result, ok := response["result"].(string)
	if !ok {
		log.Printf("Error parsing block number response: invalid format")
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
	if address == "" {
		log.Printf("Error: empty address provided")
		return nil
	}

	// Create request to get logs for the address (both incoming and outgoing)
	params := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getLogs",
		"params": []interface{}{
			map[string]interface{}{
				"fromBlock": "0x0",
				"toBlock":   "latest",
				"address":   address,
			},
		},
		"id": 1,
	}

	body, err := json.Marshal(params)
	if err != nil {
		log.Printf("Error marshaling transaction request: %v", err)
		return nil
	}

	response, err := m.ethClient.Execute(body)
	if err != nil {
		log.Printf("Error getting transactions: %v", err)
		return nil
	}

	return m.parseTransactions(response)
}

// parseTransactions converts a JSON-RPC response into a slice of transactions.
// Returns nil if the response cannot be parsed.
func (m *Manager) parseTransactions(response map[string]interface{}) []Transaction {
	result, ok := response["result"].([]interface{})
	if !ok {
		return nil
	}

	var transactions []Transaction
	for _, tx := range result {
		txMap, ok := tx.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract transaction hash from the log
		hash, ok := txMap["transactionHash"].(string)
		if !ok {
			continue
		}

		// Create additional request to get full transaction details
		params := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_getTransactionByHash",
			"params":  []interface{}{hash},
			"id":      1,
		}

		body, err := json.Marshal(params)
		if err != nil {
			continue
		}

		txResponse, err := m.ethClient.Execute(body)
		if err != nil {
			continue
		}

		txResult, ok := txResponse["result"].(map[string]interface{})
		if !ok {
			continue
		}

		tx := Transaction{
			Hash:  hash,
			From:  txResult["from"].(string),
			To:    txResult["to"].(string),
			Value: txResult["value"].(string),
		}
		transactions = append(transactions, tx)
	}

	return transactions
}

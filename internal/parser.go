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
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
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

	response, err := m.ethClient.Execute(body)
	if err != nil {
		log.Printf("Error getting current block: %v", err)
		return -1
	}

	result, ok := response.Result.(string)
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
	// Create filter with properly formatted hex values
	body, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getLogs",
		"params": []interface{}{
			map[string]interface{}{
				"address": address,
			},
		},
		"id": 1,
	})
	if err != nil {
		log.Printf("Error marshaling logs request: %v", err)
		return nil
	}

	log.Printf("Sending request: %s", string(body))

	response, err := m.ethClient.Execute(body)
	if err != nil {
		log.Printf("Error getting logs: %v", err)
		return nil
	}

	var logs []Transaction
	logsData, err := json.Marshal(response.Result)
	if err != nil {
		log.Printf("Error marshaling logs data: %v", err)
		return nil
	}

	if err := json.Unmarshal(logsData, &logs); err != nil {
		log.Printf("Error unmarshaling logs: %v", err)
		return nil
	}

	var transactions []Transaction
	for _, log := range logs {
		tx := Transaction{
			Address:          log.Address,
			BlockHash:        log.BlockHash,
			BlockNumber:      log.BlockNumber,
			Data:             log.Data,
			LogIndex:         log.LogIndex,
			Topics:           log.Topics,
			TransactionHash:  log.TransactionHash,
			TransactionIndex: log.TransactionIndex,
		}
		transactions = append(transactions, tx)
	}

	log.Printf("Found %d transactions", len(transactions))
	return []Transaction{}
}

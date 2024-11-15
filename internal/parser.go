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
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	ChainId          string `json:"chainId"`
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

type getBlockByNumResponse struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	ParentHash       string        `json:"parentHash"`
}

// GetTransactions returns all transactions for a given address.
//
// Note: In reality I would spend time adding pagination to this method.
// We are also only checking the latest block here, in reality we would
// need to go back through multiple blocks. Another thing we would need to
// handle is filtering out transactions returned that do not match the address
// specified in the search parameter.
func (m *Manager) GetTransactions(address string) []Transaction {
	currentBlock := m.GetCurrentBlock()
	if currentBlock == -1 {
		log.Printf("Error getting current block")
		return nil
	}

	// TODO: Iterate through all blocks and get transactions

	blockNumHex := fmt.Sprintf("0x%x", currentBlock)
	body, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []interface{}{
			blockNumHex,
			true,
		},
		"id": 1,
	})
	if err != nil {
		log.Printf("Error marshaling block request: %v", err)
		return nil
	}

	response, err := m.ethClient.Execute(body)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil
	}

	// Marshal and unmarshal to convert from interface{} to our struct
	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		log.Printf("Error marshaling result: %v", err)
		return nil
	}
	var getBlockByNumResponse getBlockByNumResponse
	if err := json.Unmarshal(resultBytes, &getBlockByNumResponse); err != nil {
		log.Printf("Error decoding block data: %v", err)
		return nil
	}

	if len(getBlockByNumResponse.Transactions) == 0 {
		return []Transaction{}
	}

	// TODO: Filter transactions for the specific address

	return getBlockByNumResponse.Transactions
}

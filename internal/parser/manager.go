package parser

// Public interface that defines operations that can be performed
// to interact with the ethereum blockchain.
type Parser interface {
	// Get the last parsed block number.
	GetCurrentBlock() (int, error)

	// Add an address to the observer.
	SubscribeAddress(address string) error

	// List of outbound or inbound transactions for a given address.
	GetTransactions(address string) ([]Transaction, error)
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

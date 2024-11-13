package storer

type Manager struct {
	// TODO: Add dependencies
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

func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

func (m *Manager) SaveTransaction(address string, tx Transaction) {
	// TODO: Implement
}

func (m *Manager) GetTransactions(address string) []Transaction {
	// TODO: Implement
	return nil
}

func (m *Manager) Subscribe(address string) bool {
	// TODO: Implement
	return false
}

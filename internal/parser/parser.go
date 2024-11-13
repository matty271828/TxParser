package parser

type Manager struct {
	// TODO: Add dependencies
}

// NewManager creates a new parser manager.
func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

func (m *Manager) GetCurrentBlock() (int, error) {
	// TODO: Implement
	return 0, nil
}

func (m *Manager) SubscribeAddress(address string) error {
	// TODO: Implement
	return nil
}

func (m *Manager) GetTransactions(address string) ([]Transaction, error) {
	// TODO: Implement
	return nil, nil
}

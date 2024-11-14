package parser

import (
	"txparser/internal/ethclient"
	"txparser/internal/storer"
)

type Manager struct {
	ethClient *ethclient.EthClient
	storer    storer.Storer
}

// NewManager creates a new parser manager.
func NewManager(
	ethClient *ethclient.EthClient,
	storer storer.Storer,
) (*Manager, error) {
	return &Manager{
		ethClient: ethClient,
		storer:    storer,
	}, nil
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

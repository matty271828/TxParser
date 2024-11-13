package parser

import (
	"txparser/internal/ethereum"
	"txparser/internal/storer"
)

type Manager struct {
	ethClient *ethereum.Client
	storer    storer.Storer
}

// NewManager creates a new parser manager.
func NewManager(
	ethClient *ethereum.Client,
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

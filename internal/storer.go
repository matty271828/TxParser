package internal

import "sync"

// Storer is an in-memory storer for subscribed addresses.
// This could be later repurposed to connect to a database.
type Storer struct {
	subscribedAddrs map[string]bool
	mu              sync.RWMutex // Mutex to safely access the data in concurrent scenarios
}

// NewStorer initializes a new Storer.
func NewStorer() *Storer {
	return &Storer{
		subscribedAddrs: make(map[string]bool),
	}
}

// Subscribe subscribes to notifications for a given address.
func (s *Storer) Subscribe(address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.subscribedAddrs[address]; exists {
		return false
	}
	s.subscribedAddrs[address] = true
	return true
}

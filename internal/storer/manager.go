package storer

// Storer interface defines the methods needed for storing and retrieving data
type Storer interface {
	// Save a transaction to the store.
	SaveTransaction(address string, tx Transaction)

	// Get all transactions for a given address.
	GetTransactions(address string) []Transaction

	// Subscribe marks an address as being observed.
	Subscribe(address string) bool
}

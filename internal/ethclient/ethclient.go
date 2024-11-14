package ethclient

type EthClient struct {
	// TODO: Add dependencies
}

// NewEthClient creates a new Ethereum client.
func NewEthClient() (*EthClient, error) {
	return &EthClient{}, nil
}

func (c *EthClient) GetCurrentBlock() (int, error) {
	// TODO: Implement
	return 0, nil
}

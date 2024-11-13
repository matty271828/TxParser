package ethereum

type Client struct {
	// TODO: Add dependencies
}

// NewClient creates a new Ethereum client.
func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) GetCurrentBlock() (int, error) {
	// TODO: Implement
	return 0, nil
}

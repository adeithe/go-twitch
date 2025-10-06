package api

// CharityResource handles charity related API calls.
type CharityResource struct {
	client *Client
}

// NewCharityResource creates a new CharityResource.
func NewCharityResource(client *Client) *CharityResource {
	return &CharityResource{client}
}

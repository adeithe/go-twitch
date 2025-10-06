package api

// RaidsResource represents the Twitch Raids API.
type RaidsResource struct {
	client *Client
}

// NewRaidsResource creates a new RaidsResource.
func NewRaidsResource(client *Client) *RaidsResource {
	return &RaidsResource{client}
}

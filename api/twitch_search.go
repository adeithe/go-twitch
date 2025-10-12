package api

// SearchResource represents the Twitch Search API.
type SearchResource struct {
	client *Client
}

// NewSearchResource creates a new SearchResource.
func NewSearchResource(client *Client) *SearchResource {
	return &SearchResource{client}
}

package api

// SearchResource represents the Twitch Search API.
type SearchResource struct {
	client *Client
}

const (
	// EndpointSearchCategories is the endpoint for searching categories.
	EndpointSearchCategories = TwitchAPIVersionHelix + "/search/categories"
	// EndpointSearchChannels is the endpoint for searching channels.
	EndpointSearchChannels = TwitchAPIVersionHelix + "/search/channels"
)

// NewSearchResource creates a new SearchResource.
func NewSearchResource(client *Client) *SearchResource {
	return &SearchResource{client}
}

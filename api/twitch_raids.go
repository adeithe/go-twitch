package api

// RaidsResource represents the Twitch Raids API.
type RaidsResource struct {
	client *Client
}

// EndpointRaids is the endpoint for managing raids.
const EndpointRaids = TwitchAPIVersionHelix + "/raids"

// NewRaidsResource creates a new RaidsResource.
func NewRaidsResource(client *Client) *RaidsResource {
	return &RaidsResource{client}
}

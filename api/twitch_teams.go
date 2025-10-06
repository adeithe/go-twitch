package api

// TeamsResource represents the Twitch Teams API.
type TeamsResource struct {
	client *Client
}

// NewTeamsResource creates a new TeamsResource.
func NewTeamsResource(client *Client) *TeamsResource {
	return &TeamsResource{client}
}

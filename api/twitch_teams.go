package api

// TeamsResource represents the Twitch Teams API.
type TeamsResource struct {
	client *Client
}

const (
	// EndpointTeams is the endpoint for getting teams.
	EndpointTeams = TwitchAPIVersionHelix + "/teams"
	// EndpointTeamsGetChannelTeams is the endpoint for getting a team by name.
	EndpointTeamsGetChannelTeams = TwitchAPIVersionHelix + "/teams/channel"
)

// NewTeamsResource creates a new TeamsResource.
func NewTeamsResource(client *Client) *TeamsResource {
	return &TeamsResource{client}
}

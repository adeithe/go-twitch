package api

// GoalsResource represents the Twitch Goals API.
type GoalsResource struct {
	client *Client
}

// EndpointGoals is the endpoint for managing creator goals.
const EndpointGoals = TwitchAPIVersionHelix + "/goals"

// NewGoalsResource creates a new GoalsResource.
func NewGoalsResource(client *Client) *GoalsResource {
	return &GoalsResource{client}
}

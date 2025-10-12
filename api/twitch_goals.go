package api

// GoalsResource represents the Twitch Goals API.
type GoalsResource struct {
	client *Client
}

// NewGoalsResource creates a new GoalsResource.
func NewGoalsResource(client *Client) *GoalsResource {
	return &GoalsResource{client}
}

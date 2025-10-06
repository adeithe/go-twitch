package api

// PollsResource represents the Twitch Polls API.
type PollsResource struct {
	client *Client
}

// NewPollsResource creates a new PollsResource.
func NewPollsResource(client *Client) *PollsResource {
	return &PollsResource{client}
}

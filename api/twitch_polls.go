package api

// PollsResource represents the Twitch Polls API.
type PollsResource struct {
	client *Client
}

// EndpointPolls is the endpoint for the Twitch Polls API.
const EndpointPolls = TwitchAPIVersionHelix + "/polls"

// NewPollsResource creates a new PollsResource.
func NewPollsResource(client *Client) *PollsResource {
	return &PollsResource{client}
}

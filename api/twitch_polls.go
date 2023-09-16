package api

type PollsResource struct {
	client *Client
}

func NewPollsResource(client *Client) *PollsResource {
	return &PollsResource{client}
}

package api

type StreamsResource struct {
	client *Client
}

func NewStreamsResource(client *Client) *StreamsResource {
	return &StreamsResource{client}
}

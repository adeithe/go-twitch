package api

type ChannelPointsResource struct {
	client *Client
}

func NewChannelPointsResource(client *Client) *ChannelPointsResource {
	return &ChannelPointsResource{client}
}

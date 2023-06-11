package api

type ChannelsResource struct {
	client *Client
}

func NewChannelsResource(client *Client) *ChannelsResource {
	return &ChannelsResource{client}
}

package api

type WhispersResource struct {
	client *Client
}

func NewWhispersResource(client *Client) *WhispersResource {
	return &WhispersResource{client}
}

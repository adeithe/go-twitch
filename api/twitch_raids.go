package api

type RaidsResource struct {
	client *Client
}

func NewRaidsResource(client *Client) *RaidsResource {
	return &RaidsResource{client}
}

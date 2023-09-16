package api

type CharityResource struct {
	client *Client
}

func NewCharityResource(client *Client) *CharityResource {
	return &CharityResource{client}
}

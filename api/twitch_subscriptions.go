package api

type SubscriptionsResource struct {
	client *Client
}

func NewSubscriptionsResource(client *Client) *SubscriptionsResource {
	return &SubscriptionsResource{client}
}

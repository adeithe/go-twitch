package api

// SubscriptionsResource represents the Twitch Subscriptions API.
type SubscriptionsResource struct {
	client *Client
}

// NewSubscriptionsResource creates a new SubscriptionsResource.
func NewSubscriptionsResource(client *Client) *SubscriptionsResource {
	return &SubscriptionsResource{client}
}

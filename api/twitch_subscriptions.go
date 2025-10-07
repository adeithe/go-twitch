package api

// SubscriptionsResource represents the Twitch Subscriptions API.
type SubscriptionsResource struct {
	client *Client
}

const (
	// EndpointSubscriptionsGetBroadcasterSubscriptions is the endpoint for getting broadcaster's subscriptions.
	EndpointSubscriptionsGetBroadcasterSubscriptions = TwitchAPIVersionHelix + "/subscriptions"
	// EndpointSubscriptionsCheckUserSubscription is the endpoint for getting a user's subscription.
	EndpointSubscriptionsCheckUserSubscription = TwitchAPIVersionHelix + "/subscriptions/user"
)

// NewSubscriptionsResource creates a new SubscriptionsResource.
func NewSubscriptionsResource(client *Client) *SubscriptionsResource {
	return &SubscriptionsResource{client}
}

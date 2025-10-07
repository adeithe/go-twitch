package api

// EventSubResource represents the Twitch EventSub API.
type EventSubResource struct {
	client *Client
}

// EndpointEventSubSubscriptions is the endpoint for managing EventSub subscriptions.
const EndpointEventSubSubscriptions = TwitchAPIVersionHelix + "/eventsub/subscriptions"

// NewEventSubResource creates a new EventSubResource.
func NewEventSubResource(client *Client) *EventSubResource {
	return &EventSubResource{client}
}

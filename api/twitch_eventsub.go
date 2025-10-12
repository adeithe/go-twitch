package api

// EventSubResource represents the Twitch EventSub API.
type EventSubResource struct {
	client *Client
}

// NewEventSubResource creates a new EventSubResource.
func NewEventSubResource(client *Client) *EventSubResource {
	return &EventSubResource{client}
}

package api

type EventSubResource struct {
	client *Client
}

func NewEventSubResource(client *Client) *EventSubResource {
	return &EventSubResource{client}
}

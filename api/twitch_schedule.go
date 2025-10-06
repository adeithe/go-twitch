package api

// ScheduleResource represents the Twitch Schedule API.
type ScheduleResource struct {
	client *Client
}

// NewScheduleResource creates a new ScheduleResource.
func NewScheduleResource(client *Client) *ScheduleResource {
	return &ScheduleResource{client}
}

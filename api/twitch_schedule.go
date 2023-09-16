package api

type ScheduleResource struct {
	client *Client
}

func NewScheduleResource(client *Client) *ScheduleResource {
	return &ScheduleResource{client}
}

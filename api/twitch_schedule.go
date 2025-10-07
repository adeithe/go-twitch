package api

// ScheduleResource represents the Twitch Schedule API.
type ScheduleResource struct {
	client *Client
}

const (
	// EndpointScheduleGetChannelStreamSchedule is the endpoint for getting the stream schedule for a channel.
	EndpointScheduleGetChannelStreamSchedule = TwitchAPIVersionHelix + "/schedule"
	// EndpointScheduleGetChannelCalendar is the endpoint for getting the iCalendar schedule for a channel.
	EndpointScheduleGetChannelCalendar = TwitchAPIVersionHelix + "/schedule/icalendar"
	// EndpointScheduleChannelSettings is the endpoint for getting or updating channel schedule settings.
	EndpointScheduleChannelSettings = TwitchAPIVersionHelix + "/schedule/settings"
	// EndpointScheduleCreateChannelSegment is the endpoint for managing stream schedule segments for a channel.
	EndpointScheduleCreateChannelSegment = TwitchAPIVersionHelix + "/schedule/segment"
)

// NewScheduleResource creates a new ScheduleResource.
func NewScheduleResource(client *Client) *ScheduleResource {
	return &ScheduleResource{client}
}

package api

// HypeTrainResource represents the Twitch Hype Train API.
type HypeTrainResource struct {
	client *Client
}

const (
	// EndpointHypeTrainGetEvents is the endpoint for getting Hype Train events.
	//
	// Deprecated: Scheduled for removal on December 4, 2025. Use "Get Hype Train Status" instead.
	EndpointHypeTrainGetEvents = TwitchAPIVersionHelix + "/hypetrain/events"
	// EndpointHypeTrainGetStatus is the endpoint for getting the status of a Hype Train.
	EndpointHypeTrainGetStatus = TwitchAPIVersionHelix + "/hypetrain/status"
)

// NewHypeTrainResource creates a new HypeTrainResource.
func NewHypeTrainResource(client *Client) *HypeTrainResource {
	return &HypeTrainResource{client}
}

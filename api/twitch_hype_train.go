package api

// HypeTrainResource represents the Twitch Hype Train API.
type HypeTrainResource struct {
	client *Client
}

// NewHypeTrainResource creates a new HypeTrainResource.
func NewHypeTrainResource(client *Client) *HypeTrainResource {
	return &HypeTrainResource{client}
}

package api

// PredictionsResource represents the Twitch Predictions API.
type PredictionsResource struct {
	client *Client
}

// EndpointPredictions is the endpoint for managing channel predictions.
const EndpointPredictions = TwitchAPIVersionHelix + "/predictions"

// NewPredictionsResource creates a new PredictionsResource.
func NewPredictionsResource(client *Client) *PredictionsResource {
	return &PredictionsResource{client}
}

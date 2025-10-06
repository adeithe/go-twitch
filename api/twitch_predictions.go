package api

// PredictionsResource represents the Twitch Predictions API.
type PredictionsResource struct {
	client *Client
}

// NewPredictionsResource creates a new PredictionsResource.
func NewPredictionsResource(client *Client) *PredictionsResource {
	return &PredictionsResource{client}
}

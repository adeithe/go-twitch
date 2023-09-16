package api

type PredictionsResource struct {
	client *Client
}

func NewPredictionsResource(client *Client) *PredictionsResource {
	return &PredictionsResource{client}
}

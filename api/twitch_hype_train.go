package api

type HypeTrainResource struct {
	client *Client
}

func NewHypeTrainResource(client *Client) *HypeTrainResource {
	return &HypeTrainResource{client}
}

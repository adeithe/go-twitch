package api

type AnalyticsResource struct {
	client *Client
}

func NewAnalyticsResource(client *Client) *AnalyticsResource {
	return &AnalyticsResource{client}
}

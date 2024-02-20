package api

import "time"

type AnalyticsResource struct {
	client     *Client
	Extensions *AnalyticsExtensionResource
	Games      *AnalyticsGameResource
}

type AnalyticsDateRange struct {
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

func NewAnalyticsResource(client *Client) *AnalyticsResource {
	r := &AnalyticsResource{client: client}
	r.Extensions = NewAnalyticsExtensionResource(client)
	r.Games = NewAnalyticsGameResource(client)
	return r
}

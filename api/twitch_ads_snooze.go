package api

import (
	"context"
	"net/http"
	"time"
)

type AdsSnoozeResource struct {
	client *Client
}

func NewAdsSnoozeResource(client *Client) *AdsSnoozeResource {
	return &AdsSnoozeResource{client}
}

type AdsSnoozeRequest struct {
	resource *AdsSnoozeResource
	opts     []RequestOption
}

type AdsSnoozeResponse struct {
	Header http.Header
	Data   []AdSnooze
}

type AdSnooze struct {
	SnoozeCount     int       `json:"snooze_count"`
	SnoozeRefreshAt time.Time `json:"snooze_refresh_at"`
	NextAdAt        time.Time `json:"next_ad_at"`
}

func (r *AdsSnoozeResource) Insert(broadcasterId string) *AdsSnoozeRequest {
	return &AdsSnoozeRequest{
		resource: r,
		opts: []RequestOption{
			AddQueryParameter("broadcaster_id", broadcasterId),
		},
	}
}

// Do executes the request.
func (c *AdsSnoozeRequest) Do(ctx context.Context, opts ...RequestOption) (*AdsSnoozeResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodPost, "/channels/ads/schedule/snooze", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[AdSnooze](res)
	if err != nil {
		return nil, err
	}

	return &AdsSnoozeResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

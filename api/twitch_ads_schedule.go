package api

import (
	"context"
	"net/http"
	"time"
)

type AdsScheduleResource struct {
	client *Client
}

func NewAdsScheduleResource(client *Client) *AdsScheduleResource {
	return &AdsScheduleResource{client}
}

type AdsScheduleListRequest struct {
	resource *AdsScheduleResource
	opts     []RequestOption
}

type AdsScheduleListResponse struct {
	Header http.Header
	Data   []AdSchedule
}

type AdSchedule struct {
	Duration        int       `json:"duration"`
	NextAdAt        time.Time `json:"next_ad_at"`
	LastAdAt        time.Time `json:"last_ad_at"`
	PrerollFreeTime int       `json:"preroll_free_time"`
	SnoozeCount     int       `json:"snooze_count"`
	SnoozeRefreshAt time.Time `json:"snooze_refresh_at"`
}

func (r *AdsScheduleResource) List(broadcasterId string) *AdsScheduleListRequest {
	return &AdsScheduleListRequest{
		resource: r,
		opts: []RequestOption{
			AddQueryParameter("broadcaster_id", broadcasterId),
		},
	}
}

// Do executes the request.
func (c *AdsScheduleListRequest) Do(ctx context.Context, opts ...RequestOption) (*AdsScheduleListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channels/ads", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[AdSchedule](res)
	if err != nil {
		return nil, err
	}

	return &AdsScheduleListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

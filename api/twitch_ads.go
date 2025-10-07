package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Commercial contains information about a commercial that was started.
type Commercial struct {
	Length     int    `json:"length"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

// AdsResource provides methods for the Twitch Ads API.
type AdsResource struct {
	client   *Client
	Schedule *AdsScheduleResource
	Snooze   *AdsSnoozeResource
}

// NewAdsResource creates a new AdsResource.
func NewAdsResource(client *Client) *AdsResource {
	r := &AdsResource{client: client}
	r.Schedule = NewAdsScheduleResource(client)
	r.Snooze = NewAdsSnoozeResource(client)
	return r
}

// AdsInsertRequest is a request to start a commercial.
type AdsInsertRequest struct {
	client        *Client
	broadcasterID string
	duration      int
}

// AdsInsertResponse is the response from starting a commercial.
type AdsInsertResponse struct {
	Header http.Header
	Data   []Commercial
}

const (
	// EndpointAdsStartCommercial is the endpoint for starting a commercial.
	EndpointAdsStartCommercial = TwitchAPIVersionHelix + "/channels/commercial"
	// EndpointAdsGetAdsSchedule is the endpoint for getting ad schedule information.
	EndpointAdsGetAdsSchedule = TwitchAPIVersionHelix + "/channels/ads"
	// EndpointAdsSnoozeNextAd is the endpoint for snoozing the next ad.
	EndpointAdsSnoozeNextAd = TwitchAPIVersionHelix + "/channels/ads/schedule/snooze"
)

// Insert creates a request to start a commercial for the specified broadcaster.
//
// Required Scope: channel:edit:commercial
func (r *AdsResource) Insert(broadcasterID string) *AdsInsertRequest {
	return &AdsInsertRequest{r.client, broadcasterID, 60}
}

// Duration sets the duration of the commercial in seconds.
//
// Twitch tries to serve a commercial that’s the requested length, but it may be shorter or longer. The maximum length you should request is 180 seconds.
func (c *AdsInsertRequest) Duration(seconds int) *AdsInsertRequest {
	c.duration = seconds
	return c
}

// Do executes the request.
//
//	req := client.Ads.Insert("41245072").Duration(60)
//	data, err := req.Do(ctx, api.WithBearerToken("2gbdx6oar67tqtcmt49t3wpcgycthx")
func (c *AdsInsertRequest) Do(ctx context.Context, opts ...RequestOption) (*AdsInsertResponse, error) {
	bs, err := json.Marshal(map[string]any{
		"broadcaster_id": c.broadcasterID,
		"length":         c.duration,
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.doRequest(ctx, http.MethodPost, EndpointAdsStartCommercial, bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Commercial](res)
	if err != nil {
		return nil, err
	}

	return &AdsInsertResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

// AdsScheduleResource provides methods for the Twitch Ads Schedule API.
type AdsScheduleResource struct {
	client *Client
}

// NewAdsScheduleResource creates a new AdsScheduleResource.
func NewAdsScheduleResource(client *Client) *AdsScheduleResource {
	return &AdsScheduleResource{client}
}

// AdsScheduleListRequest is a request to list ad schedule information.
type AdsScheduleListRequest struct {
	client *Client
	opts   []RequestOption
}

// AdsScheduleListResponse is the response from listing ad schedule information.
type AdsScheduleListResponse struct {
	Header http.Header
	Data   []AdSchedule
}

// AdSchedule contains information about the ad schedule for a broadcaster.
type AdSchedule struct {
	Duration        int       `json:"duration"`
	NextAdAt        time.Time `json:"next_ad_at"`
	LastAdAt        time.Time `json:"last_ad_at"`
	PrerollFreeTime int       `json:"preroll_free_time"`
	SnoozeCount     int       `json:"snooze_count"`
	SnoozeRefreshAt time.Time `json:"snooze_refresh_at"`
}

// List returns ad schedule related information, including snooze, when the last ad was run,
// when the next ad is scheduled, and if the channel is currently in pre-roll free time.
//
// NOTE: A new advertisement can NOT be run until 8 minutes running the previous.
//
// Required Scope: channel:read:ads
func (r *AdsScheduleResource) List(broadcasterID string) *AdsScheduleListRequest {
	return &AdsScheduleListRequest{
		client: r.client,
		opts: []RequestOption{
			AddQueryParameter("broadcaster_id", broadcasterID),
		},
	}
}

// Do executes the request.
func (r *AdsScheduleListRequest) Do(ctx context.Context, opts ...RequestOption) (*AdsScheduleListResponse, error) {
	res, err := r.client.doRequest(ctx, http.MethodGet, EndpointAdsGetAdsSchedule, nil, append(r.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[AdSchedule](res)
	if err != nil {
		return nil, err
	}

	return &AdsScheduleListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

// AdsSnoozeResource provides methods for the Twitch Ads Snooze API.
type AdsSnoozeResource struct {
	client *Client
}

// NewAdsSnoozeResource creates a new AdsSnoozeResource.
func NewAdsSnoozeResource(client *Client) *AdsSnoozeResource {
	return &AdsSnoozeResource{client}
}

// AdsSnoozeRequest is a request to snooze ads for a broadcaster.
type AdsSnoozeRequest struct {
	client *Client
	opts   []RequestOption
}

// AdsSnoozeResponse is the response from snoozing ads for a broadcaster.
type AdsSnoozeResponse struct {
	Header http.Header
	Data   []AdSnooze
}

// AdSnooze contains information about the ad snooze status for a broadcaster.
type AdSnooze struct {
	SnoozeCount     int       `json:"snooze_count"`
	SnoozeRefreshAt time.Time `json:"snooze_refresh_at"`
	NextAdAt        time.Time `json:"next_ad_at"`
}

// Insert if available, pushes back the timestamp of the upcoming automatic advertisement mid-roll by 5 minutes.
// This endpoint duplicates the snooze functionality in the creator dashboard's Ad Manager.
//
// Required Scope: channel:manage:ads
func (r *AdsSnoozeResource) Insert(broadcasterID string) *AdsSnoozeRequest {
	return &AdsSnoozeRequest{
		client: r.client,
		opts: []RequestOption{
			AddQueryParameter("broadcaster_id", broadcasterID),
		},
	}
}

// Do executes the request.
func (r *AdsSnoozeRequest) Do(ctx context.Context, opts ...RequestOption) (*AdsSnoozeResponse, error) {
	res, err := r.client.doRequest(ctx, http.MethodPost, EndpointAdsSnoozeNextAd, nil, append(r.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[AdSnooze](res)
	if err != nil {
		return nil, err
	}

	return &AdsSnoozeResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

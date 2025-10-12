package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// AdsResource represents the Twitch Ads API.
type AdsResource struct {
	client *Client

	// Snooze provides access to the Twitch Snooze API.
	Snooze *AdsSnoozeResource
}

// NewAdsResource creates a new AdsResource.
func NewAdsResource(client *Client) *AdsResource {
	r := &AdsResource{client: client}
	r.Snooze = NewAdsSnoozeResource(client)
	return r
}

// StartCommercialInsertCall represents a POST call to a Twitch Ads API endpoint.
type StartCommercialInsertCall struct {
	resource *AdsResource
	body     map[string]any
}

// StartCommercialResponse represents the response from a POST request to /helix/channels/commercial.
type StartCommercialResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Commercial data returned by the Twitch API.
	Data []Commercial
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/channels/commercial.
//
// Starts a commercial on the specified channel.
//
// Only Twitch Partners and Affiliates can run commercials on their channels and must be live.
//
// # Authorization
//
// Requires the broadcasters access token that includes the channel:edit:commercial scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#start-commercial
func (r *AdsResource) Insert() *StartCommercialInsertCall {
	return &StartCommercialInsertCall{resource: r, body: make(map[string]any)}
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *StartCommercialInsertCall) BroadcasterID(broadcasterID string) *StartCommercialInsertCall {
	api.body["broadcaster_id"] = broadcasterID
	return api
}

// Length sets the Length body parameter.
func (api *StartCommercialInsertCall) Length(length int) *StartCommercialInsertCall {
	api.body["length"] = length
	return api
}

// Do executes the request.
func (api *StartCommercialInsertCall) Do(ctx context.Context, opts ...RequestOption) (*StartCommercialResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/channels/commercial", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Commercial](res)
	if err != nil {
		return nil, err
	}

	return &StartCommercialResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// AdScheduleListCall represents a GET call to a Twitch Ads API endpoint.
type AdScheduleListCall struct {
	resource *AdsResource
	opts     []RequestOption
}

// AdScheduleResponse represents the response from a GET request to /helix/channels/ads.
type AdScheduleResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the AdSchedule data returned by the Twitch API.
	Data []AdSchedule
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channels/ads.
//
// Gets ad schedule related information, including snooze, when the last ad was run, when the next ad is scheduled, and if the channel is currently in pre-roll free time.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:ads scope. The user_id in the user access token must match the broadcaster_id.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-ad-schedule
func (r *AdsResource) List() *AdScheduleListCall {
	return &AdScheduleListCall{resource: r}
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *AdScheduleListCall) BroadcasterID(broadcasterID string) *AdScheduleListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *AdScheduleListCall) Do(ctx context.Context, opts ...RequestOption) (*AdScheduleResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channels/ads", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[AdSchedule](res)
	if err != nil {
		return nil, err
	}

	return &AdScheduleResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// AdsSnoozeResource represents the Twitch AdsSnooze API.
type AdsSnoozeResource struct {
	client *Client
}

// NewAdsSnoozeResource creates a new AdsSnoozeResource.
func NewAdsSnoozeResource(client *Client) *AdsSnoozeResource {
	return &AdsSnoozeResource{client}
}

// SnoozeNextAdInsertCall represents a POST call to a Twitch AdsSnooze API endpoint.
type SnoozeNextAdInsertCall struct {
	resource *AdsSnoozeResource
	opts     []RequestOption
}

// SnoozeNextAdResponse represents the response from a POST request to /helix/channels/ads/schedule/snooze.
type SnoozeNextAdResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the AdsSnoozed data returned by the Twitch API.
	Data []AdsSnoozed
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/channels/ads/schedule/snooze.
//
// If available, pushes back the timestamp of the upcoming automatic mid-roll ad by 5 minutes.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:ads scope. The user_id in the user access token must match the broadcaster_id.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#snooze-next-ad
func (r *AdsSnoozeResource) Insert() *SnoozeNextAdInsertCall {
	return &SnoozeNextAdInsertCall{resource: r}
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *SnoozeNextAdInsertCall) BroadcasterID(broadcasterID string) *SnoozeNextAdInsertCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *SnoozeNextAdInsertCall) Do(ctx context.Context, opts ...RequestOption) (*SnoozeNextAdResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/channels/ads/schedule/snooze", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[AdsSnoozed](res)
	if err != nil {
		return nil, err
	}

	return &SnoozeNextAdResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

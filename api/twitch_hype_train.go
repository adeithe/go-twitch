package api

import (
	"context"
	"net/http"
)

// HypeTrainResource represents the Twitch HypeTrain API.
type HypeTrainResource struct {
	client *Client
}

// NewHypeTrainResource creates a new HypeTrainResource.
func NewHypeTrainResource(client *Client) *HypeTrainResource {
	return &HypeTrainResource{client}
}

// HypeTrainStatusListCall represents a GET call to a Twitch HypeTrain API endpoint.
type HypeTrainStatusListCall struct {
	resource *HypeTrainResource
	opts     []RequestOption
}

// HypeTrainStatusListResponse represents the response from a GET request to /helix/hypetrain/status.
type HypeTrainStatusListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the HypeTrainStatusInfo data returned by the Twitch API.
	Data []HypeTrainStatusInfo
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/hypetrain/status.
//
// Get the status of a Hype Train for the specified broadcaster.
//
// # Authorization
//
// Requires an user access token.
//
// Requires OAuth Scope: channel:read:hype_train.
//
// Requires that the user access token belongs to BroadcasterID.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-hype-train-status
func (r *HypeTrainResource) List(broadcasterID string) *HypeTrainStatusListCall {
	c := &HypeTrainStatusListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *HypeTrainStatusListCall) BroadcasterID(broadcasterID string) *HypeTrainStatusListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *HypeTrainStatusListCall) Do(ctx context.Context, opts ...RequestOption) (*HypeTrainStatusListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/hypetrain/status", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[HypeTrainStatusInfo](res)
	if err != nil {
		return nil, err
	}

	return &HypeTrainStatusListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

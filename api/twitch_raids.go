package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// RaidsResource represents the Twitch Raids API.
type RaidsResource struct {
	client *Client
}

// NewRaidsResource creates a new RaidsResource.
func NewRaidsResource(client *Client) *RaidsResource {
	return &RaidsResource{client}
}

// RaidInsertCall represents a POST call to a Twitch Raids API endpoint.
type RaidInsertCall struct {
	resource *RaidsResource
	body     map[string]any
}

// RaidInsertResponse represents the response from a POST request to /helix/raids.
type RaidInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the InitializedRaid data returned by the Twitch API.
	Data []InitializedRaid
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/raids.
//
// Raid another channel by sending the broadcaster's viewers to the targeted channel.
//
// When you call the API from a chat bot or extension, the Twitch UX pops up a window at the top of the chat room that identifies the number of viewers in the raid.
// The raid occurs when the broadcaster clicks Raid Now or after the 90-second countdown expires.
//
// # Rate Limit
//
// The limit is 10 requests within a 10-minute window.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:raids scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#start-a-raid
func (r *RaidsResource) Insert(fromBroadcasterID string, toBroadcasterID string) *RaidInsertCall {
	c := &RaidInsertCall{resource: r, body: make(map[string]any)}
	return c.
		FromBroadcasterID(fromBroadcasterID).
		ToBroadcasterID(toBroadcasterID)
}

// FromBroadcasterID sets the FromBroadcasterID body parameter.
func (api *RaidInsertCall) FromBroadcasterID(fromBroadcasterID string) *RaidInsertCall {
	api.body["from_broadcaster_id"] = fromBroadcasterID
	return api
}

// ToBroadcasterID sets the ToBroadcasterID body parameter.
func (api *RaidInsertCall) ToBroadcasterID(toBroadcasterID string) *RaidInsertCall {
	api.body["to_broadcaster_id"] = toBroadcasterID
	return api
}

// Do executes the request.
func (api *RaidInsertCall) Do(ctx context.Context, opts ...RequestOption) (*RaidInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/raids", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[InitializedRaid](res)
	if err != nil {
		return nil, err
	}

	return &RaidInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// RaidDeleteCall represents a DELETE call to a Twitch Raids API endpoint.
type RaidDeleteCall struct {
	resource *RaidsResource
	opts     []RequestOption
}

// RaidDeleteResponse represents the response from a DELETE request to /helix/raids.
type RaidDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/raids.
//
// Cancels a pending raid that was initiated by the broadcaster.
//
// You can cancel a raid at any point up until the broadcaster clicks Raid Now in the Twitch UX or the 90-second countdown expires.
//
// # Rate Limit
//
// The limit is 10 requests within a 10-minute window.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:raids scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#cancel-a-raid
func (r *RaidsResource) Delete(broadcasterID string) *RaidDeleteCall {
	c := &RaidDeleteCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *RaidDeleteCall) BroadcasterID(broadcasterID string) *RaidDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *RaidDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*RaidDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/raids", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &RaidDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

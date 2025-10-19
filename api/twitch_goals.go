package api

import (
	"context"
	"net/http"
)

// GoalsResource represents the Twitch Goals API.
type GoalsResource struct {
	client *Client
}

// NewGoalsResource creates a new GoalsResource.
func NewGoalsResource(client *Client) *GoalsResource {
	return &GoalsResource{client}
}

// GoalsListCall represents a GET call to a Twitch Goals API endpoint.
type GoalsListCall struct {
	resource *GoalsResource
	opts     []RequestOption
}

// GoalsListResponse represents the response from a GET request to /helix/goals.
type GoalsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CreatorGoal data returned by the Twitch API.
	Data []CreatorGoal
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/goals.
//
// Gets the broadcaster's list of active goals. Use this endpoint to get the current progress of each goal.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:goals scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-creator-goals
func (r *GoalsResource) List(broadcasterID string) *GoalsListCall {
	c := &GoalsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *GoalsListCall) BroadcasterID(broadcasterID string) *GoalsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *GoalsListCall) Do(ctx context.Context, opts ...RequestOption) (*GoalsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/goals", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CreatorGoal](res)
	if err != nil {
		return nil, err
	}

	return &GoalsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

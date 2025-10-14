package api

import (
	"context"
	"net/http"
)

// VideosResource represents the Twitch Videos API.
type VideosResource struct {
	client *Client
}

// NewVideosResource creates a new VideosResource.
func NewVideosResource(client *Client) *VideosResource {
	return &VideosResource{client}
}

// VideosListCall represents a GET call to a Twitch Videos API endpoint.
type VideosListCall struct {
	resource *VideosResource
	opts     []RequestOption
}

// VideosListResponse represents the response from a GET request to /helix/videos.
type VideosListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Video data returned by the Twitch API.
	Data []Video
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/videos.
//
// Gets information about one or more published videos.
// You may get videos by ID, by user, or by game/category.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-videos
func (r *VideosResource) List() *VideosListCall {
	return &VideosListCall{resource: r}
}

// ID adds to the ID query parameter.
func (api *VideosListCall) ID(iDs ...string) *VideosListCall {
	for _, iD := range iDs {
		api.opts = append(api.opts, AddQueryParameter("id", iD))
	}
	return api
}

// UserID adds to the UserID query parameter.
func (api *VideosListCall) UserID(userIDs ...string) *VideosListCall {
	for _, userID := range userIDs {
		api.opts = append(api.opts, AddQueryParameter("user_id", userID))
	}
	return api
}

// GameID adds to the GameID query parameter.
func (api *VideosListCall) GameID(gameIDs ...string) *VideosListCall {
	for _, gameID := range gameIDs {
		api.opts = append(api.opts, AddQueryParameter("game_id", gameID))
	}
	return api
}

// Language sets the Language query parameter.
func (api *VideosListCall) Language(language string) *VideosListCall {
	api.opts = append(api.opts, SetQueryParameter("language", language))
	return api
}

// Period sets the Period query parameter.
func (api *VideosListCall) Period(period string) *VideosListCall {
	api.opts = append(api.opts, SetQueryParameter("period", period))
	return api
}

// Sort sets the Sort query parameter.
func (api *VideosListCall) Sort(sort string) *VideosListCall {
	api.opts = append(api.opts, SetQueryParameter("sort", sort))
	return api
}

// Type sets the Type query parameter.
func (api *VideosListCall) Type(t string) *VideosListCall {
	api.opts = append(api.opts, SetQueryParameter("type", t))
	return api
}

// After sets the After query parameter.
func (api *VideosListCall) After(after string) *VideosListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *VideosListCall) First(first int) *VideosListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *VideosListCall) Do(ctx context.Context, opts ...RequestOption) (*VideosListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/videos", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Video](res)
	if err != nil {
		return nil, err
	}

	return &VideosListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// VideosDeleteCall represents a DELETE call to a Twitch Videos API endpoint.
type VideosDeleteCall struct {
	resource *VideosResource
	opts     []RequestOption
}

// VideosDeleteResponse represents the response from a DELETE request to /helix/videos.
type VideosDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/videos.
//
// Deletes one or more videos.
// You may delete past broadcasts, highlights, or uploads.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:videos scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#delete-videos
func (r *VideosResource) Delete() *VideosDeleteCall {
	return &VideosDeleteCall{resource: r}
}

// ID sets the ID query parameter.
func (api *VideosDeleteCall) ID(iD string) *VideosDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("id", iD))
	return api
}

// Do executes the request.
func (api *VideosDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*VideosDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/videos", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &VideosDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

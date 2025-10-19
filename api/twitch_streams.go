package api

import (
	"context"
	"net/http"
)

// StreamsResource represents the Twitch Streams API.
type StreamsResource struct {
	client *Client
}

// NewStreamsResource creates a new StreamsResource.
func NewStreamsResource(client *Client) *StreamsResource {
	return &StreamsResource{client}
}

// StreamsListCall represents a GET call to a Twitch Streams API endpoint.
type StreamsListCall struct {
	resource *StreamsResource
	opts     []RequestOption
}

// StreamsListResponse represents the response from a GET request to /helix/streams.
type StreamsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Stream data returned by the Twitch API.
	Data []Stream
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/streams.
//
// Gets a list of all streams.
// The list is in descending order by the number of viewers watching the stream.
// Because viewers come and go during a stream, it's possible to find duplicate or missing streams in the list as you page through the results.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-streams
func (r *StreamsResource) List() *StreamsListCall {
	return &StreamsListCall{resource: r}
}

// ID adds to the ID query parameter.
func (api *StreamsListCall) ID(ids ...string) *StreamsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// UserID adds to the UserID query parameter.
func (api *StreamsListCall) UserID(userIDs ...string) *StreamsListCall {
	for _, userID := range userIDs {
		api.opts = append(api.opts, AddQueryParameter("user_id", userID))
	}
	return api
}

// GameID adds to the GameID query parameter.
func (api *StreamsListCall) GameID(gameIDs ...string) *StreamsListCall {
	for _, gameID := range gameIDs {
		api.opts = append(api.opts, AddQueryParameter("game_id", gameID))
	}
	return api
}

// Language sets the Language query parameter.
func (api *StreamsListCall) Language(language string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("language", language))
	return api
}

// Period sets the Period query parameter.
func (api *StreamsListCall) Period(period string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("period", period))
	return api
}

// Type sets the Type query parameter.
func (api *StreamsListCall) Type(t string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("type", t))
	return api
}

// Before sets the Before query parameter.
func (api *StreamsListCall) Before(before string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("before", before))
	return api
}

// After sets the After query parameter.
func (api *StreamsListCall) After(after string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *StreamsListCall) First(first int) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *StreamsListCall) Do(ctx context.Context, opts ...RequestOption) (*StreamsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/streams", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Stream](res)
	if err != nil {
		return nil, err
	}

	return &StreamsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

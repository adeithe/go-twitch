package api

import (
	"context"
	"net/http"
)

// GamesResource represents the Twitch Games API.
type GamesResource struct {
	client *Client

	// Top provides access to the Twitch Top API.
	Top *GamesTopResource
}

// NewGamesResource creates a new GamesResource.
func NewGamesResource(client *Client) *GamesResource {
	r := &GamesResource{client: client}
	r.Top = NewGamesTopResource(client)
	return r
}

// GamesTopResource represents the Twitch GamesTop API.
type GamesTopResource struct {
	client *Client
}

// NewGamesTopResource creates a new GamesTopResource.
func NewGamesTopResource(client *Client) *GamesTopResource {
	return &GamesTopResource{client}
}

// TopGamesListCall represents a GET call to a Twitch GamesTop API endpoint.
type TopGamesListCall struct {
	resource *GamesTopResource
	opts     []RequestOption
}

// TopGamesListResponse represents the response from a GET request to /helix/games/top.
type TopGamesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Game data returned by the Twitch API.
	Data []Game
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/games/top.
//
// Gets information about all broadcasts on Twitch.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-top-games
func (r *GamesTopResource) List() *TopGamesListCall {
	return &TopGamesListCall{resource: r}
}

// Before sets the Before query parameter.
func (api *TopGamesListCall) Before(before string) *TopGamesListCall {
	api.opts = append(api.opts, SetQueryParameter("before", before))
	return api
}

// After sets the After query parameter.
func (api *TopGamesListCall) After(after string) *TopGamesListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *TopGamesListCall) First(first int) *TopGamesListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *TopGamesListCall) Do(ctx context.Context, opts ...RequestOption) (*TopGamesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/games/top", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Game](res)
	if err != nil {
		return nil, err
	}

	return &TopGamesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// GamesListCall represents a GET call to a Twitch Games API endpoint.
type GamesListCall struct {
	resource *GamesResource
	opts     []RequestOption
}

// GamesListResponse represents the response from a GET request to /helix/games.
type GamesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Game data returned by the Twitch API.
	Data []Game
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/games.
//
// Gets information about one or more specified games.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-games
func (r *GamesResource) List(id string, name string) *GamesListCall {
	c := &GamesListCall{resource: r}
	return c.
		ID(id).
		Name(name)
}

// ID adds to the ID query parameter.
func (api *GamesListCall) ID(ids ...string) *GamesListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// Name adds to the Name query parameter.
func (api *GamesListCall) Name(names ...string) *GamesListCall {
	for _, name := range names {
		api.opts = append(api.opts, AddQueryParameter("name", name))
	}
	return api
}

// Do executes the request.
func (api *GamesListCall) Do(ctx context.Context, opts ...RequestOption) (*GamesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/games", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Game](res)
	if err != nil {
		return nil, err
	}

	return &GamesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

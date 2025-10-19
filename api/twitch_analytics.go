package api

import (
	"context"
	"net/http"
	"time"
)

// AnalyticsResource represents the Twitch Analytics API.
type AnalyticsResource struct {
	client *Client

	// Extensions provides access to the Twitch Extensions API.
	Extensions *AnalyticsExtensionsResource
	// Games provides access to the Twitch Games API.
	Games *AnalyticsGamesResource
}

// NewAnalyticsResource creates a new AnalyticsResource.
func NewAnalyticsResource(client *Client) *AnalyticsResource {
	r := &AnalyticsResource{client: client}
	r.Extensions = NewAnalyticsExtensionsResource(client)
	r.Games = NewAnalyticsGamesResource(client)
	return r
}

// AnalyticsExtensionsResource represents the Twitch AnalyticsExtensions API.
type AnalyticsExtensionsResource struct {
	client *Client
}

// NewAnalyticsExtensionsResource creates a new AnalyticsExtensionsResource.
func NewAnalyticsExtensionsResource(client *Client) *AnalyticsExtensionsResource {
	return &AnalyticsExtensionsResource{client}
}

// ExtensionAnalyticsListCall represents a GET call to a Twitch AnalyticsExtensions API endpoint.
type ExtensionAnalyticsListCall struct {
	resource *AnalyticsExtensionsResource
	opts     []RequestOption
}

// ExtensionAnalyticsListResponse represents the response from a GET request to /helix/analytics/extensions.
type ExtensionAnalyticsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ExtensionAnalyticsReport data returned by the Twitch API.
	Data []ExtensionAnalyticsReport
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/analytics/extensions.
//
// Gets an analytics report for one or more extensions. The response contains the URLs used to download the reports (CSV files).
//
// # Authorization
//
// Requires a user access token that includes the analytics:read:extensions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-extension-analytics
func (r *AnalyticsExtensionsResource) List() *ExtensionAnalyticsListCall {
	return &ExtensionAnalyticsListCall{resource: r}
}

// ExtensionID sets the ExtensionID query parameter.
func (api *ExtensionAnalyticsListCall) ExtensionID(extensionID string) *ExtensionAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("extension_id", extensionID))
	return api
}

// Type sets the Type query parameter.
func (api *ExtensionAnalyticsListCall) Type(t string) *ExtensionAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("type", t))
	return api
}

// After sets the After query parameter.
func (api *ExtensionAnalyticsListCall) After(after string) *ExtensionAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ExtensionAnalyticsListCall) First(first int) *ExtensionAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// StartedAt sets the StartedAt query parameter.
func (api *ExtensionAnalyticsListCall) StartedAt(startedAt time.Time) *ExtensionAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("started_at", startedAt.Format(time.RFC3339)))
	return api
}

// EndedAt sets the EndedAt query parameter.
func (api *ExtensionAnalyticsListCall) EndedAt(endedAt time.Time) *ExtensionAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("ended_at", endedAt.Format(time.RFC3339)))
	return api
}

// Do executes the request.
func (api *ExtensionAnalyticsListCall) Do(ctx context.Context, opts ...RequestOption) (*ExtensionAnalyticsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/analytics/extensions", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ExtensionAnalyticsReport](res)
	if err != nil {
		return nil, err
	}

	return &ExtensionAnalyticsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// AnalyticsGamesResource represents the Twitch AnalyticsGames API.
type AnalyticsGamesResource struct {
	client *Client
}

// NewAnalyticsGamesResource creates a new AnalyticsGamesResource.
func NewAnalyticsGamesResource(client *Client) *AnalyticsGamesResource {
	return &AnalyticsGamesResource{client}
}

// GameAnalyticsListCall represents a GET call to a Twitch AnalyticsGames API endpoint.
type GameAnalyticsListCall struct {
	resource *AnalyticsGamesResource
	opts     []RequestOption
}

// GameAnalyticsListResponse represents the response from a GET request to /helix/analytics/games.
type GameAnalyticsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the GameAnalyticsReport data returned by the Twitch API.
	Data []GameAnalyticsReport
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/analytics/games.
//
// Gets an analytics report for one or more games. The response contains the URLs used to download the reports (CSV files).
//
// # Authorization
//
// Requires a user access token that includes the analytics:read:games scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-game-analytics
func (r *AnalyticsGamesResource) List() *GameAnalyticsListCall {
	return &GameAnalyticsListCall{resource: r}
}

// GameID sets the GameID query parameter.
func (api *GameAnalyticsListCall) GameID(gameID string) *GameAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("game_id", gameID))
	return api
}

// Type sets the Type query parameter.
func (api *GameAnalyticsListCall) Type(t string) *GameAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("type", t))
	return api
}

// After sets the After query parameter.
func (api *GameAnalyticsListCall) After(after string) *GameAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *GameAnalyticsListCall) First(first int) *GameAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// StartedAt sets the StartedAt query parameter.
func (api *GameAnalyticsListCall) StartedAt(startedAt time.Time) *GameAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("started_at", startedAt.Format(time.RFC3339)))
	return api
}

// EndedAt sets the EndedAt query parameter.
func (api *GameAnalyticsListCall) EndedAt(endedAt time.Time) *GameAnalyticsListCall {
	api.opts = append(api.opts, SetQueryParameter("ended_at", endedAt.Format(time.RFC3339)))
	return api
}

// Do executes the request.
func (api *GameAnalyticsListCall) Do(ctx context.Context, opts ...RequestOption) (*GameAnalyticsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/analytics/games", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[GameAnalyticsReport](res)
	if err != nil {
		return nil, err
	}

	return &GameAnalyticsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

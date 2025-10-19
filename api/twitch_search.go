package api

import (
	"context"
	"net/http"
)

// SearchResource represents the Twitch Search API.
type SearchResource struct {
	client *Client

	// Categories provides access to the Twitch Categories API.
	Categories *SearchCategoriesResource
	// Channels provides access to the Twitch Channels API.
	Channels *SearchChannelsResource
}

// NewSearchResource creates a new SearchResource.
func NewSearchResource(client *Client) *SearchResource {
	r := &SearchResource{client: client}
	r.Categories = NewSearchCategoriesResource(client)
	r.Channels = NewSearchChannelsResource(client)
	return r
}

// SearchCategoriesResource represents the Twitch SearchCategories API.
type SearchCategoriesResource struct {
	client *Client
}

// NewSearchCategoriesResource creates a new SearchCategoriesResource.
func NewSearchCategoriesResource(client *Client) *SearchCategoriesResource {
	return &SearchCategoriesResource{client}
}

// SearchCategoriesListCall represents a GET call to a Twitch SearchCategories API endpoint.
type SearchCategoriesListCall struct {
	resource *SearchCategoriesResource
	opts     []RequestOption
}

// SearchCategoriesListResponse represents the response from a GET request to /helix/search/categories.
type SearchCategoriesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CategorySearchResult data returned by the Twitch API.
	Data []CategorySearchResult
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/search/categories.
//
// Gets the games or categories that match the specified query.
//
// To match, the category's name must contain all parts of the query string.
// For example, if the query string is 42, the response includes any category name that contains 42 in the title.
// If the query string is a phrase like love computer, the response includes any category name that contains the words love and computer anywhere in the name.
// The comparison is case insensitive.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#search-categories
func (r *SearchCategoriesResource) List(query string) *SearchCategoriesListCall {
	c := &SearchCategoriesListCall{resource: r}
	return c.
		Query(query)
}

// Query sets the Query query parameter.
func (api *SearchCategoriesListCall) Query(query string) *SearchCategoriesListCall {
	api.opts = append(api.opts, SetQueryParameter("query", query))
	return api
}

// After sets the After query parameter.
func (api *SearchCategoriesListCall) After(after string) *SearchCategoriesListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *SearchCategoriesListCall) First(first int) *SearchCategoriesListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *SearchCategoriesListCall) Do(ctx context.Context, opts ...RequestOption) (*SearchCategoriesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/search/categories", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CategorySearchResult](res)
	if err != nil {
		return nil, err
	}

	return &SearchCategoriesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// SearchChannelsResource represents the Twitch SearchChannels API.
type SearchChannelsResource struct {
	client *Client
}

// NewSearchChannelsResource creates a new SearchChannelsResource.
func NewSearchChannelsResource(client *Client) *SearchChannelsResource {
	return &SearchChannelsResource{client}
}

// SearchChannelsListCall represents a GET call to a Twitch SearchChannels API endpoint.
type SearchChannelsListCall struct {
	resource *SearchChannelsResource
	opts     []RequestOption
}

// SearchChannelsListResponse represents the response from a GET request to /helix/search/channels.
type SearchChannelsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CategorySearchResult data returned by the Twitch API.
	Data []CategorySearchResult
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/search/channels.
//
// Gets the channels that match the specified query and have streamed content within the past 6 months.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#search-channels
func (r *SearchChannelsResource) List(query string) *SearchChannelsListCall {
	c := &SearchChannelsListCall{resource: r}
	return c.
		Query(query)
}

// Query sets the Query query parameter.
func (api *SearchChannelsListCall) Query(query string) *SearchChannelsListCall {
	api.opts = append(api.opts, SetQueryParameter("query", query))
	return api
}

// After sets the After query parameter.
func (api *SearchChannelsListCall) After(after string) *SearchChannelsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *SearchChannelsListCall) First(first int) *SearchChannelsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// LiveOnly sets the LiveOnly query parameter.
func (api *SearchChannelsListCall) LiveOnly(liveOnly bool) *SearchChannelsListCall {
	api.opts = append(api.opts, SetQueryParameter("live_only", liveOnly))
	return api
}

// Do executes the request.
func (api *SearchChannelsListCall) Do(ctx context.Context, opts ...RequestOption) (*SearchChannelsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/search/channels", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CategorySearchResult](res)
	if err != nil {
		return nil, err
	}

	return &SearchChannelsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

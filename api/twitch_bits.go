package api

import (
	"context"
	"net/http"
	"time"
)

// BitsResource represents the Twitch Bits API.
type BitsResource struct {
	client *Client

	// Cheermotes provides access to the Twitch Cheermotes API.
	Cheermotes *BitsCheermotesResource
	// Extensions provides access to the Twitch Extensions API.
	Extensions *BitsExtensionsResource
	// Leaderboard provides access to the Twitch Leaderboard API.
	Leaderboard *BitsLeaderboardResource
}

// NewBitsResource creates a new BitsResource.
func NewBitsResource(client *Client) *BitsResource {
	r := &BitsResource{client: client}
	r.Cheermotes = NewBitsCheermotesResource(client)
	r.Extensions = NewBitsExtensionsResource(client)
	r.Leaderboard = NewBitsLeaderboardResource(client)
	return r
}

// BitsCheermotesResource represents the Twitch BitsCheermotes API.
type BitsCheermotesResource struct {
	client *Client
}

// NewBitsCheermotesResource creates a new BitsCheermotesResource.
func NewBitsCheermotesResource(client *Client) *BitsCheermotesResource {
	return &BitsCheermotesResource{client}
}

// BitsCheermotesListCall represents a GET call to a Twitch BitsCheermotes API endpoint.
type BitsCheermotesListCall struct {
	resource *BitsCheermotesResource
	opts     []RequestOption
}

// BitsCheermotesListResponse represents the response from a GET request to /helix/bits/cheermotes.
type BitsCheermotesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Cheermote data returned by the Twitch API.
	Data []Cheermote
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/bits/cheermotes.
//
// Gets a list of Cheermotes that users can use to cheer Bits in any Bits-enabled chat room.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-cheermotes
func (r *BitsCheermotesResource) List() *BitsCheermotesListCall {
	return &BitsCheermotesListCall{resource: r}
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *BitsCheermotesListCall) BroadcasterID(broadcasterID string) *BitsCheermotesListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *BitsCheermotesListCall) Do(ctx context.Context, opts ...RequestOption) (*BitsCheermotesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/bits/cheermotes", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Cheermote](res)
	if err != nil {
		return nil, err
	}

	return &BitsCheermotesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// BitsExtensionsResource represents the Twitch BitsExtensions API.
type BitsExtensionsResource struct {
	client *Client
}

// NewBitsExtensionsResource creates a new BitsExtensionsResource.
func NewBitsExtensionsResource(client *Client) *BitsExtensionsResource {
	return &BitsExtensionsResource{client}
}

// BitsExtensionTransactionsListCall represents a GET call to a Twitch BitsExtensions API endpoint.
type BitsExtensionTransactionsListCall struct {
	resource *BitsExtensionsResource
	opts     []RequestOption
}

// BitsExtensionTransactionsListResponse represents the response from a GET request to /helix/extensions/transactions.
type BitsExtensionTransactionsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ExtensionTransaction data returned by the Twitch API.
	Data []ExtensionTransaction
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/extensions/transactions.
//
// Gets a list of transactions for an extension. A transaction records the exchange of a currency (for example, Bits) for a digital product.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-extension-transactions
func (r *BitsExtensionsResource) List(extensionID string) *BitsExtensionTransactionsListCall {
	c := &BitsExtensionTransactionsListCall{resource: r}
	return c.
		ExtensionID(extensionID)
}

// ID adds to the ID query parameter.
func (api *BitsExtensionTransactionsListCall) ID(ids ...string) *BitsExtensionTransactionsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// ExtensionID sets the ExtensionID query parameter.
func (api *BitsExtensionTransactionsListCall) ExtensionID(extensionID string) *BitsExtensionTransactionsListCall {
	api.opts = append(api.opts, SetQueryParameter("extension_id", extensionID))
	return api
}

// After sets the After query parameter.
func (api *BitsExtensionTransactionsListCall) After(after string) *BitsExtensionTransactionsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *BitsExtensionTransactionsListCall) First(first int) *BitsExtensionTransactionsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *BitsExtensionTransactionsListCall) Do(ctx context.Context, opts ...RequestOption) (*BitsExtensionTransactionsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/extensions/transactions", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ExtensionTransaction](res)
	if err != nil {
		return nil, err
	}

	return &BitsExtensionTransactionsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// BitsLeaderboardResource represents the Twitch BitsLeaderboard API.
type BitsLeaderboardResource struct {
	client *Client
}

// NewBitsLeaderboardResource creates a new BitsLeaderboardResource.
func NewBitsLeaderboardResource(client *Client) *BitsLeaderboardResource {
	return &BitsLeaderboardResource{client}
}

// BitsLeaderboardListCall represents a GET call to a Twitch BitsLeaderboard API endpoint.
type BitsLeaderboardListCall struct {
	resource *BitsLeaderboardResource
	opts     []RequestOption
}

// BitsLeaderboardListResponse represents the response from a GET request to /helix/bits/leaderboard.
type BitsLeaderboardListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Total is the int data returned by the Twitch API.
	Total int
	// Data is the BitsLeaderboardEntry data returned by the Twitch API.
	Data []BitsLeaderboardEntry
	// DateRange is the DateRange data returned by the Twitch API.
	DateRange DateRange
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/bits/leaderboard.
//
// Gets the Bits leaderboard for the authenticated broadcaster.
//
// # Authorization
//
// Requires a user access token that includes the bits:read scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-bits-leaderboard
func (r *BitsLeaderboardResource) List() *BitsLeaderboardListCall {
	return &BitsLeaderboardListCall{resource: r}
}

// UserID sets the UserID query parameter.
func (api *BitsLeaderboardListCall) UserID(userID string) *BitsLeaderboardListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// Period sets the Period query parameter.
func (api *BitsLeaderboardListCall) Period(period string) *BitsLeaderboardListCall {
	api.opts = append(api.opts, SetQueryParameter("period", period))
	return api
}

// Count sets the Count query parameter.
func (api *BitsLeaderboardListCall) Count(count int) *BitsLeaderboardListCall {
	api.opts = append(api.opts, SetQueryParameter("count", count))
	return api
}

// StartedAt sets the StartedAt query parameter.
func (api *BitsLeaderboardListCall) StartedAt(startedAt time.Time) *BitsLeaderboardListCall {
	api.opts = append(api.opts, SetQueryParameter("started_at", startedAt.Format(time.RFC3339)))
	return api
}

// Do executes the request.
func (api *BitsLeaderboardListCall) Do(ctx context.Context, opts ...RequestOption) (*BitsLeaderboardListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/bits/leaderboard", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[BitsLeaderboardEntry](res)
	if err != nil {
		return nil, err
	}

	return &BitsLeaderboardListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Total:      data.Total,
		Data:       data.Data,
		DateRange:  data.DateRange,
		Request:    res.Request,
	}, nil
}

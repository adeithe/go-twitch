package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// EntitlementsResource represents the Twitch Entitlements API.
type EntitlementsResource struct {
	client *Client

	// Drops provides access to the Twitch Drops API.
	Drops *EntitlementsDropsResource
}

// NewEntitlementsResource creates a new EntitlementsResource.
func NewEntitlementsResource(client *Client) *EntitlementsResource {
	r := &EntitlementsResource{client: client}
	r.Drops = NewEntitlementsDropsResource(client)
	return r
}

// EntitlementsDropsResource represents the Twitch EntitlementsDrops API.
type EntitlementsDropsResource struct {
	client *Client
}

// NewEntitlementsDropsResource creates a new EntitlementsDropsResource.
func NewEntitlementsDropsResource(client *Client) *EntitlementsDropsResource {
	return &EntitlementsDropsResource{client}
}

// DropsEntitlementsListCall represents a GET call to a Twitch EntitlementsDrops API endpoint.
type DropsEntitlementsListCall struct {
	resource *EntitlementsDropsResource
	opts     []RequestOption
}

// DropsEntitlementsListResponse represents the response from a GET request to /helix/entitlements/drops.
type DropsEntitlementsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the DropEntitlement data returned by the Twitch API.
	Data []DropEntitlement
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/entitlements/drops.
//
// Gets an organization's list of entitlements that have been granted to a game, a user, or both.
//
// Entitlements returned in the response body data are not guaranteed to be sorted by any field returned by the API.
// To retrieve CLAIMED or FULFILLED entitlements, use the FulfillmentStatus method to filter results.
// To retrieve entitlements for a specific game, use the GameID method to filter results.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// The associated Client ID for the access token must be owned by a user who is a member of the organization that holds ownership of the game.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-drops-entitlements
func (r *EntitlementsDropsResource) List() *DropsEntitlementsListCall {
	return &DropsEntitlementsListCall{resource: r}
}

// ID adds to the ID query parameter.
func (api *DropsEntitlementsListCall) ID(ids ...string) *DropsEntitlementsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// UserID sets the UserID query parameter.
func (api *DropsEntitlementsListCall) UserID(userID string) *DropsEntitlementsListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// GameID sets the GameID query parameter.
func (api *DropsEntitlementsListCall) GameID(gameID string) *DropsEntitlementsListCall {
	api.opts = append(api.opts, SetQueryParameter("game_id", gameID))
	return api
}

// FulfillmentStatus sets the FulfillmentStatus query parameter.
func (api *DropsEntitlementsListCall) FulfillmentStatus(fulfillmentStatus string) *DropsEntitlementsListCall {
	api.opts = append(api.opts, SetQueryParameter("fulfillment_status", fulfillmentStatus))
	return api
}

// After sets the After query parameter.
func (api *DropsEntitlementsListCall) After(after string) *DropsEntitlementsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *DropsEntitlementsListCall) First(first int) *DropsEntitlementsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *DropsEntitlementsListCall) Do(ctx context.Context, opts ...RequestOption) (*DropsEntitlementsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/entitlements/drops", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[DropEntitlement](res)
	if err != nil {
		return nil, err
	}

	return &DropsEntitlementsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// DropsEntitlementsModifyCall represents a PATCH call to a Twitch EntitlementsDrops API endpoint.
type DropsEntitlementsModifyCall struct {
	resource *EntitlementsDropsResource
	body     map[string]any
}

// DropsEntitlementsModifyResponse represents the response from a PATCH request to /helix/entitlements/drops.
type DropsEntitlementsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the UpdatedDropEntitlement data returned by the Twitch API.
	Data []UpdatedDropEntitlement
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/entitlements/drops.
//
// Updates the Drop entitlement's fulfillment status.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// The associated Client ID for the access token must be owned by a user who is a member of the organization that holds ownership of the game.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-drops-entitlements
func (r *EntitlementsDropsResource) Modify() *DropsEntitlementsModifyCall {
	return &DropsEntitlementsModifyCall{resource: r, body: make(map[string]any)}
}

// EntitlementID sets the EntitlementID body parameter.
func (api *DropsEntitlementsModifyCall) EntitlementID(entitlementIDs ...string) *DropsEntitlementsModifyCall {
	api.body["entitlement_ids"] = entitlementIDs
	return api
}

// FulfillmentStatus sets the FulfillmentStatus body parameter.
func (api *DropsEntitlementsModifyCall) FulfillmentStatus(fulfillmentStatus string) *DropsEntitlementsModifyCall {
	api.body["fulfillment_status"] = fulfillmentStatus
	return api
}

// Do executes the request.
func (api *DropsEntitlementsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*DropsEntitlementsModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/entitlements/drops", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[UpdatedDropEntitlement](res)
	if err != nil {
		return nil, err
	}

	return &DropsEntitlementsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

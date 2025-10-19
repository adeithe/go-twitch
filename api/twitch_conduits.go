package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ConduitsResource represents the Twitch Conduits API.
type ConduitsResource struct {
	client *Client

	// Shards provides access to the Twitch Shards API.
	Shards *ConduitsShardsResource
}

// NewConduitsResource creates a new ConduitsResource.
func NewConduitsResource(client *Client) *ConduitsResource {
	r := &ConduitsResource{client: client}
	r.Shards = NewConduitsShardsResource(client)
	return r
}

// ConduitsShardsResource represents the Twitch ConduitsShards API.
type ConduitsShardsResource struct {
	client *Client
}

// NewConduitsShardsResource creates a new ConduitsShardsResource.
func NewConduitsShardsResource(client *Client) *ConduitsShardsResource {
	return &ConduitsShardsResource{client}
}

// ConduitsShardListCall represents a GET call to a Twitch ConduitsShards API endpoint.
type ConduitsShardListCall struct {
	resource *ConduitsShardsResource
	opts     []RequestOption
}

// ConduitsShardListResponse represents the response from a GET request to /helix/eventsub/conduits/shards.
type ConduitsShardListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ConduitShard data returned by the Twitch API.
	Data []ConduitShard
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/eventsub/conduits/shards.
//
// Gets a lists of all shards for a conduit.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-conduits-shards
func (r *ConduitsShardsResource) List(conduitID string) *ConduitsShardListCall {
	c := &ConduitsShardListCall{resource: r}
	return c.
		ConduitID(conduitID)
}

// ConduitID sets the ConduitID query parameter.
func (api *ConduitsShardListCall) ConduitID(conduitID string) *ConduitsShardListCall {
	api.opts = append(api.opts, SetQueryParameter("conduit_id", conduitID))
	return api
}

// Status sets the Status query parameter.
func (api *ConduitsShardListCall) Status(status string) *ConduitsShardListCall {
	api.opts = append(api.opts, SetQueryParameter("status", status))
	return api
}

// After sets the After query parameter.
func (api *ConduitsShardListCall) After(after string) *ConduitsShardListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// Do executes the request.
func (api *ConduitsShardListCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsShardListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/eventsub/conduits/shards", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ConduitShard](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsShardListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// ConduitsShardModifyCall represents a PATCH call to a Twitch ConduitsShards API endpoint.
type ConduitsShardModifyCall struct {
	resource *ConduitsShardsResource
	body     map[string]any
}

// ConduitsShardModifyResponse represents the response from a PATCH request to /helix/eventsub/conduits/shards.
type ConduitsShardModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ConduitShard data returned by the Twitch API.
	Data []ConduitShard
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/eventsub/conduits/shards.
//
// Updates a conduit shard.
//
// Shard IDs are indexed starting at 0, so a conduit with a shard_count of 5 will have shards with IDs 0 through 4.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-conduit-shards
func (r *ConduitsShardsResource) Modify(conduitID string, shards ConduitShard) *ConduitsShardModifyCall {
	c := &ConduitsShardModifyCall{resource: r, body: make(map[string]any)}
	return c.
		ConduitID(conduitID).
		Shards(shards)
}

// ConduitID sets the ConduitID body parameter.
func (api *ConduitsShardModifyCall) ConduitID(conduitID string) *ConduitsShardModifyCall {
	api.body["conduit_id"] = conduitID
	return api
}

// Shards sets the Shards body parameter.
func (api *ConduitsShardModifyCall) Shards(shardss ...ConduitShard) *ConduitsShardModifyCall {
	api.body["shards"] = shardss
	return api
}

// Do executes the request.
func (api *ConduitsShardModifyCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsShardModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/eventsub/conduits/shards", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ConduitShard](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsShardModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ConduitsListCall represents a GET call to a Twitch Conduits API endpoint.
type ConduitsListCall struct {
	resource *ConduitsResource
}

// ConduitsListResponse represents the response from a GET request to /helix/eventsub/conduits.
type ConduitsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Conduit data returned by the Twitch API.
	Data []Conduit
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/eventsub/conduits.
//
// Gets all conduits for a client ID.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-conduits
func (r *ConduitsResource) List() *ConduitsListCall {
	return &ConduitsListCall{resource: r}
}

// Do executes the request.
func (api *ConduitsListCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/eventsub/conduits", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Conduit](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ConduitsInsertCall represents a POST call to a Twitch Conduits API endpoint.
type ConduitsInsertCall struct {
	resource *ConduitsResource
	body     map[string]any
}

// ConduitsInsertResponse represents the response from a POST request to /helix/eventsub/conduits.
type ConduitsInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Conduit data returned by the Twitch API.
	Data []Conduit
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/eventsub/conduits.
//
// Creates a new conduit.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-conduits
func (r *ConduitsResource) Insert(shardCount int) *ConduitsInsertCall {
	c := &ConduitsInsertCall{resource: r, body: make(map[string]any)}
	return c.
		ShardCount(shardCount)
}

// ShardCount sets the ShardCount body parameter.
func (api *ConduitsInsertCall) ShardCount(shardCount int) *ConduitsInsertCall {
	api.body["shard_count"] = shardCount
	return api
}

// Do executes the request.
func (api *ConduitsInsertCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/eventsub/conduits", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Conduit](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ConduitsModifyCall represents a PATCH call to a Twitch Conduits API endpoint.
type ConduitsModifyCall struct {
	resource *ConduitsResource
	body     map[string]any
}

// ConduitsModifyResponse represents the response from a PATCH request to /helix/eventsub/conduits.
type ConduitsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Conduit data returned by the Twitch API.
	Data []Conduit
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/eventsub/conduits.
//
// Updates a conduit's shard count.
// To delete shards, update the count to a lower number, and the shards above the count will be deleted.
// For example, if the existing shard count is 100, by resetting shard count to 50, shards 50-99 are disabled.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-conduits
func (r *ConduitsResource) Modify(id string, shardCount int) *ConduitsModifyCall {
	c := &ConduitsModifyCall{resource: r, body: make(map[string]any)}
	return c.
		ID(id).
		ShardCount(shardCount)
}

// ID sets the ID body parameter.
func (api *ConduitsModifyCall) ID(id string) *ConduitsModifyCall {
	api.body["id"] = id
	return api
}

// ShardCount sets the ShardCount body parameter.
func (api *ConduitsModifyCall) ShardCount(shardCount int) *ConduitsModifyCall {
	api.body["shard_count"] = shardCount
	return api
}

// Do executes the request.
func (api *ConduitsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/eventsub/conduits", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Conduit](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ConduitsDeleteCall represents a DELETE call to a Twitch Conduits API endpoint.
type ConduitsDeleteCall struct {
	resource *ConduitsResource
	opts     []RequestOption
}

// ConduitsDeleteResponse represents the response from a DELETE request to /helix/eventsub/conduits.
type ConduitsDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/eventsub/conduits.
//
// Deletes a conduit.
// Note that it may take some time for Eventsub subscriptions on a deleted conduit to show as disabled when calling Get Eventsub Subscriptions.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#delete-conduit
func (r *ConduitsResource) Delete(id string) *ConduitsDeleteCall {
	c := &ConduitsDeleteCall{resource: r}
	return c.
		ID(id)
}

// ID sets the ID query parameter.
func (api *ConduitsDeleteCall) ID(id string) *ConduitsDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// Do executes the request.
func (api *ConduitsDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/eventsub/conduits", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

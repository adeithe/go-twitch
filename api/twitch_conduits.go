package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// Conduit represents a Twitch Eventsub Conduit.
type Conduit struct {
	ID         string `json:"id"`
	ShardCount int    `json:"shard_count"`
}

// ConduitsResponse represents the response from the Twitch Eventsub Conduits API.
type ConduitsResponse struct {
	Header   http.Header
	Conduits []Conduit
}

// ConduitsResource is the API resource for managing Twitch Eventsub Conduits.
type ConduitsResource struct {
	client *Client
	Shards *ConduitShardsResource
}

// NewConduitsResource creates a new ConduitsResource.
func NewConduitsResource(client *Client) *ConduitsResource {
	c := &ConduitsResource{client: client}
	c.Shards = NewConduitShardsResource(client)
	return c
}

// ConduitsListCall is the API call for listing an apps Twitch Eventsub Conduits.
type ConduitsListCall struct {
	resource *ConduitsResource
}

// EndpointConduits is the endpoint for the Twitch Eventsub Conduits API.
const EndpointConduits = TwitchAPIVersionHelix + "/eventsub/conduits"

// List creates a new ConduitsListCall.
func (r *ConduitsResource) List() *ConduitsListCall {
	return &ConduitsListCall{resource: r}
}

// Do executes the request.
func (c *ConduitsListCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, EndpointConduits, nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Conduit](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsResponse{
		Header:   res.Header,
		Conduits: data.Data,
	}, nil
}

// ConduitInsertCall is the API call for creating a new Twitch Eventsub Conduit.
type ConduitInsertCall struct {
	resource *ConduitsResource
	body     map[string]any
}

// Insert creates a new ConduitInsertCall.
func (r *ConduitsResource) Insert() *ConduitInsertCall {
	return &ConduitInsertCall{resource: r, body: make(map[string]any)}
}

// ShardCount sets the shard count for the new conduit.
func (c *ConduitInsertCall) ShardCount(n int) *ConduitInsertCall {
	c.body["shard_count"] = n
	return c
}

// Do executes the request.
func (c *ConduitInsertCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsResponse, error) {
	bs, err := json.Marshal(c.body)
	if err != nil {
		return nil, err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPost, EndpointConduits, bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Conduit](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsResponse{
		Header:   res.Header,
		Conduits: data.Data,
	}, nil
}

// ConduitUpdateCall is the API call for updating a Twitch Eventsub Conduit.
type ConduitUpdateCall struct {
	resource *ConduitsResource
	body     map[string]any
}

// Update creates a new ConduitUpdateCall.
func (r *ConduitsResource) Update(id string) *ConduitUpdateCall {
	return &ConduitUpdateCall{resource: r, body: map[string]any{"id": id}}
}

// ShardCount sets the shard count for the conduit.
func (c *ConduitUpdateCall) ShardCount(n int) *ConduitUpdateCall {
	c.body["shard_count"] = n
	return c
}

// Do executes the request.
func (c *ConduitUpdateCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitsResponse, error) {
	bs, err := json.Marshal(c.body)
	if err != nil {
		return nil, err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPatch, EndpointConduits, bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Conduit](res)
	if err != nil {
		return nil, err
	}

	return &ConduitsResponse{
		Header:   res.Header,
		Conduits: data.Data,
	}, nil
}

// ConduitDeleteCall is the API call for deleting a Twitch Eventsub Conduit.
type ConduitDeleteCall struct {
	resource *ConduitsResource
	opts     []RequestOption
}

// Delete creates a new ConduitDeleteCall.
func (r *ConduitsResource) Delete(id string) *ConduitDeleteCall {
	return &ConduitDeleteCall{
		resource: r,
		opts: []RequestOption{
			SetQueryParameter("id", id),
		},
	}
}

// Do executes the request.
func (c *ConduitDeleteCall) Do(ctx context.Context, opts ...RequestOption) error {
	res, err := c.resource.client.doRequest(ctx, http.MethodDelete, EndpointConduits, nil, append(c.opts, opts...)...)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if _, err := decodeResponse[any](res); err != nil {
		return err
	}
	return nil
}

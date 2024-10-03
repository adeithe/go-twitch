package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// ConduitShard is a shard for a Twitch Eventsub Conduit.
type ConduitShard struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Transport Transport `json:"transport"`
}

// ConduitShardsResponse represents the response from the Twitch Eventsub Conduit Shards API.
type ConduitShardsResponse struct {
	Header http.Header
	Shards []ConduitShard
	Cursor string
}

// Transport is the transport method for a Twitch Eventsub Conduit Shard.
type Transport struct {
	Method         string     `json:"method"`
	Callback       *string    `json:"callback,omitempty"`
	Secret         *string    `json:"secret,omitempty"`
	SessionID      *string    `json:"session_id,omitempty"`
	ConnectedAt    *time.Time `json:"connected_at,omitempty"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}

// NewWebhookTransport creates a new webhook transport for a Twitch Eventsub Conduit Shard.
func NewWebhookTransport(callback, secret string) Transport {
	return Transport{
		Method:   "webhook",
		Callback: &callback,
		Secret:   &secret,
	}
}

// NewWebSocketTransport creates a new websocket transport for a Twitch Eventsub Conduit Shard.
func NewWebSocketTransport(sessionID string) Transport {
	return Transport{
		Method:    "websocket",
		SessionID: &sessionID,
	}
}

// ConduitShardsResource is the API resource for managing Twitch Eventsub Conduit Shards.
type ConduitShardsResource struct {
	client *Client
}

// NewConduitShardsResource creates a new ConduitShardsResource.
func NewConduitShardsResource(client *Client) *ConduitShardsResource {
	return &ConduitShardsResource{client}
}

// ConduitShardListCall is the API call for listing Twitch Eventsub Conduit Shards.
type ConduitShardListCall struct {
	resource *ConduitShardsResource
	opts     []RequestOption
}

// List creates a new ConduitShardListCall.
func (r *ConduitShardsResource) List(conduitID string) *ConduitShardListCall {
	return &ConduitShardListCall{
		resource: r,
		opts: []RequestOption{
			SetQueryParameter("conduit_id", conduitID),
		},
	}
}

// Status filters the list of shards by status.
func (c *ConduitShardListCall) Status(status string) *ConduitShardListCall {
	c.opts = append(c.opts, SetQueryParameter("status", status))
	return c
}

// After sets the cursor for pagination.
func (c *ConduitShardListCall) After(cursor string) *ConduitShardListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Do executes the request.
func (c *ConduitShardListCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitShardsResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/eventsub/conduits/shards", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[ConduitShard](res)
	if err != nil {
		return nil, err
	}

	return &ConduitShardsResponse{
		Header: res.Header,
		Shards: data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

// ConduitShardUpdateCall is the API call for updating Twitch Eventsub Conduit Shards.
type ConduitShardUpdateCall struct {
	resource  *ConduitShardsResource
	shards    []ConduitShard
	conduitID string
}

// Update creates a new ConduitShardUpdateCall.
func (r *ConduitShardsResource) Update(conduitID string) *ConduitShardUpdateCall {
	return &ConduitShardUpdateCall{
		resource:  r,
		conduitID: conduitID,
	}
}

// Shard adds a shard to the update call.
func (c *ConduitShardUpdateCall) Shard(shardID string, transport Transport) *ConduitShardUpdateCall {
	c.shards = append(c.shards, ConduitShard{
		ID:        shardID,
		Transport: transport,
	})
	return c
}

// Do executes the request.
func (c *ConduitShardUpdateCall) Do(ctx context.Context, opts ...RequestOption) (*ConduitShardsResponse, error) {
	bs, err := json.Marshal(map[string]any{
		"conduit_id": c.conduitID,
		"shards":     c.shards,
	})
	if err != nil {
		return nil, err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPatch, "/eventsub/conduits/shards", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[ConduitShard](res)
	if err != nil {
		return nil, err
	}

	return &ConduitShardsResponse{
		Header: res.Header,
		Shards: data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

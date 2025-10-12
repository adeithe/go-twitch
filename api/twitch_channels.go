package api

import (
	"context"
	"net/http"
)

// ChannelsResource handles channel related API calls.
type ChannelsResource struct {
	client *Client
}

// NewChannelsResource creates a new ChannelsResource.
func NewChannelsResource(client *Client) *ChannelsResource {
	return &ChannelsResource{client}
}

// ChannelsListCall is a call to the channels list endpoint.
type ChannelsListCall struct {
	resource *ChannelsResource
	opts     []RequestOption
}

// ChannelsListResponse is the response from the channels list endpoint.
type ChannelsListResponse struct {
	Header http.Header
	Data   []Channel
}

// List creates a request to list channels based on the specified criteria.
func (r *ChannelsResource) List() *ChannelsListCall {
	return &ChannelsListCall{resource: r}
}

// BroadcasterID filters the results to the specified broadcaster ID.
func (c *ChannelsListCall) BroadcasterID(ids []string) *ChannelsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("broadcaster_id", id))
	}
	return c
}

// Do executes the request.
func (c *ChannelsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelsListResponse, error) {
	res, err := c.resource.client.DoRequest(ctx, http.MethodGet, EndpointChannelsGetInformation, nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Channel](res)
	if err != nil {
		return nil, err
	}

	return &ChannelsListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

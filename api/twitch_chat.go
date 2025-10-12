package api

import (
	"context"
	"fmt"
	"net/http"
)

// ChatResource handles chat related API calls.
type ChatResource struct {
	client *Client

	Chatters *ChattersResource
}

// NewChatResource creates a new ChatResource.
func NewChatResource(client *Client) *ChatResource {
	r := &ChatResource{client: client}
	r.Chatters = NewChattersResource(client)
	return r
}

// ChattersResource handles chatters related API calls.
type ChattersResource struct {
	client *Client
}

// NewChattersResource creates a new ChattersResource.
func NewChattersResource(client *Client) *ChattersResource {
	return &ChattersResource{client: client}
}

// ChattersResponse is the response from the chatters endpoint.
type ChattersResponse struct {
	Total    int
	Header   http.Header
	Chatters []UserInfo
	Cursor   string
}

// ChattersListCall is a call to the chatters list endpoint.
type ChattersListCall struct {
	resource *ChattersResource
	opts     []RequestOption
}

// List creates a new call to list chatters.
func (r *ChattersResource) List(broadcasterID, moderatorID string) *ChattersListCall {
	return &ChattersListCall{
		resource: r,
		opts: []RequestOption{
			SetQueryParameter("broadcaster_id", broadcasterID),
			SetQueryParameter("moderator_id", moderatorID),
		},
	}
}

// First filters the results to the first n chatters.
func (c *ChattersListCall) First(n int) *ChattersListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// After filters the results to those after the specified cursor.
func (c *ChattersListCall) After(cursor string) *ChattersListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Do executes the request.
func (c *ChattersListCall) Do(ctx context.Context, opts ...RequestOption) (*ChattersResponse, error) {
	res, err := c.resource.client.DoRequest(ctx, http.MethodGet, EndpointChatGetChatters, nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[UserInfo](res)
	if err != nil {
		return nil, err
	}

	return &ChattersResponse{
		Total:    data.Total,
		Header:   res.Header,
		Chatters: data.Data,
		Cursor:   data.Pagination.Cursor,
	}, nil
}

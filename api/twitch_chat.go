package api

import (
	"context"
	"fmt"
	"net/http"
)

type Chatter struct {
	ID          string `json:"user_id"`
	Username    string `json:"user_login"`
	DisplayName string `json:"user_name"`
}

type ChatResource struct {
	client *Client

	Chatters *ChattersResource
}

func NewChatResource(client *Client) *ChatResource {
	r := &ChatResource{client: client}
	r.Chatters = NewChattersResource(client)
	return r
}

type ChattersResource struct {
	client *Client
}

func NewChattersResource(client *Client) *ChattersResource {
	return &ChattersResource{client: client}
}

type ChattersResponse struct {
	Total    int
	Header   http.Header
	Chatters []Chatter
	Cursor   string
}

type ChattersListCall struct {
	resource *ChatResource
	opts     []RequestOption
}

// List creates a new call to list chatters.
func (r *ChatResource) List(broadcasterId, moderatorId string) *ChattersListCall {
	return &ChattersListCall{
		resource: r,
		opts: []RequestOption{
			SetQueryParameter("broadcaster_id", broadcasterId),
			SetQueryParameter("moderator_id", moderatorId),
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

func (c *ChattersListCall) Do(ctx context.Context, opts ...RequestOption) (*ChattersResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/chat/chatters", nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Chatter](res)
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

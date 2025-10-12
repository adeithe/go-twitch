package api

import (
	"context"
	"fmt"
	"net/http"
)

// GamesResource provides access to the Twitch Games API.
type GamesResource struct {
	client *Client

	Top *TopGamesResource
}

// NewGamesResource creates a new GamesResource.
func NewGamesResource(client *Client) *GamesResource {
	c := &GamesResource{client: client}
	c.Top = NewTopGamesResource(client)
	return c
}

// TopGamesResource provides methods for the Twitch Top Games API.
type TopGamesResource struct {
	client *Client
}

// NewTopGamesResource creates a new TopGamesResource.
func NewTopGamesResource(client *Client) *TopGamesResource {
	return &TopGamesResource{client}
}

// TopGamesListCall is a request to list top games based on the specified criteria.
type TopGamesListCall struct {
	resource *TopGamesResource
	opts     []RequestOption
}

// TopGamesListResponse is the response from listing top games.
type TopGamesListResponse struct {
	Header http.Header
	Data   []Game
	Cursor string
}

// List creates a request to list top games based on the specified criteria.
func (r *TopGamesResource) List() *TopGamesListCall {
	return &TopGamesListCall{resource: r}
}

// First limits the number of results to the specified amount.
//
// Maximum: 100 (default: 20)
func (c *TopGamesListCall) First(n int) *TopGamesListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// Before filters the results to those with a cursor value before the specified cursor.
func (c *TopGamesListCall) Before(cursor string) *TopGamesListCall {
	c.opts = append(c.opts, SetQueryParameter("before", cursor))
	return c
}

// After filters the results to those with a cursor value after the specified cursor.
func (c *TopGamesListCall) After(cursor string) *TopGamesListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Do executes the request.
func (c *TopGamesListCall) Do(ctx context.Context, opts ...RequestOption) (*TopGamesListResponse, error) {
	res, err := c.resource.client.DoRequest(ctx, http.MethodGet, EndpointGamesTop, nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Game](res)
	if err != nil {
		return nil, err
	}

	return &TopGamesListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

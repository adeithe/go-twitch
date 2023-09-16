package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Stream struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	UserLogin       string    `json:"user_login"`
	UserDisplayName string    `json:"user_name"`
	GameID          string    `json:"game_id"`
	GameName        string    `json:"game_name"`
	Type            string    `json:"type"`
	Title           string    `json:"title"`
	Tags            []string  `json:"tags"`
	ViewerCount     int       `json:"viewer_count"`
	Language        string    `json:"language"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	IsMature        bool      `json:"is_mature"`
	StartedAt       time.Time `json:"started_at"`
}

type StreamsResource struct {
	client *Client
}

func NewStreamsResource(client *Client) *StreamsResource {
	return &StreamsResource{client}
}

type StreamsListCall struct {
	resource *StreamsResource
	opts     []RequestOption
}

type StreamsListResponse struct {
	Header http.Header
	Data   []Stream
	Cursor string
}

// List creates a request to list streams based on the specified criteria.
//
// Requires an app or user access token. No scope is required.
func (r *StreamsResource) List() *StreamsListCall {
	return &StreamsListCall{resource: r}
}

// UserID filters the results to the specified user IDs.
func (c *StreamsListCall) UserID(ids ...string) *StreamsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("user_id", id))
	}
	return c
}

// Username filters the results to the specified usernames.
func (c *StreamsListCall) Username(usernames ...string) *StreamsListCall {
	for _, username := range usernames {
		c.opts = append(c.opts, AddQueryParameter("user_login", username))
	}
	return c
}

// GameID filters the results to the specified game IDs.
func (c *StreamsListCall) GameID(ids ...string) *StreamsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("game_id", id))
	}
	return c
}

// Type filters the results to the specified stream types.
//
// Possible values: "all", "live" (Default: "all")
func (c *StreamsListCall) Type(t string) *StreamsListCall {
	c.opts = append(c.opts, SetQueryParameter("type", t))
	return c
}

// Language filters the results to the specified languages.
func (c *StreamsListCall) Languages(languages ...string) *StreamsListCall {
	for _, language := range languages {
		c.opts = append(c.opts, AddQueryParameter("language", language))
	}
	return c
}

// First limits the number of results to the specified amount.
//
// Maximum: 100 (default: 20)
func (c *StreamsListCall) First(n int) *StreamsListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// Before filters the results to streams that started before the specified cursor.
func (c *StreamsListCall) Before(cursor string) *StreamsListCall {
	c.opts = append(c.opts, SetQueryParameter("before", cursor))
	return c
}

// After filters the results to streams that started after the specified cursor.
func (c *StreamsListCall) After(cursor string) *StreamsListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Do executes the request.
func (c *StreamsListCall) Do(ctx context.Context, opts ...RequestOption) (*StreamsListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/streams", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Stream](res)
	if err != nil {
		return nil, err
	}

	return &StreamsListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

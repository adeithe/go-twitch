package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Video struct {
	ID              string        `json:"id"`
	StreamID        string        `json:"stream_id"`
	UserID          string        `json:"user_id"`
	UserLogin       string        `json:"user_login"`
	UserDisplayName string        `json:"user_name"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	URL             string        `json:"url"`
	ThumbnailURL    string        `json:"thumbnail_url"`
	Viewable        string        `json:"viewable"`
	ViewCount       int           `json:"view_count"`
	Language        string        `json:"language"`
	Type            string        `json:"type"`
	Duration        VideoDuration `json:"duration"`
	PublishedAt     time.Time     `json:"published_at"`
	CreatedAt       time.Time     `json:"created_at"`
}

type VideoDuration time.Duration

type VideosResource struct {
	client *Client
}

func NewVideosResource(client *Client) *VideosResource {
	return &VideosResource{client}
}

type VideosListCall struct {
	resource *VideosResource
	opts     []RequestOption
}

type VideosListResponse struct {
	Header http.Header
	Data   []Video
	Cursor string
}

// List creates a new call to list videos.
//
// One of ID, UserID, or GameID must be specified.
func (r *VideosResource) List() *VideosListCall {
	return &VideosListCall{resource: r}
}

// ID filters the results to those with the specified ID.
func (c *VideosListCall) ID(ids []string) *VideosListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

// UserID filters the results to those with the specified user ID.
func (c *VideosListCall) UserID(id string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("user_id", id))
	return c
}

// GameID filters the results to those with the specified game ID.
func (c *VideosListCall) GameID(id string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("game_id", id))
	return c
}

// Language filters the results to those with the specified language.
func (c *VideosListCall) Language(language string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("language", language))
	return c
}

// Period filters the results to those with a specified period.
//
// Possible values: "all", "day", "week", "month" (default: all)
func (c *VideosListCall) Period(p string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("period", p))
	return c
}

// Sort sets the order in which to list videos.
//
// Possible values: "time", "trending", "views" (default: time)
func (c *VideosListCall) Sort(s string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("sort", s))
	return c
}

func (c *VideosListCall) First(n int) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// Type filters the results to those with the specified type.
//
// Possible values: "all", "upload", "archive", "highlight" (default: all)
func (c *VideosListCall) Type(t string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("type", t))
	return c
}

// Before filters the results to those with a cursor value before the specified cursor.
func (c *VideosListCall) Before(cursor string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("before", cursor))
	return c
}

// After filters the results to those with a cursor value after the specified cursor.
func (c *VideosListCall) After(cursor string) *VideosListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

func (c *VideosListCall) Do(ctx context.Context, opts ...RequestOption) (*VideosListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/videos", nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Video](res)
	if err != nil {
		return nil, err
	}

	return &VideosListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

type VideosDeleteCall struct {
	resource *VideosResource
	opts     []RequestOption
}

type VideosDeleteResponse struct {
	Header http.Header
	Data   []string
}

// Delete creates a new call to delete videos.
func (r *VideosResource) Delete(ids []string) *VideosDeleteCall {
	c := &VideosDeleteCall{resource: r}
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

func (c *VideosDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*VideosDeleteResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodDelete, "/videos", nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[string](res)
	if err != nil {
		return nil, err
	}

	return &VideosDeleteResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

func (d *VideoDuration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = VideoDuration(parsed)
	return nil
}

func (d VideoDuration) AsDuration() time.Duration {
	return time.Duration(d)
}

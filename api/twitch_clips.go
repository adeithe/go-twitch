package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Clip struct {
	ID              string       `json:"id"`
	URL             string       `json:"url"`
	EmbedURL        string       `json:"embed_url"`
	BroadcasterID   string       `json:"broadcaster_id"`
	BroadcasterName string       `json:"broadcaster_name"`
	CreatorID       string       `json:"creator_id"`
	CreatorName     string       `json:"creator_name"`
	VideoID         string       `json:"video_id"`
	GameID          string       `json:"game_id"`
	Language        string       `json:"language"`
	Title           string       `json:"title"`
	ViewCount       int          `json:"view_count"`
	ThumbnailURL    string       `json:"thumbnail_url"`
	Duration        ClipDuration `json:"duration"`
	VODOffset       int          `json:"vod_offset"`
	Featured        bool         `json:"is_featured"`
	CreatedAt       time.Time    `json:"created_at"`
}

type ClipDuration time.Duration

type ClipsResource struct {
	client *Client
}

func NewClipsResource(client *Client) *ClipsResource {
	return &ClipsResource{client}
}

type ClipsListCall struct {
	resource *ClipsResource
	opts     []RequestOption
}

type ClipsListResponse struct {
	Header http.Header
	Data   []Clip
	Cursor string
}

// List creates a new call to list clips.
//
// One or more of ID, BroadcasterID, or GameID must be specified.
func (r *ClipsResource) List() *ClipsListCall {
	return &ClipsListCall{resource: r}
}

// ID filters the results to those with the specified clip IDs.
func (c *ClipsListCall) ID(ids []string) *ClipsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

// BroadcasterID filters the results to those with the specified broadcaster ID.
func (c *ClipsListCall) BroadcasterID(id string) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", id))
	return c
}

// GameID filters the results to those with the specified game ID.
func (c *ClipsListCall) GameID(id string) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("game_id", id))
	return c
}

// StartedAt filters the results to those created after the specified time.
func (c *ClipsListCall) StartedAt(t time.Time) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("started_at", t.Format(time.RFC3339)))
	return c
}

// EndedAt filters the results to those created before the specified time.
func (c *ClipsListCall) EndedAt(t time.Time) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("ended_at", t.Format(time.RFC3339)))
	return c
}

// First filters the results to the first n clips.
func (c *ClipsListCall) First(n int) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// Before filters the results to those before the specified cursor.
func (c *ClipsListCall) Before(cursor string) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("before", cursor))
	return c
}

// After filters the results to those after the specified cursor.
func (c *ClipsListCall) After(cursor string) *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Featured filters the results to only those that are featured.
func (c *ClipsListCall) Featured() *ClipsListCall {
	c.opts = append(c.opts, SetQueryParameter("is_featured", "true"))
	return c
}

// Do executes the call.
func (c *ClipsListCall) Do(ctx context.Context, opts ...RequestOption) (*ClipsListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/clips", nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Clip](res)
	if err != nil {
		return nil, err
	}

	return &ClipsListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

func (d *ClipDuration) UnmarshalJSON(data []byte) error {
	var duration float64
	if err := json.Unmarshal(data, &duration); err != nil {
		return err
	}

	*d = ClipDuration(time.Duration(float64(time.Second) * duration))
	return nil
}

func (d ClipDuration) AsDuration() time.Duration {
	return time.Duration(d)
}

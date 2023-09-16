package api

import (
	"context"
	"net/http"
)

type Channel struct {
	ID                          string   `json:"broadcaster_id"`
	Login                       string   `json:"broadcaster_login"`
	DisplayName                 string   `json:"broadcaster_name"`
	GameID                      string   `json:"game_id"`
	GameName                    string   `json:"game_name"`
	Title                       string   `json:"title"`
	Delay                       int      `json:"delay"`
	Tags                        []string `json:"tags"`
	ContentClassificationLabels []string `json:"content_classification_labels"`
	IsBrandedContent            bool     `json:"is_branded_content"`
}

type ChannelsResource struct {
	client *Client
}

func NewChannelsResource(client *Client) *ChannelsResource {
	return &ChannelsResource{client}
}

type ChannelsListCall struct {
	resource *ChannelsResource
	opts     []RequestOption
}

// List creates a request to list channels based on the specified criteria.
func (r *ChannelsResource) List() *ChannelsListCall {
	return &ChannelsListCall{resource: r}
}

// BroadcasterID filters the results to the specified broadcaster ID.
func (c *ChannelsListCall) BroadcasterID(ids ...string) *ChannelsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("broadcaster_id", id))
	}
	return c
}

// Do executes the request.
func (c *ChannelsListCall) Do(ctx context.Context, opts ...RequestOption) ([]Channel, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channels", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Channel](res)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

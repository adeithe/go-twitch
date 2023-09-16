package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Commercial struct {
	Length     int    `json:"length"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

type AdsResource struct {
	client *Client
}

func NewAdsResource(client *Client) *AdsResource {
	return &AdsResource{client}
}

type AdsInsertRequest struct {
	resource      *AdsResource
	broadcasterID string
	duration      int
}

// Insert creates a request to start a commercial for the specified broadcaster.
//
// Required Scope: channel:edit:commercial
func (r *AdsResource) Insert(broadcasterId string) *AdsInsertRequest {
	return &AdsInsertRequest{r, broadcasterId, 60}
}

// Duration sets the duration of the commercial in seconds.
//
// Twitch tries to serve a commercial that’s the requested length, but it may be shorter or longer. The maximum length you should request is 180 seconds.
func (c *AdsInsertRequest) Duration(seconds int) *AdsInsertRequest {
	c.duration = seconds
	return c
}

// Do executes the request.
//
//	req := client.Ads.Insert("41245072").Duration(60)
//	data, err := req.Do(ctx, api.WithBearerToken("2gbdx6oar67tqtcmt49t3wpcgycthx")
func (c *AdsInsertRequest) Do(ctx context.Context, opts ...RequestOption) ([]Commercial, error) {
	bs, err := json.Marshal(map[string]interface{}{
		"broadcaster_id": c.broadcasterID,
		"length":         c.duration,
	})
	if err != nil {
		return nil, err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPost, "/channels/commercial", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Commercial](res)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

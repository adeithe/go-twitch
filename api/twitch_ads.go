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

type StartCommercialRequest struct {
	resource      *AdsResource
	broadcasterID string
	duration      int
}

// StartCommercial creates a request to start a commercial for the specified broadcaster.
//
// Required Scope: channel:edit:commercial
func (r *AdsResource) StartCommercial(broadcasterId string) *StartCommercialRequest {
	return &StartCommercialRequest{r, broadcasterId, 60}
}

// BroadcasterID sets the broadcaster ID for the request.
//
// This ID must match the user ID found in the OAuth token.
func (c *StartCommercialRequest) BroadcasterID(broadcasterId string) *StartCommercialRequest {
	c.broadcasterID = broadcasterId
	return c
}

// Duration sets the duration of the commercial in seconds.
//
// Twitch tries to serve a commercial that’s the requested length, but it may be shorter or longer. The maximum length you should request is 180 seconds.
func (c *StartCommercialRequest) Duration(seconds int) *StartCommercialRequest {
	c.duration = seconds
	return c
}

// Do executes the request.
//
//	req := client.Ads.StartCommercial("41245072").Duration(60)
//	data, err := req.Do(ctx, api.WithBearerToken("2gbdx6oar67tqtcmt49t3wpcgycthx")
func (c *StartCommercialRequest) Do(ctx context.Context, opts ...RequestOption) ([]Commercial, error) {
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

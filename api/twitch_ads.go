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

func (r *AdsResource) StartCommercial(broadcasterID string) *StartCommercialRequest {
	return &StartCommercialRequest{r, broadcasterID, 60}
}

func (c *StartCommercialRequest) Duration(seconds int) *StartCommercialRequest {
	c.duration = seconds
	return c
}

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

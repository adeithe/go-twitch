package api

import (
	"context"
	"net/http"
	"time"
)

type Cheermote struct {
	Prefix       string          `json:"prefix"`
	Tiers        []CheermoteTier `json:"tiers"`
	Type         string          `json:"type"`
	Order        int             `json:"order"`
	IsCharitable bool            `json:"is_charitable"`
	LastUpdated  time.Time       `json:"last_updated"`
}

type CheermoteTier struct {
	ID             string            `json:"id"`
	MinBits        int               `json:"min_bits"`
	Color          string            `json:"color"`
	Images         map[string]string `json:"images"`
	CanCheer       bool              `json:"can_cheer"`
	ShowInBitsCard bool              `json:"show_in_bits_card"`
}

type CheermotesResource struct {
	client *Client
}

func NewCheermotesResource(client *Client) *CheermotesResource {
	return &CheermotesResource{client}
}

type CheermotesListCall struct {
	resource *CheermotesResource
	opts     []RequestOption
}

// List creates a request to list cheermotes based on the specified criteria.
//
// Requires an app or user access token. No scope is required.
func (r *CheermotesResource) List() *CheermotesListCall {
	return &CheermotesListCall{resource: r}
}

// BroadcasterID filters the results to the specified broadcaster ID.
func (c *CheermotesListCall) BroadcasterID(id string) *CheermotesListCall {
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", id))
	return c
}

// BroadcasterName filters the results to the specified broadcaster name.
func (c *CheermotesListCall) Do(ctx context.Context, opts ...RequestOption) ([]Cheermote, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/bits/cheermotes", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Cheermote](res)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

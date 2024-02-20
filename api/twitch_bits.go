package api

import (
	"context"
	"fmt"
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

type BitsLeaderboardEntry struct {
	ID          string `json:"user_id"`
	Login       string `json:"user_login"`
	DisplayName string `json:"user_name"`
	Rank        int    `json:"rank"`
	Score       int    `json:"score"`
}

type BitsResource struct {
	client *Client

	Cheermotes            *CheermotesResource
	Leaderboard           *BitsLeaderboardResource
	ExtensionTransactions *BitsExtensionTransactionsResource
}

func NewBitsResource(client *Client) *BitsResource {
	r := &BitsResource{client: client}
	r.Cheermotes = NewCheermotesResource(client)
	r.Leaderboard = NewBitsLeaderboardResource(client)
	r.ExtensionTransactions = NewBitsExtensionTransactionsResource(client)
	return r
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

type CheermotesListResponse struct {
	Header http.Header
	Data   []Cheermote
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
func (c *CheermotesListCall) Do(ctx context.Context, opts ...RequestOption) (*CheermotesListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/bits/cheermotes", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Cheermote](res)
	if err != nil {
		return nil, err
	}

	return &CheermotesListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

type BitsLeaderboardResource struct {
	client *Client
}

func NewBitsLeaderboardResource(client *Client) *BitsLeaderboardResource {
	return &BitsLeaderboardResource{client}
}

type BitsLeaderboardListCall struct {
	resource *BitsLeaderboardResource
	opts     []RequestOption
}

type BitsLeaderboardListResponse struct {
	Header http.Header
	Total  int
	Data   []BitsLeaderboardEntry
}

// List creates a request to list users from the authenticated users Bits leaderboard.
func (r *BitsLeaderboardResource) List() *BitsLeaderboardListCall {
	return &BitsLeaderboardListCall{resource: r}
}

// Count limits the number of results to return.
//
// Maximum: 100 (default: 10)
func (c *BitsLeaderboardListCall) Count(n int) *BitsLeaderboardListCall {
	c.opts = append(c.opts, SetQueryParameter("count", fmt.Sprint(n)))
	return c
}

// Period sets the time period over which data is aggregated.
//
// Possible values: "day", "week", "month", "year", "all" (default: "all")
func (c *BitsLeaderboardListCall) Period(period string) *BitsLeaderboardListCall {
	c.opts = append(c.opts, SetQueryParameter("period", period))
	return c
}

// StartedAt the start date used for determining the aggregation period.
func (c *BitsLeaderboardListCall) StartedAt(t time.Time) *BitsLeaderboardListCall {
	c.opts = append(c.opts, SetQueryParameter("started_at", t.Format(time.RFC3339)))
	return c
}

// UserID limits the aggregated results to the specified user ID.
// If count is greater than 1, the response may include users ranked above and below the specified user.
//
// To get the leaderboard's top leaders, don't specify this.
func (c *BitsLeaderboardListCall) UserID(userId string) *BitsLeaderboardListCall {
	c.opts = append(c.opts, SetQueryParameter("user_id", userId))
	return c
}

// Do executes the request.
func (c *BitsLeaderboardListCall) Do(ctx context.Context, opts ...RequestOption) (*BitsLeaderboardListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/bits/leaderboard", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[BitsLeaderboardEntry](res)
	if err != nil {
		return nil, err
	}

	return &BitsLeaderboardListResponse{
		Header: res.Header,
		Total:  data.Total,
		Data:   data.Data,
	}, nil
}

type BitsExtensionTransactionsResource struct {
	client *Client
}

func NewBitsExtensionTransactionsResource(client *Client) *BitsExtensionTransactionsResource {
	return &BitsExtensionTransactionsResource{client}
}

type BitsTransactionsListCall struct {
	client *Client
	opts   []RequestOption
}

type BitsTransactionsListResponse struct {
	Header     http.Header
	Data       []ExtensionTransaction
	Pagination Pagination
}

type ExtensionTransaction struct {
	ID               string           `json:"id"`
	BroadcasterID    string           `json:"broadcaster_id"`
	BroadcasterLogin string           `json:"broadcaster_login"`
	BroadcasterName  string           `json:"broadcaster_name"`
	UserID           string           `json:"user_id"`
	UserLogin        string           `json:"user_login"`
	UserName         string           `json:"user_name"`
	ProductType      string           `json:"product_type"`
	Product          ExtensionProduct `json:"product_data"`
	Timestamp        time.Time        `json:"timestamp"`
}

type ExtensionProduct struct {
	Sku           string               `json:"sku"`
	Domain        string               `json:"domain"`
	Cost          ExtensionProductCost `json:"cost"`
	InDevelopment bool                 `json:"inDevelopment"`
	DisplayName   string               `json:"displayName"`
	Broadcast     bool                 `json:"broadcast"`
}

type ExtensionProductCost struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

func (r *BitsExtensionTransactionsResource) List(extensionId string) *BitsTransactionsListCall {
	return &BitsTransactionsListCall{
		client: r.client,
		opts: []RequestOption{
			SetQueryParameter("extension_id", extensionId),
		},
	}
}

func (c *BitsTransactionsListCall) TransactionID(ids ...string) *BitsTransactionsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

func (c *BitsTransactionsListCall) First(n int) *BitsTransactionsListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

func (c *BitsTransactionsListCall) After(cursor string) *BitsTransactionsListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

func (c *BitsTransactionsListCall) Do(ctx context.Context, opts ...RequestOption) (*BitsTransactionsListResponse, error) {
	res, err := c.client.doRequest(ctx, http.MethodGet, "/extensions/transactions", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[ExtensionTransaction](res)
	if err != nil {
		return nil, err
	}

	return &BitsTransactionsListResponse{
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
	}, nil
}

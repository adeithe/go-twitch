package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type CustomReward struct {
	ID                          string    `json:"id"`
	BroadcasterID               string    `json:"broadcaster_id"`
	BroadcasterLogin            string    `json:"broadcaster_login"`
	BroadcasterDisplayName      string    `json:"broadcaster_name"`
	BackgroundColor             string    `json:"background_color"`
	Title                       string    `json:"title"`
	Prompt                      string    `json:"prompt"`
	Cost                        int64     `json:"cost"`
	Enabled                     bool      `json:"is_enabled"`
	Paused                      bool      `json:"is_paused"`
	InStock                     bool      `json:"is_in_stock"`
	IsUserInputRequired         bool      `json:"is_user_input_required"`
	RedemptionsSkipRequestQueue bool      `json:"should_redemptions_skip_request_queue"`
	CooldownExpiresAt           time.Time `json:"cooldown_expires_at"`
}

type CustomRewardRedemption struct {
	ID                     string `json:"id"`
	BroadcasterID          string `json:"broadcaster_id"`
	BroadcasterLogin       string `json:"broadcaster_login"`
	BroadcasterDisplayName string `json:"broadcaster_name"`
	UserID                 string `json:"user_id"`
	UserLogin              string `json:"user_login"`
	UserDisplayName        string `json:"user_name"`
	UserInput              string `json:"user_input"`
	Status                 string `json:"status"`
	Reward                 struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Prompt string `json:"prompt"`
		Cost   int64  `json:"cost"`
	} `json:"reward"`
	RedeemedAt time.Time `json:"redeemed_at"`
}

type ChannelPointsResource struct {
	client *Client

	CustomRewards *CustomRewardsResource
}

func NewChannelPointsResource(client *Client) *ChannelPointsResource {
	r := &ChannelPointsResource{client: client}
	r.CustomRewards = NewCustomRewardsResource(client)
	return r
}

type CustomRewardsResource struct {
	client *Client

	Redemption *CustomRewardsRedemptionResource
}

func NewCustomRewardsResource(client *Client) *CustomRewardsResource {
	return &CustomRewardsResource{client: client}
}

type CustomRewardsListCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
}

type CustomRewardsListResponse struct {
	Header http.Header
	Data   []CustomReward
}

// List creates a reqyest to list custom channel point rewards for a given broadcaster.
func (r *CustomRewardsResource) List(broadcasterId string) *CustomRewardsListCall {
	c := &CustomRewardsListCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterId))
	return c
}

// ID filters the results to the specified reward IDs.
func (c *CustomRewardsListCall) ID(ids ...string) *CustomRewardsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, SetQueryParameter("id", id))
	}
	return c
}

// OnlyManageable filters the results to only rewards that the app may manage.
func (c *CustomRewardsListCall) OnlyManageable() *CustomRewardsListCall {
	c.opts = append(c.opts, SetQueryParameter("only_manageable_rewards", "true"))
	return c
}

// Do executes the request.
func (c *CustomRewardsListCall) Do(ctx context.Context, opts ...RequestOption) (*CustomRewardsListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channel_points/custom_rewards", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

type CustomRewardsInsertCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
	body     map[string]interface{}
}

type CustomRewardsInsertResponse struct {
	Header http.Header
	Data   []CustomReward
}

func (r *CustomRewardsResource) Insert(broadcasterId string) *CustomRewardsInsertCall {
	c := &CustomRewardsInsertCall{resource: r, body: make(map[string]interface{})}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterId))
	return c
}

func (c *CustomRewardsInsertCall) Title(title string) *CustomRewardsInsertCall {
	c.body["title"] = title
	return c
}

func (c *CustomRewardsInsertCall) Prompt(prompt string) *CustomRewardsInsertCall {
	c.body["prompt"] = prompt
	return c
}

func (c *CustomRewardsInsertCall) Cost(cost int64) *CustomRewardsInsertCall {
	c.body["cost"] = cost
	return c
}

func (c *CustomRewardsInsertCall) BackgroundColor(hexCode string) *CustomRewardsInsertCall {
	c.body["background_color"] = hexCode
	return c
}

func (c *CustomRewardsInsertCall) IsEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_enabled"] = enabled
	return c
}

func (c *CustomRewardsInsertCall) IsUserInputRequired(required bool) *CustomRewardsInsertCall {
	c.body["is_user_input_required"] = required
	return c
}

func (c *CustomRewardsInsertCall) IsMaxPerStreamEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_max_per_stream_enabled"] = enabled
	return c
}

func (c *CustomRewardsInsertCall) MaxPerStream(max int64) *CustomRewardsInsertCall {
	c.body["max_per_stream"] = max
	return c
}

func (c *CustomRewardsInsertCall) IsMaxPerUserPerStreamEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_max_per_user_per_stream_enabled"] = enabled
	return c
}

func (c *CustomRewardsInsertCall) MaxPerUserPerStream(max int64) *CustomRewardsInsertCall {
	c.body["max_per_user_per_stream"] = max
	return c
}

func (c *CustomRewardsInsertCall) IsGlobalCooldownEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_global_cooldown_enabled"] = enabled
	return c
}

func (c *CustomRewardsInsertCall) GlobalCooldown(d time.Duration) *CustomRewardsInsertCall {
	c.body["global_cooldown_seconds"] = d.Seconds()
	return c
}

func (c *CustomRewardsInsertCall) IsPaused(paused bool) *CustomRewardsInsertCall {
	c.body["is_paused"] = paused
	return c
}

func (c *CustomRewardsInsertCall) ShouldRedemptionsSkipRequestQueue(b bool) *CustomRewardsInsertCall {
	c.body["should_redemptions_skip_request_queue"] = b
	return c
}

// Do executes the request.
func (c *CustomRewardsInsertCall) Do(ctx context.Context, opts ...RequestOption) (*CustomRewardsInsertResponse, error) {
	bs, err := json.Marshal(c.body)
	if err != nil {
		return nil, err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPost, "/channel_points/custom_rewards", bytes.NewReader(bs), append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsInsertResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

type CustomRewardsUpdateCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
	body     map[string]interface{}
}

type CustomRewardsUpdateResponse struct {
	Header http.Header
	Data   []CustomReward
}

func (r *CustomRewardsResource) Update(broadcasterId, id string) *CustomRewardsUpdateCall {
	c := &CustomRewardsUpdateCall{resource: r, body: make(map[string]interface{})}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterId))
	c.opts = append(c.opts, SetQueryParameter("id", id))
	return c
}

func (c *CustomRewardsUpdateCall) Title(title string) *CustomRewardsUpdateCall {
	c.body["title"] = title
	return c
}

func (c *CustomRewardsUpdateCall) Prompt(prompt string) *CustomRewardsUpdateCall {
	c.body["prompt"] = prompt
	return c
}

func (c *CustomRewardsUpdateCall) Cost(cost int64) *CustomRewardsUpdateCall {
	c.body["cost"] = cost
	return c
}

func (c *CustomRewardsUpdateCall) BackgroundColor(hexCode string) *CustomRewardsUpdateCall {
	c.body["background_color"] = hexCode
	return c
}

func (c *CustomRewardsUpdateCall) IsEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_enabled"] = enabled
	return c
}

func (c *CustomRewardsUpdateCall) IsUserInputRequired(required bool) *CustomRewardsUpdateCall {
	c.body["is_user_input_required"] = required
	return c
}

func (c *CustomRewardsUpdateCall) IsMaxPerStreamEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_max_per_stream_enabled"] = enabled
	return c
}

func (c *CustomRewardsUpdateCall) MaxPerStream(max int64) *CustomRewardsUpdateCall {
	c.body["max_per_stream"] = max
	return c
}

func (c *CustomRewardsUpdateCall) IsMaxPerUserPerStreamEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_max_per_user_per_stream_enabled"] = enabled
	return c
}

func (c *CustomRewardsUpdateCall) MaxPerUserPerStream(max int64) *CustomRewardsUpdateCall {
	c.body["max_per_user_per_stream"] = max
	return c
}

func (c *CustomRewardsUpdateCall) IsGlobalCooldownEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_global_cooldown_enabled"] = enabled
	return c
}

func (c *CustomRewardsUpdateCall) GlobalCooldown(d time.Duration) *CustomRewardsUpdateCall {
	c.body["global_cooldown_seconds"] = d.Seconds()
	return c
}

func (c *CustomRewardsUpdateCall) IsPaused(paused bool) *CustomRewardsUpdateCall {
	c.body["is_paused"] = paused
	return c
}

func (c *CustomRewardsUpdateCall) ShouldRedemptionsSkipRequestQueue(b bool) *CustomRewardsUpdateCall {
	c.body["should_redemptions_skip_request_queue"] = b
	return c
}

// Do executes the request.
func (c *CustomRewardsUpdateCall) Do(ctx context.Context, opts ...RequestOption) (*CustomRewardsUpdateResponse, error) {
	bs, err := json.Marshal(c.body)
	if err != nil {
		return nil, err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPatch, "/channel_points/custom_rewards", bytes.NewReader(bs), append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsUpdateResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

type CustomRewardsDeleteCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
}

func (r *CustomRewardsResource) Delete(broadcasterId, id string) *CustomRewardsDeleteCall {
	c := &CustomRewardsDeleteCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterId))
	c.opts = append(c.opts, SetQueryParameter("id", id))
	return c
}

// Do executes the request.
func (c *CustomRewardsDeleteCall) Do(ctx context.Context, opts ...RequestOption) error {
	res, err := c.resource.client.doRequest(ctx, http.MethodDelete, "/channel_points/custom_rewards", nil, append(opts, c.opts...)...)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = decodeResponse[CustomRewardRedemption](res)
	return err
}

type CustomRewardsRedemptionResource struct {
	client *Client
}

func NewCustomRewardsRedemptionResource(client *Client) *CustomRewardsRedemptionResource {
	return &CustomRewardsRedemptionResource{client: client}
}

type CustomRewardsRedemptionListCall struct {
	resource *CustomRewardsRedemptionResource
	opts     []RequestOption
}

type CustomRewardsRedemptionListResponse struct {
	Header http.Header
	Data   []CustomRewardRedemption
	Cursor string
}

// List creates a request to list custom channel point reward redemptions for a given broadcaster.
func (r *CustomRewardsRedemptionResource) List(broadcasterId, rewardId string) *CustomRewardsRedemptionListCall {
	c := &CustomRewardsRedemptionListCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterId))
	c.opts = append(c.opts, SetQueryParameter("reward_id", rewardId))
	c.opts = append(c.opts, SetQueryParameter("status", "UNFULFILLED"))
	return c
}

// Status filters the results to the specified statuses.
//
// Possible values: "UNFULFILLED", "FULFILLED", "CANCELED" (default: UNFULFILLED)
func (c *CustomRewardsRedemptionListCall) Status(status string) *CustomRewardsRedemptionListCall {
	c.opts = append(c.opts, SetQueryParameter("status", status))
	return c
}

// ID filters the results to the specified reward redemption IDs.
func (c *CustomRewardsRedemptionListCall) ID(ids ...string) *CustomRewardsRedemptionListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

// Sort specifies the order in which to sort the results.
//
// Possible values: "OLDEST", "NEWEST" (default: OLDEST)
func (c *CustomRewardsRedemptionListCall) Sort(sort string) *CustomRewardsRedemptionListCall {
	c.opts = append(c.opts, SetQueryParameter("sort", sort))
	return c
}

// Before filters the results to those with a cursor value before the specified cursor.
func (c *CustomRewardsRedemptionListCall) Before(cursor string) *CustomRewardsRedemptionListCall {
	c.opts = append(c.opts, SetQueryParameter("before", cursor))
	return c
}

// After filters the results to those with a cursor value after the specified cursor.
func (c *CustomRewardsRedemptionListCall) After(cursor string) *CustomRewardsRedemptionListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Do executes the request.
func (c *CustomRewardsRedemptionListCall) Do(ctx context.Context, opts ...RequestOption) (*CustomRewardsRedemptionListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channel_points/custom_rewards/redemptions", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[CustomRewardRedemption](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsRedemptionListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

type CustomRewardsRedemptionUpdateCall struct {
	resource *CustomRewardsRedemptionResource
	opts     []RequestOption
}

type CustomRewardsRedemptionUpdateResponse struct {
	Header http.Header
	Data   []CustomRewardRedemption
}

func (r *CustomRewardsRedemptionResource) Update(broadcasterId, rewardId string, id []string) *CustomRewardsRedemptionUpdateCall {
	c := &CustomRewardsRedemptionUpdateCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterId))
	c.opts = append(c.opts, SetQueryParameter("reward_id", rewardId))
	for _, id := range id {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

func (c *CustomRewardsRedemptionUpdateCall) Cancel() *CustomRewardsRedemptionUpdateCall {
	c.opts = append(c.opts, SetQueryParameter("status", "CANCELED"))
	return c
}

func (c *CustomRewardsRedemptionUpdateCall) Fulfill() *CustomRewardsRedemptionUpdateCall {
	c.opts = append(c.opts, SetQueryParameter("status", "FULFILLED"))
	return c
}

// Do executes the request.
func (c *CustomRewardsRedemptionUpdateCall) Do(ctx context.Context, opts ...RequestOption) (*CustomRewardsRedemptionUpdateResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodPatch, "/channel_points/custom_rewards/redemptions", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[CustomRewardRedemption](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsRedemptionUpdateResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

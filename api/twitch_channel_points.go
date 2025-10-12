package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// ChannelPointsResource provides methods for the Twitch Channel Points API.
type ChannelPointsResource struct {
	client *Client

	CustomRewards *CustomRewardsResource
}

// NewChannelPointsResource creates a new ChannelPointsResource.
func NewChannelPointsResource(client *Client) *ChannelPointsResource {
	r := &ChannelPointsResource{client: client}
	r.CustomRewards = NewCustomRewardsResource(client)
	return r
}

// CustomRewardsResource provides methods for the Twitch Channel Points Custom Rewards API.
type CustomRewardsResource struct {
	client *Client

	Redemption *CustomRewardsRedemptionResource
}

// NewCustomRewardsResource creates a new CustomRewardsResource.
func NewCustomRewardsResource(client *Client) *CustomRewardsResource {
	return &CustomRewardsResource{client: client}
}

// CustomRewardsListCall is a call to the custom rewards list endpoint.
type CustomRewardsListCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
}

// CustomRewardsListResponse is the response from the custom rewards list endpoint.
type CustomRewardsListResponse struct {
	Header http.Header
	Data   []CustomReward
}

// List creates a reqyest to list custom channel point rewards for a given broadcaster.
func (r *CustomRewardsResource) List(broadcasterID string) *CustomRewardsListCall {
	c := &CustomRewardsListCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return c
}

// ID filters the results to the specified reward IDs.
func (c *CustomRewardsListCall) ID(ids []string) *CustomRewardsListCall {
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
	res, err := c.resource.client.DoRequest(ctx, http.MethodGet, EndpointChannelPointsGetCustomRewards, nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

// CustomRewardsInsertCall is a call to the custom rewards insert endpoint.
type CustomRewardsInsertCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
	body     map[string]any
}

// CustomRewardsInsertResponse is the response from the custom rewards insert endpoint.
type CustomRewardsInsertResponse struct {
	Header http.Header
	Data   []CustomReward
}

// Insert creates a request to create a custom channel point reward for a given broadcaster.
func (r *CustomRewardsResource) Insert(broadcasterID string) *CustomRewardsInsertCall {
	c := &CustomRewardsInsertCall{resource: r, body: make(map[string]any)}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return c
}

// Title sets the title of the reward.
func (c *CustomRewardsInsertCall) Title(title string) *CustomRewardsInsertCall {
	c.body["title"] = title
	return c
}

// Prompt sets the prompt of the reward.
func (c *CustomRewardsInsertCall) Prompt(prompt string) *CustomRewardsInsertCall {
	c.body["prompt"] = prompt
	return c
}

// Cost sets the cost of the reward.
func (c *CustomRewardsInsertCall) Cost(cost int64) *CustomRewardsInsertCall {
	c.body["cost"] = cost
	return c
}

// BackgroundColor sets the background color of the reward in hex code format.
func (c *CustomRewardsInsertCall) BackgroundColor(hexCode string) *CustomRewardsInsertCall {
	c.body["background_color"] = hexCode
	return c
}

// IsEnabled sets whether the reward is enabled.
func (c *CustomRewardsInsertCall) IsEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_enabled"] = enabled
	return c
}

// IsUserInputRequired sets whether user input is required for the reward.
func (c *CustomRewardsInsertCall) IsUserInputRequired(required bool) *CustomRewardsInsertCall {
	c.body["is_user_input_required"] = required
	return c
}

// IsMaxPerStreamEnabled sets whether the max per stream limit is enabled for the reward.
func (c *CustomRewardsInsertCall) IsMaxPerStreamEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_max_per_stream_enabled"] = enabled
	return c
}

// MaxPerStream sets the max per stream limit for the reward.
func (c *CustomRewardsInsertCall) MaxPerStream(limit int64) *CustomRewardsInsertCall {
	c.body["max_per_stream"] = limit
	return c
}

// IsMaxPerUserPerStreamEnabled sets whether the max per user per stream limit is enabled for the reward.
func (c *CustomRewardsInsertCall) IsMaxPerUserPerStreamEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_max_per_user_per_stream_enabled"] = enabled
	return c
}

// MaxPerUserPerStream sets the max per user per stream limit for the reward.
func (c *CustomRewardsInsertCall) MaxPerUserPerStream(limit int64) *CustomRewardsInsertCall {
	c.body["max_per_user_per_stream"] = limit
	return c
}

// IsGlobalCooldownEnabled sets whether the global cooldown is enabled for the reward.
func (c *CustomRewardsInsertCall) IsGlobalCooldownEnabled(enabled bool) *CustomRewardsInsertCall {
	c.body["is_global_cooldown_enabled"] = enabled
	return c
}

// GlobalCooldown sets the global cooldown duration for the reward.
func (c *CustomRewardsInsertCall) GlobalCooldown(d time.Duration) *CustomRewardsInsertCall {
	c.body["global_cooldown_seconds"] = d.Seconds()
	return c
}

// IsPaused sets whether the reward is paused.
func (c *CustomRewardsInsertCall) IsPaused(paused bool) *CustomRewardsInsertCall {
	c.body["is_paused"] = paused
	return c
}

// ShouldRedemptionsSkipRequestQueue sets whether redemptions should skip the request queue.
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

	res, err := c.resource.client.DoRequest(ctx, http.MethodPost, EndpointChannelPointsCreateCustomRewards, bytes.NewReader(bs), append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsInsertResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

// CustomRewardsUpdateCall is a call to the custom rewards update endpoint.
type CustomRewardsUpdateCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
	body     map[string]any
}

// CustomRewardsUpdateResponse is the response from the custom rewards update endpoint.
type CustomRewardsUpdateResponse struct {
	Header http.Header
	Data   []CustomReward
}

// Update creates a request to update a custom channel point reward for a given broadcaster.
func (r *CustomRewardsResource) Update(broadcasterID, id string) *CustomRewardsUpdateCall {
	c := &CustomRewardsUpdateCall{resource: r, body: make(map[string]any)}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	c.opts = append(c.opts, SetQueryParameter("id", id))
	return c
}

// Title sets the title of the reward.
func (c *CustomRewardsUpdateCall) Title(title string) *CustomRewardsUpdateCall {
	c.body["title"] = title
	return c
}

// Prompt sets the prompt of the reward.
func (c *CustomRewardsUpdateCall) Prompt(prompt string) *CustomRewardsUpdateCall {
	c.body["prompt"] = prompt
	return c
}

// Cost sets the cost of the reward.
func (c *CustomRewardsUpdateCall) Cost(cost int64) *CustomRewardsUpdateCall {
	c.body["cost"] = cost
	return c
}

// BackgroundColor sets the background color of the reward in hex code format.
func (c *CustomRewardsUpdateCall) BackgroundColor(hexCode string) *CustomRewardsUpdateCall {
	c.body["background_color"] = hexCode
	return c
}

// IsEnabled sets whether the reward is enabled.
func (c *CustomRewardsUpdateCall) IsEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_enabled"] = enabled
	return c
}

// IsUserInputRequired sets whether user input is required for the reward.
func (c *CustomRewardsUpdateCall) IsUserInputRequired(required bool) *CustomRewardsUpdateCall {
	c.body["is_user_input_required"] = required
	return c
}

// IsMaxPerStreamEnabled sets whether the max per stream limit is enabled for the reward.
func (c *CustomRewardsUpdateCall) IsMaxPerStreamEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_max_per_stream_enabled"] = enabled
	return c
}

// MaxPerStream sets the max per stream limit for the reward.
func (c *CustomRewardsUpdateCall) MaxPerStream(limit int64) *CustomRewardsUpdateCall {
	c.body["max_per_stream"] = limit
	return c
}

// IsMaxPerUserPerStreamEnabled sets whether the max per user per stream limit is enabled for the reward.
func (c *CustomRewardsUpdateCall) IsMaxPerUserPerStreamEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_max_per_user_per_stream_enabled"] = enabled
	return c
}

// MaxPerUserPerStream sets the max per user per stream limit for the reward.
func (c *CustomRewardsUpdateCall) MaxPerUserPerStream(limit int64) *CustomRewardsUpdateCall {
	c.body["max_per_user_per_stream"] = limit
	return c
}

// IsGlobalCooldownEnabled sets whether the global cooldown is enabled for the reward.
func (c *CustomRewardsUpdateCall) IsGlobalCooldownEnabled(enabled bool) *CustomRewardsUpdateCall {
	c.body["is_global_cooldown_enabled"] = enabled
	return c
}

// GlobalCooldown sets the global cooldown duration for the reward.
func (c *CustomRewardsUpdateCall) GlobalCooldown(d time.Duration) *CustomRewardsUpdateCall {
	c.body["global_cooldown_seconds"] = d.Seconds()
	return c
}

// IsPaused sets whether the reward is paused.
func (c *CustomRewardsUpdateCall) IsPaused(paused bool) *CustomRewardsUpdateCall {
	c.body["is_paused"] = paused
	return c
}

// ShouldRedemptionsSkipRequestQueue sets whether redemptions should skip the request queue.
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

	res, err := c.resource.client.DoRequest(ctx, http.MethodPatch, EndpointChannelPointsUpdateCustomReward, bytes.NewReader(bs), append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsUpdateResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

// CustomRewardsDeleteCall is a call to the custom rewards delete endpoint.
type CustomRewardsDeleteCall struct {
	resource *CustomRewardsResource
	opts     []RequestOption
}

// Delete creates a request to delete a custom channel point reward for a given broadcaster.
func (r *CustomRewardsResource) Delete(broadcasterID, id string) *CustomRewardsDeleteCall {
	c := &CustomRewardsDeleteCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	c.opts = append(c.opts, SetQueryParameter("id", id))
	return c
}

// Do executes the request.
func (c *CustomRewardsDeleteCall) Do(ctx context.Context, opts ...RequestOption) error {
	res, err := c.resource.client.DoRequest(ctx, http.MethodDelete, EndpointChannelPointsDeleteCustomReward, nil, append(opts, c.opts...)...)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[CustomRewardRedemption](res)
	return err
}

// CustomRewardsRedemptionResource provides methods for the Twitch Channel Points Custom Rewards Redemptions API.
type CustomRewardsRedemptionResource struct {
	client *Client
}

// NewCustomRewardsRedemptionResource creates a new CustomRewardsRedemptionResource.
func NewCustomRewardsRedemptionResource(client *Client) *CustomRewardsRedemptionResource {
	return &CustomRewardsRedemptionResource{client: client}
}

// CustomRewardsRedemptionListCall is a call to the custom rewards redemptions list endpoint.
type CustomRewardsRedemptionListCall struct {
	resource *CustomRewardsRedemptionResource
	opts     []RequestOption
}

// CustomRewardsRedemptionListResponse is the response from the custom rewards redemptions list endpoint.
type CustomRewardsRedemptionListResponse struct {
	Header http.Header
	Data   []CustomRewardRedemption
	Cursor string
}

// List creates a request to list custom channel point reward redemptions for a given broadcaster.
func (r *CustomRewardsRedemptionResource) List(broadcasterID, rewardID string) *CustomRewardsRedemptionListCall {
	c := &CustomRewardsRedemptionListCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	c.opts = append(c.opts, SetQueryParameter("reward_id", rewardID))
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
func (c *CustomRewardsRedemptionListCall) ID(ids []string) *CustomRewardsRedemptionListCall {
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
	res, err := c.resource.client.DoRequest(ctx, http.MethodGet, EndpointChannelPointsGetCustomRewardRedemptions, nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

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

// CustomRewardsRedemptionUpdateCall is a call to the custom rewards redemptions update endpoint.
type CustomRewardsRedemptionUpdateCall struct {
	resource *CustomRewardsRedemptionResource
	opts     []RequestOption
}

// CustomRewardsRedemptionUpdateResponse is the response from the custom rewards redemptions update endpoint.
type CustomRewardsRedemptionUpdateResponse struct {
	Header http.Header
	Data   []CustomRewardRedemption
}

// Update creates a request to update the status of one or more custom channel point reward redemptions for a given broadcaster and reward.
func (r *CustomRewardsRedemptionResource) Update(broadcasterID, rewardID string, id []string) *CustomRewardsRedemptionUpdateCall {
	c := &CustomRewardsRedemptionUpdateCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	c.opts = append(c.opts, SetQueryParameter("reward_id", rewardID))
	for _, id := range id {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

// Cancel sets the status of the redemption(s) to "CANCELED".
func (c *CustomRewardsRedemptionUpdateCall) Cancel() *CustomRewardsRedemptionUpdateCall {
	c.opts = append(c.opts, SetQueryParameter("status", "CANCELED"))
	return c
}

// Fulfill sets the status of the redemption(s) to "FULFILLED".
func (c *CustomRewardsRedemptionUpdateCall) Fulfill() *CustomRewardsRedemptionUpdateCall {
	c.opts = append(c.opts, SetQueryParameter("status", "FULFILLED"))
	return c
}

// Do executes the request.
func (c *CustomRewardsRedemptionUpdateCall) Do(ctx context.Context, opts ...RequestOption) (*CustomRewardsRedemptionUpdateResponse, error) {
	res, err := c.resource.client.DoRequest(ctx, http.MethodPatch, EndpointChannelPointsUpdateRedemptionStatus, nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomRewardRedemption](res)
	if err != nil {
		return nil, err
	}

	return &CustomRewardsRedemptionUpdateResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

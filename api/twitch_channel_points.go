package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ChannelPointsResource represents the Twitch ChannelPoints API.
type ChannelPointsResource struct {
	client *Client

	// Redemptions provides access to the Twitch Redemptions API.
	Redemptions *ChannelPointsRedemptionsResource
	// Rewards provides access to the Twitch Rewards API.
	Rewards *ChannelPointsRewardsResource
}

// NewChannelPointsResource creates a new ChannelPointsResource.
func NewChannelPointsResource(client *Client) *ChannelPointsResource {
	r := &ChannelPointsResource{client: client}
	r.Redemptions = NewChannelPointsRedemptionsResource(client)
	r.Rewards = NewChannelPointsRewardsResource(client)
	return r
}

// ChannelPointsRedemptionsResource represents the Twitch ChannelPointsRedemptions API.
type ChannelPointsRedemptionsResource struct {
	client *Client
}

// NewChannelPointsRedemptionsResource creates a new ChannelPointsRedemptionsResource.
func NewChannelPointsRedemptionsResource(client *Client) *ChannelPointsRedemptionsResource {
	return &ChannelPointsRedemptionsResource{client}
}

// ChannelPointRedemptionsListCall represents a GET call to a Twitch ChannelPointsRedemptions API endpoint.
type ChannelPointRedemptionsListCall struct {
	resource *ChannelPointsRedemptionsResource
	opts     []RequestOption
}

// ChannelPointRedemptionsListResponse represents the response from a GET request to /helix/channel_points/custom_rewards/redemptions.
type ChannelPointRedemptionsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CustomRewardRedemption data returned by the Twitch API.
	Data []CustomRewardRedemption
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channel_points/custom_rewards/redemptions.
//
// Gets a list of redemptions for the specified custom reward.
// The app used to create the reward is the only app that may get the redemptions.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:redemptions or channel:manage:redemptions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-custom-reward-redemption
func (r *ChannelPointsRedemptionsResource) List(broadcasterID string, rewardID string, status string) *ChannelPointRedemptionsListCall {
	c := &ChannelPointRedemptionsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		RewardID(rewardID).
		Status(status)
}

// ID sets the ID query parameter.
func (api *ChannelPointRedemptionsListCall) ID(id string) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelPointRedemptionsListCall) BroadcasterID(broadcasterID string) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// RewardID sets the RewardID query parameter.
func (api *ChannelPointRedemptionsListCall) RewardID(rewardID string) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("reward_id", rewardID))
	return api
}

// Status sets the Status query parameter.
func (api *ChannelPointRedemptionsListCall) Status(status string) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("status", status))
	return api
}

// Sort sets the Sort query parameter.
func (api *ChannelPointRedemptionsListCall) Sort(sort string) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("sort", sort))
	return api
}

// After sets the After query parameter.
func (api *ChannelPointRedemptionsListCall) After(after string) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ChannelPointRedemptionsListCall) First(first int) *ChannelPointRedemptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *ChannelPointRedemptionsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelPointRedemptionsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channel_points/custom_rewards/redemptions", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomRewardRedemption](res)
	if err != nil {
		return nil, err
	}

	return &ChannelPointRedemptionsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelPointRedemptionsModifyCall represents a PATCH call to a Twitch ChannelPointsRedemptions API endpoint.
type ChannelPointRedemptionsModifyCall struct {
	resource *ChannelPointsRedemptionsResource
	opts     []RequestOption
}

// ChannelPointRedemptionsModifyResponse represents the response from a PATCH request to /helix/channel_points/custom_rewards/redemptions.
type ChannelPointRedemptionsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CustomRewardRedemption data returned by the Twitch API.
	Data []CustomRewardRedemption
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/channel_points/custom_rewards/redemptions.
//
// Updates the status for a redemption. You may update a redemption only if its status is UNFULFILLED.
// The app used to create the reward is the only app that may update the redemption.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:redemptions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-redemption-status
func (r *ChannelPointsRedemptionsResource) Modify(id string, broadcasterID string, rewardID string) *ChannelPointRedemptionsModifyCall {
	c := &ChannelPointRedemptionsModifyCall{resource: r}
	return c.
		ID(id).
		BroadcasterID(broadcasterID).
		RewardID(rewardID)
}

// ID adds to the ID query parameter.
func (api *ChannelPointRedemptionsModifyCall) ID(ids ...string) *ChannelPointRedemptionsModifyCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelPointRedemptionsModifyCall) BroadcasterID(broadcasterID string) *ChannelPointRedemptionsModifyCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// RewardID sets the RewardID query parameter.
func (api *ChannelPointRedemptionsModifyCall) RewardID(rewardID string) *ChannelPointRedemptionsModifyCall {
	api.opts = append(api.opts, SetQueryParameter("reward_id", rewardID))
	return api
}

// Do executes the request.
func (api *ChannelPointRedemptionsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelPointRedemptionsModifyResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/channel_points/custom_rewards/redemptions", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomRewardRedemption](res)
	if err != nil {
		return nil, err
	}

	return &ChannelPointRedemptionsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelPointsRewardsResource represents the Twitch ChannelPointsRewards API.
type ChannelPointsRewardsResource struct {
	client *Client
}

// NewChannelPointsRewardsResource creates a new ChannelPointsRewardsResource.
func NewChannelPointsRewardsResource(client *Client) *ChannelPointsRewardsResource {
	return &ChannelPointsRewardsResource{client}
}

// ChannelPointRewardsInsertCall represents a POST call to a Twitch ChannelPointsRewards API endpoint.
type ChannelPointRewardsInsertCall struct {
	resource *ChannelPointsRewardsResource
	body     map[string]any
	opts     []RequestOption
}

// ChannelPointRewardsInsertResponse represents the response from a POST request to /helix/channel_points/custom_rewards.
type ChannelPointRewardsInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CustomReward data returned by the Twitch API.
	Data []CustomReward
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/channel_points/custom_rewards.
//
// Creates a new Custom Reward in a channel.
//
// The maximum number of custom rewards per channel is 50, which includes both enabled and disabled rewards.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:redemptions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-custom-rewards
func (r *ChannelPointsRewardsResource) Insert(broadcasterID string, title string, cost int64) *ChannelPointRewardsInsertCall {
	c := &ChannelPointRewardsInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		Title(title).
		Cost(cost)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelPointRewardsInsertCall) BroadcasterID(broadcasterID string) *ChannelPointRewardsInsertCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Title sets the Title body parameter.
func (api *ChannelPointRewardsInsertCall) Title(title string) *ChannelPointRewardsInsertCall {
	api.body["title"] = title
	return api
}

// Prompt sets the Prompt body parameter.
func (api *ChannelPointRewardsInsertCall) Prompt(prompt string) *ChannelPointRewardsInsertCall {
	api.body["prompt"] = prompt
	return api
}

// BackgroundColor sets the BackgroundColor body parameter.
func (api *ChannelPointRewardsInsertCall) BackgroundColor(backgroundColor string) *ChannelPointRewardsInsertCall {
	api.body["prompt"] = backgroundColor
	return api
}

// Cost sets the Cost body parameter.
func (api *ChannelPointRewardsInsertCall) Cost(cost int64) *ChannelPointRewardsInsertCall {
	api.body["cost"] = cost
	return api
}

// MaxPerStream sets the MaxPerStream body parameter.
func (api *ChannelPointRewardsInsertCall) MaxPerStream(maxPerStream int) *ChannelPointRewardsInsertCall {
	api.body["max_per_stream"] = maxPerStream
	return api
}

// MaxPerUserPerStream sets the MaxPerUserPerStream body parameter.
func (api *ChannelPointRewardsInsertCall) MaxPerUserPerStream(maxPerUserPerStream int) *ChannelPointRewardsInsertCall {
	api.body["max_per_user_per_stream"] = maxPerUserPerStream
	return api
}

// GlobalCooldownSeconds sets the GlobalCooldownSeconds body parameter.
func (api *ChannelPointRewardsInsertCall) GlobalCooldownSeconds(globalCooldownSeconds int) *ChannelPointRewardsInsertCall {
	api.body["global_cooldown_seconds"] = globalCooldownSeconds
	return api
}

// IsEnabled sets the IsEnabled body parameter.
func (api *ChannelPointRewardsInsertCall) IsEnabled(isEnabled bool) *ChannelPointRewardsInsertCall {
	api.body["is_enabled"] = isEnabled
	return api
}

// IsUserInputRequired sets the IsUserInputRequired body parameter.
func (api *ChannelPointRewardsInsertCall) IsUserInputRequired(isUserInputRequired bool) *ChannelPointRewardsInsertCall {
	api.body["is_user_input_required"] = isUserInputRequired
	return api
}

// IsMaxPerStreamEnabled sets the IsMaxPerStreamEnabled body parameter.
func (api *ChannelPointRewardsInsertCall) IsMaxPerStreamEnabled(isMaxPerStreamEnabled bool) *ChannelPointRewardsInsertCall {
	api.body["is_max_per_stream_enabled"] = isMaxPerStreamEnabled
	return api
}

// IsMaxPerUserPerStreamEnabled sets the IsMaxPerUserPerStreamEnabled body parameter.
func (api *ChannelPointRewardsInsertCall) IsMaxPerUserPerStreamEnabled(isMaxPerUserPerStreamEnabled bool) *ChannelPointRewardsInsertCall {
	api.body["is_max_per_user_per_stream_enabled"] = isMaxPerUserPerStreamEnabled
	return api
}

// IsGlobalCooldownEnabled sets the IsGlobalCooldownEnabled body parameter.
func (api *ChannelPointRewardsInsertCall) IsGlobalCooldownEnabled(isGlobalCooldownEnabled bool) *ChannelPointRewardsInsertCall {
	api.body["is_global_cooldown_enabled"] = isGlobalCooldownEnabled
	return api
}

// ShouldRedemptionsSkipRequestQueue sets the ShouldRedemptionsSkipRequestQueue body parameter.
func (api *ChannelPointRewardsInsertCall) ShouldRedemptionsSkipRequestQueue(shouldRedemptionsSkipRequestQueue bool) *ChannelPointRewardsInsertCall {
	api.body["should_redemptions_skip_request_queue"] = shouldRedemptionsSkipRequestQueue
	return api
}

// Do executes the request.
func (api *ChannelPointRewardsInsertCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelPointRewardsInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/channel_points/custom_rewards", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &ChannelPointRewardsInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelPointRewardsDeleteCall represents a DELETE call to a Twitch ChannelPointsRewards API endpoint.
type ChannelPointRewardsDeleteCall struct {
	resource *ChannelPointsRewardsResource
	opts     []RequestOption
}

// ChannelPointRewardsDeleteResponse represents the response from a DELETE request to /helix/channel_points/custom_rewards.
type ChannelPointRewardsDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/channel_points/custom_rewards.
//
// Deletes a custom reward that the broadcaster created.
//
// The app used to create the reward is the only app that may delete it.
// If the redemption status for the reward is UNFULFILLED at the time the reward is deleted, its redemption status is marked as FULFILLED.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:redemptions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#delete-custom-reward
func (r *ChannelPointsRewardsResource) Delete(id string, broadcasterID string) *ChannelPointRewardsDeleteCall {
	c := &ChannelPointRewardsDeleteCall{resource: r}
	return c.
		ID(id).
		BroadcasterID(broadcasterID)
}

// ID sets the ID query parameter.
func (api *ChannelPointRewardsDeleteCall) ID(id string) *ChannelPointRewardsDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelPointRewardsDeleteCall) BroadcasterID(broadcasterID string) *ChannelPointRewardsDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *ChannelPointRewardsDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelPointRewardsDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/channel_points/custom_rewards", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &ChannelPointRewardsDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

// ChannelPointRewardsListCall represents a GET call to a Twitch ChannelPointsRewards API endpoint.
type ChannelPointRewardsListCall struct {
	resource *ChannelPointsRewardsResource
	opts     []RequestOption
}

// ChannelPointRewardsListResponse represents the response from a GET request to /helix/channel_points/custom_rewards.
type ChannelPointRewardsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CustomReward data returned by the Twitch API.
	Data []CustomReward
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channel_points/custom_rewards.
//
// Gets a list of custom rewards that the specified broadcaster created.
//
// A channel may offer a maximum of 50 rewards, which includes both enabled and disabled rewards.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:redemptions or channel:manage:redemptions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-custom-reward
func (r *ChannelPointsRewardsResource) List(broadcasterID string) *ChannelPointRewardsListCall {
	c := &ChannelPointRewardsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// ID adds to the ID query parameter.
func (api *ChannelPointRewardsListCall) ID(ids ...string) *ChannelPointRewardsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelPointRewardsListCall) BroadcasterID(broadcasterID string) *ChannelPointRewardsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// OnlyManageableRewards sets the OnlyManageableRewards query parameter.
func (api *ChannelPointRewardsListCall) OnlyManageableRewards(onlyManageableRewards bool) *ChannelPointRewardsListCall {
	api.opts = append(api.opts, SetQueryParameter("only_manageable_rewards", onlyManageableRewards))
	return api
}

// Do executes the request.
func (api *ChannelPointRewardsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelPointRewardsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channel_points/custom_rewards", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &ChannelPointRewardsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelPointRewardsModifyCall represents a PATCH call to a Twitch ChannelPointsRewards API endpoint.
type ChannelPointRewardsModifyCall struct {
	resource *ChannelPointsRewardsResource
	body     map[string]any
	opts     []RequestOption
}

// ChannelPointRewardsModifyResponse represents the response from a PATCH request to /helix/channel_points/custom_rewards.
type ChannelPointRewardsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CustomReward data returned by the Twitch API.
	Data []CustomReward
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/channel_points/custom_rewards.
//
// Updates a custom reward. The app used to create the reward is the only app that may update the reward.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:redemptions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-custom-reward
func (r *ChannelPointsRewardsResource) Modify(id string, broadcasterID string) *ChannelPointRewardsModifyCall {
	c := &ChannelPointRewardsModifyCall{resource: r, body: make(map[string]any)}
	return c.
		ID(id).
		BroadcasterID(broadcasterID)
}

// ID sets the ID query parameter.
func (api *ChannelPointRewardsModifyCall) ID(id string) *ChannelPointRewardsModifyCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelPointRewardsModifyCall) BroadcasterID(broadcasterID string) *ChannelPointRewardsModifyCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Title sets the Title body parameter.
func (api *ChannelPointRewardsModifyCall) Title(title string) *ChannelPointRewardsModifyCall {
	api.body["title"] = title
	return api
}

// Prompt sets the Prompt body parameter.
func (api *ChannelPointRewardsModifyCall) Prompt(prompt string) *ChannelPointRewardsModifyCall {
	api.body["prompt"] = prompt
	return api
}

// BackgroundColor sets the BackgroundColor body parameter.
func (api *ChannelPointRewardsModifyCall) BackgroundColor(backgroundColor string) *ChannelPointRewardsModifyCall {
	api.body["backgroundColor"] = backgroundColor
	return api
}

// Cost sets the Cost body parameter.
func (api *ChannelPointRewardsModifyCall) Cost(cost int) *ChannelPointRewardsModifyCall {
	api.body["cost"] = cost
	return api
}

// MaxPerStream sets the MaxPerStream body parameter.
func (api *ChannelPointRewardsModifyCall) MaxPerStream(maxPerStream int) *ChannelPointRewardsModifyCall {
	api.body["maxPerStream"] = maxPerStream
	return api
}

// MaxPerUserPerStream sets the MaxPerUserPerStream body parameter.
func (api *ChannelPointRewardsModifyCall) MaxPerUserPerStream(maxPerUserPerStream int) *ChannelPointRewardsModifyCall {
	api.body["maxPerUserPerStream"] = maxPerUserPerStream
	return api
}

// GlobalCooldownSeconds sets the GlobalCooldownSeconds body parameter.
func (api *ChannelPointRewardsModifyCall) GlobalCooldownSeconds(globalCooldownSeconds int) *ChannelPointRewardsModifyCall {
	api.body["globalCooldownSeconds"] = globalCooldownSeconds
	return api
}

// IsPaused sets the IsPaused body parameter.
func (api *ChannelPointRewardsModifyCall) IsPaused(isPaused bool) *ChannelPointRewardsModifyCall {
	api.body["isPaused"] = isPaused
	return api
}

// IsEnabled sets the IsEnabled body parameter.
func (api *ChannelPointRewardsModifyCall) IsEnabled(isEnabled bool) *ChannelPointRewardsModifyCall {
	api.body["isEnabled"] = isEnabled
	return api
}

// IsUserInputRequired sets the IsUserInputRequired body parameter.
func (api *ChannelPointRewardsModifyCall) IsUserInputRequired(isUserInputRequired bool) *ChannelPointRewardsModifyCall {
	api.body["isUserInputRequired"] = isUserInputRequired
	return api
}

// IsMaxPerStreamEnabled sets the IsMaxPerStreamEnabled body parameter.
func (api *ChannelPointRewardsModifyCall) IsMaxPerStreamEnabled(isMaxPerStreamEnabled bool) *ChannelPointRewardsModifyCall {
	api.body["isMaxPerStreamEnabled"] = isMaxPerStreamEnabled
	return api
}

// IsMaxPerUserPerStreamEnabled sets the IsMaxPerUserPerStreamEnabled body parameter.
func (api *ChannelPointRewardsModifyCall) IsMaxPerUserPerStreamEnabled(isMaxPerUserPerStreamEnabled bool) *ChannelPointRewardsModifyCall {
	api.body["isMaxPerUserPerStreamEnabled"] = isMaxPerUserPerStreamEnabled
	return api
}

// IsGlobalCooldownEnabled sets the IsGlobalCooldownEnabled body parameter.
func (api *ChannelPointRewardsModifyCall) IsGlobalCooldownEnabled(isGlobalCooldownEnabled bool) *ChannelPointRewardsModifyCall {
	api.body["isGlobalCooldownEnabled"] = isGlobalCooldownEnabled
	return api
}

// ShouldRedemptionsSkipRequestQueue sets the ShouldRedemptionsSkipRequestQueue body parameter.
func (api *ChannelPointRewardsModifyCall) ShouldRedemptionsSkipRequestQueue(shouldRedemptionsSkipRequestQueue bool) *ChannelPointRewardsModifyCall {
	api.body["shouldRedemptionsSkipRequestQueue"] = shouldRedemptionsSkipRequestQueue
	return api
}

// Do executes the request.
func (api *ChannelPointRewardsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelPointRewardsModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/channel_points/custom_rewards", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CustomReward](res)
	if err != nil {
		return nil, err
	}

	return &ChannelPointRewardsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ChannelsResource represents the Twitch Channels API.
type ChannelsResource struct {
	client *Client

	// Editors provides access to the Twitch Editors API.
	Editors *ChannelsEditorsResource
	// Followed provides access to the Twitch Followed API.
	Followed *ChannelsFollowedResource
	// Followers provides access to the Twitch Followers API.
	Followers *ChannelsFollowersResource
}

// NewChannelsResource creates a new ChannelsResource.
func NewChannelsResource(client *Client) *ChannelsResource {
	r := &ChannelsResource{client: client}
	r.Editors = NewChannelsEditorsResource(client)
	r.Followed = NewChannelsFollowedResource(client)
	r.Followers = NewChannelsFollowersResource(client)
	return r
}

// ChannelsEditorsResource represents the Twitch ChannelsEditors API.
type ChannelsEditorsResource struct {
	client *Client
}

// NewChannelsEditorsResource creates a new ChannelsEditorsResource.
func NewChannelsEditorsResource(client *Client) *ChannelsEditorsResource {
	return &ChannelsEditorsResource{client}
}

// ChannelEditorsListCall represents a GET call to a Twitch ChannelsEditors API endpoint.
type ChannelEditorsListCall struct {
	resource *ChannelsEditorsResource
	opts     []RequestOption
}

// ChannelEditorsListResponse represents the response from a GET request to /helix/channels/editors.
type ChannelEditorsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ChannelEditor data returned by the Twitch API.
	Data []ChannelEditor
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channels/editors.
//
// Gets the list of editors for a broadcaster.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:editors scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-editors
func (r *ChannelsEditorsResource) List(broadcasterID string) *ChannelEditorsListCall {
	c := &ChannelEditorsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelEditorsListCall) BroadcasterID(broadcasterID string) *ChannelEditorsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *ChannelEditorsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelEditorsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channels/editors", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChannelEditor](res)
	if err != nil {
		return nil, err
	}

	return &ChannelEditorsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelsFollowedResource represents the Twitch ChannelsFollowed API.
type ChannelsFollowedResource struct {
	client *Client
}

// NewChannelsFollowedResource creates a new ChannelsFollowedResource.
func NewChannelsFollowedResource(client *Client) *ChannelsFollowedResource {
	return &ChannelsFollowedResource{client}
}

// ChannelsFollowedListCall represents a GET call to a Twitch ChannelsFollowed API endpoint.
type ChannelsFollowedListCall struct {
	resource *ChannelsFollowedResource
	opts     []RequestOption
}

// ChannelsFollowedListResponse represents the response from a GET request to /helix/channels/followed.
type ChannelsFollowedListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Total is the int data returned by the Twitch API.
	Total int
	// Data is the Followed data returned by the Twitch API.
	Data []Followed
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channels/followed.
//
// Gets a list of broadcasters that the specified user follows.
// You can also use this endpoint to see whether a user follows a specific broadcaster.
//
// # Authorization
//
// Requires a user access token that includes the user:read:follows scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-followed-channels
func (r *ChannelsFollowedResource) List(userID string) *ChannelsFollowedListCall {
	c := &ChannelsFollowedListCall{resource: r}
	return c.
		UserID(userID)
}

// UserID sets the UserID query parameter.
func (api *ChannelsFollowedListCall) UserID(userID string) *ChannelsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelsFollowedListCall) BroadcasterID(broadcasterID string) *ChannelsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// After sets the After query parameter.
func (api *ChannelsFollowedListCall) After(after string) *ChannelsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ChannelsFollowedListCall) First(first int) *ChannelsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *ChannelsFollowedListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelsFollowedListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channels/followed", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Followed](res)
	if err != nil {
		return nil, err
	}

	return &ChannelsFollowedListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Total:      data.Total,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// ChannelsFollowersResource represents the Twitch ChannelsFollowers API.
type ChannelsFollowersResource struct {
	client *Client
}

// NewChannelsFollowersResource creates a new ChannelsFollowersResource.
func NewChannelsFollowersResource(client *Client) *ChannelsFollowersResource {
	return &ChannelsFollowersResource{client}
}

// ChannelFollowersListCall represents a GET call to a Twitch ChannelsFollowers API endpoint.
type ChannelFollowersListCall struct {
	resource *ChannelsFollowersResource
	opts     []RequestOption
}

// ChannelFollowersListResponse represents the response from a GET request to /helix/channels/followers.
type ChannelFollowersListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Total is the int data returned by the Twitch API.
	Total int
	// Data is the Follower data returned by the Twitch API.
	Data []Follower
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channels/followers.
//
// Gets a list of users that follow the specified broadcaster.
// You can also use this endpoint to see whether a specific user follows the broadcaster.
//
// # Authorization
//
// Requires a user access token that includes the moderator:read:followers scope.
//
// The ID in the broadcaster_id query parameter must match the user ID in the access token or the user ID in the access token must be a moderator for the specified broadcaster.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-followers
func (r *ChannelsFollowersResource) List(broadcasterID string) *ChannelFollowersListCall {
	c := &ChannelFollowersListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// UserID sets the UserID query parameter.
func (api *ChannelFollowersListCall) UserID(userID string) *ChannelFollowersListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelFollowersListCall) BroadcasterID(broadcasterID string) *ChannelFollowersListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// After sets the After query parameter.
func (api *ChannelFollowersListCall) After(after string) *ChannelFollowersListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ChannelFollowersListCall) First(first int) *ChannelFollowersListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *ChannelFollowersListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelFollowersListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channels/followers", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Follower](res)
	if err != nil {
		return nil, err
	}

	return &ChannelFollowersListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Total:      data.Total,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// ChannelInformationListCall represents a GET call to a Twitch Channels API endpoint.
type ChannelInformationListCall struct {
	resource *ChannelsResource
	opts     []RequestOption
}

// ChannelInformationListResponse represents the response from a GET request to /helix/channels.
type ChannelInformationListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Channel data returned by the Twitch API.
	Data []Channel
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/channels.
//
// Gets information about one or more channels.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-information
func (r *ChannelsResource) List(broadcasterID string) *ChannelInformationListCall {
	c := &ChannelInformationListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID adds to the BroadcasterID query parameter.
func (api *ChannelInformationListCall) BroadcasterID(broadcasterIDs ...string) *ChannelInformationListCall {
	for _, broadcasterID := range broadcasterIDs {
		api.opts = append(api.opts, AddQueryParameter("broadcaster_id", broadcasterID))
	}
	return api
}

// Do executes the request.
func (api *ChannelInformationListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelInformationListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/channels", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Channel](res)
	if err != nil {
		return nil, err
	}

	return &ChannelInformationListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelInformationModifyCall represents a PATCH call to a Twitch Channels API endpoint.
type ChannelInformationModifyCall struct {
	resource *ChannelsResource
	body     map[string]any
	opts     []RequestOption
}

// ChannelInformationModifyResponse represents the response from a PATCH request to /helix/channels.
type ChannelInformationModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/channels.
//
// Updates properties for a channel.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:broadcast scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#modify-channel-information
func (r *ChannelsResource) Modify(broadcasterID string) *ChannelInformationModifyCall {
	c := &ChannelInformationModifyCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelInformationModifyCall) BroadcasterID(broadcasterID string) *ChannelInformationModifyCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// GameID sets the GameID body parameter.
func (api *ChannelInformationModifyCall) GameID(gameID string) *ChannelInformationModifyCall {
	api.body["game_id"] = gameID
	return api
}

// BroadcasterLanguage sets the BroadcasterLanguage body parameter.
func (api *ChannelInformationModifyCall) BroadcasterLanguage(broadcasterLanguage string) *ChannelInformationModifyCall {
	api.body["broadcaster_language"] = broadcasterLanguage
	return api
}

// Title sets the Title body parameter.
func (api *ChannelInformationModifyCall) Title(title string) *ChannelInformationModifyCall {
	api.body["title"] = title
	return api
}

// Tags sets the Tags body parameter.
func (api *ChannelInformationModifyCall) Tags(tagss ...string) *ChannelInformationModifyCall {
	api.body["tags"] = tagss
	return api
}

// ContentClassificationLabels sets the ContentClassificationLabels body parameter.
func (api *ChannelInformationModifyCall) ContentClassificationLabels(contentClassificationLabelss ...ChannelContentClassificationLabel) *ChannelInformationModifyCall {
	api.body["content_classification_labels"] = contentClassificationLabelss
	return api
}

// Delay sets the Delay body parameter.
func (api *ChannelInformationModifyCall) Delay(delay int) *ChannelInformationModifyCall {
	api.body["delay"] = delay
	return api
}

// IsBrandedContent sets the IsBrandedContent body parameter.
func (api *ChannelInformationModifyCall) IsBrandedContent(isBrandedContent bool) *ChannelInformationModifyCall {
	api.body["is_branded_content"] = isBrandedContent
	return api
}

// Do executes the request.
func (api *ChannelInformationModifyCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelInformationModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/channels", bytes.NewReader(bs), append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &ChannelInformationModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ChatResource represents the Twitch Chat API.
type ChatResource struct {
	client *Client

	// Badges provides access to the Twitch Badges API.
	Badges *ChatBadgesResource
	// Chatters provides access to the Twitch Chatters API.
	Chatters *ChatChattersResource
	// Emotes provides access to the Twitch Emotes API.
	Emotes *ChatEmotesResource
	// EmoteSets provides access to the Twitch EmoteSets API.
	EmoteSets *ChatEmoteSetsResource
	// Settings provides access to the Twitch Settings API.
	Settings *ChatSettingsResource
	// Shared provides access to the Twitch Shared API.
	Shared *ChatSharedResource
	// Announcement provides access to the Twitch Announcement API.
	Announcement *ChatAnnouncementResource
	// Shoutout provides access to the Twitch Shoutout API.
	Shoutout *ChatShoutoutResource
}

// NewChatResource creates a new ChatResource.
func NewChatResource(client *Client) *ChatResource {
	r := &ChatResource{client: client}
	r.Badges = NewChatBadgesResource(client)
	r.Chatters = NewChatChattersResource(client)
	r.Emotes = NewChatEmotesResource(client)
	r.EmoteSets = NewChatEmoteSetsResource(client)
	r.Settings = NewChatSettingsResource(client)
	r.Shared = NewChatSharedResource(client)
	r.Announcement = NewChatAnnouncementResource(client)
	r.Shoutout = NewChatShoutoutResource(client)
	return r
}

// ChatBadgesResource represents the Twitch ChatBadges API.
type ChatBadgesResource struct {
	client *Client

	// Global provides access to the Twitch Global API.
	Global *ChatBadgesGlobalResource
}

// NewChatBadgesResource creates a new ChatBadgesResource.
func NewChatBadgesResource(client *Client) *ChatBadgesResource {
	r := &ChatBadgesResource{client: client}
	r.Global = NewChatBadgesGlobalResource(client)
	return r
}

// ChatBadgesGlobalResource represents the Twitch ChatBadgesGlobal API.
type ChatBadgesGlobalResource struct {
	client *Client
}

// NewChatBadgesGlobalResource creates a new ChatBadgesGlobalResource.
func NewChatBadgesGlobalResource(client *Client) *ChatBadgesGlobalResource {
	return &ChatBadgesGlobalResource{client}
}

// GlobalChatBadgesListCall represents a GET call to a Twitch ChatBadgesGlobal API endpoint.
type GlobalChatBadgesListCall struct {
	resource *ChatBadgesGlobalResource
}

// GlobalChatBadgesListResponse represents the response from a GET request to /helix/chat/badges/global.
type GlobalChatBadgesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ChatBadge data returned by the Twitch API.
	Data []ChatBadge
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/badges/global.
//
// Gets list of chat badges on Twitch, which users may use in any chat room.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-global-chat-badges
func (r *ChatBadgesGlobalResource) List() *GlobalChatBadgesListCall {
	return &GlobalChatBadgesListCall{resource: r}
}

// Do executes the request.
func (api *GlobalChatBadgesListCall) Do(ctx context.Context, opts ...RequestOption) (*GlobalChatBadgesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/badges/global", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChatBadge](res)
	if err != nil {
		return nil, err
	}

	return &GlobalChatBadgesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChannelChatBadgesListCall represents a GET call to a Twitch ChatBadges API endpoint.
type ChannelChatBadgesListCall struct {
	resource *ChatBadgesResource
	opts     []RequestOption
}

// ChannelChatBadgesListResponse represents the response from a GET request to /helix/chat/badges.
type ChannelChatBadgesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ChatBadge data returned by the Twitch API.
	Data []ChatBadge
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/badges.
//
// Gets the list of custom chat badges for a broadcaster.
//
// The list is empty if the broadcaster hasnt created custom chat badges.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-chat-badges
func (r *ChatBadgesResource) List(broadcasterID string) *ChannelChatBadgesListCall {
	c := &ChannelChatBadgesListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelChatBadgesListCall) BroadcasterID(broadcasterID string) *ChannelChatBadgesListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *ChannelChatBadgesListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelChatBadgesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/badges", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChatBadge](res)
	if err != nil {
		return nil, err
	}

	return &ChannelChatBadgesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatChattersResource represents the Twitch ChatChatters API.
type ChatChattersResource struct {
	client *Client

	// User provides access to the Twitch User API.
	User *ChatChattersUserResource
}

// NewChatChattersResource creates a new ChatChattersResource.
func NewChatChattersResource(client *Client) *ChatChattersResource {
	r := &ChatChattersResource{client: client}
	r.User = NewChatChattersUserResource(client)
	return r
}

// ChatChattersUserResource represents the Twitch ChatChattersUser API.
type ChatChattersUserResource struct {
	client *Client

	// Color provides access to the Twitch Color API.
	Color *ChatChattersUserColorResource
}

// NewChatChattersUserResource creates a new ChatChattersUserResource.
func NewChatChattersUserResource(client *Client) *ChatChattersUserResource {
	r := &ChatChattersUserResource{client: client}
	r.Color = NewChatChattersUserColorResource(client)
	return r
}

// ChatChattersUserColorResource represents the Twitch ChatChattersUserColor API.
type ChatChattersUserColorResource struct {
	client *Client
}

// NewChatChattersUserColorResource creates a new ChatChattersUserColorResource.
func NewChatChattersUserColorResource(client *Client) *ChatChattersUserColorResource {
	return &ChatChattersUserColorResource{client}
}

// UserChatColorListCall represents a GET call to a Twitch ChatChattersUserColor API endpoint.
type UserChatColorListCall struct {
	resource *ChatChattersUserColorResource
	opts     []RequestOption
}

// UserChatColorListResponse represents the response from a GET request to /helix/chat/color.
type UserChatColorListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the UserChatColor data returned by the Twitch API.
	Data []UserChatColor
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/color.
//
// Gets the chat color for a user.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-user-chat-color
func (r *ChatChattersUserColorResource) List(userID string) *UserChatColorListCall {
	c := &UserChatColorListCall{resource: r}
	return c.
		UserID(userID)
}

// UserID adds to the UserID query parameter.
func (api *UserChatColorListCall) UserID(userIDs ...string) *UserChatColorListCall {
	for _, userID := range userIDs {
		api.opts = append(api.opts, AddQueryParameter("user_id", userID))
	}
	return api
}

// Do executes the request.
func (api *UserChatColorListCall) Do(ctx context.Context, opts ...RequestOption) (*UserChatColorListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/color", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[UserChatColor](res)
	if err != nil {
		return nil, err
	}

	return &UserChatColorListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// UserChatColorUpdateCall represents a PUT call to a Twitch ChatChattersUserColor API endpoint.
type UserChatColorUpdateCall struct {
	resource *ChatChattersUserColorResource
	opts     []RequestOption
}

// UserChatColorUpdateResponse represents the response from a PUT request to /helix/chat/color.
type UserChatColorUpdateResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Update creates a new PUT request to /helix/chat/color.
//
// Updates the display color for the user in chat.
//
// # Authorization
//
// Requires a user access token that includes the user:manage:chat_color scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-user-chat-color
func (r *ChatChattersUserColorResource) Update(userID string, color string) *UserChatColorUpdateCall {
	c := &UserChatColorUpdateCall{resource: r}
	return c.
		UserID(userID).
		Color(color)
}

// UserID sets the UserID query parameter.
func (api *UserChatColorUpdateCall) UserID(userID string) *UserChatColorUpdateCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// Color sets the Color query parameter.
func (api *UserChatColorUpdateCall) Color(color string) *UserChatColorUpdateCall {
	api.opts = append(api.opts, SetQueryParameter("color", color))
	return api
}

// Do executes the request.
func (api *UserChatColorUpdateCall) Do(ctx context.Context, opts ...RequestOption) (*UserChatColorUpdateResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "PUT", "/helix/chat/color", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &UserChatColorUpdateResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

// ChatChattersListCall represents a GET call to a Twitch ChatChatters API endpoint.
type ChatChattersListCall struct {
	resource *ChatChattersResource
	opts     []RequestOption
}

// ChatChattersListResponse represents the response from a GET request to /helix/chat/chatters.
type ChatChattersListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Total is the int data returned by the Twitch API.
	Total int
	// Data is the UserInfo data returned by the Twitch API.
	Data []UserInfo
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/chatters.
//
// Gets the list of users that are connected to a chat session.
//
// To determine whether a user is a moderator or VIP, use the Get Moderators and Get VIPs endpoints.
// You can check the roles of up to 100 users.
//
// # Authorization
//
// Requires a user access token that includes the moderator:read:chatters scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-chatters
func (r *ChatChattersResource) List(broadcasterID string, moderatorID string) *ChatChattersListCall {
	c := &ChatChattersListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChatChattersListCall) BroadcasterID(broadcasterID string) *ChatChattersListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *ChatChattersListCall) ModeratorID(moderatorID string) *ChatChattersListCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// After sets the After query parameter.
func (api *ChatChattersListCall) After(after string) *ChatChattersListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ChatChattersListCall) First(first int) *ChatChattersListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *ChatChattersListCall) Do(ctx context.Context, opts ...RequestOption) (*ChatChattersListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/chatters", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[UserInfo](res)
	if err != nil {
		return nil, err
	}

	return &ChatChattersListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Total:      data.Total,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// ChatEmotesResource represents the Twitch ChatEmotes API.
type ChatEmotesResource struct {
	client *Client

	// Channel provides access to the Twitch Channel API.
	Channel *ChatEmotesChannelResource
	// Global provides access to the Twitch Global API.
	Global *ChatEmotesGlobalResource
	// User provides access to the Twitch User API.
	User *ChatEmotesUserResource
}

// NewChatEmotesResource creates a new ChatEmotesResource.
func NewChatEmotesResource(client *Client) *ChatEmotesResource {
	r := &ChatEmotesResource{client: client}
	r.Channel = NewChatEmotesChannelResource(client)
	r.Global = NewChatEmotesGlobalResource(client)
	r.User = NewChatEmotesUserResource(client)
	return r
}

// ChatEmotesChannelResource represents the Twitch ChatEmotesChannel API.
type ChatEmotesChannelResource struct {
	client *Client
}

// NewChatEmotesChannelResource creates a new ChatEmotesChannelResource.
func NewChatEmotesChannelResource(client *Client) *ChatEmotesChannelResource {
	return &ChatEmotesChannelResource{client}
}

// ChannelEmotesListCall represents a GET call to a Twitch ChatEmotesChannel API endpoint.
type ChannelEmotesListCall struct {
	resource *ChatEmotesChannelResource
	opts     []RequestOption
}

// ChannelEmotesListResponse represents the response from a GET request to /helix/chat/emotes.
type ChannelEmotesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Emote data returned by the Twitch API.
	Data []Emote
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/emotes.
//
// Gets the list of custom emotes for a broadcaster.
// Broadcasters create these custom emotes for users who subscribe to or follow the channel or cheer Bits in the chat window.
//
// With the exception of custom follower emotes, users may use custom emotes in any Twitch chat.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-emotes
func (r *ChatEmotesChannelResource) List(broadcasterID string) *ChannelEmotesListCall {
	c := &ChannelEmotesListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelEmotesListCall) BroadcasterID(broadcasterID string) *ChannelEmotesListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *ChannelEmotesListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelEmotesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/emotes", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Emote](res)
	if err != nil {
		return nil, err
	}

	return &ChannelEmotesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatEmotesGlobalResource represents the Twitch ChatEmotesGlobal API.
type ChatEmotesGlobalResource struct {
	client *Client
}

// NewChatEmotesGlobalResource creates a new ChatEmotesGlobalResource.
func NewChatEmotesGlobalResource(client *Client) *ChatEmotesGlobalResource {
	return &ChatEmotesGlobalResource{client}
}

// GlobalEmotesListCall represents a GET call to a Twitch ChatEmotesGlobal API endpoint.
type GlobalEmotesListCall struct {
	resource *ChatEmotesGlobalResource
}

// GlobalEmotesListResponse represents the response from a GET request to /helix/chat/emotes/global.
type GlobalEmotesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Emote data returned by the Twitch API.
	Data []Emote
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/emotes/global.
//
// Gets the list of global emotes.
// Global emotes are Twitch-created emotes that users can use in any Twitch chat.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-global-emotes
func (r *ChatEmotesGlobalResource) List() *GlobalEmotesListCall {
	return &GlobalEmotesListCall{resource: r}
}

// Do executes the request.
func (api *GlobalEmotesListCall) Do(ctx context.Context, opts ...RequestOption) (*GlobalEmotesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/emotes/global", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Emote](res)
	if err != nil {
		return nil, err
	}

	return &GlobalEmotesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatEmotesUserResource represents the Twitch ChatEmotesUser API.
type ChatEmotesUserResource struct {
	client *Client
}

// NewChatEmotesUserResource creates a new ChatEmotesUserResource.
func NewChatEmotesUserResource(client *Client) *ChatEmotesUserResource {
	return &ChatEmotesUserResource{client}
}

// UserEmotesListCall represents a GET call to a Twitch ChatEmotesUser API endpoint.
type UserEmotesListCall struct {
	resource *ChatEmotesUserResource
	opts     []RequestOption
}

// UserEmotesListResponse represents the response from a GET request to /helix/chat/emotes/user.
type UserEmotesListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Emote data returned by the Twitch API.
	Data []Emote
	// Template is the string data returned by the Twitch API.
	Template string
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/emotes/user.
//
// Retrieves emotes available to the user across all channels.
//
// # Authorization
//
// Requires a user access token that includes the user:read:emotes scope.
//
// Query parameter user_id must match the user_id in the user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-user-emotes
func (r *ChatEmotesUserResource) List(userID string) *UserEmotesListCall {
	c := &UserEmotesListCall{resource: r}
	return c.
		UserID(userID)
}

// UserID sets the UserID query parameter.
func (api *UserEmotesListCall) UserID(userID string) *UserEmotesListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *UserEmotesListCall) BroadcasterID(broadcasterID string) *UserEmotesListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// After sets the After query parameter.
func (api *UserEmotesListCall) After(after string) *UserEmotesListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// Do executes the request.
func (api *UserEmotesListCall) Do(ctx context.Context, opts ...RequestOption) (*UserEmotesListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/emotes/user", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Emote](res)
	if err != nil {
		return nil, err
	}

	return &UserEmotesListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// ChatEmoteSetsResource represents the Twitch ChatEmoteSets API.
type ChatEmoteSetsResource struct {
	client *Client
}

// NewChatEmoteSetsResource creates a new ChatEmoteSetsResource.
func NewChatEmoteSetsResource(client *Client) *ChatEmoteSetsResource {
	return &ChatEmoteSetsResource{client}
}

// EmoteSetsListCall represents a GET call to a Twitch ChatEmoteSets API endpoint.
type EmoteSetsListCall struct {
	resource *ChatEmoteSetsResource
	opts     []RequestOption
}

// EmoteSetsListResponse represents the response from a GET request to /helix/chat/emotes/set.
type EmoteSetsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Emote data returned by the Twitch API.
	Data []Emote
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/emotes/set.
//
// Gets emotes for one or more specified emote sets.
//
// An emote set groups emotes that have a similar context.
// For example, Twitch places all the subscriber emotes that a broadcaster uploads for their channel in the same emote set.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-emote-sets
func (r *ChatEmoteSetsResource) List(emoteSetID string) *EmoteSetsListCall {
	c := &EmoteSetsListCall{resource: r}
	return c.
		EmoteSetID(emoteSetID)
}

// EmoteSetID adds to the EmoteSetID query parameter.
func (api *EmoteSetsListCall) EmoteSetID(emoteSetIDs ...string) *EmoteSetsListCall {
	for _, emoteSetID := range emoteSetIDs {
		api.opts = append(api.opts, AddQueryParameter("emote_set_id", emoteSetID))
	}
	return api
}

// Do executes the request.
func (api *EmoteSetsListCall) Do(ctx context.Context, opts ...RequestOption) (*EmoteSetsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/emotes/set", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Emote](res)
	if err != nil {
		return nil, err
	}

	return &EmoteSetsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatSettingsResource represents the Twitch ChatSettings API.
type ChatSettingsResource struct {
	client *Client
}

// NewChatSettingsResource creates a new ChatSettingsResource.
func NewChatSettingsResource(client *Client) *ChatSettingsResource {
	return &ChatSettingsResource{client}
}

// ChatSettingsListCall represents a GET call to a Twitch ChatSettings API endpoint.
type ChatSettingsListCall struct {
	resource *ChatSettingsResource
	opts     []RequestOption
}

// ChatSettingsListResponse represents the response from a GET request to /helix/chat/settings.
type ChatSettingsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ChatSettings data returned by the Twitch API.
	Data []ChatSettings
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/chat/settings.
//
// Gets the chat settings for a channel.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-chat-settings
func (r *ChatSettingsResource) List(broadcasterID string, moderatorID string) *ChatSettingsListCall {
	c := &ChatSettingsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChatSettingsListCall) BroadcasterID(broadcasterID string) *ChatSettingsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *ChatSettingsListCall) ModeratorID(moderatorID string) *ChatSettingsListCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// Do executes the request.
func (api *ChatSettingsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChatSettingsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/chat/settings", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChatSettings](res)
	if err != nil {
		return nil, err
	}

	return &ChatSettingsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatSettingsModifyCall represents a PATCH call to a Twitch ChatSettings API endpoint.
type ChatSettingsModifyCall struct {
	resource *ChatSettingsResource
	body     map[string]any
	opts     []RequestOption
}

// ChatSettingsModifyResponse represents the response from a PATCH request to /helix/chat/settings.
type ChatSettingsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ChatSettings data returned by the Twitch API.
	Data []ChatSettings
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/chat/settings.
//
// Updates the chat settings for a channel.
//
// # Authorization
//
// Requires a user access token that includes the moderator:manage:chat_settings scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#update-chat-settings
func (r *ChatSettingsResource) Modify(broadcasterID string, moderatorID string) *ChatSettingsModifyCall {
	c := &ChatSettingsModifyCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChatSettingsModifyCall) BroadcasterID(broadcasterID string) *ChatSettingsModifyCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *ChatSettingsModifyCall) ModeratorID(moderatorID string) *ChatSettingsModifyCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// EmoteMode sets the EmoteMode body parameter.
func (api *ChatSettingsModifyCall) EmoteMode(emoteMode bool) *ChatSettingsModifyCall {
	api.body["emote_mode"] = emoteMode
	return api
}

// FollowerMode sets the FollowerMode body parameter.
func (api *ChatSettingsModifyCall) FollowerMode(followerMode bool) *ChatSettingsModifyCall {
	api.body["follower_mode"] = followerMode
	return api
}

// FollowerModeDuration sets the FollowerModeDuration body parameter.
func (api *ChatSettingsModifyCall) FollowerModeDuration(followerModeDuration int) *ChatSettingsModifyCall {
	api.body["follower_mode_duration"] = followerModeDuration
	return api
}

// NonModeratorChatDelay sets the NonModeratorChatDelay body parameter.
func (api *ChatSettingsModifyCall) NonModeratorChatDelay(nonModeratorChatDelay bool) *ChatSettingsModifyCall {
	api.body["non_moderator_chat_delay"] = nonModeratorChatDelay
	return api
}

// NonModeratorChatDelayDuration sets the NonModeratorChatDelayDuration body parameter.
func (api *ChatSettingsModifyCall) NonModeratorChatDelayDuration(nonModeratorChatDelayDuration int) *ChatSettingsModifyCall {
	api.body["non_moderator_chat_delay_duration"] = nonModeratorChatDelayDuration
	return api
}

// SlowMode sets the SlowMode body parameter.
func (api *ChatSettingsModifyCall) SlowMode(slowMode bool) *ChatSettingsModifyCall {
	api.body["slow_mode"] = slowMode
	return api
}

// SlowModeWaitTime sets the SlowModeWaitTime body parameter.
func (api *ChatSettingsModifyCall) SlowModeWaitTime(slowModeWaitTime int) *ChatSettingsModifyCall {
	api.body["slow_mode_wait_time"] = slowModeWaitTime
	return api
}

// SubscriberMode sets the SubscriberMode body parameter.
func (api *ChatSettingsModifyCall) SubscriberMode(subscriberMode bool) *ChatSettingsModifyCall {
	api.body["subscriber_mode"] = subscriberMode
	return api
}

// UniqueChatMode sets the UniqueChatMode body parameter.
func (api *ChatSettingsModifyCall) UniqueChatMode(uniqueChatMode bool) *ChatSettingsModifyCall {
	api.body["unique_chat_mode"] = uniqueChatMode
	return api
}

// Do executes the request.
func (api *ChatSettingsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*ChatSettingsModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/chat/settings", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChatSettings](res)
	if err != nil {
		return nil, err
	}

	return &ChatSettingsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatSharedResource represents the Twitch ChatShared API.
type ChatSharedResource struct {
	client *Client
}

// NewChatSharedResource creates a new ChatSharedResource.
func NewChatSharedResource(client *Client) *ChatSharedResource {
	return &ChatSharedResource{client}
}

// SharedChatSessionListCall represents a GET call to a Twitch ChatShared API endpoint.
type SharedChatSessionListCall struct {
	resource *ChatSharedResource
	opts     []RequestOption
}

// SharedChatSessionListResponse represents the response from a GET request to /helix/shared_chat/session.
type SharedChatSessionListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the SharedChatSession data returned by the Twitch API.
	Data []SharedChatSession
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/shared_chat/session.
//
// Retrieves the active shared chat session for a channel.
//
// # Authorization
//
// Requires a user access token that includes the moderator:read:chat_settings scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-shared-chat-session
func (r *ChatSharedResource) List(broadcasterID string) *SharedChatSessionListCall {
	c := &SharedChatSessionListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *SharedChatSessionListCall) BroadcasterID(broadcasterID string) *SharedChatSessionListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *SharedChatSessionListCall) Do(ctx context.Context, opts ...RequestOption) (*SharedChatSessionListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/shared_chat/session", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[SharedChatSession](res)
	if err != nil {
		return nil, err
	}

	return &SharedChatSessionListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ChatAnnouncementResource represents the Twitch ChatAnnouncement API.
type ChatAnnouncementResource struct {
	client *Client
}

// NewChatAnnouncementResource creates a new ChatAnnouncementResource.
func NewChatAnnouncementResource(client *Client) *ChatAnnouncementResource {
	return &ChatAnnouncementResource{client}
}

// ChatAnnouncementInsertCall represents a POST call to a Twitch ChatAnnouncement API endpoint.
type ChatAnnouncementInsertCall struct {
	resource *ChatAnnouncementResource
	body     map[string]any
	opts     []RequestOption
}

// ChatAnnouncementInsertResponse represents the response from a POST request to /helix/chat/announcements.
type ChatAnnouncementInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/chat/announcements.
//
// Sends an announcement to a chat room.
//
// # Rate Limits
//
// One announcement may be sent every 2 seconds.
//
// # Authorization
//
// Requires a user access token that includes the moderator:manage:announcements scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#send-chat-announcement
func (r *ChatAnnouncementResource) Insert(broadcasterID string, moderatorID string, message string) *ChatAnnouncementInsertCall {
	c := &ChatAnnouncementInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID).
		Message(message)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChatAnnouncementInsertCall) BroadcasterID(broadcasterID string) *ChatAnnouncementInsertCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *ChatAnnouncementInsertCall) ModeratorID(moderatorID string) *ChatAnnouncementInsertCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// Message sets the Message body parameter.
func (api *ChatAnnouncementInsertCall) Message(message string) *ChatAnnouncementInsertCall {
	api.body["message"] = message
	return api
}

// Color sets the Color body parameter.
func (api *ChatAnnouncementInsertCall) Color(color string) *ChatAnnouncementInsertCall {
	api.body["color"] = color
	return api
}

// Do executes the request.
func (api *ChatAnnouncementInsertCall) Do(ctx context.Context, opts ...RequestOption) (*ChatAnnouncementInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/chat/announcements", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &ChatAnnouncementInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

// ChatShoutoutResource represents the Twitch ChatShoutout API.
type ChatShoutoutResource struct {
	client *Client
}

// NewChatShoutoutResource creates a new ChatShoutoutResource.
func NewChatShoutoutResource(client *Client) *ChatShoutoutResource {
	return &ChatShoutoutResource{client}
}

// ChatShoutoutInsertCall represents a POST call to a Twitch ChatShoutout API endpoint.
type ChatShoutoutInsertCall struct {
	resource *ChatShoutoutResource
	opts     []RequestOption
}

// ChatShoutoutInsertResponse represents the response from a POST request to /helix/chat/shoutouts.
type ChatShoutoutInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/chat/shoutouts.
//
// Sends a shoutout in chat to another channel.
//
// # Rate Limits
//
// The broadcaster may send a Shoutout once every 2 minutes.
//
// They may send the same broadcaster a Shoutout once every 60 minutes.
//
// # Authorization
//
// Requires a user access token that includes the moderator:manage:shoutouts scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#send-a-shoutout
func (r *ChatShoutoutResource) Insert(fromBroadcasterID string, toBroadcasterID string, moderatorID string) *ChatShoutoutInsertCall {
	c := &ChatShoutoutInsertCall{resource: r}
	return c.
		FromBroadcasterID(fromBroadcasterID).
		ToBroadcasterID(toBroadcasterID).
		ModeratorID(moderatorID)
}

// FromBroadcasterID sets the FromBroadcasterID query parameter.
func (api *ChatShoutoutInsertCall) FromBroadcasterID(fromBroadcasterID string) *ChatShoutoutInsertCall {
	api.opts = append(api.opts, SetQueryParameter("from_broadcaster_id", fromBroadcasterID))
	return api
}

// ToBroadcasterID sets the ToBroadcasterID query parameter.
func (api *ChatShoutoutInsertCall) ToBroadcasterID(toBroadcasterID string) *ChatShoutoutInsertCall {
	api.opts = append(api.opts, SetQueryParameter("to_broadcaster_id", toBroadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *ChatShoutoutInsertCall) ModeratorID(moderatorID string) *ChatShoutoutInsertCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// Do executes the request.
func (api *ChatShoutoutInsertCall) Do(ctx context.Context, opts ...RequestOption) (*ChatShoutoutInsertResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/chat/shoutouts", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &ChatShoutoutInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

// SendMessageInsertCall represents a POST call to a Twitch Chat API endpoint.
type SendMessageInsertCall struct {
	resource *ChatResource
	body     map[string]any
}

// SendMessageInsertResponse represents the response from a POST request to /helix/chat/messages.
type SendMessageInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the OutboundChatMessage data returned by the Twitch API.
	Data []OutboundChatMessage
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/chat/messages.
//
// Sends a message to the specified chat room.
//
// # Rate Limits
//
// A user may send 20 messages every 30 seconds per channel.
//
// # Authorization
//
// Requires an app access token or user access token that includes the user:write:chat scope.
//
// If app access token used, then additionally requires user:bot scope from chatting user, and either channel:bot scope from broadcaster or moderator status.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#send-chat-message
func (r *ChatResource) Insert(broadcasterID string, senderID string, message string) *SendMessageInsertCall {
	c := &SendMessageInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		SenderID(senderID).
		Message(message)
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *SendMessageInsertCall) BroadcasterID(broadcasterID string) *SendMessageInsertCall {
	api.body["broadcaster_id"] = broadcasterID
	return api
}

// SenderID sets the SenderID body parameter.
func (api *SendMessageInsertCall) SenderID(senderID string) *SendMessageInsertCall {
	api.body["sender_id"] = senderID
	return api
}

// Message sets the Message body parameter.
func (api *SendMessageInsertCall) Message(message string) *SendMessageInsertCall {
	api.body["message"] = message
	return api
}

// ReplayParentMessageID sets the ReplayParentMessageID body parameter.
func (api *SendMessageInsertCall) ReplayParentMessageID(replayParentMessageID string) *SendMessageInsertCall {
	api.body["replay_parent_message_id"] = replayParentMessageID
	return api
}

// ForSourceOnly sets the ForSourceOnly body parameter.
func (api *SendMessageInsertCall) ForSourceOnly(forSourceOnly bool) *SendMessageInsertCall {
	api.body["for_source_only"] = forSourceOnly
	return api
}

// Do executes the request.
func (api *SendMessageInsertCall) Do(ctx context.Context, opts ...RequestOption) (*SendMessageInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/chat/messages", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[OutboundChatMessage](res)
	if err != nil {
		return nil, err
	}

	return &SendMessageInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

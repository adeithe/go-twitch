package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ModerationResource represents the Twitch Moderation API.
type ModerationResource struct {
	client *Client

	// Bans provides access to the Twitch Bans API.
	Bans *ModerationBansResource
	// ClearChat provides access to the Twitch ClearChat API.
	ClearChat *ModerationClearChatResource
}

// NewModerationResource creates a new ModerationResource.
func NewModerationResource(client *Client) *ModerationResource {
	r := &ModerationResource{client: client}
	r.Bans = NewModerationBansResource(client)
	r.ClearChat = NewModerationClearChatResource(client)
	return r
}

// ModerationBansResource represents the Twitch ModerationBans API.
type ModerationBansResource struct {
	client *Client
}

// NewModerationBansResource creates a new ModerationBansResource.
func NewModerationBansResource(client *Client) *ModerationBansResource {
	return &ModerationBansResource{client}
}

// BanUserInsertCall represents a POST call to a Twitch ModerationBans API endpoint.
type BanUserInsertCall struct {
	resource *ModerationBansResource
	body     map[string]any
	opts     []RequestOption
}

// BanUserInsertResponse represents the response from a POST request to /helix/moderation/bans.
type BanUserInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the IssuedBan data returned by the Twitch API.
	Data []IssuedBan
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/moderation/bans.
//
// Bans a user from participating in the specified broadcaster's chat room or puts them in a timeout.
//
// If the user is currently in a timeout, you can call this endpoint to change the duration of the timeout or ban them altogether.
// If the user is currently banned, you cannot call this method to put them in a timeout instead.
//
// # Authorization
//
// Requires a user access token that includes the moderator:manage:banned_users scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#ban-user
func (r *ModerationBansResource) Insert(broadcasterID string, moderatorID string, ban OutboundBan) *BanUserInsertCall {
	c := &BanUserInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID).
		Ban(ban)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *BanUserInsertCall) BroadcasterID(broadcasterID string) *BanUserInsertCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *BanUserInsertCall) ModeratorID(moderatorID string) *BanUserInsertCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// Ban sets the Ban body parameter.
func (api *BanUserInsertCall) Ban(bans ...OutboundBan) *BanUserInsertCall {
	api.body["data"] = bans
	return api
}

// Do executes the request.
func (api *BanUserInsertCall) Do(ctx context.Context, opts ...RequestOption) (*BanUserInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/moderation/bans", bytes.NewReader(bs), append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[IssuedBan](res)
	if err != nil {
		return nil, err
	}

	return &BanUserInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// UnbanUserDeleteCall represents a DELETE call to a Twitch ModerationBans API endpoint.
type UnbanUserDeleteCall struct {
	resource *ModerationBansResource
	opts     []RequestOption
}

// UnbanUserDeleteResponse represents the response from a DELETE request to /helix/moderation/bans.
type UnbanUserDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/moderation/bans.
//
// Removes the ban or timeout that was placed on the specified user.
//
// # Authorization
//
// Requires a user access token that includes the moderator:manage:banned_users scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#unban-user
func (r *ModerationBansResource) Delete(broadcasterID string, moderatorID string, userID string) *UnbanUserDeleteCall {
	c := &UnbanUserDeleteCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID).
		UserID(userID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *UnbanUserDeleteCall) BroadcasterID(broadcasterID string) *UnbanUserDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *UnbanUserDeleteCall) ModeratorID(moderatorID string) *UnbanUserDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// UserID sets the UserID query parameter.
func (api *UnbanUserDeleteCall) UserID(userID string) *UnbanUserDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// Do executes the request.
func (api *UnbanUserDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*UnbanUserDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/moderation/bans", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &UnbanUserDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

// ModerationClearChatResource represents the Twitch ModerationClearChat API.
type ModerationClearChatResource struct {
	client *Client
}

// NewModerationClearChatResource creates a new ModerationClearChatResource.
func NewModerationClearChatResource(client *Client) *ModerationClearChatResource {
	return &ModerationClearChatResource{client}
}

// ChatMessagesDeleteCall represents a DELETE call to a Twitch ModerationClearChat API endpoint.
type ChatMessagesDeleteCall struct {
	resource *ModerationClearChatResource
	opts     []RequestOption
}

// ChatMessagesDeleteResponse represents the response from a DELETE request to /helix/moderation/chat.
type ChatMessagesDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/moderation/chat.
//
// Removes a single chat message or all chat messages from the broadcaster's chat room.
//
// # Authorization
//
// Requires a user access token that includes the moderator:manage:chat_messages scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#delete-chat-messages
func (r *ModerationClearChatResource) Delete(broadcasterID string, moderatorID string) *ChatMessagesDeleteCall {
	c := &ChatMessagesDeleteCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChatMessagesDeleteCall) BroadcasterID(broadcasterID string) *ChatMessagesDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *ChatMessagesDeleteCall) ModeratorID(moderatorID string) *ChatMessagesDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// MessageID sets the MessageID query parameter.
func (api *ChatMessagesDeleteCall) MessageID(messageID string) *ChatMessagesDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("message_id", messageID))
	return api
}

// Do executes the request.
func (api *ChatMessagesDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*ChatMessagesDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/moderation/chat", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &ChatMessagesDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

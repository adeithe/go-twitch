package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ChatterBan represents a ban or timeout on a user in a channel.
type ChatterBan struct {
	BroadcasterID string     `json:"broadcaster_id"`
	ModeratorID   string     `json:"moderator_id"`
	UserID        string     `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	EndsAt        *time.Time `json:"ends_at,omitempty"`
}

// ModerationResource provides access to the Twitch Moderation API.
type ModerationResource struct {
	client *Client
}

// NewModerationResource creates a new ModerationResource.
func NewModerationResource(client *Client) *ModerationResource {
	return &ModerationResource{client}
}

// CreateBanRequest is a request to ban or put a user in a timeout from a channel.
type CreateBanRequest struct {
	resource      *ModerationResource
	broadcasterID string
	moderatorID   string
	userID        string
	duration      *time.Duration
	reason        string
}

// CreateBan creates a request to ban a user from a channel.
//
// Required Scope: moderator:manage:banned_users
func (r *ModerationResource) CreateBan(broadcasterID, moderatorID, userID string) *CreateBanRequest {
	return &CreateBanRequest{r, broadcasterID, moderatorID, userID, nil, ""}
}

// TargetID the ID of the user to ban or put in a timeout.
func (c *CreateBanRequest) TargetID(userID string) *CreateBanRequest {
	c.userID = userID
	return c
}

// Duration the duration of the timeout, in seconds. If omitted, the ban is permanent.
//
// The minimum timeout is 1 second and the maximum is 1,209,600 seconds (2 weeks).
//
// To end a user's timeout early, set this field to 1, or use the Unban user endpoint.
func (c *CreateBanRequest) Duration(duration time.Duration) *CreateBanRequest {
	c.duration = &duration
	return c
}

// Reason the reason the you're banning the user or putting them in a timeout. This is optional and may be an empty string.
//
// Reason is limited to a maximum of 500 characters.
func (c *CreateBanRequest) Reason(reason string) *CreateBanRequest {
	c.reason = reason
	return c
}

// Do executes the request.
func (c *CreateBanRequest) Do(ctx context.Context, opts ...RequestOption) ([]ChatterBan, error) {
	bs, err := json.Marshal(map[string][]map[string]any{
		"data": {{
			"user_id":  c.userID,
			"duration": c.duration,
			"reason":   c.reason,
		}},
	})
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("broadcaster_id", c.broadcasterID)
	query.Set("moderator_id", c.moderatorID)
	res, err := c.resource.client.DoRequest(ctx, http.MethodPost, fmt.Sprintf("%s?%s", EndpointModerationBans, query.Encode()), bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChatterBan](res)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

// RemoveBanRequest is a request to remove a ban on a user from a channel.
type RemoveBanRequest struct {
	resource      *ModerationResource
	broadcasterID string
	moderatorID   string
	userID        string
}

// RemoveBan creates a request to remove a ban on a user from a channel.
//
// Required Scope: moderator:manage:banned_users
func (r *ModerationResource) RemoveBan(broadcasterID, moderatorID, userID string) *RemoveBanRequest {
	return &RemoveBanRequest{r, broadcasterID, moderatorID, userID}
}

// TargetID the ID of the user to unban.
func (c *RemoveBanRequest) TargetID(userID string) *RemoveBanRequest {
	c.userID = userID
	return c
}

// Do executes the request.
func (c *RemoveBanRequest) Do(ctx context.Context, opts ...RequestOption) error {
	query := url.Values{}
	query.Set("broadcaster_id", c.broadcasterID)
	query.Set("moderator_id", c.moderatorID)
	query.Set("user_id", c.userID)

	res, err := c.resource.client.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("%s?%s", EndpointModerationBans, query.Encode()), nil, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	_, err = decodeResponse[any](res)
	return err
}

// ClearChatRequest is a request to clear all messages from a channel.
type ClearChatRequest struct {
	resource      *ModerationResource
	broadcasterID string
	moderatorID   string
	messageID     string
}

// ClearChat creates a request to clear all messages from a channel.
//
// Required Scope: moderator:manage:chat_messages
func (r *ModerationResource) ClearChat(broadcasterID, moderatorID string) *ClearChatRequest {
	return &ClearChatRequest{r, broadcasterID, moderatorID, ""}
}

// MessageID the ID of the message to delete.
func (c *ClearChatRequest) MessageID(messageID string) *ClearChatRequest {
	c.messageID = messageID
	return c
}

// Do executes the request.
func (c *ClearChatRequest) Do(ctx context.Context, opts ...RequestOption) error {
	query := url.Values{}
	query.Set("broadcaster_id", c.broadcasterID)
	query.Set("moderator_id", c.moderatorID)
	query.Set("message_id", c.messageID)

	res, err := c.resource.client.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("%s?%s", EndpointModerationDeleteChatMessages, query.Encode()), nil, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	_, err = decodeResponse[any](res)
	return err
}

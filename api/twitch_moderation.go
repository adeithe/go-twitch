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

type ChatterBan struct {
	BroadcasterID string     `json:"broadcaster_id"`
	ModeratorID   string     `json:"moderator_id"`
	UserID        string     `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	EndsAt        *time.Time `json:"ends_at,omitempty"`
}

type ModerationResource struct {
	client *Client
}

func NewModerationResource(client *Client) *ModerationResource {
	return &ModerationResource{client}
}

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
func (r *ModerationResource) CreateBan(broadcasterId, moderatorId, userId string) *CreateBanRequest {
	return &CreateBanRequest{r, broadcasterId, moderatorId, userId, nil, ""}
}

// UserID the ID of the user to ban or put in a timeout.
func (c *CreateBanRequest) TargetID(userId string) *CreateBanRequest {
	c.userID = userId
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
	res, err := c.resource.client.doRequest(ctx, http.MethodPost, fmt.Sprintf("/moderation/bans?%s", query.Encode()), bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[ChatterBan](res)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

type RemoveBanRequest struct {
	resource      *ModerationResource
	broadcasterID string
	moderatorID   string
	userID        string
}

// RemoveBan creates a request to remove a ban on a user from a channel.
//
// Required Scope: moderator:manage:banned_users
func (r *ModerationResource) RemoveBan(broadcasterId, moderatorId, userId string) *RemoveBanRequest {
	return &RemoveBanRequest{r, broadcasterId, moderatorId, userId}
}

// UserID the ID of the user to unban.
func (c *RemoveBanRequest) TargetID(userId string) *RemoveBanRequest {
	c.userID = userId
	return c
}

// Do executes the request.
func (c *RemoveBanRequest) Do(ctx context.Context, opts ...RequestOption) error {
	query := url.Values{}
	query.Set("broadcaster_id", c.broadcasterID)
	query.Set("moderator_id", c.moderatorID)
	query.Set("user_id", c.userID)

	res, err := c.resource.client.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/moderation/bans?%s", query.Encode()), nil, opts...)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = decodeResponse[any](res)
	return err
}

type ClearChatRequest struct {
	resource      *ModerationResource
	broadcasterId string
	moderatorId   string
	messageId     string
}

// ClearChat creates a request to clear all messages from a channel.
//
// Required Scope: moderator:manage:chat_messages
func (r *ModerationResource) ClearChat(broadcasterId, moderatorId string) *ClearChatRequest {
	return &ClearChatRequest{r, broadcasterId, moderatorId, ""}
}

// MessageID the ID of the message to delete.
func (c *ClearChatRequest) MessageID(messageId string) *ClearChatRequest {
	c.messageId = messageId
	return c
}

// Do executes the request.
func (c *ClearChatRequest) Do(ctx context.Context, opts ...RequestOption) error {
	query := url.Values{}
	query.Set("broadcaster_id", c.broadcasterId)
	query.Set("moderator_id", c.moderatorId)
	query.Set("message_id", c.messageId)

	res, err := c.resource.client.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/moderation/chat?%s", query.Encode()), nil, opts...)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = decodeResponse[any](res)
	return err
}

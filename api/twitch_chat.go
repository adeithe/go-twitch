package api

import (
	"context"
	"fmt"
	"net/http"
)

// Chatter represents a user in Twitch chat.
type Chatter struct {
	ID          string `json:"user_id"`
	Username    string `json:"user_login"`
	DisplayName string `json:"user_name"`
}

// ChatResource handles chat related API calls.
type ChatResource struct {
	client *Client

	Chatters *ChattersResource
}

// NewChatResource creates a new ChatResource.
func NewChatResource(client *Client) *ChatResource {
	r := &ChatResource{client: client}
	r.Chatters = NewChattersResource(client)
	return r
}

// ChattersResource handles chatters related API calls.
type ChattersResource struct {
	client *Client
}

// NewChattersResource creates a new ChattersResource.
func NewChattersResource(client *Client) *ChattersResource {
	return &ChattersResource{client: client}
}

// ChattersResponse is the response from the chatters endpoint.
type ChattersResponse struct {
	Total    int
	Header   http.Header
	Chatters []Chatter
	Cursor   string
}

// ChattersListCall is a call to the chatters list endpoint.
type ChattersListCall struct {
	resource *ChattersResource
	opts     []RequestOption
}

const (
	// EndpointChatGetChatters is the endpoint for getting chatters in a channel.
	EndpointChatGetChatters = TwitchAPIVersionHelix + "/chat/chatters"
	// EndpointChatGetChannelEmotes is the endpoint for getting channel emotes.
	EndpointChatGetChannelEmotes = TwitchAPIVersionHelix + "/chat/emotes"
	// EndpointChatGetGlobalEmotes is the endpoint for getting global emotes.
	EndpointChatGetGlobalEmotes = TwitchAPIVersionHelix + "/chat/emotes/global"
	// EndpointChatGetEmoteSets is the endpoint for getting emote sets.
	EndpointChatGetEmoteSets = TwitchAPIVersionHelix + "/chat/emotes/set"
	// EndpointChatGetChannelBadges is the endpoint for getting a channels chat badges.
	EndpointChatGetChannelBadges = TwitchAPIVersionHelix + "/chat/badges"
	// EndpointChatGetGlobalBadges is the endpoint for getting global chat badges.
	EndpointChatGetGlobalBadges = TwitchAPIVersionHelix + "/chat/badges/global"
	// EndpointChatGetSettings is the endpoint for getting chat settings.
	EndpointChatGetSettings = TwitchAPIVersionHelix + "/chat/settings"
	// EndpointChatGetSharedChatSession is the endpoint for getting a shared chat session.
	EndpointChatGetSharedChatSession = TwitchAPIVersionHelix + "/shared_chat/session"
	// EndpointChatGetUserEmotes is the endpoint for getting user emotes.
	EndpointChatGetUserEmotes = TwitchAPIVersionHelix + "/chat/emotes/user"
	// EndpointChatUpdateSettings is the endpoint for updating a chatrooms settings.
	EndpointChatUpdateSettings = TwitchAPIVersionHelix + "/chat/settings"
	// EndpointChatSendAnnouncement is the endpoint for sending an announcement to a chatroom.
	EndpointChatSendAnnouncement = TwitchAPIVersionHelix + "/chat/announcements"
	// EndpointChatSendShoutout is the endpoint for sending a shoutout in chat.
	EndpointChatSendShoutout = TwitchAPIVersionHelix + "/chat/shoutouts"
	// EndpointChatSendMessage is the endpoint for sending a message in chat.
	EndpointChatSendMessage = TwitchAPIVersionHelix + "/chat/messages"
	// EndpointChatGetUserColor is the endpoint for getting a user's chat color.
	EndpointChatGetUserColor = TwitchAPIVersionHelix + "/chat/color"
	// EndpointChatUpdateUserColor is the endpoint for updating a user's chat color.
	EndpointChatUpdateUserColor = TwitchAPIVersionHelix + "/chat/color"
)

// List creates a new call to list chatters.
func (r *ChattersResource) List(broadcasterID, moderatorID string) *ChattersListCall {
	return &ChattersListCall{
		resource: r,
		opts: []RequestOption{
			SetQueryParameter("broadcaster_id", broadcasterID),
			SetQueryParameter("moderator_id", moderatorID),
		},
	}
}

// First filters the results to the first n chatters.
func (c *ChattersListCall) First(n int) *ChattersListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// After filters the results to those after the specified cursor.
func (c *ChattersListCall) After(cursor string) *ChattersListCall {
	c.opts = append(c.opts, SetQueryParameter("after", cursor))
	return c
}

// Do executes the request.
func (c *ChattersListCall) Do(ctx context.Context, opts ...RequestOption) (*ChattersResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, EndpointChatGetChatters, nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Chatter](res)
	if err != nil {
		return nil, err
	}

	return &ChattersResponse{
		Total:    data.Total,
		Header:   res.Header,
		Chatters: data.Data,
		Cursor:   data.Pagination.Cursor,
	}, nil
}

package irc

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrShardIDOutOfBounds returned when an invalid shard id is provided
	ErrShardIDOutOfBounds = errors.New("shard id out of bounds")
	// ErrShardedMessageSend returned when a sharded connection tries to send a message
	ErrShardedMessageSend = errors.New("messages can not be sent on a sharded connection")
	// ErrNotConnected returned when the connection is closed
	ErrNotConnected = errors.New("connection is closed")
	// ErrNotAuthenticated returned when authentication is required but not set
	ErrNotAuthenticated = errors.New("not authenticated")
	// ErrAlreadyConnected returned when a connection tries to connect but is already running
	ErrAlreadyConnected = errors.New("client is already connected")
	// ErrPingTimeout returned when the server takes too long to respond to a ping message
	ErrPingTimeout = errors.New("server took too long to respond to ping")
	// ErrPartialMessage returned when a raw message is not complete
	ErrPartialMessage = errors.New("parseError: partial message")
	// ErrNoCommand returned when a raw message has no command
	ErrNoCommand = errors.New("parseError: no command")
)

// GlobalUserState is the state of the authenticated user that does not change across channels
type GlobalUserState struct {
	ID          int
	DisplayName string
	Color       string
	Badges      map[string]string
	BadgeInfo   map[string]string
	EmoteSets   []string
	Type        string
}

// UserState is the state of the authenticated user in a channels chatroom
type UserState struct {
	ID            int
	DisplayName   string
	Color         string
	Badges        map[string]string
	BadgeInfo     map[string]string
	EmoteSets     []string
	IsBroadcaster bool
	IsModerator   bool
	IsVIP         bool
	IsSubscriber  bool
	Type          string
}

// RoomState is the current state of a channels chatroom
type RoomState struct {
	UserState         UserState
	ID                int
	Name              string
	isEmoteOnly       bool
	isSubscribersOnly bool
	isRitualsEnabled  bool
	isR9KModeEnabled  bool
	followersOnly     float64
	slowMode          float64
}

// ChatSender is a user that sent a message in a channels chatroom
type ChatSender struct {
	ID            int
	Username      string
	DisplayName   string
	Color         string
	Badges        map[string]string
	BadgeInfo     map[string]string
	IsBroadcaster bool
	IsModerator   bool
	IsVIP         bool
	IsSubscriber  bool
	Type          string
}

// ChatMessage is a message sent in a channels chatroom
type ChatMessage struct {
	IRCMessage Message
	Sender     ChatSender
	ID         string
	Channel    string
	ChannelID  int
	Text       string
	IsCheer    bool
	IsAction   bool
	CreatedAt  time.Time
}

// NewGlobalUserState parses the state of an authenticated user
func NewGlobalUserState(msg Message) (state GlobalUserState) {
	user := NewChatSender(msg)
	state = GlobalUserState{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Color:       user.Color,
		Badges:      user.Badges,
		BadgeInfo:   user.BadgeInfo,
		EmoteSets:   strings.Split(msg.Tags["emote-sets"], ","),
		Type:        user.Type,
	}
	return
}

// NewUserState parses the state of an authenticated user in a channels chatroom
func NewUserState(msg Message) (state UserState) {
	user := NewChatSender(msg)
	state = UserState{
		ID:            user.ID,
		DisplayName:   user.DisplayName,
		Color:         user.Color,
		Badges:        user.Badges,
		BadgeInfo:     user.BadgeInfo,
		EmoteSets:     strings.Split(msg.Tags["emote-sets"], ","),
		IsBroadcaster: user.IsBroadcaster,
		IsModerator:   user.IsModerator,
		IsVIP:         user.IsVIP,
		IsSubscriber:  user.IsSubscriber,
		Type:          user.Type,
	}
	return
}

// NewRoomState parses the state of a channels chatroom and stores it to the provided pointer
func NewRoomState(msg Message, state *RoomState) *RoomState {
	if state == nil {
		state = &RoomState{}
	}
	state.Name = strings.TrimPrefix(msg.Params[0], "#")
	if val, ok := msg.Tags["room-id"]; ok {
		if id, err := strconv.Atoi(val); err == nil {
			state.ID = id
		}
	}
	if val, ok := msg.Tags["emote-only"]; ok {
		state.isEmoteOnly = val == "1"
	}
	if val, ok := msg.Tags["subs-only"]; ok {
		state.isSubscribersOnly = val == "1"
	}
	if val, ok := msg.Tags["rituals"]; ok {
		state.isRitualsEnabled = val == "1"
	}
	if val, ok := msg.Tags["r9k"]; ok {
		state.isR9KModeEnabled = val == "1"
	}
	if val, ok := msg.Tags["followers-only"]; ok {
		state.followersOnly, _ = strconv.ParseFloat(val, 64)
	}
	if val, ok := msg.Tags["slow"]; ok {
		state.slowMode, _ = strconv.ParseFloat(val, 64)
	}
	return state
}

// NewChatSender parses the sender for a message in a channels chatroom
func NewChatSender(msg Message) ChatSender {
	sender := ChatSender{
		Username:    msg.Sender.Nickname,
		DisplayName: msg.Sender.Nickname,
		Color:       msg.Tags["color"],
		Badges:      make(map[string]string),
		BadgeInfo:   make(map[string]string),
		Type:        msg.Tags["user-type"],
	}
	if name, ok := msg.Tags["display-name"]; ok {
		sender.DisplayName = name
	}
	if id, err := strconv.Atoi(msg.Tags["user-id"]); err == nil {
		sender.ID = id
	}

	if len(msg.Tags["badges"]) > 0 {
		badges := strings.Split(msg.Tags["badges"], ",")
		if len(badges) > 0 {
			for _, badge := range badges {
				data := strings.Split(badge, "/")
				sender.Badges[data[0]] = data[1]
			}
		}
	}

	if len(msg.Tags["badge-info"]) > 0 {
		badges := strings.Split(msg.Tags["badge-info"], ",")
		if len(badges) > 0 {
			for _, badge := range badges {
				data := strings.Split(badge, "/")
				sender.BadgeInfo[data[0]] = data[1]
			}
		}
	}

	_, sender.IsBroadcaster = sender.Badges["broadcaster"]
	_, sender.IsModerator = sender.Badges["moderator"]
	_, sender.IsVIP = sender.Badges["vip"]
	_, sender.IsSubscriber = sender.Badges["subscriber"]

	return sender
}

// NewChatMessage parses a message sent in a channels chatroom
func NewChatMessage(msg Message) ChatMessage {
	chatMsg := ChatMessage{
		IRCMessage: msg,
		Sender:     NewChatSender(msg),
		ID:         msg.Tags["id"],
		Channel:    strings.TrimPrefix(msg.Params[0], "#"),
		Text:       msg.Message,
	}
	if id, err := strconv.Atoi(msg.Tags["room-id"]); err == nil {
		chatMsg.ChannelID = id
	}

	_, isCheer := msg.Tags["bits"]
	chatMsg.IsCheer = isCheer

	if strings.HasPrefix(msg.Message, "\u0001ACTION") && strings.HasSuffix(msg.Message, "\u0001") {
		chatMsg.Text = chatMsg.Text[8 : len(chatMsg.Text)-1]
		chatMsg.IsAction = true
	}
	if ts, err := strconv.ParseInt(msg.Tags["tmi-sent-ts"], 10, 64); err == nil {
		chatMsg.CreatedAt = time.Unix(0, ts*1e6)
	}
	return chatMsg
}

// IsEmoteOnly returns true if users without VIP or moderator privileges are only permitted to send emotes
func (state RoomState) IsEmoteOnly() bool {
	return state.isEmoteOnly
}

// IsSubscribersOnly returns true if users without VIP or moderator privileges must be subscribed to send messages
func (state RoomState) IsSubscribersOnly() bool {
	return state.isSubscribersOnly
}

// IsRitualsEnabled returns true if users are able to signal to the chatroom that they are new to the community
func (state RoomState) IsRitualsEnabled() bool {
	return state.isRitualsEnabled
}

// IsR9KModeEnabled returns true if messages must contain more than 9 unique characters to be sent successfully
func (state RoomState) IsR9KModeEnabled() bool {
	return state.isR9KModeEnabled
}

// IsFollowersOnly returns true if users must be following for a set duration before sending messages
func (state RoomState) IsFollowersOnly() (enabled bool, duration time.Duration) {
	enabled = state.followersOnly >= 0
	if state.followersOnly > 0 {
		duration = time.Duration(state.followersOnly) * time.Minute
	}
	return
}

// IsSlowModeEnabled returns true if non-moderators must wait a set duration between sending messages
func (state RoomState) IsSlowModeEnabled() (enabled bool, duration time.Duration) {
	enabled = state.slowMode > 0
	if state.slowMode > 0 {
		duration = time.Duration(state.slowMode) * time.Minute
	}
	return
}

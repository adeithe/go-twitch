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

// ServerNotice is a message sent from the IRC server with general notices
type ServerNotice struct {
	IRCMessage Message
	Channel    string
	Message    string
	Type       string
}

// UserState is a generic state of an authenticated user
type UserState struct {
	ID          int64
	DisplayName string
	Color       string
	Badges      map[string]string
	BadgeInfo   map[string]string
	EmoteSets   []int64
	Type        string
}

// GlobalUserState is the state of the authenticated user that does not change across channels
type GlobalUserState struct {
	UserState
}

// ChannelUserState is the state of the authenticated user in a channels chatroom
type ChannelUserState struct {
	UserState
	IsBroadcaster bool
	IsModerator   bool
	IsVIP         bool
	IsSubscriber  bool
}

// RoomState is the current state of a channels chatroom
type RoomState struct {
	UserState         ChannelUserState
	ID                int64
	Name              string
	isEmoteOnly       bool
	isSubscribersOnly bool
	isRitualsEnabled  bool
	isR9KModeEnabled  bool
	followersOnly     float64
	slowMode          float64
}

// UserNotice is a generic user based event in a channels chatroom
type UserNotice struct {
	IRCMessage Message
	Sender     ChatSender
	ID         string
	Message    string
	Type       string
	CreatedAt  time.Time
}

// ChatSender is a user that sent a message in a channels chatroom
type ChatSender struct {
	ID            int64
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
	ChannelID  int64
	Text       string
	IsCheer    bool
	IsAction   bool
	CreatedAt  time.Time
}

// ChatBan is a timeout or permanent ban that was issued on a user in a chatroom
type ChatBan struct {
	IRCMessage  Message
	ChannelName string
	ChannelID   int64
	TargetName  string
	TargetID    int64
	duration    int
	CreatedAt   time.Time
}

// ChatMessageDelete is a message that was deleted from the chatroom
type ChatMessageDelete struct {
	IRCMessage       Message
	ChannelName      string
	TargetID         string
	TargetSenderName string
	Text             string
	CreatedAt        time.Time
}

// NewServerNotice parses a notice message sent from the IRC server
func NewServerNotice(msg Message) ServerNotice {
	notice := ServerNotice{
		IRCMessage: msg,
		Channel:    strings.TrimPrefix(msg.Params[0], "#"),
		Message:    msg.Text,
		Type:       msg.Tags["msg-id"],
	}
	return notice
}

// NewUserState parses a generic state of an authenticated user
func NewUserState(msg Message) (state UserState) {
	user := NewChatSender(msg)
	state = UserState{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Color:       user.Color,
		Badges:      user.Badges,
		BadgeInfo:   user.BadgeInfo,
		Type:        user.Type,
	}
	for _, set := range strings.Split(msg.Tags["emote-sets"], ",") {
		if i, err := toParsedID(set); err == nil {
			state.EmoteSets = append(state.EmoteSets, i)
		}
	}
	return
}

// NewGlobalUserState parses the global state of an authenticated user
func NewGlobalUserState(msg Message) (state GlobalUserState) {
	state = GlobalUserState{
		UserState: NewUserState(msg),
	}
	return
}

// NewChannelUserState parses the state of an authenticated user in a channels chatroom
func NewChannelUserState(msg Message) (state ChannelUserState) {
	user := NewChatSender(msg)
	state = ChannelUserState{
		UserState:     NewUserState(msg),
		IsBroadcaster: user.IsBroadcaster,
		IsModerator:   user.IsModerator,
		IsVIP:         user.IsVIP,
		IsSubscriber:  user.IsSubscriber,
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
		if id, err := toParsedID(val); err == nil {
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

// NewUserNotice parses a generic user based event in a channels chatroom
func NewUserNotice(msg Message) UserNotice {
	notice := UserNotice{
		IRCMessage: msg,
		Sender:     NewChatSender(msg),
		ID:         msg.Tags["id"],
		Message:    msg.Tags["system-msg"],
		Type:       msg.Tags["msg-id"],
	}
	if ts, err := toParsedTimestamp(msg.Tags["tmi-sent-ts"]); err == nil {
		notice.CreatedAt = ts
	}
	for _, escape := range escapeChars {
		notice.Message = strings.ReplaceAll(notice.Message, escape.from, escape.to)
	}
	return notice
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
	if name, ok := msg.Tags["login"]; ok {
		sender.Username = name
	}
	if name, ok := msg.Tags["display-name"]; ok {
		sender.DisplayName = name
	}
	if id, err := toParsedID(msg.Tags["user-id"]); err == nil {
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
		Text:       msg.Text,
	}
	if id, err := toParsedID(msg.Tags["room-id"]); err == nil {
		chatMsg.ChannelID = id
	}

	if strings.HasPrefix(msg.Text, "\u0001ACTION") && strings.HasSuffix(msg.Text, "\u0001") {
		chatMsg.Text = chatMsg.Text[8 : len(chatMsg.Text)-1]
		chatMsg.IsAction = true
	}
	if ts, err := toParsedTimestamp(msg.Tags["tmi-sent-ts"]); err == nil {
		chatMsg.CreatedAt = ts
	}

	_, chatMsg.IsCheer = msg.Tags["bits"]

	return chatMsg
}

// NewChatBan parses a ban or timeout in a channels chatroom
func NewChatBan(msg Message) ChatBan {
	ban := ChatBan{
		IRCMessage:  msg,
		ChannelName: strings.TrimPrefix(msg.Params[0], "#"),
		TargetName:  msg.Text,
	}
	if id, err := toParsedID(msg.Tags["room-id"]); err == nil {
		ban.ChannelID = id
	}
	if id, err := toParsedID(msg.Tags["target-user-id"]); err == nil {
		ban.TargetID = id
	}
	if n, err := strconv.Atoi(msg.Tags["ban-duration"]); err == nil {
		ban.duration = n
	}
	if ts, err := toParsedTimestamp(msg.Tags["tmi-sent-ts"]); err == nil {
		ban.CreatedAt = ts
	}
	return ban
}

// NewChatMessageDelete parses a notice that a message was deleted
func NewChatMessageDelete(msg Message) ChatMessageDelete {
	delete := ChatMessageDelete{
		IRCMessage:       msg,
		ChannelName:      strings.TrimPrefix(msg.Params[0], "#"),
		TargetID:         msg.Tags["target-msg-id"],
		TargetSenderName: msg.Tags["login"],
		Text:             msg.Text,
	}
	if ts, err := toParsedTimestamp(msg.Tags["tmi-sent-ts"]); err == nil {
		delete.CreatedAt = ts
	}
	return delete
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

// IsTemporary returns true if a ban is set to expire after a set amount of time
func (ban ChatBan) IsTemporary() bool {
	return ban.duration > 0
}

// Duration returns the duration of a temporary ban
func (ban ChatBan) Duration() time.Duration {
	if ban.duration < 1 {
		return 0
	}
	return time.Duration(ban.duration) * time.Second
}

// Expiration returns the time that a temporary ban will expire
func (ban ChatBan) Expiration() time.Time {
	return ban.CreatedAt.Add(time.Duration(ban.duration) * time.Second)
}

func toParsedID(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func toParsedTimestamp(ts string) (time.Time, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(0, i*1e6), nil
}

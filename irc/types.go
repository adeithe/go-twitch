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

// UserState is the state of the authenticated user
type UserState struct {
	UserID      int
	DisplayName string
	Color       string
	EmoteSets   []string
}

// ChatSender is a user that sent a message in a channels chatroom
type ChatSender struct {
	Username     string
	DisplayName  string
	UserID       int
	Color        string
	Badges       map[string]string
	IsSubscriber bool
	IsModerator  bool
	Type         string
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

// NewUserState parses the state of an authenticated user
func NewUserState(msg Message) UserState {
	state := UserState{
		DisplayName: msg.Tags["display-name"],
		Color:       msg.Tags["color"],
		EmoteSets:   strings.Split(msg.Tags["emote-sets"], ","),
	}
	if id, err := strconv.Atoi(msg.Tags["user-id"]); err == nil {
		state.UserID = id
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
		Type:        msg.Tags["user-type"],
	}
	if name, ok := msg.Tags["display-name"]; ok {
		sender.DisplayName = name
	}
	if id, err := strconv.Atoi(msg.Tags["user-id"]); err == nil {
		sender.UserID = id
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

	_, subBadge := sender.Badges["subscriber"]
	sender.IsSubscriber = subBadge || msg.Tags["subscriber"] == "1"

	_, modBadge := sender.Badges["moderator"]
	sender.IsModerator = modBadge || msg.Tags["mod"] == "1"
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

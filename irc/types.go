package irc

import (
	"strconv"
	"strings"
	"time"
)

type UserState struct {
	UserID      int
	DisplayName string
	Color       string
	EmoteSets   []string
}

// ChatSender stores information about the sender of a PRIVMSG
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

// ChatMessage stores information about a parsed IRC message
type ChatMessage struct {
	IRCMessage IRCMessage
	Sender     ChatSender
	ID         string
	Channel    string
	ChannelID  int
	Message    string
	IsCheer    bool
	IsAction   bool
	CreatedAt  time.Time
}

// NewUserState uses an IRCMessage to parse a UserState.
func NewUserState(ircMsg IRCMessage) UserState {
	state := &UserState{
		DisplayName: ircMsg.Tags["display-name"],
		Color:       ircMsg.Tags["color"],
		EmoteSets:   strings.Split(ircMsg.Tags["emote-sets"], ","),
	}
	if id, err := strconv.Atoi(ircMsg.Tags["user-id"]); err == nil {
		state.UserID = id
	}
	return *state
}

// NewChatSender uses an IRCMessage to get information about a PRIVMSG sender
func NewChatSender(ircMsg IRCMessage) ChatSender {
	sender := &ChatSender{
		Username:    ircMsg.Sender.Nickname,
		DisplayName: ircMsg.Sender.Nickname,
		Color:       ircMsg.Tags["color"],
		Badges:      make(map[string]string),
		Type:        ircMsg.Tags["user-type"],
	}
	if displayName, ok := ircMsg.Tags["display-name"]; ok {
		sender.DisplayName = displayName
	}
	if id, err := strconv.Atoi(ircMsg.Tags["user-id"]); err == nil {
		sender.UserID = id
	}

	badges := strings.Split(ircMsg.Tags["badges"], ",")
	if len(badges[0]) > 0 {
		for _, badge := range badges {
			data := strings.Split(badge, "/")
			sender.Badges[data[0]] = data[1]
		}
	}

	_, subBadge := sender.Badges["subscriber"]
	sender.IsSubscriber = subBadge || ircMsg.Tags["subscriber"] == "1"

	_, modBadge := sender.Badges["moderator"]
	sender.IsModerator = modBadge || ircMsg.Tags["mod"] == "1"
	return *sender
}

// NewChatMessage uses IRCMessage to get information about a PRIVMSG
func NewChatMessage(ircMsg IRCMessage) ChatMessage {
	msg := &ChatMessage{
		IRCMessage: ircMsg,
		Sender:     NewChatSender(ircMsg),
		ID:         ircMsg.Tags["id"],
		Channel:    strings.TrimPrefix(ircMsg.Params[0], "#"),
		Message:    ircMsg.Message,
	}
	if id, err := strconv.Atoi(ircMsg.Tags["room-id"]); err == nil {
		msg.ChannelID = id
	}
	if _, ok := ircMsg.Tags["bits"]; ok {
		msg.IsCheer = ok
	}
	if strings.HasPrefix(msg.Message, "\u0001ACTION") && strings.HasSuffix(msg.Message, "\u0001") {
		msg.Message = msg.Message[8 : len(msg.Message)-1]
		msg.IsAction = true
	}
	if ts, err := strconv.ParseInt(ircMsg.Tags["tmi-sent-ts"], 10, 64); err == nil {
		msg.CreatedAt = time.Unix(0, ts*1e6)
	}
	return *msg
}

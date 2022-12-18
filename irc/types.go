package irc

import (
	"strconv"
	"strings"
	"time"
)

type Notice struct {
	MessageID    string
	TargetUserID string
	Channel      *Channel
	raw          *RawMessage
}

type Badge struct {
	Version  string
	Metadata string
}

// Source the source of the incoming message.
type Source struct {
	Nickname string
	Username string
	Host     string
}

// RawMessage a parsed IRC message.
type RawMessage struct {
	Tags    Tags
	Source  Source
	Command IRCCommand
	Params  []string
	Text    string
	Raw     string
}

// Tags is a map of Twitch IRC tags.
type Tags map[string]string

// Command is the command sent in an IRCMessage
type IRCCommand string

const (
	// CMDReady is a 376 command
	CMDReady IRCCommand = "376"
	// CMDPing is a PING command
	CMDPing IRCCommand = "PING"
	// CMDPong is a PONG command
	CMDPong IRCCommand = "PONG"
	// CMDJoin is a JOIN command
	CMDJoin IRCCommand = "JOIN"
	// CMDPart is a PART command
	CMDPart IRCCommand = "PART"
	// CMDReconnect is a RECONNECT command
	CMDReconnect IRCCommand = "RECONNECT"
	// CMDGlobalUserState is a GLOBALUSERSTATE command
	CMDGlobalUserState IRCCommand = "GLOBALUSERSTATE"
	// CMDNotice is a NOTICE command
	CMDNotice IRCCommand = "NOTICE"
	// CMDRoomState is a ROOMSTATE command
	CMDRoomState IRCCommand = "ROOMSTATE"
	// CMDUserState is a USERSTATE command
	CMDUserState IRCCommand = "USERSTATE"
	// CMDUserNotice is a USERNOTICE command
	CMDUserNotice IRCCommand = "USERNOTICE"
	// CMDPrivMessage is a PRIVMSG command
	CMDPrivMessage IRCCommand = "PRIVMSG"
	// CMDWhisper is a WHISPER command
	CMDWhisper IRCCommand = "WHISPER"
	// CMDClearChat is a CLEARCHAT command
	CMDClearChat IRCCommand = "CLEARCHAT"
	// CMDClearMessage is a CLEARMSG command
	CMDClearMessage IRCCommand = "CLEARMSG"
	// Deprecated: CMDHostTarget is a HOSTTARGET command
	CMDHostTarget IRCCommand = "HOSTTARGET"
)

func (m RawMessage) String() string {
	return m.Raw
}

func (m Notice) String() string {
	return m.raw.Text
}

func toUsername(channel string) string {
	return strings.TrimPrefix(strings.ToLower(channel), "#")
}

func toBadges(msg *RawMessage) map[string]Badge {
	badges := make(map[string]Badge)
	for _, badge := range strings.Split(msg.Tags["badges"], ",") {
		parts := strings.SplitN(badge, "/", 2)
		if len(parts) > 1 {
			badges[parts[0]] = Badge{Version: parts[1]}
		}
	}
	for _, badge := range strings.Split(msg.Tags["badge-info"], ",") {
		parts := strings.SplitN(badge, "/", 2)
		if len(parts) > 1 {
			if badge, ok := badges[parts[0]]; ok {
				badge.Metadata = parts[1]
			}
		}
	}
	return badges
}

func toParsedTimestamp(ts string) time.Time {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Now()
	}
	return time.Unix(0, i*1e6)
}

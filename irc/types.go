package irc

import "strings"

// Source the source of the incoming message.
type Source struct {
	Nickname string
	Username string
	Host     string
}

// Message a parsed IRC message.
type Message struct {
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
	// CMDNotice is a NOTICE command
	CMDNotice IRCCommand = "NOTICE"
	// CMDReconnect is a RECONNECT command
	CMDReconnect IRCCommand = "RECONNECT"
	// CMDRoomState is a ROOMSTATE command
	CMDRoomState IRCCommand = "ROOMSTATE"
	// CMDUserNotice is a USERNOTICE command
	CMDUserNotice IRCCommand = "USERNOTICE"
	// CMDUserState is a USERSTATE command
	CMDUserState IRCCommand = "USERSTATE"
	// CMDGlobalUserState is a GLOBALUSERSTATE command
	CMDGlobalUserState IRCCommand = "GLOBALUSERSTATE"
	// CMDJoin is a JOIN command
	CMDJoin IRCCommand = "JOIN"
	// CMDPart is a PART command
	CMDPart IRCCommand = "PART"
	// CMDPing is a PING command
	CMDPing IRCCommand = "PING"
	// CMDPong is a PONG command
	CMDPong IRCCommand = "PONG"
	// CMDReady is a 376 command
	CMDReady IRCCommand = "376"
)

func (m Message) String() string {
	return m.Raw
}

func sanitizeUsername(username string) string {
	return strings.TrimPrefix(strings.ToLower(username), "#")
}

package irc

import (
	"regexp"
	"strings"
)

// Command is the command sent in an IRCMessage
type Command string

const (
	// CMDPrivMessage is a PRIVMSG command
	CMDPrivMessage Command = "PRIVMSG"
	// CMDClearChat is a CLEARCHAT command
	CMDClearChat Command = "CLEARCHAT"
	// CMDClearMessage is a CLEARMSG command
	CMDClearMessage Command = "CLEARMSG"
	// CMDHostTarget is a HOSTTARGET command
	CMDHostTarget Command = "HOSTTARGET"
	// CMDNotice is a NOTICE command
	CMDNotice Command = "NOTICE"
	// CMDReconnect is a RECONNECT command
	CMDReconnect Command = "RECONNECT"
	// CMDRoomState is a ROOMSTATE command
	CMDRoomState Command = "ROOMSTATE"
	// CMDUserNotice is a USERNOTICE command
	CMDUserNotice Command = "USERNOTICE"
	// CMDUserState is a USERSTATE command
	CMDUserState Command = "USERSTATE"
	// CMDGlobalUserState is a GLOBALUSERSTATE command
	CMDGlobalUserState Command = "GLOBALUSERSTATE"
	// CMDJoin is a JOIN command
	CMDJoin Command = "JOIN"
	// CMDPart is a PART command
	CMDPart Command = "PART"
	// CMDPing is a PING command
	CMDPing Command = "PING"
	// CMDPong is a PONG command
	CMDPong Command = "PONG"
	// CMDReady is a 376 command
	CMDReady Command = "376"
)

// Source is basic info about in IRC message
type Source struct {
	Nickname string
	Username string
	Host     string
}

// Message is an IRC message received by the socket
type Message struct {
	Raw     string
	Command Command
	Sender  Source
	Tags    map[string]string
	Params  []string
	Text    string
}

// IMessageParser is a generic parser for an IRC message
type IMessageParser interface {
	Parse() error

	tags(string)
	sender(string)
}

var escapeChars = []struct {
	from string
	to   string
}{
	{`\s`, ` `},
	{`\n`, ``},
	{`\r`, ``},
	{`\:`, `;`},
	{`\\`, `\`},
}

// NewParsedMessage parses raw data from the IRC server and returns an IRCMessage
func NewParsedMessage(raw string) (Message, error) {
	msg := &Message{Raw: raw}
	if err := msg.Parse(); err != nil {
		return *msg, err
	}
	return *msg, nil
}

// Parse takes the raw data in an IRCMessage and stores it accordingly
//
// This is done automatically when running NewParsedMessage but can be run again at any time
func (msg *Message) Parse() error {
	var index int
	parts := strings.Split(msg.Raw, " ")

	if strings.HasPrefix(parts[index], "@") {
		msg.tags(strings.TrimPrefix(parts[index], "@"))
		index++
	}

	if index >= len(parts) {
		return ErrPartialMessage
	}

	if strings.HasPrefix(parts[index], ":") {
		msg.sender(strings.TrimPrefix(parts[index], ":"))
		index++
	}

	if index >= len(parts) {
		return ErrNoCommand
	}

	msg.Command = Command(parts[index])
	index++

	if index >= len(parts) {
		return nil
	}

	var params []string
	for i, v := range parts[index:] {
		if strings.HasPrefix(v, ":") {
			msg.Text = strings.TrimPrefix(strings.Join(parts[index+i:], " "), ":")
			break
		}
		params = append(params, v)
	}
	msg.Params = params
	return nil
}

func (msg *Message) tags(raw string) {
	tags := make(map[string]string)
	for _, tag := range strings.Split(raw, ";") {
		pair := strings.SplitN(tag, "=", 2)
		var value string
		if len(pair) == 2 {
			rawValue := pair[1]
			for _, escape := range escapeChars {
				rawValue = strings.ReplaceAll(rawValue, escape.from, escape.to)
			}
			value = strings.TrimSpace(strings.TrimSuffix(rawValue, "\\"))
		}
		tags[pair[0]] = value
	}
	msg.Tags = tags
}

func (msg *Message) sender(raw string) {
	regex := regexp.MustCompile(`!|@`)
	sourceData := regex.Split(raw, -1)
	sender := Source{}
	if len(sourceData) > 0 {
		switch len(sourceData) {
		case 1:
			sender.Host = sourceData[0]
		case 2:
			sender.Nickname = sourceData[0]
			sender.Username = sourceData[0]
			sender.Host = sourceData[1]
		case 3:
			sender.Nickname = sourceData[0]
			sender.Username = sourceData[1]
			sender.Host = sourceData[2]
		}
	}
	msg.Sender = sender
}

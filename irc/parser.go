package irc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Adeithe/go-twitch/irc/cmd"
)

type IRCSource struct {
	Nickname string
	Username string
	Host     string
}

type IRCMessage struct {
	Raw     string
	Command cmd.IRCCommand
	Sender  IRCSource
	Tags    map[string]string
	Params  []string
	Message string
}

// IMessageParser is an interface for IRCMessage
type IMessageParser interface {
	Parse() error

	tags(string)
	sender(string)
}

var _ IMessageParser = &IRCMessage{}

func NewParsedMessage(raw string) (IRCMessage, error) {
	msg := &IRCMessage{Raw: raw}
	if err := msg.Parse(); err != nil {
		return IRCMessage{}, err
	}
	return *msg, nil
}

func (msg *IRCMessage) Parse() error {
	var index int
	parts := strings.Split(msg.Raw, " ")

	if strings.HasPrefix(parts[index], "@") {
		msg.tags(strings.TrimPrefix(parts[index], "@"))
		index++
	}

	if index >= len(parts) {
		return fmt.Errorf("parseError: partial message")
	}

	if strings.HasPrefix(parts[index], ":") {
		msg.sender(strings.TrimPrefix(parts[index], ":"))
		index++
	}

	if index >= len(parts) {
		return fmt.Errorf("parseError: no command")
	}

	msg.Command = cmd.IRCCommand(parts[index])
	index++

	if index >= len(parts) {
		return nil
	}

	var params []string
	for i, v := range parts[index:] {
		if strings.HasPrefix(v, ":") {
			msg.Message = strings.TrimPrefix(strings.Join(parts[index+i:], " "), ":")
			break
		}
		params = append(params, v)
	}
	msg.Params = params

	return nil
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

func (msg *IRCMessage) tags(raw string) {
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

func (msg *IRCMessage) sender(raw string) {
	regex := regexp.MustCompile(`!|@`)
	sourceData := regex.Split(raw, -1)
	sender := &IRCSource{}
	if len(sourceData) > 0 {
		switch len(sourceData) {
		case 1:
			sender.Host = sourceData[0]
		case 2:
			sender.Nickname = sourceData[0]
			sender.Username = sourceData[0]
			sender.Host = sourceData[1]
		default:
			sender.Nickname = sourceData[0]
			sender.Username = sourceData[1]
			sender.Host = sourceData[2]
		}
	}
	msg.Sender = *sender
}

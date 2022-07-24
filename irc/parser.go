package irc

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidTags returned when the tags string is invalid.
	ErrInvalidTags = errors.New("invalid tags string")
	// ErrInvalidSource returned when the source string is invalid.
	ErrInvalidSource = errors.New("invalid source")
	// ErrPartialMessage returned when the raw message is incomplete.
	ErrPartialMessage = errors.New("partial message")
	// ErrNoCommand returned when the raw message has no command.
	ErrNoCommand = errors.New("no irc command")

	senderRegexp = regexp.MustCompile(`!|@`)
	escapeChars  = []struct {
		from string
		to   string
	}{
		{`\s`, ` `},
		{`\n`, ``},
		{`\r`, ``},
		{`\:`, `;`},
		{`\\`, `\`},
	}
)

// ParseMessage parses an IRC message as received by Twitch.
// In most cases this is handled for you and received in the form of an event.
//
// See: https://dev.twitch.tv/docs/irc/send-receive-messages#receiving-chat-messages
func ParseMessage(raw string) (*Message, error) {
	var i int
	parts := strings.Split(raw, " ")

	msg := &Message{Raw: raw}
	if strings.HasPrefix(raw, "@") {
		tags, err := ParseTags(parts[i])
		if err != nil {
			return nil, err
		}
		msg.Tags = tags
		i++
	}

	if i >= len(parts) {
		return nil, ErrPartialMessage
	}

	if strings.HasPrefix(parts[i], ":") {
		source, err := ParseSource(parts[i])
		if err != nil {
			return nil, err
		}
		msg.Source = *source
		i++
	}

	if i >= len(parts) {
		return nil, ErrNoCommand
	}

	msg.Command = IRCCommand(parts[i])
	i++

	if i >= len(parts) {
		return msg, nil
	}

	var params []string
	for n, param := range parts[i:] {
		if strings.HasPrefix(param, ":") {
			text := strings.Join(parts[i+n:], " ")
			msg.Text = strings.TrimPrefix(text, ":")
			break
		}
		params = append(params, param)
	}
	msg.Params = params
	return msg, nil
}

// ParseTags parses a tags string into an object.
// In most cases this is handled for you and received in the form of an event.
//
// Tag strings are provided by Twitch in a form of "@tag-name-1=<tag-value-1>;tag-name-2=<tag-value-2>;..."
//
// See: https://dev.twitch.tv/docs/irc/tags
func ParseTags(tags string) (Tags, error) {
	if len(tags) < 2 || !strings.HasPrefix(tags, "@") {
		return nil, ErrInvalidTags
	}
	tags = strings.TrimPrefix(tags, "@")

	data := make(Tags)
	for _, pair := range strings.Split(tags, ";") {
		parts := strings.SplitN(pair, "=", 2)
		value := parts[1]
		if len(parts) > 1 {
			for _, escape := range escapeChars {
				value = strings.ReplaceAll(value, escape.from, escape.to)
			}
		}
		data[parts[0]] = strings.TrimSpace(strings.TrimSuffix(value, "\\"))
	}
	return data, nil
}

// ParseSource parses a source string.
// In most cases this is handled for you and received in the form of an event.
//
// Source strings are provided by Twitch in a format such as ":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv"
// This message format may vary depending on the type of message.
//
// See: https://dev.twitch.tv/docs/irc#parsing-messages
func ParseSource(source string) (*Source, error) {
	if len(source) < 2 || !strings.HasPrefix(source, ":") {
		return nil, ErrInvalidSource
	}
	source = strings.TrimPrefix(source, ":")

	data := &Source{}
	sourceData := senderRegexp.Split(source, -1)
	if len(sourceData) > 0 {
		switch len(sourceData) {
		case 1:
			data.Host = sourceData[0]
		case 2:
			data.Nickname = sourceData[0]
			data.Username = sourceData[0]
			data.Host = sourceData[1]
		case 3:
			data.Nickname = sourceData[0]
			data.Username = sourceData[1]
			data.Host = sourceData[2]
		}
	}
	return data, nil
}

package irc

import (
	"errors"
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

	escapeChars = []struct {
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
func ParseMessage(raw string) (*RawMessage, error) {
	var i int
	parts := strings.Split(raw, " ")

	msg := &RawMessage{Raw: raw}
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

	data := make(Tags)
	for _, pair := range strings.Split(tags[1:], ";") {
		parts := strings.SplitN(pair, "=", 2)
		for _, escape := range escapeChars {
			if strings.Contains(parts[1], escape.from) {
				parts[1] = strings.ReplaceAll(parts[1], escape.from, escape.to)
			}
		}
		data[parts[0]] = strings.TrimSpace(strings.TrimSuffix(parts[1], "\\"))
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

	hostParts := strings.SplitN(source, "@", 2)
	if len(hostParts) < 2 {
		return &Source{Host: hostParts[0]}, nil
	}

	data := &Source{Host: hostParts[1]}
	userParts := strings.SplitN(hostParts[0], "!", 2)
	if len(userParts) < 2 {
		data.Nickname = userParts[0]
		data.Username = userParts[0]
		return data, nil
	}

	data.Nickname = userParts[0]
	data.Username = userParts[1]
	return data, nil
}

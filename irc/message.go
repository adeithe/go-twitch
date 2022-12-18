package irc

import (
	"strconv"
	"strings"
	"time"
)

type ChatMessage struct {
	Channel *Channel
	Sender  User
	SentAt  time.Time

	tags Tags
	text string
}

type ParentMessage struct {
	ID          string
	UserID      string
	Login       string
	DisplayName string
	Text        string
}

func (c *Conn) handlePrivMessage(msg *RawMessage) {
	if len(msg.Params) < 1 {
		return
	}

	channelName := toUsername(msg.Params[0])
	channel, ok := c.GetChannel(channelName)
	if !ok {
		return
	}

	msg.Tags["login"] = msg.Source.Username
	doEvent(c.events.OnChatMessage)(&ChatMessage{
		Channel: channel,
		Sender: User{
			badges: toBadges(msg),
			tags:   msg.Tags,
		},
		SentAt: toParsedTimestamp(msg.Tags["tmi-sent-ts"]),

		tags: msg.Tags,
		text: msg.Text,
	})
}

func (m ChatMessage) ID() string {
	return m.tags["id"]
}

func (m ChatMessage) Parent() (ParentMessage, bool) {
	if id, ok := m.tags["reply-parent-msg-id"]; ok {
		return ParentMessage{
			ID:          id,
			UserID:      m.tags["reply-parent-user-id"],
			Login:       m.tags["reply-parent-user-login"],
			DisplayName: m.tags["reply-parent-display-name"],
			Text:        m.tags["reply-parent-msg-body"],
		}, ok
	}
	return ParentMessage{}, false
}

func (m ChatMessage) Bits() int {
	bits, _ := strconv.Atoi(m.tags["bits"])
	return bits
}

func (m ChatMessage) IsAction() bool {
	return strings.HasPrefix(m.text, "\001ACTION") && strings.HasSuffix(m.text, "\001")
}

func (m ChatMessage) Text() string {
	return strings.TrimSpace(strings.TrimPrefix("\001ACTION", strings.TrimSuffix(m.text, "\001")))
}

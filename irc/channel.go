package irc

import (
	"fmt"
	"time"
)

type Channel struct {
	userState             UserState
	conn                  *Conn
	name                  string
	roomId                string
	r9kMode               bool
	subOnly               bool
	emoteOnly             bool
	followersOnly         bool
	slowMode              bool
	acknowledged          bool
	followersOnlyDuration time.Duration
	slowModeDuration      time.Duration
	ackC                  chan error
}

func (c *Conn) handleRoomState(msg *RawMessage) {
	if len(msg.Params) < 1 {
		return
	}

	channelName := toUsername(msg.Params[0])
	channel, ok := c.GetChannel(channelName)
	if !ok {
		c.channelsMx.Lock()
		channel = &Channel{conn: c, name: channelName, ackC: make(chan error), acknowledged: true}
		c.channels[channelName] = channel
		c.channelsMx.Unlock()
	}

	for key, val := range msg.Tags {
		switch key {
		case "room-id":
			channel.roomId = val
		case "emote-only":
			channel.emoteOnly = val != "0"
		case "followers-only":
			channel.followersOnly = val != "-1"
			channel.followersOnlyDuration, _ = time.ParseDuration(val + "m")
		case "r9k":
			channel.r9kMode = val != "0"
		case "slow":
			channel.slowMode = val != "0"
			channel.slowModeDuration, _ = time.ParseDuration(val + "s")
		case "subs-only":
			channel.subOnly = val != "0"
		}
	}
	channel.ack(nil)
}

// Username returns the username of the channel.
func (c Channel) Username() string {
	return toUsername(c.name)
}

// RoomID returns the room id of the channel.
func (c Channel) RoomID() string {
	return c.roomId
}

// IsR9KMode returns whether r9k mode is enabled.
func (c Channel) IsR9KMode() bool {
	return c.r9kMode
}

// IsSubOnly returns whether sub only mode is enabled.
func (c Channel) IsSubOnly() bool {
	return c.subOnly
}

// IsEmoteOnly returns whether emote only mode is enabled.
func (c Channel) IsEmoteOnly() bool {
	return c.emoteOnly
}

// IsFollowersOnly returns the duration of followers only mode and whether followers only mode is enabled.
func (c Channel) IsFollowersOnly() (time.Duration, bool) {
	return c.followersOnlyDuration, c.followersOnly
}

// IsSlowMode returns the duration of slow mode and whether slow mode is enabled.
func (c Channel) IsSlowMode() (time.Duration, bool) {
	return c.slowModeDuration, c.slowMode
}

// IsJoined returns true if the channel has been acknowledged by the server.
func (c Channel) IsJoined() bool {
	return c.acknowledged
}

// Say sends a message to the channel.
func (c *Channel) Say(message string) error {
	return c.conn.SendRaw(fmt.Sprintf("PRIVMSG #%s :%s", c.Username(), message))
}

// Sayf sends a formatted message to the channel.
func (c *Channel) Sayf(format string, v ...interface{}) error {
	return c.Say(fmt.Sprintf(format, v...))
}

func (c *Channel) ack(err error) {
	if !c.acknowledged {
		c.acknowledged = true
		c.ackC <- err
	}
}

func (c *Channel) unack() {
	c.acknowledged = false
}

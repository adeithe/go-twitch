package irc

import "strings"

type UserState struct {
	User
}

type User struct {
	badges map[string]Badge
	tags   Tags
}

type UserType string

const (
	UserTypeNone      UserType = ""
	UserTypeGlobalMod UserType = "global_mod"
	UserTypeStaff     UserType = "staff"
	UserTypeAdmin     UserType = "admin"
)

func (c *Conn) handleUserState(msg *RawMessage) {
	if len(msg.Params) < 1 {
		return
	}

	channelName := toUsername(msg.Params[0])
	c.channelsMx.RLock()
	channel, ok := c.channels[channelName]
	c.channelsMx.RUnlock()
	if !ok {
		c.channelsMx.Lock()
		channel = &Channel{conn: c, name: channelName, ackC: make(chan error), acknowledged: true}
		c.channels[channelName] = channel
		c.channelsMx.Unlock()
	}

	msg.Tags["login"] = c.username
	msg.Tags["user-id"] = c.userId
	channel.userState = UserState{User: User{badges: toBadges(msg), tags: msg.Tags}}
}

// ID returns the users ID.
func (s User) ID() string {
	return s.tags["user-id"]
}

// Username returns the users login.
func (s User) Username() string {
	return s.tags["login"]
}

// DisplayName returns the users display name.
func (s User) DisplayName() string {
	if name, ok := s.tags["display-name"]; ok {
		return name
	}
	return s.Username()
}

// EmoteSets returns the emote sets the user has permission to use.
func (s UserState) EmoteSets() []string {
	return strings.Split(s.tags["emote-sets"], ",")
}

// Badge returns the badge with the given name, if it exists.
func (s User) Badge(name string) (Badge, bool) {
	badge, ok := s.badges[strings.ToLower(name)]
	return badge, ok
}

// Color returns the users color, if the user has set one. Otherwise, it returns an empty string.
func (s User) Color() string {
	return s.tags["color"]
}

// IsBroadcaster returns true if the user is the broadcaster of the channel.
func (s User) IsBroadcaster() bool {
	_, ok := s.badges["broadcaster"]
	return ok
}

// IsMod returns true if the user is allowed to perform moderator actions in the channel.
func (s User) IsMod() bool {
	return s.tags["mod"] == "1" || s.IsBroadcaster()
}

// IsSubscriber returns true if the user is a subscriber to the channel.
func (s User) IsSubscriber() bool {
	return s.tags["subscriber"] == "1"
}

// IsTurbo returns true if the user has a turbo subscription.
// Turbo subscriptions allow users to bypass ads regardless of sub status.
func (s User) IsTurbo() bool {
	return s.tags["turbo"] == "1"
}

// UserType returns the users type. For normal users, this is usually empty.
func (s User) UserType() UserType {
	return UserType(s.tags["user-type"])
}

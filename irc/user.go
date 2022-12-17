package irc

import "strings"

type UserState struct {
	badges      map[string]Badge
	color       string
	displayName string
	emoteSets   []string
	mod         bool
	subscriber  bool
	turbo       bool
	userType    UserType
}

type Badge struct {
	Version  string
	Metadata string
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

	state := UserState{badges: make(map[string]Badge)}
	for key, val := range msg.Tags {
		switch key {
		case "badges":
			for _, badge := range strings.Split(val, ",") {
				parts := strings.SplitN(badge, "/", 2)
				if len(parts) > 1 {
					state.badges[parts[0]] = Badge{Version: parts[1]}
				}
			}
			for _, badge := range strings.Split(msg.Tags["badge-info"], ",") {
				parts := strings.SplitN(badge, "/", 2)
				if len(parts) > 1 {
					if badge, ok := state.badges[parts[0]]; ok {
						badge.Metadata = parts[1]
					}
				}
			}
		case "color":
			state.color = val
		case "display-name":
			state.displayName = val
		case "emote-sets":
			state.emoteSets = strings.Split(val, ",")
		case "mod":
			state.mod = val == "1"
		case "subscriber":
			state.subscriber = val == "1"
		case "turbo":
			state.turbo = val == "1"
		case "user-type":
			state.userType = UserType(val)
		}
	}
	channel.userState = state
}

// Badge returns the badge with the given name, if it exists.
func (s UserState) Badge(name string) (Badge, bool) {
	badge, ok := s.badges[strings.ToLower(name)]
	return badge, ok
}

// Color returns the users color, if the user has set one. Otherwise, it returns an empty string.
func (s UserState) Color() string {
	return s.color
}

// DisplayName returns the users display name.
func (s UserState) DisplayName() string {
	return s.displayName
}

// EmoteSets returns the emote sets the user has permission to use.
func (s UserState) EmoteSets() []string {
	return s.emoteSets
}

// IsBroadcaster returns true if the user is the broadcaster of the channel.
func (s UserState) IsBroadcaster() bool {
	_, ok := s.badges["broadcaster"]
	return ok
}

// IsMod returns true if the user is allowed to perform moderator actions in the channel.
func (s UserState) IsMod() bool {
	return s.mod || s.IsBroadcaster()
}

// IsSubscriber returns true if the user is a subscriber to the channel.
func (s UserState) IsSubscriber() bool {
	return s.subscriber
}

// IsTurbo returns true if the user has a turbo subscription.
// Turbo subscriptions allow users to bypass ads regardless of sub status.
func (s UserState) IsTurbo() bool {
	return s.turbo
}

// UserType returns the users type. For normal users, this is usually empty.
func (s UserState) UserType() UserType {
	return s.userType
}

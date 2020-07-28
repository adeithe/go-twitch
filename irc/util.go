package irc

import "strings"

// ToChannelName returns a lowercase string and removes prefixed #.
func ToChannelName(str string) string {
	return strings.ToLower(strings.TrimPrefix(str, "#"))
}

package pubsub

import "fmt"

// Topic formattable string based on known Twitch PubSub Topics.
// Visit https://dev.twitch.tv/docs/pubsub#topics for more information.
type Topic string

const (
	// ChannelBitEventsV1 - Anyone cheers in a specified channel.
	// Arguments: ChannelID
	// Scopes: bits:read
	ChannelBitEventsV1 Topic = "channel-bits-events-v1.%v"
	// ChannelBitEventsV2 - Anyone cheers in a specified channel.
	// Arguments: ChannelID
	// Scopes: bits:read
	ChannelBitEventsV2 Topic = "channel-bits-events-v2.%v"
	// BitsBadgeNotification - Message sent when a user earns a new Bits badge in a particular channel, and chooses to share the notification with chat.
	// Arguments: ChannelID
	// Scopes: bits:read
	BitsBadgeNotification Topic = "channel-bits-badge-unlocks.%v"
	// ChannelPoints - A custom reward is redeemed in a channel.
	// Arguments: ChannelID
	// Scopes: channel:read:redemptions
	ChannelPoints Topic = "channel-points-channel-v1.%v"
	// ChannelSubscriptions - Anyone subscribes (first month), resubscribes (subsequent months), or gifts a subscription to a channel. Subgift subscription messages contain recipient information.
	// Arguments: ChannelID
	// Scopes: channel_subscriptions
	ChannelSubscriptions Topic = "channel-subscribe-events-v1.%v"
	// ChatModeratorActions - A moderator performs an action in the channel.
	// Arguments: ChannelID
	// Scopes: channel:moderate
	ChatModeratorActions Topic = "chat_moderator_actions.%v"
	// Whispers - Anyone whispers the specified user.
	// Arguments: UserID
	// Scopes: whispers:read
	Whispers Topic = "whispers.%v"
)

// ParseTopic will parse the given Topic with the arguments provided.
func ParseTopic(topic Topic, args ...interface{}) string {
	return ParseTopicFromString(string(topic), args...)
}

// ParseTopicFromString allows you to parse a topic manually if it isn't available in Topic.
// This is the same as using fmt.Sprintf
func ParseTopicFromString(topic string, args ...interface{}) string {
	return fmt.Sprintf(topic, args...)
}

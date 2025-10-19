package api

import (
	"encoding/json"
	"math"
	"time"
)

// Pagination represents pagination information in a Twitch API response.
type Pagination struct {
	Cursor string `json:"cursor,omitempty"`
}

// DateRange represents a range of dates with a start and end time.
type DateRange struct {
	// StartedAt is the reporting windows start date.
	StartedAt time.Time `json:"started_at"`
	// EndedAt is the reporting windows end date.
	EndedAt time.Time `json:"ended_at"`
}

// SizedImage represents an image with multiple sizes.
type SizedImage struct {
	// Size1x is the URL to the 1x size image.
	Size1x string `json:"url_1x"`
	// Size2x is the URL to the 2x size image.
	Size2x string `json:"url_2x"`
	// Size4x is the URL to the 4x size image.
	Size4x string `json:"url_4x"`
}

// Commercial represents a commercial that was started.
type Commercial struct {
	// Length is the duration of the commercial in seconds.
	Length int `json:"length"`
	// Message is a message which indicates whether Twitch was able to serve an ad.
	Message string `json:"message"`
	// RetryAfter is the number of seconds the client should wait before trying to start another commercial.
	RetryAfter int `json:"retry_after"`
}

// AdSchedule represents the ad schedule for a broadcaster.
type AdSchedule struct {
	// DurationSeconds is the number of seconds the upcoming commercial will run.
	DurationSeconds int `json:"duration"`
	// PrerollFreeTime is the number of seconds of pre-roll free time remaining for the channel in seconds.
	PrerollFreeTime int `json:"preroll_free_time"`
	// SnoozeCount is the number of snoozes the broadcaster has available.
	SnoozeCount int `json:"snooze_count"`
	// SnoozeRefreshAt is the UTC timestamp for when the broadcaster will receive another snooze.
	SnoozeRefreshAt time.Time `json:"snooze_refresh_at"`
	// NextAdAt is the UTC timestamp for when the next commercial is scheduled to run.
	NextAdAt time.Time `json:"next_ad_at"`
	// LastAdAt is the UTC timestamp for when the last commercial ran.
	LastAdAt time.Time `json:"last_ad_at"`
}

// AdsSnoozed represents the result of a successful ad snooze request.
type AdsSnoozed struct {
	// SnoozeCount is the number of snoozes the broadcaster has available.
	SnoozeCount int `json:"snooze_count"`
	// SnoozeRefreshAt is the UTC timestamp for when the broadcaster will receive another snooze.
	SnoozeRefreshAt time.Time `json:"snooze_refresh_at"`
	// NextAdAt is the UTC timestamp for when the next commercial is scheduled to run.
	NextAdAt time.Time `json:"next_ad_at"`
}

// ExtensionAnalyticsReport represents analytics data for a Twitch extension.
type ExtensionAnalyticsReport struct {
	// ExtensionID is the ID of the extension the analytics report was generated for.
	ExtensionID string `json:"extension_id"`
	// URL is a temporary link to download the analytics report. Expires after 5 minutes.
	URL string `json:"URL"`
	// Type is the type of analytics report.
	Type string `json:"type"`
	// DateRange is the date range for which the analytics report was generated.
	DateRange DateRange `json:"date_range"`
}

// GameAnalyticsReport represents analytics data for a Twitch game.
type GameAnalyticsReport struct {
	// GameID is the ID of the game the analytics report was generated for.
	GameID string `json:"game_id"`
	// URL is a temporary link to download the analytics report. Expires after 5 minutes.
	URL string `json:"URL"`
	// Type is the type of analytics report.
	Type string `json:"type"`
	// DateRange is the date range for which the analytics report was generated.
	DateRate DateRange `json:"date_range"`
}

// BitsLeaderboardEntry represents a leaderboard of users who have spent the most bits in a channel.
type BitsLeaderboardEntry struct {
	// UserID is the ID of the user.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user.
	UserName string `json:"user_name"`
	// Rank is the leaderboard rank for the user in the channel.
	Rank int `json:"rank"`
	// Score is the total number of bits the user has spent in the channel.
	Score int `json:"score"`
}

// Cheermote represents a Twitch cheermote.
type Cheermote struct {
	// Prefix is the prefix of the cheermote.
	Prefix string `json:"prefix"`
	// Tiers is the list of tiers for the cheermote.
	Tiers []CheermoteTier `json:"tiers"`
	// Type is the type of the cheermote.
	Type string `json:"type"`
	// Order is the order of the cheermote in the bits card.
	Order int `json:"order"`
	// IsCharitable indicates whether the cheermote provides a charitable contribution match during charity campaigns.
	IsCharitable bool `json:"is_charitable"`
	// LastUpdated is the last time the cheermote was updated.
	LastUpdated time.Time `json:"last_updated"`
}

// CheermoteTier represents a tier of a Twitch cheermote.
type CheermoteTier struct {
	// ID is the tier for the cheermote.
	ID string `json:"id"`
	// MinBits is the minimum number of bits needed to use the cheermote.
	MinBits int `json:"min_bits"`
	// Color is the color associated with this tier.
	Color string `json:"color"`
	// Images is a map of image URLs for the cheermote in different sizes.
	Images map[string]string `json:"images"`
	// CanCheer indicates whether the cheermote can be used to cheer.
	CanCheer bool `json:"can_cheer"`
	// ShowInBitsCard indicates whether the cheermote is shown in the bits card.
	ShowInBitsCard bool `json:"show_in_bits_card"`
}

// ExtensionTransaction contains information about a transaction for an extension.
type ExtensionTransaction struct {
	// ID is the ID of the transaction.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster the transaction occurred in.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster the transaction occurred in.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster the transaction occurred in.
	BroadcasterName string `json:"broadcaster_name"`
	// UserID is the ID of the user that made the transaction.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user that made the transaction.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user that made the transaction.
	UserName string `json:"user_name"`
	// ProductType is the type of product that was purchased.
	ProductType string `json:"product_type"`
	// Product is information about the product that was purchased.
	Product ExtensionProduct `json:"product_data"`
	// Timestamp is the UTC timestamp of when the transaction occurred.
	Timestamp time.Time `json:"timestamp"`
}

// ExtensionProduct contains information about a product for an extension.
type ExtensionProduct struct {
	// SKU is the stock keeping unit that identifies the product.
	SKU string `json:"sku"`
	// Domain is the extension domain of the product.
	Domain string `json:"domain"`
	// Cost is the cost of the product.
	Cost ExtensionProductCost `json:"cost"`
	// InDevelopment indicates whether the product is in development and not currently available to the public.
	InDevelopment bool `json:"inDevelopment"`
	// DisplayName is the display name of the product.
	DisplayName string `json:"displayName"`
	// Broadcast indicates whether the data was broadcast to all instances of the extension.
	Broadcast bool `json:"broadcast"`
	// Expiration is the UTC timestamp of when the product expires.
	Expiration time.Time `json:"expiration"`
}

// ExtensionProductCost contains information about the cost of a product for an extension.
type ExtensionProductCost struct {
	// Amount is the cost of the product in bits.
	Amount int `json:"amount"`
	// Type is the type of currency used to purchase the product.
	Type string `json:"type"`
}

// Channel represents a Twitch channel.
type Channel struct {
	// GameID is the ID of the game being played on the channel.
	GameID string `json:"game_id"`
	// GameName is the name of the game being played on the channel.
	GameName string `json:"game_name"`
	// BroadcasterID is the ID of the broadcaster.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster.
	BroadcasterName string `json:"broadcaster_name"`
	// Title is stream title for the channel.
	Title string `json:"title"`
	// Delay is the delay of the stream in seconds.
	Delay int `json:"delay"`
	// Tags are the tags applied to the channel.
	Tags []string `json:"tags"`
	// ContentClassificationLabels are the CCLs applied to the channel.
	ContentClassificationLabels []string `json:"content_classification_labels"`
	// IsBrandedContent indicates whether the channel is marked as branded content such as sponorships.
	IsBrandedContent bool `json:"is_branded_content"`
}

// ChannelContentClassificationLabel represents a content classification label.
type ChannelContentClassificationLabel struct {
	// ID is the ID of the content classification label.
	ID string `json:"id"`
	// IsEnabled indicates whether the content classification label is enabled.
	IsEnabled bool `json:"is_enabled"`
}

// ChannelEditor represents a user who is an editor of a Twitch channel.
type ChannelEditor struct {
	// UserID is the ID of the user who is an editor of the channel.
	UserID string `json:"user_id"`
	// UserName is the display name of the user who is an editor of the channel.
	UserName string `json:"user_name"`
	// AddedAt is the UTC timestamp of when the user was added as an editor of the channel.
	AddedAt time.Time `json:"created_at"`
}

// Followed represents a Twitch channel followed by a channel.
type Followed struct {
	// BroadcasterID is the ID of the broadcaster being followed.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster being followed.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster being followed.
	BroadcasterName string `json:"broadcaster_name"`
	// FollowedAt is the UTC timestamp of when the channel was followed.
	FollowedAt time.Time `json:"followed_at"`
}

// Follower represents a user who follows a Twitch channel.
type Follower struct {
	// UserID is the ID of the user who follows the channel.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user who follows the channel.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user who follows the channel.
	UserName string `json:"user_name"`
	// FollowedAt is the UTC timestamp of when the user followed the channel.
	FollowedAt time.Time `json:"followed_at"`
}

// CustomReward represents a Twitch Channel Point custom reward.
type CustomReward struct {
	// RewardID the ID that uniquely identifies the custom reward.
	RewardID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that created the custom reward.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that created the custom reward.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that created the custom reward.
	BroadcasterName string `json:"broadcaster_name"`
	// BackgroundColor is a hex formatted color code representing the background color of the reward.
	BackgroundColor string `json:"background_color"`
	// Title is the title of the custom reward.
	Title string `json:"title"`
	// Prompt is the prompt for the custom reward.
	Prompt string `json:"prompt"`
	// Image is a set of images for the custom reward.
	Image *SizedImage `json:"image"`
	// DefaultImage is a set of default images for the custom reward.
	DefaultImage SizedImage `json:"default_image"`
	// MaxPerStreamSetting is the setting for the custom rewards maximum number of redemptions per stream.
	MaxPerStreamSetting CustomRewardMaxPerStreamSetting `json:"max_per_stream_setting"`
	// MaxPerUserPerStreamSetting is the setting for the custom rewards maximum number of redemptions per user per stream.
	MaxPerUserPerStreamSetting CustomRewardMaxPerUserPerStreamSetting `json:"max_per_user_per_stream_setting"`
	// GlobalCooldownSetting is the setting for the custom rewards global cooldown.
	GlobalCooldownSetting CustomRewardGlobalCooldownSetting `json:"global_cooldown_setting"`
	// Cost is the cost of the custom reward.
	Cost int64 `json:"cost"`
	// TimesRedeemedThisStream is the number of redemptions of the reward in the current stream.
	TimesRedeemedThisStream int `json:"redemptions_redeemed_current_stream"`
	// Enabled indicates whether the custom reward is enabled.
	Enabled bool `json:"is_enabled"`
	// Paused indicates whether the custom reward is paused.
	Paused bool `json:"is_paused"`
	// InStock indicates whether the custom reward is currently in stock.
	InStock bool `json:"is_in_stock"`
	// IsUserInputRequired indicates whether the custom reward requires user input.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// RedemptionsSkipRequestQueue indicates whether redemptions for the reward skip the request queue.
	RedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// CooldownExpiresAt is the UTC timestamp of when the cooldown for the reward expires.
	CooldownExpiresAt *time.Time `json:"cooldown_expires_at,omitempty"`
}

// CustomRewardMaxPerStreamSetting represents the maximum redemptions per stream setting for a custom reward.
type CustomRewardMaxPerStreamSetting struct {
	// Enabled  indicates whether the max redemptions per stream setting is enabled.
	Enabled bool `json:"is_enabled"`
	// Value is the maximum number of redemptions a user can make per stream.
	Value int64 `json:"max_per_stream"`
}

// CustomRewardMaxPerUserPerStreamSetting represents the per user per stream setting for a custom reward.
type CustomRewardMaxPerUserPerStreamSetting struct {
	// Enabled indicates whether the max per user per stream setting is enabled.
	Enabled bool `json:"is_enabled"`
	// Value is the maximum number of redemptions a user can make per stream.
	Value int64 `json:"max_per_user_per_stream"`
}

// CustomRewardGlobalCooldownSetting represents the global cooldown setting for a custom reward.
type CustomRewardGlobalCooldownSetting struct {
	// Enabled indicates whether the global cooldown setting is enabled.
	Enabled bool `json:"is_enabled"`
	// Value is the length of the global cooldown in seconds.
	Value int64 `json:"global_cooldown_seconds"`
}

// CustomRewardRedemption represents a Twitch Channel Point custom reward redemption.
type CustomRewardRedemption struct {
	// ID is the ID that uniquely identifies the custom reward redemption.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that created the custom reward.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that created the custom reward.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that created the custom reward.
	BroadcasterName string `json:"broadcaster_name"`
	// UserID is the ID of the user that redeemed the custom reward.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user that redeemed the custom reward.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user that redeemed the custom reward.
	UserName string `json:"user_name"`
	// UserInput is the user input provided by the user when redeeming the custom reward.
	UserInput string `json:"user_input"`
	// Status is the current status of the custom reward redemption.
	Status string `json:"status"`
	// Reward is basic information about the custom reward that was redeemed.
	Reward RedemptionRewardInfo `json:"reward"`
	// RedeemedAt is the UTC timestamp of when the custom reward was redeemed.
	RedeemedAt time.Time `json:"redeemed_at"`
}

// RedemptionRewardInfo represents basic reward info for a custom reward redemption.
type RedemptionRewardInfo struct {
	// RewardID is the ID of the custom reward.
	RewardID string `json:"id"`
	// Title is the title of the custom reward.
	Title string `json:"title"`
	// Prompt is the prompt of the custom reward.
	Prompt string `json:"prompt"`
	// Cost is the cost of the custom reward.
	Cost int64 `json:"cost"`
}

// CharityCampaign represents a charity campaign on a Twitch channel.
type CharityCampaign struct {
	// ID is the ID that uniquely identifies the charity campaign.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster running the charity campaign.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster running the charity campaign.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster running the charity campaign.
	BroadcasterName string `json:"broadcaster_name"`
	// CharityName is the name of the charity.
	CharityName string `json:"charity_name"`
	// CharityDescription is the description of the charity.
	CharityDescription string `json:"charity_description"`
	// CharityLogo is the URL to the logo of the charity.
	CharityLogo string `json:"charity_logo"`
	// CharityWebsite is the URL to the website of the charity.
	CharityWebsite string `json:"charity_website"`
	// CurrentAmount is the current amount of money raised by the charity campaign.
	CurrentAmount CharityCampaignAmount `json:"current_amount"`
	// TargetAmount is the target amount of money to be raised by the charity campaign.
	TargetAmount CharityCampaignAmount `json:"target_amount"`
}

// CharityCampaignDonation represents a donation made to a charity campaign.
type CharityCampaignDonation struct {
	// ID is the ID that uniquely identifies the charity campaign donation.
	ID string `json:"id"`
	// CampaignID is the ID of the charity campaign the donation was made to.
	CampaignID string `json:"campaign_id"`
	// UserID is the ID of the user that made the donation.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user that made the donation.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user that made the donation.
	UserName string `json:"user_name"`
	// Amount is the amount of money that was donated.
	Amount CharityCampaignAmount `json:"amount"`
}

// CharityCampaignAmount represents a monetary amount in a specific currency.
type CharityCampaignAmount struct {
	// Value is the amount in the currency's minor unit (e.g., cents for USD).
	Value int `json:"value"`
	// Decimal is the number of decimal places used by the currency. This number is used to translate value from minor units to major units.
	Decimal int `json:"decimal"`
	// Currency is the ISO-4217 three-letter currency code that identifies the type of currency (e.g., "USD").
	Currency string `json:"currency"`
}

// Emote represents a Twitch emote.
type Emote struct {
	// ID is the ID of the emote.
	ID string `json:"id"`
	// OwnerID is the ID of the user who owns the emote. May be empty.
	OwnerID string `json:"owner_id,omitempty"`
	// Name is the name of the emote.
	Name string `json:"name"`
	// Tier is the tier of the emote.
	Tier string `json:"tier,omitempty"`
	// Template is the template of the emote.
	Template string `json:"template"`
	// EmoteType is the type of the emote.
	EmoteType string `json:"emote_type,omitempty"`
	// EmoteSetID is the ID of the emote set the emote belongs to.
	EmoteSetID string `json:"emote_set_id,omitempty"`
	// Format is a list of supported image formats.
	Format []string `json:"format"`
	// Scale is a list of supported image scales.
	Scale []string `json:"scale"`
	// ThemeMode is a list of supported theme modes.
	ThemeMode []string `json:"theme_mode"`
	// Images is a set of sized images for the emote.
	Images SizedImage `json:"images"`
}

// ChatBadge represents a Twitch chat badge.
type ChatBadge struct {
	// SetID is the ID of the badge set the badge belongs to.
	SetID string `json:"set_id"`
	// Versions is a list of versions of the badge.
	Versions []ChatBadgeVersion `json:"versions"`
}

// ChatBadgeVersion represents a version of a Twitch chat badge.
type ChatBadgeVersion struct {
	// ID is the ID of the badge version.
	ID string `json:"id"`
	// ImageSize1X is the URL to the 1x size image of the badge.
	ImageSize1X string `json:"image_url_1x"`
	// ImageSize2X is the URL to the 2x size image of the badge.
	ImageSize2X string `json:"image_url_2x"`
	// ImageSize4X is the URL to the 4x size image of the badge.
	ImageSize4X string `json:"image_url_4x"`
	// Description is the description of the badge.
	Description string `json:"description"`
	// Title is the title of the badge.
	Title string `json:"title"`
	// ClickAction is the click action of the badge.
	ClickAction string `json:"click_action"`
	// ClickURL is the URL to open when the badge is clicked.
	ClickURL string `json:"click_url"`
}

// ChatSettings represents the chat settings for a Twitch channel.
type ChatSettings struct {
	// BroadcasterID is the ID of the broadcaster whose chat settings are being managed.
	BroadcasterID string `json:"broadcaster_id"`
	// EmoteMode indicates whether emote-only mode is enabled.
	EmoteMode bool `json:"emote_mode"`
	// FollowMode indicates whether follow-only mode is enabled.
	FollowMode bool `json:"follow_mode"`
	// FollowModeDurationMinutes is the duration in minutes that a user must follow the channel before they can chat.
	FollowModeDurationMinutes int `json:"follow_mode_duration"`
	// NonModeratorChatDelay indicates whether non-moderator chat delay is enabled.
	NonModeratorChatDelay bool `json:"non_moderator_chat_delay"`
	// NonModeratorChatDelaySeconds is the duration in seconds of the non-moderator chat delay.
	NonModeratorChatDelaySeconds int `json:"non_moderator_chat_delay_duration"`
	// SlowMode indicates whether slow mode is enabled.
	SlowMode bool `json:"slow_mode"`
	// SlowModeWaitTimeSeconds is the duration in seconds that a user must wait between sending messages in slow mode.
	SlowModeWaitTimeSeconds int `json:"slow_mode_wait_time"`
	// SubscriberMode indicates whether subscriber-only mode is enabled.
	SubscriberMode bool `json:"subscriber_mode"`
	// UniqueChatMode indicates whether unique chat mode is enabled.
	UniqueChatMode bool `json:"unique_chat_mode"`
}

// SharedChatSession represents a shared chat session.
type SharedChatSession struct {
	// SessionID is the ID of the shared chat session.
	SessionID string `json:"session_id"`
	// HostBroadcasterID is the ID of the broadcaster that created the shared chat session.
	HostBroadcasterID string `json:"host_broadcaster_id"`
	// Participants is a list of participants in the shared chat session.
	Participants []SharedChatSessionParticipant `json:"participants"`
	// CreatedAt is the UTC timestamp of when the shared chat session was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the UTC timestamp of when the shared chat session was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// SharedChatSessionParticipant represents a participant in a shared chat session.
type SharedChatSessionParticipant struct {
	// BroadcasterID is the ID of the broadcaster that is a participant in the shared chat session.
	BroadcasterID string `json:"broadcaster_id"`
}

// OutboundChatMessage represents the result of an attempt to send a message in chat.
type OutboundChatMessage struct {
	// MessageID is the ID of the message that was sent.
	MessageID string `json:"message_id"`
	// IsSent indicates whether the message was successfully sent.
	IsSent bool `json:"is_sent"`
	// DropReason is the reason the message was not sent, if applicable.
	DropReason *OutboundChatMessageDropReason `json:"drop_reason,omitempty"`
}

// OutboundChatMessageDropReason represents the reason a chat message was not sent.
type OutboundChatMessageDropReason struct {
	// Code is the code representing the reason the message was not sent.
	Code string `json:"code"`
	// Message is a message describing the reason the message was not sent.
	Message string `json:"message"`
}

// UserChatColor represents a user's chat color.
type UserChatColor struct {
	// UserID is the ID of the user.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user.
	UserName string `json:"user_name"`
	// Color is the hex color code of the user's chat color.
	Color string `json:"color"`
}

// EditableClip represents links to edit a Twitch clip.
type EditableClip struct {
	// ID is the ID that uniquely identifies the clip.
	ID string `json:"id"`
	// EditURL is a URL that can be used to edit the clip's title and visibility settings.
	EditURL string `json:"edit_url"`
}

// Clip represents a Twitch clip.
type Clip struct {
	// ID is the ID that uniquely identifies the clip.
	ID string `json:"id"`
	// CreatorID is the ID of the user who created the clip.
	CreatorID string `json:"creator_id"`
	// CreatorName is the display name of the user who created the clip.
	CreatorName string `json:"creator_name"`
	// BroadcasterID is the ID of the broadcaster on whose channel the clip was created.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterName is the display name of the broadcaster on whose channel the clip was created.
	BroadcasterName string `json:"broadcaster_name"`
	// VideoID is the ID of the video from which the clip was created.
	VideoID string `json:"video_id"`
	// GameID is the ID of the game being played on the channel when the clip was created.
	GameID string `json:"game_id"`
	// Language is the language of the channel when the clip was created.
	Language string `json:"language"`
	// Title is the title of the clip.
	Title string `json:"title"`
	// URL is the URL of the clip.
	URL string `json:"url"`
	// EmbedURL is the URL that can be used to embed the clip.
	EmbedURL string `json:"embed_url"`
	// ThumbnailURL is the URL of the clip's thumbnail image.
	ThumbnailURL string `json:"thumbnail_url"`
	// VODOffset is the offset in seconds from the start of the VOD where the clip was created.
	VODOffset int `json:"vod_offset"`
	// ViewCount is the number of times the clip has been viewed.
	ViewCount int `json:"view_count"`
	// Duration is the duration of the clip. Value has 0.1 second precision.
	Duration float64 `json:"duration"`
	// CreatedAt is the UTC timestamp of when the clip was created.
	Featured bool `json:"featured"`
	// CreatedAt is the UTC timestamp of when the clip was created.
	CreatedAt time.Time `json:"created_at"`
}

// DownloadableClip represents downloadable URLs for a Twitch clip.
type DownloadableClip struct {
	// ClipID is the ID that uniquely identifies the clip.
	ClipID string `json:"clip_id"`
	// LandscapeDownloadURL is the URL to download the clip in landscape format.
	LandscapeDownloadURL string `json:"landscape_download_url"`
	// PortraitDownloadURL is the URL to download the clip in portrait format.
	PortraitDownloadURL string `json:"portrait_download_url"`
}

// Conduit represents a Twitch Eventsub Conduit.
type Conduit struct {
	// ID is the ID of the conduit.
	ID string `json:"id"`
	// ShardCount is the number of shards associated with the conduit.
	ShardCount int `json:"shard_count"`
}

// ConduitShard is a shard for a Twitch Eventsub Conduit.
type ConduitShard struct {
	// ID is the ID of the conduit shard.
	ID string `json:"id"`
	// Status is the current status of the conduit shard.
	Status string `json:"status,omitempty"`
	// Transport is the transport method for the conduit shard.
	Transport Transport `json:"transport"`
}

// ConduitError represents an error returned by an EventSub Conduits API endpoint.
type ConduitError struct {
	// ShardID is the ID of the shard the error is for.
	ShardID string `json:"id,omitempty"`
	// Code is used to represent a specific error condition while attempting to update shards.
	Code string `json:"code,omitempty"`
	// Message is a human-readable description of the error condition.
	Message string `json:"message,omitempty"`
}

// Transport is the transport method for a Twitch Eventsub Conduit Shard.
type Transport struct {
	// Method is the transport method. Possible values are "webhook" and "websocket".
	Method string `json:"method"`
	// Callback is the callback URL for webhook transports.
	Callback *string `json:"callback,omitempty"`
	// Secret is the secret for webhook transports.
	Secret *string `json:"secret,omitempty"`
	// SessionID is the session ID for websocket transports.
	SessionID *string `json:"session_id,omitempty"`
	// ConnectedAt is the UTC timestamp of when the websocket transport connected.
	ConnectedAt *time.Time `json:"connected_at,omitempty"`
	// DisconnectedAt is the UTC timestamp of when the websocket transport disconnected.
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty"`
}

// ContentClassificationLabel represents a content classification label.
type ContentClassificationLabel struct {
	// ID is the ID of the content classification label.
	ID string `json:"id"`
	// Name is the name of the content classification label.
	Name string `json:"name"`
	// Description is the description of the content classification label.
	Description string `json:"description"`
}

// DropEntitlement represents a Twitch drop entitlement.
type DropEntitlement struct {
	// ID is the ID that uniquely identifies the drop entitlement.
	ID string `json:"id"`
	// GameID is the ID of the game associated with the drop entitlement.
	GameID string `json:"game_id"`
	// BenefitID is the ID of the benefit associated with the drop entitlement.
	BenefitID string `json:"benefit_id"`
	// UserID is the ID of the user who has the drop entitlement.
	UserID string `json:"user_id"`
	// FulfillmentStatus is the fulfillment status of the drop entitlement.
	FulfillmentStatus string `json:"fulfillment_status"`
	// LastUpdatedAt is the UTC timestamp of when the drop entitlement was last updated.
	LastUpdatedAt time.Time `json:"last_updated"`
	// Timestamp is the UTC timestamp of when the drop entitlement was granted.
	Timestamp time.Time `json:"timestamp"`
}

// UpdatedDropEntitlement represents the result of an attempt to update drop entitlements.
type UpdatedDropEntitlement struct {
	// IDs is a list of IDs that were successfully updated.
	IDs []string `json:"ids"`
	// Status is the status of the update request.
	Status string `json:"status"`
}

// ExtensionConfigurationSegment represents a configuration segment for a Twitch extension.
type ExtensionConfigurationSegment struct {
	// BroadcasterID is the ID of the broadcaster that owns the extension.
	BroadcasterID string `json:"broadcaster_id"`
	// Segment is the segment of the extension configuration. Possible values are "broadcaster", "developer", and "global".
	Segment string `json:"segment"`
	// Content is the content of the extension configuration segment.
	Content string `json:"content"`
	// Version is the version of the extension configuration segment.
	Version string `json:"version"`
}

// ExtensionLiveChannel represents a live channel for a Twitch extension.
type ExtensionLiveChannel struct {
	// BroadcasterID is the ID of the broadcaster that has the extension live.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that has the extension live.
	BroadcasterName string `json:"broadcaster_name"`
	// GameID is the ID of the game being played on the channel.
	GameID string `json:"game_id"`
	// GameName is the name of the game being played on the channel.
	GameName string `json:"game_name"`
	// Title is the stream title for the channel.
	Title string `json:"title"`
}

// ExtensionSecret represents the secrets for a Twitch extension.
type ExtensionSecret struct {
	// FormatVersion is the version of the secret format.
	FormatVersion int `json:"format_version"`
	// Secrets is a list of secrets for the extension.
	Secrets []ExtensionSecretItem `json:"secrets"`
}

// ExtensionSecretItem represents a single secret for a Twitch extension.
type ExtensionSecretItem struct {
	// Content is the content of the secret.
	Content string `json:"content"`
	// ActiveAt is the UTC timestamp of when the secret was activated.
	ActiveAt time.Time `json:"active_at"`
	// ExpiresAt is the UTC timestamp of when the secret expires.
	ExpiresAt time.Time `json:"expires_at"`
}

// Extension represents a Twitch extension.
type Extension struct {
	// ID is the ID that uniquely identifies the extension.
	ID string `json:"id"`
	// Name is the name of the extension.
	Name string `json:"name"`
	// IconURL is the URL to the icon of the extension.
	IconURL string `json:"icon_url"`
	// AuthorName is the name of the extension author.
	AuthorName string `json:"author_name"`
	// Description is the description of the extension.
	Description string `json:"description"`
	// PrivacyPolicyURL is the URL to the privacy policy of the extension.
	PrivacyPolicyURL string `json:"privacy_policy_url"`
	// TermsURL is the URL to the terms of service of the extension.
	TermsURL string `json:"eula_tos_url"`
	// ConfigurationLocation is the location where the extension configuration is stored.
	ConfigurationLocation string `json:"configuration_location"`
	// SubscriptionsSupportLevel is the level of subscription support for the extension.
	SubscriptionsSupportLevel string `json:"subscriptions_support_level"`
	// Summary is the summary of the extension.
	Summary string `json:"summary"`
	// SupportEmail is the support email for the extension.
	SupportEmail string `json:"support_email"`
	// State is the current state of the extension.
	State string `json:"state"`
	// Version is the current version of the extension.
	Version string `json:"version"`
	// ViewerSummary is the summary of the extension for viewers.
	ViewerSummary string `json:"viewer_summary"`
	// AllowedConfigURLs is a list of allowed configuration URLs for the extension.
	AllowedConfigURLs []string `json:"allowlisted_config_urls"`
	// AllowedPanelURLs is a list of allowed panel URLs for the extension.
	AllowedPanelURLs []string `json:"allowlisted_panel_urls"`
	// ScreenshotURLs is a list of screenshot URLs for the extension.
	ScreenshotURLs []string `json:"screenshot_urls"`
	// IconURLs is a map of icon URLs for the extension in different sizes.
	IconURLs map[string]string `json:"icon_urls"`
	// Views is the different views for the extension.
	Views ExtensionViews `json:"views"`
	// RequestIdentityLink indicates whether the extension requests identity linking for users.
	RequestIdentityLink bool `json:"request_identity_link"`
	// HasChatSupport indicates whether the extension has chat support.
	HasChatSupport bool `json:"has_chat_support"`
	// BitsEnabled indicates whether the extension has bits enabled.
	BitsEnabled bool `json:"bits_enabled"`
	// CanInstall indicates whether the extension can be installed by broadcasters.
	CanInstall bool `json:"can_install"`
}

// ExtensionViews represents the different views for a Twitch extension.
type ExtensionViews struct {
	// Mobile is the mobile view for the extension.
	Mobile *ExtensionMobileView `json:"mobile,omitempty"`
	// Panel is the panel view for the extension.
	Panel *ExtensionPanelView `json:"panel,omitempty"`
	// VideoOverlay is the video overlay view for the extension.
	VideoOverlay *ExtensionVideoOverlayView `json:"video_overlay,omitempty"`
	// Component is the component view for the extension.
	Component *ExtensionComponentView `json:"component,omitempty"`
}

// ExtensionMobileView represents the different views for a Twitch extension.
type ExtensionMobileView struct {
	// ViewerURL is the URL to load the extension in a mobile view.
	ViewerURL string `json:"viewer_url"`
}

// ExtensionPanelView represents the panel view for a Twitch extension.
type ExtensionPanelView struct {
	// ViewerURL is the URL to load the extension in a panel view.
	ViewerURL string `json:"viewer_url"`
	// Width is the width of the panel in pixels.
	Height int `json:"height"`
	// CanLinkExternalContent indicates whether the panel can link to external content.
	CanLinkExternalContent bool `json:"can_link_external_content"`
}

// ExtensionVideoOverlayView represents the video overlay view for a Twitch extension.
type ExtensionVideoOverlayView struct {
	// ViewerURL is the URL to load the extension in a video overlay view.
	ViewerURL string `json:"viewer_url"`
	// CanLinkExternalContent indicates whether the video overlay can link to external content.
	CanLinkExternalContent bool `json:"can_link_external_content"`
}

// ExtensionComponentView represents the component view for a Twitch extension.
type ExtensionComponentView struct {
	// ViewerURL is the URL to load the extension in a component view.
	ViewerURL string `json:"viewer_url"`
	// AspectWidth is the aspect width of the component.
	AspectWidth int `json:"aspect_width"`
	// AspectHeight is the aspect height of the component.
	AspectHeight int `json:"aspect_height"`
	// AspectRatioX is the aspect ratio X of the component.
	AspectRatioX int `json:"aspect_ratio_x"`
	// AspectRatioY is the aspect ratio Y of the component.
	AspectRatioY int `json:"aspect_ratio_y"`
	// ScalePixels is the scale in pixels of the component.
	ScalePixels int `json:"scale_pixels"`
	// TargetWidth is the target width of the component in pixels.
	TargetHeight int `json:"target_height"`
	// Size is the size of the component in pixels.
	Size int `json:"size"`
	// TargetHeight is the target height of the component in pixels.
	ZoomPixels int `json:"zoom_pixels"`
	// Zoom is whether the component is zoomed.
	Zoom bool `json:"zoom"`
	// Autoscale is whether the component autoscales.
	Autoscale bool `json:"autoscale"`
	// CanLinkExternalContent indicates whether the component can link to external content.
	CanLinkExternalContent bool `json:"can_link_external_content"`
}

// ExtensionBitsProduct represents a bits product for a Twitch extension.
type ExtensionBitsProduct struct {
	// SKU is the stock keeping unit that identifies the product.
	SKU string `json:"sku"`
	// DisplayName is the display name of the product.
	DisplayName string `json:"display_name"`
	// Cost is the cost of the product in bits.
	Cost ExtensionBitsProductCost `json:"cost"`
	// InDevelopment indicates whether the product is in development and not currently available to the public.
	InDevelopment bool `json:"in_development"`
	// IsBroadcast indicates whether the product is broadcast to all instances of the extension.
	IsBroadcast bool `json:"is_broadcast"`
	// ExpiresAt is the UTC timestamp of when the product expires.
	ExpiresAt *time.Time `json:"expiration,omitempty"`
}

// ExtensionBitsProductCost contains information about the cost of a bits product for a Twitch extension.
type ExtensionBitsProductCost struct {
	// Amount is the cost of the product in bits.
	Amount int `json:"amount"`
	// Type is the type of currency used to purchase the product.
	Type string `json:"type"`
}

// EventSubSubscription represents a Twitch EventSub subscription.
type EventSubSubscription struct {
	// ID is the ID that uniquely identifies the subscription.
	ID string `json:"id"`
	// Status is the current status of the subscription.
	Status string `json:"status"`
	// Type is the type of event the subscription is for.
	Type string `json:"type"`
	// Version is the version of the subscription type.
	Version string `json:"version"`
	// Condition is the condition that triggers the event.
	Condition map[string]any `json:"condition"`
	// Transport is the transport information for the subscription.
	Transport EventSubTransport `json:"transport"`
	// Cost is the cost of the subscription in points.
	Cost int `json:"cost"`
	// TotalCost is the total cost of the subscription in points.
	TotalCost int `json:"total_cost"`
	// MaxTotalCost is the maximum total cost of the subscription in points.
	MaxTotalCost int `json:"max_total_cost"`
	// CreatedAt is the UTC timestamp of when the subscription was created.
	CreatedAt time.Time `json:"created_at"`
}

// EventSubTransport represents the transport information for a Twitch EventSub subscription.
type EventSubTransport struct {
	// Method is the transport method.
	Method string `json:"method"`
	// Callback is the callback URL for webhook transports.
	Callback string `json:"callback,omitempty"`
	// SessionID is the session ID for websocket transports.
	SessionID string `json:"session_id,omitempty"`
	// ConduitID is the ID of the conduit being used for an EventSub transport.
	ConduitID string `json:"conduit_id,omitempty"`
}

// Game represents a game on Twitch.
type Game struct {
	// ID is the ID that uniquely identifies the game.
	ID string `json:"id"`
	// Name is the name of the game.
	Name string `json:"name"`
	// BoxArtURL is the URL to the box art of the game.
	BoxArtURL string `json:"box_art_url"`
	// IGDB is the ID of the game in the IGDB database.
	IGDB string `json:"igdb_id"`
}

// CreatorGoal represents a creator goal on a Twitch channel.
type CreatorGoal struct {
	// ID is the ID that uniquely identifies the creator goal.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that created the creator goal.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that created the creator goal.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that created the creator goal.
	BroadcasterName string `json:"broadcaster_name"`
	// Type is the type of the creator goal.
	Type string `json:"type"`
	// Description is the description of the creator goal.
	Description string `json:"description"`
	// CurrentAmount is the current amount of progress towards the creator goal.
	CurrentAmount int `json:"current_amount"`
	// TargetAmount is the target amount of the creator goal.
	TargetAmount int `json:"target_amount"`
	// CreatedAt is the UTC timestamp of when the creator goal was created.
	CreatedAt time.Time `json:"created_at"`
}

// ChannelGuestStarSettings represents the guest star settings for a Twitch channel.
type ChannelGuestStarSettings struct {
	// GroupLayout is the layout of the guest star group.
	GroupLayout string `json:"group_layout"`
	// BrowserSourceToken is the token used to authenticate the browser source for guest stars.
	BrowserSourceToken string `json:"browser_source_token"`
	// IsBrowserSourceEnabled indicates whether the browser source for guest stars is enabled.
	IsBrowserSourceAudioEnabled bool `json:"is_browser_source_audio_enabled"`
	// IsModeratorSendLiveEnabled indicates whether moderators can send live invitations to guest stars.
	IsModeratorSendLiveEnabled bool `json:"is_moderator_send_live_enabled"`
	// SlotCount is the number of slots available for guest stars.
	SlotCount int `json:"slot_count"`
}

// GuestStarSession represents a guest star session on a Twitch channel.
type GuestStarSession struct {
	// ID is the ID that uniquely identifies the guest star session.
	ID string `json:"id"`
	// Guests is a list of guests in the guest star session.
	Guests []GuestStarGuest `json:"guests"`
}

// GuestStarGuest represents a guest in a guest star session.
type GuestStarGuest struct {
	// SlotID is the ID of the slot assigned to the guest.
	SlotID string `json:"slot_id"`
	// IsLive indicates whether the guest is currently live in the session.
	IsLive bool `json:"is_live"`
	// UserID is the ID of the guest user.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the guest user.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the guest user.
	UserName string `json:"user_display_name"`
	// Volume is the volume level of the guest user.
	Volume int `json:"volume"`
	// AudioSettings is the audio settings for the guest user.
	AudioSettings MediaSettings `json:"audio_settings"`
	// VideoSettings is the video settings for the guest user.
	VideoSettings MediaSettings `json:"video_settings"`
	// AssignedAt is the UTC timestamp of when the guest was assigned to the slot.
	AssignedAt time.Time `json:"assigned_at"`
	// JoinedAt is the UTC timestamp of when the guest joined the session.
	JoinedAt time.Time `json:"joined_at"`
}

// MediaSettings represents the media settings.
type MediaSettings struct {
	// IsHostEnabled indicates whether the host is enabled.
	IsHostEnabled bool `json:"is_host_enabled"`
	// IsGuestEnabled indicates whether the guest is enabled.
	IsGuestEnabled bool `json:"is_guest_enabled"`
	// IsAvailable indicates whether the media is available.
	IsAvailable bool `json:"is_available"`
}

// GuestStarInvite represents an invite to a guest star session.
type GuestStarInvite struct {
	// UserID is the ID of the user who was invited.
	UserID string `json:"user_id"`
	// Status is the status of the invite.
	Status string `json:"status"`
	// IsAudioAvailable indicates whether audio is available for the invite.
	IsAudioAvailable bool `json:"is_audio_available"`
	// IsVideoAvailable indicates whether video is available for the invite.
	IsVideoAvailable bool `json:"is_video_available"`
	// IsAudioEnabled indicates whether audio is enabled for the invite.
	IsAudioEnabled bool `json:"is_audio_enabled"`
	// IsVideoEnabled indicates whether video is enabled for the invite.
	IsVideoEnabled bool `json:"is_video_enabled"`
	// InvitedAt is the UTC timestamp of when the user was invited.
	InvitedAt time.Time `json:"invited_at"`
}

// HypeTrainEvent represents a Hype Train event on a Twitch channel.
type HypeTrainEvent struct {
	// ID is the ID that uniquely identifies the Hype Train event.
	ID string `json:"id"`
	// EventType is the type of the Hype Train event.
	EventType string `json:"event_type"`
	// Version is the version of the Hype Train event.
	Version string `json:"version"`
	// EventData provides information about the Hype Train event.
	EventData HypeTrainEventData `json:"event_data"`
	// EventTimestamp is the UTC timestamp of when the Hype Train event occurred.
	EventTimestamp time.Time `json:"event_timestamp"`
}

// HypeTrainEventData provides information about a Hype Train event.
type HypeTrainEventData struct {
	// ID is the ID that uniquely identifies the Hype Train.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that the Hype Train event occurred on.
	BroadcasterID string `json:"broadcaster_id"`
	// Goal is the number of points required to reach the next level of the Hype Train.
	Goal int `json:"goal"`
	// Total is the total number of points contributed to the Hype Train so far.
	Total int `json:"total"`
	// Level is the current level of the Hype Train.
	Level int `json:"level"`
	// LastContribution is information about the last contribution made to the Hype Train.
	LastContribution HypeTrainEventContribution `json:"last_contribution"`
	// TopContributions is a list of the top contributions made to the Hype Train.
	TopContributions []HypeTrainEventContribution `json:"top_contributions"`
	// CooldownEndTime is the UTC timestamp of when the Hype Train cooldown ends.
	CooldownEndTime time.Time `json:"cooldown_end_time"`
	// StartedAt is the UTC timestamp of when the Hype Train started.
	StartedAt time.Time `json:"started_at"`
	// ExpiresAt is the UTC timestamp of when the Hype Train expires.
	ExpiresAt time.Time `json:"expires_at"`
}

// HypeTrainEventContribution represents a contribution made to a Hype Train.
type HypeTrainEventContribution struct {
	// Type is the type of contribution.
	Type string `json:"type"`
	// UserID is the ID of the user who made the contribution.
	UserID string `json:"user"`
	// Total is the total number of points contributed by the user.
	Total int `json:"total"`
}

// HypeTrainStatusInfo represents the status information of a Hype Train on a Twitch channel.
type HypeTrainStatusInfo struct {
	// Current is the current status of the Hype Train, if one is active.
	Current *HypeTrainStatus `json:"current,omitempty"`
	// AllTimeHigh is the highest level reached during a Hype Train on the channel.
	AllTimeHigh *HypeTrainRecord `json:"all_time_high,omitempty"`
	// SharedAllTimeHigh is the highest level reached during a shared Hype Train involving multiple channels.
	SharedAllTimeHigh *HypeTrainRecord `json:"shared_all_time_high,omitempty"`
}

// HypeTrainStatus represents the status of a Hype Train on a Twitch channel.
type HypeTrainStatus struct {
	// ID is the ID that uniquely identifies the Hype Train.
	ID string `json:"id"`
	// Type is the type of the Hype Train.
	Type string `json:"type"`
	// BroadcasterID is the ID of the broadcaster that the Hype Train is occurring on.
	BroadcasterID string `json:"broadcaster_user_id"`
	// BroadcasterLogin is the login name of the broadcaster that the Hype Train is occurring on.
	BroadcasterLogin string `json:"broadcaster_user_login"`
	// BroadcasterName is the display name of the broadcaster that the Hype Train is occurring on.
	BroadcasterName string `json:"broadcaster_user_name"`
	// Level is the level of the Hype Train.
	Level int `json:"level"`
	// Total is the total number of points contributed to the Hype Train so far.
	Total int `json:"total"`
	// Progress is the progress towards the next level of the Hype Train.
	Progress int `json:"progress"`
	// Goal is the number of points required to reach the next level of the Hype Train.
	Goal int `json:"goal"`
	// TopContributions is a list of the top contributions made to the Hype Train.
	TopContributions []HypeTrainContribution `json:"top_contributions"`
	// SharedTrainParticipants is a list of broadcaster IDs that are participants in a shared Hype Train.
	SharedTrainParticipants []string `json:"shared_train_participants"`
	// StartedAt is the UTC timestamp of when the Hype Train started.
	StartedAt time.Time `json:"started_at"`
	// ExpiresAt is the UTC timestamp of when the Hype Train expires.
	ExpiresAt time.Time `json:"expires_at"`
}

// HypeTrainContribution represents a contribution made to a Hype Train.
type HypeTrainContribution struct {
	// UserID is the ID of the user who made the contribution.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user who made the contribution.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user who made the contribution.
	UserName string `json:"user_name"`
	// Type is the type of contribution.
	Type string `json:"type"`
	// Total is the total number of points contributed by the user.
	Total int `json:"total"`
}

// HypeTrainParticipant represents a participant in a shared Hype Train.
type HypeTrainParticipant struct {
	// BroadcasterID is the ID of the broadcaster that is a participant in the shared Hype Train.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that is a participant in the shared Hype Train.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that is a participant in the shared Hype Train.
	BroadcasterName string `json:"broadcaster_name"`
}

// HypeTrainRecord represents a record of a completed Hype Train on a Twitch channel.
type HypeTrainRecord struct {
	// Level is the highest level reached during the Hype Train.
	Level int `json:"level"`
	// Total is the total number of points contributed during the Hype Train.
	Total int `json:"total"`
	// AchievedAt is the UTC timestamp of when the Hype Train was completed.
	AchievedAt time.Time `json:"achieved_at"`
}

// AutoModStatus represents the status of an AutoMod message review.
type AutoModStatus struct {
	// MessageID is the ID of the message being reviewed.
	MessageID string `json:"msg_id"`
	// Permitted indicates whether the message would be permitted or blocked by AutoMod.
	Permitted bool `json:"is_permitted"`
}

// AutoModSettings represents the AutoMod settings for a Twitch channel.
type AutoModSettings struct {
	// BroadcasterID is the ID of the broadcaster whose AutoMod settings are being managed.
	BroadcasterID string `json:"broadcaster_id"`
	// ModeratorID is the ID of the moderator managing the AutoMod settings.
	ModeratorID string `json:"moderator_id"`
	// OverallLevel is the overall AutoMod level setting.
	OverallLevel int `json:"overall_level"`
	// Ableism is the AutoMod level setting for messages containing references to disability.
	Ableism int `json:"disability"`
	// Aggressive is the AutoMod level setting for messages containing aggressive content.
	Aggressive int `json:"aggression"`
	// SexualDiscrimination is the AutoMod level setting for messages containing discrimination based on sexuality, sex, or gender.
	SexualDiscrimination int `json:"sexuality_sex_or_gender"`
	// MysogynySexism is the AutoMod level setting for messages containing misogynistic or sexist language.
	MisogynySexism int `json:"misogyny"`
	// Bullying is the AutoMod level setting for messages containing bullying or harassment.
	Bullying int `json:"bullying"`
	// Swearing is the AutoMod level setting for messages containing swearing or profanity.
	Swearing int `json:"swearing"`
	// RacismEthnicity is the AutoMod level setting for messages containing racism or ethnic discrimination.
	RacismEthnicity int `json:"race_ethnicity_or_religion"`
	// SexualTerms is the AutoMod level setting for messages containing sexual content or innuendo.
	SexualTerms int `json:"sex_based_terms"`
}

// BannedUser represents a user who has been banned from a Twitch channel.
type BannedUser struct {
	// ModeratorID is the ID of the moderator who issued the ban.
	ModeratorID string `json:"moderator_id"`
	// ModeratorLogin is the login name of the moderator who issued the ban.
	ModeratorLogin string `json:"moderator_login"`
	// ModeratorName is the display name of the moderator who issued the ban.
	ModeratorName string `json:"moderator_name"`
	// UserID is the ID of the user who was banned.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user who was banned.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user who was banned.
	UserName string `json:"user_name"`
	// Reason is the reason the user was banned.
	Reason string `json:"reason"`
	// CreatedAt is the UTC timestamp of when the ban was issued.
	CreatedAt time.Time `json:"created_at"`
	// ExpiresAt is the UTC timestamp of when the ban expires, if applicable.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// OutboundBan represents a request to ban a user from a Twitch channel.
type OutboundBan struct {
	// UserID is the ID of the user to be banned.
	UserID string `json:"user_id"`
	// Reason is the reason for the ban. Maximum length: 500 characters.
	Reason string `json:"reason,omitempty"`
	// Duration is the duration of the ban in seconds. If not specified, the ban is permanent.
	Duration *int `json:"duration,omitempty"`
}

// IssuedBan represents a ban that has been issued on a Twitch channel.
type IssuedBan struct {
	// BroadcasterID is the ID of the broadcaster that issued the ban.
	BroadcasterID string `json:"broadcaster_id"`
	// ModeratorID is the ID of the moderator that issued the ban.
	ModeratorID string `json:"moderator_id"`
	// UserID is the ID of the user that was banned.
	UserID string `json:"user_id"`
	// CreatedAt is the UTC timestamp of when the ban was issued.
	CreatedAt time.Time `json:"created_at"`
	// ExpiresAt is the UTC timestamp of when the ban expires, if applicable.
	ExpiresAt *time.Time `json:"end_time,omitempty"`
}

// UnbanRequest represents a users request to be unbanned from a Twitch channel.
type UnbanRequest struct {
	// ID is the ID that uniquely identifies the unban request.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that the unban request was sent to.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that the unban request was sent to.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that the unban request was sent to.
	BroadcasterName string `json:"broadcaster_name"`
	// ModeratorID is the ID of the moderator.
	ModeratorID string `json:"moderator_id"`
	// ModeratorLogin is the login name of the moderator.
	ModeratorLogin string `json:"moderator_login"`
	// ModeratorName is the display name of the moderator.
	ModeratorName string `json:"moderator_name"`
	// UserID is the ID of the user that sent the unban request.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user that sent the unban request.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user that sent the unban request.
	UserName string `json:"user_name"`
	// Text is the message from the user that sent the unban request.
	Text string `json:"text"`
	// Status is the status of the unban request.
	Status string `json:"status"`
	// ResolutionText is the message from the moderator that resolved the unban request.
	ResolutionText string `json:"resolution_text,omitempty"`
	// CreatedAt is the UTC timestamp of when the unban request was created.
	CreatedAt time.Time `json:"created_at"`
	// ResolvedAt is the UTC timestamp of when the unban request was resolved, if applicable.
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

// BlockedTerm represents a term that has been blocked in a Twitch channel.
type BlockedTerm struct {
	// ID is the ID that uniquely identifies the blocked term.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster the blocked term belongs to.
	BroadcasterID string `json:"broadcaster_id"`
	// ModeratorID is the ID of the moderator that blocked the term.
	ModeratorID string `json:"moderator_id"`
	// Text is the blocked term.
	Text string `json:"text"`
	// CreatedAt is the UTC timestamp of when the blocked term was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the UTC timestamp of when the blocked term was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// ExpiresAt is the UTC timestamp of when the blocked term expires, if applicable.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ModeratedChannel represents a Twitch channel a user moderates.
type ModeratedChannel struct {
	// BroadcasterID is the ID of the broadcaster that the user moderates.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that the user moderates.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that the user moderates.
	BroadcasterName string `json:"broadcaster_name"`
}

// UserInfo represents basic information about a Twitch user.
type UserInfo struct {
	// UserID is the ID of the user.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user.
	UserName string `json:"user_name"`
}

// ShieldModeStatus represents the status of Shield Mode on a Twitch channel.
type ShieldModeStatus struct {
	// ModeratorID is the ID of the moderator that changed the Shield Mode status.
	ModeratorID string `json:"moderator_id"`
	// ModeratorLogin is the login name of the moderator that changed the Shield Mode status.
	ModeratorLogin string `json:"moderator_login"`
	// ModeratorName is the display name of the moderator that changed the Shield Mode status.
	ModeratorName string `json:"moderator_name"`
	// Active indicates whether Shield Mode is currently active on the channel.
	Active bool `json:"is_active"`
	// LastActivatedAt is the UTC timestamp of when Shield Mode was last activated.
	LastActivatedAt time.Time `json:"last_activated_at"`
}

// UserWarning represents a warning issued to a user on a Twitch channel.
type UserWarning struct {
	// BroadcasterID is the ID of the channel in which the warning will take effect.
	BroadcasterID string `json:"broadcaster_id"`
	// ModeratorID is the ID of the moderator that issued the warning.
	ModeratorID string `json:"moderator_id"`
	// UserID is the ID of the user that was warned.
	UserID string `json:"user_id"`
	// Reason is the reason the user was warned.
	Reason string `json:"reason"`
}

// Poll represents a Twitch poll on a channel.
type Poll struct {
	// ID is the ID that uniquely identifies the poll.
	ID string `json:"id"`
	// Status is the current status of the poll.
	Status string `json:"status"`
	// BroadcasterID is the ID of the broadcaster that created the poll.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that created the poll.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that created the poll.
	BroadcasterName string `json:"broadcaster_name"`
	// Title is the title of the poll.
	Title string `json:"title"`
	// Choices is a list of choices in the poll.
	Choices []PollChoice `json:"choices"`
	// BitsVotingEnabled indicates whether voting with bits is enabled for the poll.
	BitsVotingEnabled bool `json:"bits_voting_enabled"`
	// ChannelPointsVotingEnabled indicates whether voting with channel points is enabled for the poll.
	ChannelPointsVotingEnabled bool `json:"channel_points_voting_enabled"`
	// BitsPerVote is the number of bits required to cast a vote, if bits voting is enabled.
	BitsPerVote int `json:"bits_per_vote"`
	// ChannelPointsPerVote is the number of channel points required to cast a vote, if channel points voting is enabled.
	ChannelPointsPerVote int `json:"channel_points_per_vote"`
	// Duration is the duration of the poll in seconds.
	Duration int `json:"duration"`
	// StartedAt is the UTC timestamp of when the poll started.
	StartedAt time.Time `json:"started_at"`
	// EndedAt is the UTC timestamp of when the poll ended, if applicable.
	EndedAt *time.Time `json:"ended_at,omitempty"`
}

// OutboundChoice represents a choice when creating a new Twitch poll or prediction.
type OutboundChoice struct {
	// Title is the title of the choice.
	Title string `json:"title"`
}

// PollChoice represents a choice in a Twitch poll.
type PollChoice struct {
	// ID is the ID that uniquely identifies the choice.
	ID string `json:"id"`
	// Title is the title of the choice.
	Title string `json:"title"`
	// Votes is the number of votes the choice has received.
	Votes int `json:"votes"`
	// ChannelPointsVotes is the number of channel points votes the choice has received.
	ChannelPointsVotes int `json:"channel_points_votes"`
	// BitsVotes is the number of bits votes the choice has received.
	BitsVotes int `json:"bits_votes"`
}

// Prediction represents a Twitch prediction on a channel.
type Prediction struct {
	// ID is the ID that uniquely identifies the prediction.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that created the prediction.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that created the prediction.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that created the prediction.
	BroadcasterName string `json:"broadcaster_name"`
	// Title is the title of the prediction.
	Title string `json:"title"`
	// WinningOutcomeID is the ID of the winning outcome, if the prediction has been resolved.
	WinningOutcomeID string `json:"winning_outcome_id,omitempty"`
	// Status is the current status of the prediction.
	Status string `json:"status"`
	// PredictionWindowSeconds is the duration of the prediction in seconds.
	PredictionWindowSeconds int `json:"prediction_window"`
	// Outcomes is a list of possible outcomes for the prediction.
	Outcomes []PredictionOutcome `json:"outcomes"`
	// CreatedAt is the UTC timestamp of when the prediction was created.
	CreatedAt time.Time `json:"created_at"`
	// EndedAt is the UTC timestamp of when the prediction ended, if applicable.
	EndedAt *time.Time `json:"ended_at,omitempty"`
	// LockedAt is the UTC timestamp of when the prediction was locked, if applicable.
	LockedAt *time.Time `json:"locked_at,omitempty"`
}

// OutboundPredictionOutcome represents an outcome when creating a new Twitch prediction.
type OutboundPredictionOutcome struct {
	// Title is the title of the outcome.
	Title string `json:"title"`
}

// PredictionOutcome represents an outcome in a Twitch prediction.
type PredictionOutcome struct {
	// ID is the ID that uniquely identifies the outcome.
	ID string `json:"id"`
	// Title is the title of the outcome.
	Title string `json:"title"`
	// Color is the color associated with the outcome.
	Color string `json:"color"`
	// Users is the number of users who voted for the outcome.
	Users int `json:"users"`
	// ChannelPoints is the total number of channel points wagered on the outcome.
	ChannelPoints int `json:"channel_points"`
	// TopPredictors is a list of the top predictors for the outcome.
	TopPredictors []PredictionParticipant `json:"top_predictors"`
}

// PredictionParticipant represents a participant in a Twitch prediction.
type PredictionParticipant struct {
	// UserID is the ID of the user who participated in the prediction.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user who participated in the prediction.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user who participated in the prediction.
	UserName string `json:"user_name"`
	// ChannelPointsUsed is the number of channel points the user wagered.
	ChannelPointsUsed int `json:"channel_points_used"`
	// ChannelPointsWon is the number of channel points the user won, if applicable.
	ChannelPointsWon int `json:"channel_points_won"`
}

// InitializedRaid represents a raid that has been initialized on a Twitch channel.
type InitializedRaid struct {
	// Mature indicates whether the channel being raided contains mature content.
	Mature bool `json:"is_mature"`
	// CreatedAt is the UTC timestamp of when the raid was created.
	CreatedAt time.Time `json:"created_at"`
}

// StreamSchedule represents a Twitch stream schedule for a channel.
type StreamSchedule struct {
	// BroadcasterID is the ID of the broadcaster that the stream schedule belongs to.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that the stream schedule belongs to.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that the stream schedule belongs to.
	BroadcasterName string `json:"broadcaster_name"`
	// Vacation is the vacation period in the stream schedule, if applicable.
	Vacation *StreamScheduleVacation `json:"vacation"`
	// Segments is a list of segments in the stream schedule.
	Segments []StreamScheduleSegment `json:"segments"`
}

// StreamScheduleVacation represents a vacation period in a Twitch stream schedule.
type StreamScheduleVacation struct {
	// StartsAt is the UTC timestamp of when the vacation starts.
	StartsAt time.Time `json:"start_time"`
	// EndsAt is the UTC timestamp of when the vacation ends.
	EndsAt time.Time `json:"end_time"`
}

// StreamScheduleSegment represents a segment in a Twitch stream schedule.
type StreamScheduleSegment struct {
	// ID is the ID that uniquely identifies the segment.
	ID string `json:"id"`
	// Title is the title of the segment.
	Title string `json:"title"`
	// IsRecurring indicates whether the segment is recurring.
	IsRecurring bool `json:"is_recurring"`
	// Categories is a list of categories for the segment.
	Categories []StreamScheduleCategory `json:"categories"`
	// CanceledUntil is the UTC timestamp of when the segment is canceled until, if applicable.
	CanceledUntil *time.Time `json:"canceled_until"`
	// StartsAt is the UTC timestamp of when the segment starts.
	StartsAt time.Time `json:"start_time"`
	// EndsAt is the UTC timestamp of when the segment ends.
	EndsAt time.Time `json:"end_time"`
}

// StreamScheduleCategory represents a category in a Twitch stream schedule.
type StreamScheduleCategory struct {
	// ID is the ID that uniquely identifies the category.
	ID string `json:"id"`
	// Name is the name of the category.
	Name string `json:"name"`
}

// CategorySearchResult represents a search result for a Twitch category.
type CategorySearchResult struct {
	// ID is the ID that uniquely identifies the category.
	ID string
	// Name is the name of the category.
	Name string
	// BoxArtURL is the URL to the box art of the category.
	BoxArtURL string
}

// ChannelSearchResult represents a search result for a Twitch channel.
type ChannelSearchResult struct {
	// BroadcasterID is the ID of the broadcaster.
	BroadcasterID string `json:"id"`
	// BroadcasterLogin is the login name of the broadcaster.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster.
	BroadcasterName string `json:"display_name"`
	// BroadcasterLanguage is the language of the broadcaster.
	BroadcasterLanguage string `json:"broadcaster_language"`
	// GameID is the ID of the game being played on the channel.
	GameID string `json:"game_id"`
	// GameName is the name of the game being played on the channel.
	GameName string `json:"game_name"`
	// Title is the title of the channel.
	Title string `json:"title"`
	// ThumbnailURL is the URL to the thumbnail of the channel.
	ThumbnailURL string `json:"thumbnail_url"`
	// Tags is a list of tags associated with the channel.
	Tags []string `json:"tag"`
	// Live indicates whether the channel is currently live.
	Live bool `json:"is_live"`
	// Started is the UTC timestamp of when the channel went live, if applicable.
	StartedAt time.Time `json:"started_at"`
}

// StreamKey represents a Twitch stream key.
type StreamKey struct {
	// StreamKey is the stream key to use in order to stream to Twitch.
	//
	// # This value is sensitive and should be treated like a password. Do not share it with anyone.
	Key string `json:"stream_key"`
}

// Stream represents a Twitch stream.
type Stream struct {
	// ID is the ID that uniquely identifies the stream.
	ID string `json:"id"`
	// UserID is the ID of the user who is streaming.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user who is streaming.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user who is streaming.
	UserName string `json:"user_name"`
	// GameID is the ID of the game being played on the stream.
	GameID string `json:"game_id"`
	// GameName is the name of the game being played on the stream.
	GameName string `json:"game_name"`
	// Type is the type of the stream (e.g., "live").
	Type string `json:"type"`
	// Title is the title of the stream.
	Title string `json:"title"`
	// Language is the language of the stream.
	Language string `json:"language"`
	// ThumbnailURL is the URL to the thumbnail of the stream.
	ThumbnailURL string `json:"thumbnail_url"`
	// Tags is a list of tags associated with the stream.
	Tags []string `json:"tags"`
	// ViewerCount is the number of viewers currently watching the stream.
	ViewerCount int `json:"viewer_count"`
	// IsMature indicates whether the stream is marked as mature.
	IsMature bool `json:"is_mature"`
	// StartedAt is the UTC timestamp of when the stream started.
	StartedAt time.Time `json:"started_at"`
}

// StreamMarkerData represents a stream marker on a Twitch video.
type StreamMarkerData struct {
	// ID is the ID that uniquely identifies the stream marker.
	ID string `json:"id"`
	// Description is the description of the stream marker.
	Description string `json:"description"`
	// URL is the URL that opens the video in Twitch Highlighter
	URL string `json:"url,omitempty"`
	// PositionSeconds is the position of the stream marker in seconds from the start of the stream.
	PositionSeconds int `json:"position_seconds"`
	// CreatedAt is the UTC timestamp of when the stream marker was created.
	CreatedAt time.Time `json:"created_at"`
}

// StreamMarker represents a Twitch stream marker.
type StreamMarker struct {
	// UserID is the ID of the user who created the stream marker.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user who created the stream marker.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user who created the stream marker.
	UserName string `json:"user_name"`
	// Videos is a list of videos associated with the stream marker.
	Videos []StreamMarkerVideo `json:"videos"`
}

// StreamMarkerVideo represents a video associated with a Twitch stream marker.
type StreamMarkerVideo struct {
	// VideoID is the ID of the video.
	VideoID string `json:"video_id"`
	// Markers is a list of markers in the video.
	Markers []StreamMarkerData `json:"markers"`
}

// ChannelSubscription represents a subscription to a Twitch channel.
type ChannelSubscription struct {
	// BroadcasterID is the ID of the broadcaster that the subscription belongs to.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that the subscription belongs to.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that the subscription belongs to.
	BroadcasterName string `json:"broadcaster_name"`
	// UserID is the ID of the user that owns the subscription.
	UserID string `json:"user_id"`
	// UserLogin is the login name of the user that owns the subscription.
	UserLogin string `json:"user_login"`
	// UserName is the display name of the user that owns the subscription.
	UserName string `json:"user_name"`
	// GifterID is the ID of the user that gifted the subscription, if applicable.
	GifterID string `json:"gifter_id,omitempty"`
	// GifterLogin is the login name of the user that gifted the subscription, if applicable.
	GifterLogin string `json:"gifter_login,omitempty"`
	// GifterName is the display name of the user that gifted the subscription, if applicable.
	GifterName string `json:"gifter_name,omitempty"`
	// PlanName is the name of the subscription plan.
	PlanName string `json:"plan_name"`
	// Tier is the subscription tier.
	Tier string `json:"tier"`
	// IsGift indicates whether the subscription is a gift.
	IsGift bool `json:"is_gift"`
}

// UserSubscriptionStatus represents the subscription status of a user to a Twitch channel.
type UserSubscriptionStatus struct {
	// BroadcasterID is the ID of the broadcaster that the subscription status belongs to.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that the subscription status belongs to.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that the subscription status belongs to.
	BroadcasterName string `json:"broadcaster_name"`
	// GifterID is the ID of the user that gifted the subscription, if applicable.
	GifterID string `json:"gifter_id,omitempty"`
	// GifterLogin is the login name of the user that gifted the subscription, if applicable.
	GifterLogin string `json:"gifter_login,omitempty"`
	// GifterName is the display name of the user that gifted the subscription, if applicable.
	GifterName string `json:"gifter_name,omitempty"`
	// Tier is the subscription tier.
	Tier string `json:"tier"`
	// IsGift indicates whether the subscription is a gift.
	IsGift bool `json:"is_gift"`
}

// StreamTag represents a tag that can be associated with a Twitch stream.
type StreamTag struct {
	// TagID is the ID that uniquely identifies the tag.
	TagID string `json:"tag_id"`
	// IsAutomatic indicates whether the tag was automatically applied by Twitch.
	IsAutomatic bool `json:"is_auto"`
	// LocalizedNames is a map of locale codes to localized tag names. The key is in the form, <locale>-<country/region>. For example, en-us.
	LocalizedNames map[string]string `json:"localization_names,omitempty"`
	// LocalizedDescriptions is a map of locale codes to localized tag descriptions. The key is in the form, <locale>-<country/region>. For example, en-us.
	LocalizedDescriptions map[string]string `json:"localization_descriptions,omitempty"`
}

// ChannelTeam represents a Twitch channel team.
type ChannelTeam struct {
	// ID is the ID that uniquely identifies the team.
	ID string `json:"id"`
	// BroadcasterID is the ID of the broadcaster that the team belongs to.
	BroadcasterID string `json:"broadcaster_id"`
	// BroadcasterLogin is the login name of the broadcaster that the team belongs to.
	BroadcasterLogin string `json:"broadcaster_login"`
	// BroadcasterName is the display name of the broadcaster that the team belongs to.
	BroadcasterName string `json:"broadcaster_name"`
	// BroadcasterImageURL is the URL to the profile image of the broadcaster that the team belongs to.
	BackgroundImageURL string `json:"background_image_url"`
	// ThumbnailURL is the URL to the thumbnail of the team.
	ThumbnailURL string `json:"thumbnail_url"`
	// TeamName is the name of the team.
	TeamName string `json:"team_name"`
	// TeamDisplayName is the display name of the team.
	TeamDisplayName string `json:"team_display_name"`
	// Banner is the URL to the banner of the team.
	Banner string `json:"banner"`
	// Info is the description of the team.
	Info string `json:"info"`
	// CreatedAt is the UTC timestamp of when the team was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the UTC timestamp of when the team was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// Team represents a Twitch team.
type Team struct {
	// ID is the ID that uniquely identifies the team.
	ID string `json:"id"`
	// TeamName is the name of the team.
	TeamName string `json:"team_name"`
	// ThumbnailURL is the URL to the thumbnail of the team.
	ThumbnailURL string `json:"thumbnail_url"`
	// BackgroundImageURL is the URL to the background image of the team.
	BackgroundImageURL string `json:"background_image_url"`
	// Banner is the URL to the banner of the team.
	Banner string `json:"banner"`
	// Info is the description of the team.
	Info string `json:"info"`
	// Users is a list of users that are members of the team.
	Users []UserInfo `json:"users"`
	// CreatedAt is the UTC timestamp of when the team was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the UTC timestamp of when the team was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents a Twitch user.
type User struct {
	// UserID is the ID of the user.
	UserID string `json:"id"`
	// UserLogin is the login name of the user.
	UserLogin string `json:"login"`
	// UserName is the display name of the user.
	UserName string `json:"display_name"`
	// Type is the type of user. Possible values are "staff", "admin", "global_mod", or "" (empty string) for regular users.
	Type string `json:"type"`
	// BroadcasterType is the broadcaster type of the user. Possible values are "partner", "affiliate", or "" (empty string) for regular users.
	BroadcasterType string `json:"broadcaster_type"`
	// Description is the description of the user.
	Description string `json:"description"`
	// ProfileImageURL is the URL to the profile image of the user.
	ProfileImageURL string `json:"profile_image_url"`
	// OfflineImageURL is the URL to the offline image of the user.
	OfflineImageURL string `json:"offline_image_url"`
	// Email is the email address of the user. This field is only included if app has the user:read:email scope for the user.
	Email string `json:"email,omitempty"`
	// CreatedAt is the UTC timestamp of when the user was created.
	CreatedAt time.Time `json:"created_at"`
}

// UserExtension represents a Twitch user extension.
type UserExtension struct {
	// ID is the ID that uniquely identifies the extension.
	ID string `json:"id"`
	// Name is the name of the extension.
	Name string `json:"name"`
	// Version is the version of the extension.
	Version string `json:"version"`
	// Type is the extension types that can be activated for this extension.
	Types []string `json:"type"`
	// CanActivate indicates whether the extension can be activated by the user.
	CanActivate bool `json:"can_activate"`
}

// UserActiveExtension represents the active extensions on a Twitch channel.
type UserActiveExtension struct {
	// Panel is a map of active panel extensions on the channel.
	Panel map[string]UserActiveExtensionData `json:"panel"`
	// Overlay is a map of active overlay extensions on the channel.
	Overlay map[string]UserActiveExtensionData `json:"overlay"`
	// Component is a map of active component extensions on the channel.
	Component map[string]UserActiveExtensionPositionalData `json:"component"`
}

// UserActiveExtensionData represents data about an active extension on a Twitch channel.
type UserActiveExtensionData struct {
	// ID is the ID that uniquely identifies the extension.
	ID string `json:"id"`
	// Name is the name of the extension.
	Name string `json:"name"`
	// Version is the version of the extension.
	Version string `json:"version"`
	// Active indicates whether the extension is active on the channel.
	Active bool `json:"active"`
}

// UserActiveExtensionPositionalData represents data about an active extension on a Twitch channel.
type UserActiveExtensionPositionalData struct {
	// ID is the ID that uniquely identifies the extension.
	ID string `json:"id"`
	// Name is the name of the extension.
	Name string `json:"name"`
	// Version is the version of the extension.
	Version string `json:"version"`
	// X is the x position of the extension on the channel.
	X int `json:"x"`
	// Y is the y position of the extension on the channel.
	Y int `json:"y"`
	// Active indicates whether the extension is active on the channel.
	Active bool `json:"active"`
}

// Video represents a Twitch video.
type Video struct {
	// ID is the ID that uniquely identifies the video.
	ID string `json:"id"`
	// StreamID is the ID of the stream that the video is associated with.
	StreamID string `json:"stream_id"`
	// BroadcasterID is the ID of the user who created the video.
	BroadcasterID string `json:"user_id"`
	// BroadcasterLogin is the login name of the user who created the video.
	BroadcasterLogin string `json:"user_login"`
	// BroadcasterName is the display name of the user who created the video.
	BroadcasterName string `json:"user_name"`
	// Title is the title of the video.
	Title string `json:"title"`
	// Description is the description of the video.
	Description string `json:"description"`
	// URL is the URL to the video.
	URL string `json:"url"`
	// ThumbnailURL is the URL to the thumbnail of the video.
	ThumbnailURL string `json:"thumbnail_url"`
	// Viewable indicates the viewability of the video.
	Viewable string `json:"viewable"`
	// Language is the language of the video.
	Language string `json:"language"`
	// Type is the type of the video (e.g., "archive", "highlight", "upload").
	Type string `json:"type"`
	// Duration is the duration of the video.
	Duration VideoDuration `json:"duration"`
	// MutedSegments is a list of muted segments in the video.
	MutedSegments []VideoMutedSegment `json:"muted_segments"`
	// ViewCount is the number of views the video has received.
	ViewCount int `json:"view_count"`
	// PublishedAt is the UTC timestamp of when the video was published.
	PublishedAt time.Time `json:"published_at"`
	// CreatedAt is the UTC timestamp of when the video was created.
	CreatedAt time.Time `json:"created_at"`
}

// VideoMutedSegment represents a muted segment in a Twitch video.
type VideoMutedSegment struct {
	// Duration is the duration of the muted segment in seconds.
	Duration int `json:"duration"`
	// Offset is the offset of the muted segment in seconds from the start of the video.
	Offset int `json:"offset"`
}

// VideoDuration represents the duration of a video.
type VideoDuration time.Duration

// Amount returns the amount in major currency units.
func (a CharityCampaignAmount) Amount() float64 {
	return float64(a.Value) / float64(math.Pow10(a.Decimal))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *VideoDuration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = VideoDuration(parsed)
	return nil
}

// AsDuration converts the VideoDuration to a time.Duration.
func (d VideoDuration) AsDuration() time.Duration {
	return time.Duration(d)
}

// WithWebhookTransport creates a new webhook transport for a Twitch Eventsub Conduit Shard.
func WithWebhookTransport(callback, secret string) Transport {
	return Transport{
		Method:   "webhook",
		Callback: &callback,
		Secret:   &secret,
	}
}

// WithWebSocketTransport creates a new websocket transport for a Twitch Eventsub Conduit Shard.
func WithWebSocketTransport(sessionID string) Transport {
	return Transport{
		Method:    "websocket",
		SessionID: &sessionID,
	}
}

// WithChoice creates a new outbound choice with the given title.
func WithChoice(title string) OutboundChoice {
	return OutboundChoice{Title: title}
}

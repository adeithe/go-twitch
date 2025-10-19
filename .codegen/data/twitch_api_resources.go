package data

var (
	// AdsResource is the resource for the Twitch Ads API
	AdsResource = NewTwitchAPIResource("Ads", AdSnoozeResource)
	// AdSnoozeResource is the resource for the Twitch AdSnooze API
	AdSnoozeResource = NewTwitchAPIResource("Snooze")
	// AnalyticsResource is the resource for the Twitch Analytics API
	AnalyticsResource = NewTwitchAPIResource("Analytics", AnalyticsExtensionsResource, AnalyticsGamesResource)
	// AnalyticsExtensionsResource is the resource for the Twitch Extensions Analytics API
	AnalyticsExtensionsResource = NewTwitchAPIResource("Extensions")
	// AnalyticsGamesResource is the resource for the Twitch Games Analytics API
	AnalyticsGamesResource = NewTwitchAPIResource("Games")
	// BitsResource is the resource for the Twitch Bits API
	BitsResource = NewTwitchAPIResource("Bits", BitsCheermotesResource, BitsExtensionsResource, BitsLeaderboardResource)
	// BitsLeaderboardResource is the resource for the Twitch Bits Leaderboard API
	BitsLeaderboardResource = NewTwitchAPIResource("Leaderboard")
	// BitsCheermotesResource is the resource for the Twitch Bits Cheermotes API
	BitsCheermotesResource = NewTwitchAPIResource("Cheermotes")
	// BitsExtensionsResource is the resource for the Twitch Bits Extensions API
	BitsExtensionsResource = NewTwitchAPIResource("Extensions")
	// ChannelsResource is the resource for the Twitch Channels API
	ChannelsResource = NewTwitchAPIResource("Channels", ChannelEditorsResource, ChannelFollowedResource, ChannelFollowersResource)
	// ChannelEditorsResource is the resource for the Twitch Channel Editors API
	ChannelEditorsResource = NewTwitchAPIResource("Editors")
	// ChannelFollowedResource is the resource for the Twitch Channel Followed API
	ChannelFollowedResource = NewTwitchAPIResource("Followed")
	// ChannelFollowersResource is the resource for the Twitch Channel Followers API
	ChannelFollowersResource = NewTwitchAPIResource("Followers")
	// ChannelPointsResource is the resource for the Twitch Channel Points API
	ChannelPointsResource = NewTwitchAPIResource("ChannelPoints", ChannelPointsRedemptionsResource, ChannelPointsRewardsResource)
	// ChannelPointsRewardsResource is the resource for the Twitch Channel Points Rewards API
	ChannelPointsRewardsResource = NewTwitchAPIResource("Rewards")
	// ChannelPointsRedemptionsResource is the resource for the Twitch Channel Points Redemptions API
	ChannelPointsRedemptionsResource = NewTwitchAPIResource("Redemptions")
	// CharityResource is the resource for the Twitch Charity API
	CharityResource = NewTwitchAPIResource("Charity", CharityCampaignResource, CharityDonationsResource)
	// CharityCampaignResource is the resource for the Twitch Charity Campaign API
	CharityCampaignResource = NewTwitchAPIResource("Campaign")
	// CharityDonationsResource is the resource for the Twitch Charity Donations API
	CharityDonationsResource = NewTwitchAPIResource("Donations")
	// ChatResource is the resource for the Twitch Chat API
	ChatResource = NewTwitchAPIResource("Chat", ChatBadgesResource, ChatChattersResource, ChatEmotesResource, ChatEmoteSetsResource, ChatSettingsResource, ChatSharedResource, ChatAnnouncementResource, ChatShoutoutResource)
	// ChatChattersResource is the resource for the Twitch Chat Chatters API
	ChatChattersResource = NewTwitchAPIResource("Chatters", ChatChattersUserResource)
	// ChatChattersUserResource is the resource for the Twitch Chat Chatters User API
	ChatChattersUserResource = NewTwitchAPIResource("User", ChatChattersUserColorResource)
	// ChatChattersUserColorResource is the resource for the Twitch Chat Chatters User Color API
	ChatChattersUserColorResource = NewTwitchAPIResource("Color")
	// ChatEmotesResource is the resource for the Twitch ChatEmotes API
	ChatEmotesResource = NewTwitchAPIResource("Emotes", ChatChannelEmotesResource, ChatEmotesGlobalResource, ChatUserEmotesResource)
	// ChatChannelEmotesResource is the resource for the Twitch Chat Channel Emotes API
	ChatChannelEmotesResource = NewTwitchAPIResource("Channel")
	// ChatUserEmotesResource is the resource for the Twitch Chat User Emotes API
	ChatUserEmotesResource = NewTwitchAPIResource("User")
	// ChatEmotesGlobalResource is the resource for the Twitch Chat Emotes Global API
	ChatEmotesGlobalResource = NewTwitchAPIResource("Global")
	// ChatEmoteSetsResource is the resource for the Twitch Chat Emote Sets API
	ChatEmoteSetsResource = NewTwitchAPIResource("EmoteSets")
	// ChatBadgesResource is the resource for the Twitch Chat Badges API
	ChatBadgesResource = NewTwitchAPIResource("Badges", ChatBadgesGlobalResource)
	// ChatBadgesGlobalResource is the resource for the Twitch Chat Badges Global API
	ChatBadgesGlobalResource = NewTwitchAPIResource("Global")
	// ChatSettingsResource is the resource for the Twitch Chat Settings API
	ChatSettingsResource = NewTwitchAPIResource("Settings")
	// ChatSharedResource is the resource for the Twitch Chat Shared API
	ChatSharedResource = NewTwitchAPIResource("Shared")
	// ChatAnnouncementResource is the resource for the Twitch Chat Announcement API
	ChatAnnouncementResource = NewTwitchAPIResource("Announcement")
	// ChatShoutoutResource is the resource for the Twitch Chat Shoutout API
	ChatShoutoutResource = NewTwitchAPIResource("Shoutout")
	// ClipsResource is the resource for the Twitch Clips API
	ClipsResource = NewTwitchAPIResource("Clips", ClipsDownloadResource)
	// ClipsDownloadResource is the resource for the Twitch Clips Download API
	ClipsDownloadResource = NewTwitchAPIResource("Download")
	// ConduitsResource is the resource for the Twitch Conduits API
	ConduitsResource = NewTwitchAPIResource("Conduits", ConduitsShardsResource)
	// ConduitsShardsResource is the resource for the Twitch Conduits Shards API
	ConduitsShardsResource = NewTwitchAPIResource("Shards")
	// CCLResource is the resource for the Twitch Content Classification Labels API
	CCLResource = NewTwitchAPIResource("ContentLabels")
	// EntitlementsResource is the resource for the Twitch Entitlements API
	EntitlementsResource = NewTwitchAPIResource("Entitlements", EntitlementsDropsResource)
	// EntitlementsDropsResource is the resource for the Twitch Entitlements Drops API
	EntitlementsDropsResource = NewTwitchAPIResource("Drops")
	// ExtensionsResource is the resource for the Twitch Extensions API
	ExtensionsResource = NewTwitchAPIResource("Extensions")
	// EventSubResource is the resource for the Twitch EventSub API
	EventSubResource = NewTwitchAPIResource("EventSub")
	// GamesResource is the resource for the Twitch Games API
	GamesResource = NewTwitchAPIResource("Games", GamesTopResource)
	// GamesTopResource is the resource for the Twitch Games Top API
	GamesTopResource = NewTwitchAPIResource("Top")
	// GoalsResource is the resource for the Twitch Goals API
	GoalsResource = NewTwitchAPIResource("Goals")
	// GuestStarResource is the resource for the Twitch Guest Star API
	GuestStarResource = NewTwitchAPIResource("GuestStar", GuestStarSessionResource)
	// GuestStarSessionResource is the resource for the Twitch Guest Star Session API
	GuestStarSessionResource = NewTwitchAPIResource("Session")
	// HypeTrainResource is the resource for the Twitch Hype Train API
	HypeTrainResource = NewTwitchAPIResource("HypeTrain")
	// ModerationResource is the resource for the Twitch Moderation API
	ModerationResource = NewTwitchAPIResource("Moderation", ModerationBansResource, ModerationClearChatResource)
	// ModerationBansResource is the resource for the Twitch Moderation Bans API
	ModerationBansResource = NewTwitchAPIResource("Bans")
	// ModerationClearChatResource is the resource for the Twitch Moderation Clear Chat API
	ModerationClearChatResource = NewTwitchAPIResource("ClearChat")
	// PollsResource is the resource for the Twitch Polls API
	PollsResource = NewTwitchAPIResource("Polls")
	// PredictionsResource is the resource for the Twitch Predictions API
	PredictionsResource = NewTwitchAPIResource("Predictions")
	// RaidsResource is the resource for the Twitch Raids API
	RaidsResource = NewTwitchAPIResource("Raids")
	// ScheduleResource is the resource for the Twitch Schedule API
	ScheduleResource = NewTwitchAPIResource("Schedule")
	// SearchResource is the resource for the Twitch Search API
	SearchResource = NewTwitchAPIResource("Search", SearchCategoriesResource, SearchChannelsResource)
	// SearchCategoriesResource is the resource for the Twitch Search Categories API
	SearchCategoriesResource = NewTwitchAPIResource("Categories")
	// SearchChannelsResource is the resource for the Twitch Search Channels API
	SearchChannelsResource = NewTwitchAPIResource("Channels")
	// StreamsResource is the resource for the Twitch Streams API
	StreamsResource = NewTwitchAPIResource("Streams", StreamKeyResource, StreamsFollowedResource, StreamsMarkersResource)
	// StreamKeyResource is the resource for the Twitch Stream Key API
	StreamKeyResource = NewTwitchAPIResource("StreamKey")
	// StreamsFollowedResource is the resource for the Twitch Streams Followed API
	StreamsFollowedResource = NewTwitchAPIResource("Followed")
	// StreamsMarkersResource is the resource for the Twitch Streams Markers API
	StreamsMarkersResource = NewTwitchAPIResource("Markers")
	// SubscriptionsResource is the resource for the Twitch Subscriptions API
	SubscriptionsResource = NewTwitchAPIResource("Subscriptions", SubscriptionsUserResource)
	// SubscriptionsUserResource is the resource for the Twitch Subscriptions User API
	SubscriptionsUserResource = NewTwitchAPIResource("Subscribed")
	// TeamsResource is the resource for the Twitch Teams API
	TeamsResource = NewTwitchAPIResource("Teams", TeamsChannelsResource)
	// TeamsChannelsResource is the resource for the Twitch Teams Channels API
	TeamsChannelsResource = NewTwitchAPIResource("Channels")
	// UsersResource is the resource for the Twitch Users API
	UsersResource = NewTwitchAPIResource("Users")
	// VideosResource is the resource for the Twitch Videos API
	VideosResource = NewTwitchAPIResource("Videos")
	// WhispersResource is the resource for the Twitch Whispers API
	WhispersResource = NewTwitchAPIResource("Whispers")

	// Resources is the list of top-level API resources to generate. Subresources are included automatically.
	Resources = []*TwitchAPIResource{
		AdsResource, AnalyticsResource, BitsResource, ChannelsResource, ChannelPointsResource,
		CharityResource, ChatResource, ClipsResource, ConduitsResource, CCLResource, EntitlementsResource,
		ExtensionsResource, EventSubResource, GamesResource, GoalsResource, GuestStarResource, HypeTrainResource,
		ModerationResource, PollsResource, PredictionsResource, RaidsResource, ScheduleResource, SearchResource,
		StreamsResource, SubscriptionsResource, TeamsResource, UsersResource, VideosResource, WhispersResource,
	}
)

// ResourceMapping maps each TwitchAPIResource to its associated endpoints.
var ResourceMapping map[*TwitchAPIResource][]*TwitchAPIEndpoint

func init() {
	ResourceMapping = make(map[*TwitchAPIResource][]*TwitchAPIEndpoint)
	for _, e := range Endpoints {
		ResourceMapping[e.Resource] = append(ResourceMapping[e.Resource], e)
	}
}

// NewTwitchAPIResource creates a new TwitchAPIResource with the given name and sub-resources.
func NewTwitchAPIResource(name string, subresources ...*TwitchAPIResource) *TwitchAPIResource {
	return &TwitchAPIResource{Name: name, SubResources: subresources}
}

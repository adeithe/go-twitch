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
	// GamesResource is the resource for the Twitch Games API
	GamesResource = NewTwitchAPIResource("Games", GamesTopResource)
	// GamesTopResource is the resource for the Twitch Games Top API
	GamesTopResource = NewTwitchAPIResource("Top")
	// ModerationResource is the resource for the Twitch Moderation API
	ModerationResource = NewTwitchAPIResource("Moderation", ModerationBansResource, ModerationClearChatResource)
	// ModerationBansResource is the resource for the Twitch Moderation Bans API
	ModerationBansResource = NewTwitchAPIResource("Bans")
	// ModerationClearChatResource is the resource for the Twitch Moderation Clear Chat API
	ModerationClearChatResource = NewTwitchAPIResource("ClearChat")
	// StreamsResource is the resource for the Twitch Streams API
	StreamsResource = NewTwitchAPIResource("Streams")
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
		CharityResource, ChatResource, ClipsResource, ConduitsResource, CCLResource, GamesResource,
		ModerationResource, StreamsResource, TeamsResource, UsersResource, VideosResource, WhispersResource,
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

package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/adeithe/go-twitch/api"
)

type TwitchAPIResource struct {
	Name         string
	SubResources []*TwitchAPIResource
}

type TwitchAPIEndpoint struct {
	Resource     *TwitchAPIResource
	Name         string
	Method, Path string
	DocsURL      string
	Comments     []string
	Params       any
	Response     any
}

type BasicResponse[T any] struct {
	Data T
}

const (
	TemplateAPIHeader   = "template_twitch_api_header.go.tmpl"
	TemplateAPIResource = "template_twitch_api_resource.go.tmpl"
	TemplateAPICall     = "template_twitch_api_call.go.tmpl"
)

var (
	AdsResource                      = NewTwitchAPIResource("Ads", AdSnoozeResource)
	AdSnoozeResource                 = NewTwitchAPIResource("Snooze")
	AnalyticsResource                = NewTwitchAPIResource("Analytics", AnalyticsExtensionsResource, AnalyticsGamesResource)
	AnalyticsExtensionsResource      = NewTwitchAPIResource("Extensions")
	AnalyticsGamesResource           = NewTwitchAPIResource("Games")
	BitsResource                     = NewTwitchAPIResource("Bits", BitsCheermotesResource, BitsExtensionsResource, BitsLeaderboardResource)
	BitsLeaderboardResource          = NewTwitchAPIResource("Leaderboard")
	BitsCheermotesResource           = NewTwitchAPIResource("Cheermotes")
	BitsExtensionsResource           = NewTwitchAPIResource("Extensions")
	ChannelsResource                 = NewTwitchAPIResource("Channels", ChannelEditorsResource, ChannelFollowedResource, ChannelFollowersResource)
	ChannelEditorsResource           = NewTwitchAPIResource("Editors")
	ChannelFollowedResource          = NewTwitchAPIResource("Followed")
	ChannelFollowersResource         = NewTwitchAPIResource("Followers")
	ChannelPointsResource            = NewTwitchAPIResource("ChannelPoints", ChannelPointsRedemptionsResource, ChannelPointsRewardsResource)
	ChannelPointsRewardsResource     = NewTwitchAPIResource("Rewards")
	ChannelPointsRedemptionsResource = NewTwitchAPIResource("Redemptions")
	CharityResource                  = NewTwitchAPIResource("Charity", CharityCampaignResource, CharityDonationsResource)
	CharityCampaignResource          = NewTwitchAPIResource("Campaign")
	CharityDonationsResource         = NewTwitchAPIResource("Donations")
	ChatResource                     = NewTwitchAPIResource("Chat", ChatBadgesResource, ChatChattersResource, ChatEmotesResource, ChatEmoteSetsResource, ChatSettingsResource, ChatSharedResource, ChatAnnouncementResource, ChatShoutoutResource)
	ChatChattersResource             = NewTwitchAPIResource("Chatters", ChatChattersUserResource)
	ChatChattersUserResource         = NewTwitchAPIResource("User", ChatChattersUserColorResource)
	ChatChattersUserColorResource    = NewTwitchAPIResource("Color")
	ChatEmotesResource               = NewTwitchAPIResource("Emotes", ChatChannelEmotesResource, ChatEmotesGlobalResource, ChatUserEmotesResource)
	ChatChannelEmotesResource        = NewTwitchAPIResource("Channel")
	ChatUserEmotesResource           = NewTwitchAPIResource("User")
	ChatEmotesGlobalResource         = NewTwitchAPIResource("Global")
	ChatEmoteSetsResource            = NewTwitchAPIResource("EmoteSets")
	ChatBadgesResource               = NewTwitchAPIResource("Badges", ChatBadgesGlobalResource)
	ChatBadgesGlobalResource         = NewTwitchAPIResource("Global")
	ChatSettingsResource             = NewTwitchAPIResource("Settings")
	ChatSharedResource               = NewTwitchAPIResource("Shared")
	ChatAnnouncementResource         = NewTwitchAPIResource("Announcement")
	ChatShoutoutResource             = NewTwitchAPIResource("Shoutout")
	ClipsResource                    = NewTwitchAPIResource("Clips", ClipsDownloadResource)
	ClipsDownloadResource            = NewTwitchAPIResource("Download")
	ConduitsResource                 = NewTwitchAPIResource("Conduits", ConduitsShardsResource)
	ConduitsShardsResource           = NewTwitchAPIResource("Shards")
	GamesResource                    = NewTwitchAPIResource("Games", GamesTopResource)
	GamesTopResource                 = NewTwitchAPIResource("Top")
	ModerationResource               = NewTwitchAPIResource("Moderation", ModerationBansResource, ModerationClearChatResource)
	ModerationBansResource           = NewTwitchAPIResource("Bans")
	ModerationClearChatResource      = NewTwitchAPIResource("ClearChat")
	StreamsResource                  = NewTwitchAPIResource("Streams")
	UsersResource                    = NewTwitchAPIResource("Users")
	VideosResource                   = NewTwitchAPIResource("Videos")
	WhispersResource                 = NewTwitchAPIResource("Whispers")

	// Resources is the list of top-level API resources to generate. Subresources are included automatically.
	Resources = []*TwitchAPIResource{
		AdsResource, AnalyticsResource, BitsResource, ChannelsResource, ChannelPointsResource,
		CharityResource, ChatResource, ClipsResource, ConduitsResource, GamesResource,
		ModerationResource, StreamsResource, UsersResource, VideosResource, WhispersResource,
	}

	// Endpoints is the list of all API endpoints to generate. Resource mapping is done automatically using the Resource field.
	Endpoints = []*TwitchAPIEndpoint{
		// Ads
		{
			// https://dev.twitch.tv/docs/api/reference/#start-commercial
			Resource: AdsResource,
			Name:     "StartCommercial",
			Method:   http.MethodPost,
			Path:     api.EndpointAdsStartCommercial,
			DocsURL:  "#start-commercial",
			Comments: []string{
				"Starts a commercial on the specified channel.", "",
				"Only Twitch Partners and Affiliates can run commercials on their channels and must be live.", "",
				"# Authorization", "", "Requires the broadcasters access token that includes the channel:edit:commercial scope.",
			},
			Params: struct {
				BroadcasterID string `body:"-"`
				Length        int    `body:"-"`
			}{},
			Response: BasicResponse[api.Commercial]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-ad-schedule
			Resource: AdsResource,
			Name:     "AdSchedule",
			Method:   http.MethodGet,
			Path:     api.EndpointAdsGetAdsSchedule,
			DocsURL:  "#get-ad-schedule",
			Comments: []string{
				"Gets ad schedule related information, including snooze, when the last ad was run, when the next ad is scheduled, and if the channel is currently in pre-roll free time.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:read:ads scope. The user_id in the user access token must match the broadcaster_id.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
			}{},
			Response: BasicResponse[api.AdSchedule]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#snooze-next-ad
			Resource: AdSnoozeResource,
			Name:     "SnoozeNextAd",
			Method:   http.MethodPost,
			Path:     api.EndpointAdsSnoozeNextAd,
			DocsURL:  "#snooze-next-ad",
			Comments: []string{
				"If available, pushes back the timestamp of the upcoming automatic mid-roll ad by 5 minutes.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:ads scope. The user_id in the user access token must match the broadcaster_id.",
			},
			Params: struct {
				BroadcasterID string `query:"-"`
			}{},
			Response: BasicResponse[api.AdsSnoozed]{},
		},
		// Analytics
		{
			// https://dev.twitch.tv/docs/api/reference/#get-extension-analytics
			Resource: AnalyticsExtensionsResource,
			Name:     "ExtensionAnalytics",
			Method:   http.MethodGet,
			Path:     api.EndpointAnalyticsGetExtensionAnalytics,
			DocsURL:  "#get-extension-analytics",
			Comments: []string{
				"Gets an analytics report for one or more extensions. The response contains the URLs used to download the reports (CSV files).", "",
				"# Authorization", "", "Requires a user access token that includes the analytics:read:extensions scope.",
			},
			Params: struct {
				ExtensionID string    `query:"-"`
				Type        string    `query:"-"`
				After       string    `query:"-"`
				First       int       `query:"-"`
				StartedAt   time.Time `query:"-"`
				EndedAt     time.Time `query:"-"`
			}{},
			Response: struct {
				Data       api.ExtensionAnalyticsReport
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-game-analytics
			Resource: AnalyticsGamesResource,
			Name:     "GameAnalytics",
			Method:   http.MethodGet,
			Path:     api.EndpointAnalyticsGetGameAnalytics,
			DocsURL:  "#get-game-analytics",
			Comments: []string{
				"Gets an analytics report for one or more games. The response contains the URLs used to download the reports (CSV files).", "",
				"# Authorization", "", "Requires a user access token that includes the analytics:read:games scope.",
			},
			Params: struct {
				GameID    string    `query:"-"`
				Type      string    `query:"-"`
				StartedAt time.Time `query:"-"`
				EndedAt   time.Time `query:"-"`
				First     int       `query:"-"`
				After     string    `query:"-"`
			}{},
			Response: struct {
				Data       api.GameAnalyticsReport
				Pagination api.Pagination
			}{},
		},
		// Bits
		{
			// https://dev.twitch.tv/docs/api/reference/#get-bits-leaderboard
			Resource: BitsLeaderboardResource,
			Name:     "BitsLeaderboard",
			Method:   http.MethodGet,
			Path:     api.EndpointBitsGetLeaderboard,
			DocsURL:  "#get-bits-leaderboard",
			Comments: []string{
				"Gets the Bits leaderboard for the authenticated broadcaster.", "",
				"# Authorization", "", "Requires a user access token that includes the bits:read scope.",
			},
			Params: struct {
				UserID    string    `query:"-"`
				Period    string    `query:"-"`
				Count     int       `query:"-"`
				StartedAt time.Time `query:"-"`
			}{},
			Response: struct {
				Total     int
				Data      api.BitsLeaderboardEntry
				DateRange api.DateRange
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-cheermotes
			Resource: BitsCheermotesResource,
			Name:     "BitsCheermotes",
			Method:   http.MethodGet,
			Path:     api.EndpointBitsGetCheermotes,
			DocsURL:  "#get-cheermotes",
			Comments: []string{
				"Gets a list of Cheermotes that users can use to cheer Bits in any Bits-enabled chat room.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				BroadcasterID string `query:"-"`
			}{},
			Response: BasicResponse[api.Cheermote]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-extension-transactions
			Resource: BitsExtensionsResource,
			Name:     "BitsExtensionTransactions",
			Method:   http.MethodGet,
			Path:     api.EndpointBitsGetExtensionTransactions,
			DocsURL:  "#get-extension-transactions",
			Comments: []string{
				"Gets a list of transactions for an extension. A transaction records the exchange of a currency (for example, Bits) for a digital product.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Params: struct {
				ID          string `query:"-"`
				ExtensionID string `query:"-"`
				After       string `query:"-"`
				First       int    `query:"-"`
			}{},
			Response: BasicResponse[api.ExtensionTransaction]{},
		},
		// Channels
		{
			// https://dev.twitch.tv/docs/api/reference/#get-channel-information
			Resource: ChannelsResource,
			Name:     "ChannelInformation",
			Method:   http.MethodGet,
			Path:     api.EndpointChannelsGetInformation,
			DocsURL:  "#get-channel-information",
			Comments: []string{
				"Gets information about one or more channels.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				BroadcasterID []string `query:"-"`
			}{},
			Response: BasicResponse[api.Channel]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#modify-channel-information
			Resource: ChannelsResource,
			Name:     "ChannelInformation",
			Method:   http.MethodPatch,
			Path:     api.EndpointChannelsModifyInformation,
			DocsURL:  "#modify-channel-information",
			Comments: []string{
				"Updates properties for a channel.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:broadcast scope.",
			},
			Params: struct {
				BroadcasterID               string                                  `query:"-"`
				GameID                      string                                  `body:"-"`
				BroadcasterLanguage         string                                  `body:"-"`
				Title                       string                                  `body:"-"`
				Tags                        []string                                `body:"-"`
				ContentClassificationLabels []api.ChannelContentClassificationLabel `body:"-"`
				Delay                       int                                     `body:"-"`
				IsBrandedContent            bool                                    `body:"-"`
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-channel-editors
			Resource: ChannelEditorsResource,
			Name:     "ChannelEditors",
			Method:   http.MethodGet,
			Path:     api.EndpointChannelsGetEditors,
			DocsURL:  "#get-channel-editors",
			Comments: []string{
				"Gets the list of editors for a broadcaster.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:read:editors scope.",
			},
			Params: struct {
				BroadcasterID string `query:"-"`
			}{},
			Response: BasicResponse[api.ChannelEditor]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-followed-channels
			Resource: ChannelFollowedResource,
			Name:     "ChannelsFollowed",
			Method:   http.MethodGet,
			Path:     api.EndpointChannelsGetFollowedChannels,
			DocsURL:  "#get-followed-channels",
			Comments: []string{
				"Gets a list of broadcasters that the specified user follows.", "You can also use this endpoint to see whether a user follows a specific broadcaster.", "",
				"# Authorization", "", "Requires a user access token that includes the user:read:follows scope.",
			},
			Params: struct {
				UserID        string `query:"-"`
				BroadcasterID string `query:"-"`
				After         string `query:"-"`
				First         int    `query:"-"`
			}{},
			Response: struct {
				Total      int
				Data       api.Followed
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-channel-followers
			Resource: ChannelFollowersResource,
			Name:     "ChannelFollowers",
			Method:   http.MethodGet,
			Path:     api.EndpointChannelsGetFollowers,
			DocsURL:  "#get-channel-followers",
			Comments: []string{
				"Gets a list of users that follow the specified broadcaster.", "You can also use this endpoint to see whether a specific user follows the broadcaster.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:read:followers scope.", "",
				"The ID in the broadcaster_id query parameter must match the user ID in the access token or the user ID in the access token must be a moderator for the specified broadcaster.",
			},
			Params: struct {
				UserID        string `query:"-"`
				BroadcasterID string `query:"-"`
				After         string `query:"-"`
				First         int    `query:"-"`
			}{},
			Response: struct {
				Total      int
				Data       api.Follower
				Pagination api.Pagination
			}{},
		},
		// Channel Points
		{
			// https://dev.twitch.tv/docs/api/reference/#create-custom-rewards
			Resource: ChannelPointsRewardsResource,
			Name:     "ChannelPointRewards",
			Method:   http.MethodPost,
			Path:     api.EndpointChannelPointsGetCustomRewards,
			DocsURL:  "#create-custom-rewards",
			Comments: []string{
				"Creates a new Custom Reward in a channel.", "",
				"The maximum number of custom rewards per channel is 50, which includes both enabled and disabled rewards.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:redemptions scope.",
			},
			Params: struct {
				BroadcasterID                     string `query:"-"`
				Title                             string `body:"-"`
				Prompt                            string `body:"-"`
				BackgroundColor                   string `body:"-"`
				Cost                              int64  `body:"-"`
				MaxPerStream                      int    `body:"-"`
				MaxPerUserPerStream               int    `body:"-"`
				GlobalCooldownSeconds             int    `body:"-"`
				IsEnabled                         bool   `body:"-"`
				IsUserInputRequired               bool   `body:"-"`
				IsMaxPerStreamEnabled             bool   `body:"-"`
				IsMaxPerUserPerStreamEnabled      bool   `body:"-"`
				IsGlobalCooldownEnabled           bool   `body:"-"`
				ShouldRedemptionsSkipRequestQueue bool   `body:"-"`
			}{},
			Response: BasicResponse[api.CustomReward]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#delete-custom-reward
			Resource: ChannelPointsRewardsResource,
			Name:     "ChannelPointRewards",
			Method:   http.MethodDelete,
			Path:     api.EndpointChannelPointsDeleteCustomReward,
			DocsURL:  "#delete-custom-reward",
			Comments: []string{
				"Deletes a custom reward that the broadcaster created.", "",
				"The app used to create the reward is the only app that may delete it.",
				"If the redemption status for the reward is UNFULFILLED at the time the reward is deleted, its redemption status is marked as FULFILLED.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:redemptions scope.",
			},
			Params: struct {
				ID            string `query:"required"`
				BroadcasterID string `query:"required"`
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-custom-reward
			Resource: ChannelPointsRewardsResource,
			Name:     "ChannelPointRewards",
			Method:   http.MethodGet,
			Path:     api.EndpointChannelPointsGetCustomRewards,
			DocsURL:  "#get-custom-reward",
			Comments: []string{
				"Gets a list of custom rewards that the specified broadcaster created.", "",
				"A channel may offer a maximum of 50 rewards, which includes both enabled and disabled rewards.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:read:redemptions or channel:manage:redemptions scope.",
			},
			Params: struct {
				ID                    []string `query:"-"`
				BroadcasterID         string   `query:"required"`
				OnlyManageableRewards bool     `query:"-"`
			}{},
			Response: BasicResponse[api.CustomReward]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-custom-reward-redemption
			Resource: ChannelPointsRedemptionsResource,
			Name:     "ChannelPointRedemptions",
			Method:   http.MethodGet,
			Path:     api.EndpointChannelPointsGetCustomRewardRedemptions,
			DocsURL:  "#get-custom-reward-redemption",
			Comments: []string{
				"Gets a list of redemptions for the specified custom reward.", "The app used to create the reward is the only app that may get the redemptions.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:read:redemptions or channel:manage:redemptions scope.",
			},
			Params: struct {
				ID            string `query:"-"`
				BroadcasterID string `query:"required"`
				RewardID      string `query:"required"`
				Status        string `query:"required"`
				Sort          string `query:"-"`
				After         string `query:"-"`
				First         int    `query:"-"`
			}{},
			Response: BasicResponse[api.CustomRewardRedemption]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#update-custom-reward
			Resource: ChannelPointsRewardsResource,
			Name:     "ChannelPointRewards",
			Method:   http.MethodPatch,
			Path:     api.EndpointChannelPointsUpdateCustomReward,
			DocsURL:  "#update-custom-reward",
			Comments: []string{
				"Updates a custom reward. The app used to create the reward is the only app that may update the reward.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:redemptions scope.",
			},
			Params: struct {
				ID                                string `query:"required"`
				BroadcasterID                     string `query:"required"`
				Title                             string `body:"-"`
				Prompt                            string `body:"-"`
				BackgroundColor                   string `body:"-"`
				Cost                              int    `body:"-"`
				MaxPerStream                      int    `body:"-"`
				MaxPerUserPerStream               int    `body:"-"`
				GlobalCooldownSeconds             int    `body:"-"`
				IsPaused                          bool   `body:"-"`
				IsEnabled                         bool   `body:"-"`
				IsUserInputRequired               bool   `body:"-"`
				IsMaxPerStreamEnabled             bool   `body:"-"`
				IsMaxPerUserPerStreamEnabled      bool   `body:"-"`
				IsGlobalCooldownEnabled           bool   `body:"-"`
				ShouldRedemptionsSkipRequestQueue bool   `body:"-"`
			}{},
			Response: BasicResponse[api.CustomReward]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#update-redemption-status
			Resource: ChannelPointsRedemptionsResource,
			Name:     "ChannelPointRedemptions",
			Method:   http.MethodPatch,
			Path:     api.EndpointChannelPointsGetCustomRewardRedemptions,
			DocsURL:  "#update-redemption-status",
			Comments: []string{
				"Updates the status for a redemption. You may update a redemption only if its status is UNFULFILLED.",
				"The app used to create the reward is the only app that may update the redemption.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:redemptions scope.",
			},
			Params: struct {
				ID            string `query:"required"`
				BroadcasterID string `query:"required"`
				RewardID      string `query:"required"`
			}{},
			Response: BasicResponse[api.CustomRewardRedemption]{},
		},
		// Charity
		{
			// https://dev.twitch.tv/docs/api/reference/#get-charity-campaign
			Resource: CharityCampaignResource,
			Name:     "CharityCampaign",
			Method:   http.MethodGet,
			Path:     api.EndpointCharityGetCampaign,
			DocsURL:  "#get-charity-campaign",
			Comments: []string{
				"Gets information about the charity campaign that a broadcaster is running.", "For example, the  fundraising goal for a campaign and the current amount of donations.", "",
				"To receive events when progress is made towards the campaign goal or the broadcaster changes the fundraising goal, subscribe to the channel.charity_campaign.progress subscription type.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:read:charity scope.",
			},
			Params: struct {
				BroadcasterID string `query:"-"`
			}{},
			Response: BasicResponse[api.CharityCampaign]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-charity-campaign-donations
			Resource: CharityDonationsResource,
			Name:     "CharityDonations",
			Method:   http.MethodGet,
			Path:     api.EndpointCharityGetCampaign,
			DocsURL:  "#get-charity-campaign-donations",
			Comments: []string{
				"Gets the list of donations that users have made to the broadcasters active charity campaign.", "",
				"To receive events as donations occur, subscribe to the channel.charity_campaign.donate subscription type.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:read:charity scope.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
				After         string `query:"-"`
				First         int    `query:"-"`
			}{},
			Response: struct {
				Data       api.CharityCampaignDonation
				Pagination api.Pagination
			}{},
		},
		// Chat
		{
			// https://dev.twitch.tv/docs/api/reference/#get-chatters
			Resource: ChatChattersResource,
			Name:     "ChatChatters",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetChatters,
			DocsURL:  "#get-chatters",
			Comments: []string{
				"Gets the list of users that are connected to a chat session.", "",
				"To determine whether a user is a moderator or VIP, use the Get Moderators and Get VIPs endpoints.",
				"You can check the roles of up to 100 users.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:read:chatters scope.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
				ModeratorID   string `query:"required"`
				After         string `query:"-"`
				First         int    `query:"-"`
			}{},
			Response: struct {
				Total      int
				Data       api.UserInfo
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-channel-emotes
			Resource: ChatChannelEmotesResource,
			Name:     "ChannelEmotes",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetChannelEmotes,
			DocsURL:  "#get-channel-emotes",
			Comments: []string{
				"Gets the list of custom emotes for a broadcaster.",
				"Broadcasters create these custom emotes for users who subscribe to or follow the channel or cheer Bits in the chat window.", "",
				"With the exception of custom follower emotes, users may use custom emotes in any Twitch chat.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
			}{},
			Response: BasicResponse[api.Emote]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-global-emotes
			Resource: ChatEmotesGlobalResource,
			Name:     "GlobalEmotes",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetGlobalEmotes,
			DocsURL:  "#get-global-emotes",
			Comments: []string{
				"Gets the list of global emotes.", "Global emotes are Twitch-created emotes that users can use in any Twitch chat.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Response: BasicResponse[api.Emote]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-emote-sets
			Resource: ChatEmoteSetsResource,
			Name:     "EmoteSets",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetEmoteSets,
			DocsURL:  "#get-emote-sets",
			Comments: []string{
				"Gets emotes for one or more specified emote sets.", "",
				"An emote set groups emotes that have a similar context.",
				"For example, Twitch places all the subscriber emotes that a broadcaster uploads for their channel in the same emote set.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				EmoteSetID []string `query:"required"`
			}{},
			Response: BasicResponse[api.Emote]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-channel-chat-badges
			Resource: ChatBadgesResource,
			Name:     "ChannelChatBadges",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetChannelBadges,
			DocsURL:  "#get-channel-chat-badges",
			Comments: []string{
				"Gets the list of custom chat badges for a broadcaster.", "",
				"The list is empty if the broadcaster hasnt created custom chat badges.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
			}{},
			Response: BasicResponse[api.ChatBadge]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-global-chat-badges
			Resource: ChatBadgesGlobalResource,
			Name:     "GlobalChatBadges",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetGlobalBadges,
			DocsURL:  "#get-global-chat-badges",
			Comments: []string{"Gets list of chat badges on Twitch, which users may use in any chat room."},
			Response: BasicResponse[api.ChatBadge]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-chat-settings
			Resource: ChatSettingsResource,
			Name:     "ChatSettings",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetSettings,
			DocsURL:  "#get-chat-settings",
			Comments: []string{
				"Gets the chat settings for a channel.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
				ModeratorID   string `query:"required"`
			}{},
			Response: BasicResponse[api.ChatSettings]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-shared-chat-session
			Resource: ChatSharedResource,
			Name:     "SharedChatSession",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetSharedChatSession,
			DocsURL:  "#get-shared-chat-session",
			Comments: []string{
				"Retrieves the active shared chat session for a channel.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:read:chat_settings scope.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
			}{},
			Response: BasicResponse[api.SharedChatSession]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-user-emotes
			Resource: ChatUserEmotesResource,
			Name:     "UserEmotes",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetUserEmotes,
			DocsURL:  "#get-user-emotes",
			Comments: []string{
				"Retrieves emotes available to the user across all channels.", "",
				"# Authorization", "",
				"Requires a user access token that includes the user:read:emotes scope.", "",
				"Query parameter user_id must match the user_id in the user access token.",
			},
			Params: struct {
				UserID        string `query:"required"`
				BroadcasterID string `query:"-"`
				After         string `query:"-"`
			}{},
			Response: struct {
				Data       api.Emote
				Template   string
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#update-chat-settings
			Resource: ChatSettingsResource,
			Name:     "ChatSettings",
			Method:   http.MethodPatch,
			Path:     api.EndpointChatUpdateSettings,
			DocsURL:  "#update-chat-settings",
			Comments: []string{
				"Updates the chat settings for a channel.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:manage:chat_settings scope.",
			},
			Params: struct {
				BroadcasterID                 string `query:"required"`
				ModeratorID                   string `query:"required"`
				EmoteMode                     bool   `body:"-"`
				FollowerMode                  bool   `body:"-"`
				FollowerModeDuration          int    `body:"-"`
				NonModeratorChatDelay         bool   `body:"-"`
				NonModeratorChatDelayDuration int    `body:"-"`
				SlowMode                      bool   `body:"-"`
				SlowModeWaitTime              int    `body:"-"`
				SubscriberMode                bool   `body:"-"`
				UniqueChatMode                bool   `body:"-"`
			}{},
			Response: BasicResponse[api.ChatSettings]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#send-chat-announcement
			Resource: ChatAnnouncementResource,
			Name:     "ChatAnnouncement",
			Method:   http.MethodPost,
			Path:     api.EndpointChatSendAnnouncement,
			DocsURL:  "#send-chat-announcement",
			Comments: []string{
				"Sends an announcement to a chat room.", "",
				"# Rate Limits", "", "One announcement may be sent every 2 seconds.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:manage:announcements scope.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
				ModeratorID   string `query:"required"`
				Message       string `body:"required"`
				Color         string `body:"-"`
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#send-a-shoutout
			Resource: ChatShoutoutResource,
			Name:     "ChatShoutout",
			Method:   http.MethodPost,
			Path:     api.EndpointChatSendShoutout,
			DocsURL:  "#send-a-shoutout",
			Comments: []string{
				"Sends a shoutout in chat to another channel.", "",
				"# Rate Limits", "",
				"The broadcaster may send a Shoutout once every 2 minutes.", "",
				"They may send the same broadcaster a Shoutout once every 60 minutes.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:manage:shoutouts scope.",
			},
			Params: struct {
				FromBroadcasterID string `query:"required"`
				ToBroadcasterID   string `query:"required"`
				ModeratorID       string `query:"required"`
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#send-chat-message
			Resource: ChatResource,
			Name:     "SendMessage",
			Method:   http.MethodPost,
			Path:     api.EndpointChatSendMessage,
			DocsURL:  "#send-chat-message",
			Comments: []string{
				"Sends a message to the specified chat room.", "",
				"# Rate Limits", "", "A user may send 20 messages every 30 seconds per channel.", "",
				"# Authorization", "", "Requires an app access token or user access token that includes the user:write:chat scope.", "",
				"If app access token used, then additionally requires user:bot scope from chatting user, and either channel:bot scope from broadcaster or moderator status.",
			},
			Params: struct {
				BroadcasterID         string `body:"required"`
				SenderID              string `body:"required"`
				Message               string `body:"required"`
				ReplayParentMessageID string `body:"-"`
				ForSourceOnly         bool   `body:"-"`
			}{},
			Response: BasicResponse[api.OutboundChatMessage]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-user-chat-color
			Resource: ChatChattersUserColorResource,
			Name:     "UserChatColor",
			Method:   http.MethodGet,
			Path:     api.EndpointChatGetUserColor,
			DocsURL:  "#get-user-chat-color",
			Comments: []string{
				"Gets the chat color for a user.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				UserID []string `query:"required"`
			}{},
			Response: BasicResponse[api.UserChatColor]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#update-user-chat-color
			Resource: ChatChattersUserColorResource,
			Name:     "UserChatColor",
			Method:   http.MethodPut,
			Path:     api.EndpointChatUpdateUserColor,
			DocsURL:  "#update-user-chat-color",
			Comments: []string{
				"Updates the display color for the user in chat.", "",
				"# Authorization", "", "Requires a user access token that includes the user:manage:chat_color scope.",
			},
			Params: struct {
				UserID string `query:"required"`
				Color  string `query:"required"`
			}{},
		},
		// Clips
		{
			// https://dev.twitch.tv/docs/api/reference/#create-clip
			Resource: ClipsResource,
			Name:     "CreateClip",
			Method:   http.MethodPost,
			Path:     api.EndpointClips,
			DocsURL:  "#create-clip",
			Comments: []string{
				"Creates a clip for a stream.", "",
				"# Authorization", "", "Requires a user access token that includes the clips:edit scope.",
			},
			Params: struct {
				BroadcasterID string `query:"-"`
				HasDelay      bool   `query:"-"`
			}{},
			Response: BasicResponse[api.Clip]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-clips
			Resource: ClipsResource,
			Name:     "Clips",
			Method:   http.MethodGet,
			Path:     api.EndpointClips,
			DocsURL:  "#get-clips",
			Comments: []string{},
			Params: struct {
				ID            []string  `query:"-"`
				BroadcasterID string    `query:"-"`
				GameID        string    `query:"-"`
				Before        string    `query:"-"`
				After         string    `query:"-"`
				First         int       `query:"-"`
				IsFeatured    bool      `query:"-"`
				StartedAt     time.Time `query:"-"`
				EndedAt       time.Time `query:"-"`
			}{},
			Response: struct {
				Data       api.Clip
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-clips-download
			Resource: ClipsDownloadResource,
			Name:     "ClipsDownload",
			Method:   http.MethodGet,
			Path:     api.EndpointClipsGetClipsDownload,
			DocsURL:  "#get-clips-download",
			Comments: []string{
				"Provides URLs to download the video file for the specified clips.", "",
				"# Rate Limits", "", "Limited to 100 requests per minute.", "",
				"# Authorization", "", "Requires an app access token or user access token that includes the editor:manage:clips or channel:manage:clips scope.",
			},
			Params: struct {
				ClipID        string `query:"required"`
				EditorID      string `query:"-"`
				BroadcasterID string `query:"-"`
			}{},
			Response: BasicResponse[api.DownloadableClip]{},
		},
		// Conduits
		{
			// https://dev.twitch.tv/docs/api/reference/#get-conduits
			Resource: ConduitsResource,
			Name:     "Conduits",
			Method:   http.MethodGet,
			Path:     api.EndpointConduits,
			DocsURL:  "#get-conduits",
			Comments: []string{
				"Gets all conduits for a client ID.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Response: BasicResponse[api.Conduit]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#create-conduits
			Resource: ConduitsResource,
			Name:     "Conduits",
			Method:   http.MethodPost,
			Path:     api.EndpointConduits,
			DocsURL:  "#create-conduits",
			Comments: []string{
				"Creates a new conduit.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Params: struct {
				ShardCount int `body:"required"`
			}{},
			Response: BasicResponse[api.Conduit]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#update-conduits
			Resource: ConduitsResource,
			Name:     "Conduits",
			Method:   http.MethodPatch,
			Path:     api.EndpointConduits,
			DocsURL:  "#update-conduits",
			Comments: []string{
				"Updates a conduit’s shard count.", "To delete shards, update the count to a lower number, and the shards above the count will be deleted.",
				"For example, if the existing shard count is 100, by resetting shard count to 50, shards 50-99 are disabled.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Params: struct {
				ID         string `body:"required"`
				ShardCount int    `body:"required"`
			}{},
			Response: BasicResponse[api.Conduit]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#delete-conduit
			Resource: ConduitsResource,
			Name:     "Conduits",
			Method:   http.MethodDelete,
			Path:     api.EndpointConduits,
			DocsURL:  "#delete-conduit",
			Comments: []string{
				"Deletes a conduit.",
				"Note that it may take some time for Eventsub subscriptions on a deleted conduit to show as disabled when calling Get Eventsub Subscriptions.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Params: struct {
				ID string `query:"required"`
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-conduit-shards
			Resource: ConduitsShardsResource,
			Name:     "ConduitsShard",
			Method:   http.MethodGet,
			Path:     api.EndpointConduitsShards,
			DocsURL:  "#get-conduits-shards",
			Comments: []string{
				"Gets a lists of all shards for a conduit.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Params: struct {
				ConduitID string `query:"required"`
				Status    string `query:"-"`
				After     string `query:"-"`
			}{},
			Response: struct {
				Data       api.ConduitShard
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#update-conduit-shards
			Resource: ConduitsShardsResource,
			Name:     "ConduitsShard",
			Method:   http.MethodPatch,
			Path:     api.EndpointConduitsShards,
			DocsURL:  "#update-conduit-shards",
			Comments: []string{
				"Updates a conduit shard.", "",
				"Shard IDs are indexed starting at 0, so a conduit with a shard_count of 5 will have shards with IDs 0 through 4.", "",
				"# Authorization", "", "Requires an app access token.",
			},
			Params: struct {
				ConduitID string             `body:"required"`
				Shards    []api.ConduitShard `body:"required"`
			}{},
			Response: BasicResponse[api.ConduitShard]{},
		},
		// Entitlements
		// EventSub
		// Games
		{
			// https://dev.twitch.tv/docs/api/reference/#get-top-games
			Resource: GamesTopResource,
			Name:     "TopGames",
			Method:   http.MethodGet,
			Path:     api.EndpointGamesTop,
			DocsURL:  "#get-top-games",
			Comments: []string{
				"Gets information about all broadcasts on Twitch.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				Before string `query:"-"`
				After  string `query:"-"`
				First  int    `query:"-"`
			}{},
			Response: struct {
				Data       api.Game
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#get-games
			Resource: GamesResource,
			Name:     "Games",
			Method:   http.MethodGet,
			Path:     api.EndpointGames,
			DocsURL:  "#get-games",
			Comments: []string{"Gets information about one or more specified games."},
			Params: struct {
				ID   []string `query:"required"`
				Name []string `query:"required"`
			}{},
			Response: BasicResponse[api.Game]{},
		},
		// Goals
		// Guest Star
		// Hype Train
		// Moderation
		{
			// https://dev.twitch.tv/docs/api/reference/#ban-user
			Resource: ModerationBansResource,
			Name:     "BanUser",
			Method:   http.MethodPost,
			Path:     api.EndpointModerationBans,
			DocsURL:  "#ban-user",
			Comments: []string{
				"Bans a user from participating in the specified broadcaster's chat room or puts them in a timeout.", "",
				"If the user is currently in a timeout, you can call this endpoint to change the duration of the timeout or ban them altogether.",
				"If the user is currently banned, you cannot call this method to put them in a timeout instead.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:manage:banned_users scope.",
			},
			Params: struct {
				BroadcasterID string            `query:"required"`
				ModeratorID   string            `query:"required"`
				Data          []api.OutboundBan `body:"required"`
			}{},
			Response: BasicResponse[api.IssuedBan]{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#unban-user
			Resource: ModerationBansResource,
			Name:     "UnbanUser",
			Method:   http.MethodDelete,
			Path:     api.EndpointModerationBans,
			DocsURL:  "#unban-user",
			Comments: []string{
				"Removes the ban or timeout that was placed on the specified user.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:manage:banned_users scope.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
				ModeratorID   string `query:"required"`
				UserID        string `query:"required"`
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#delete-chat-messages
			Resource: ModerationClearChatResource,
			Name:     "ChatMessages",
			Method:   http.MethodDelete,
			Path:     api.EndpointModerationDeleteChatMessages,
			DocsURL:  "#delete-chat-messages",
			Comments: []string{
				"Removes a single chat message or all chat messages from the broadcaster's chat room.", "",
				"# Authorization", "", "Requires a user access token that includes the moderator:manage:chat_messages scope.",
			},
			Params: struct {
				BroadcasterID string `query:"required"`
				ModeratorID   string `query:"required"`
				MessageID     string `query:"-"`
			}{},
		},
		// Polls
		// Predictions
		// Raids
		// Schedule
		// Search
		// Streams
		{
			// https://dev.twitch.tv/docs/api/reference/#get-streams
			Resource: StreamsResource,
			Name:     "Streams",
			Method:   http.MethodGet,
			Path:     api.EndpointStreams,
			DocsURL:  "#get-streams",
			Comments: []string{
				"Gets a list of all streams.", "The list is in descending order by the number of viewers watching the stream.",
				"Because viewers come and go during a stream, it's possible to find duplicate or missing streams in the list as you page through the results.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				ID       []string `query:"-"`
				UserID   []string `query:"-"`
				GameID   []string `query:"-"`
				Language string   `query:"-"`
				Period   string   `query:"-"`
				Type     string   `query:"-"`
				Before   string   `query:"-"`
				After    string   `query:"-"`
				First    int      `query:"-"`
			}{},
			Response: struct {
				Data       api.Stream
				Pagination api.Pagination
			}{},
		},
		// Subscriptions
		// Tags
		// Teams
		// Users
		{
			// https://dev.twitch.tv/docs/api/reference/#get-users
			Resource: UsersResource,
			Name:     "Users",
			Method:   http.MethodGet,
			Path:     api.EndpointUsers,
			DocsURL:  "#get-users",
			Comments: []string{
				"Gets information about one or more users.", "You may specify users by ID or by login name.", "",
				"You may look up users using their user ID, login name, or both but the sum total of the number of users you may look up is 100.", "For example, you may specify 50 IDs and 50 names or 100 IDs or names, but you cannot specify 100 IDs and 100 names.", "",
				"If you don't specify IDs or login names, the request returns information about the user in the access token if you specify a user access token.", "",
				"# Authorization", "", "Requires an app access token or user access token.", "",
				"To include the user's verified email address in the response, you must use a user access token that includes the user:read:email scope.",
			},
			Params: struct {
				ID    []string `query:"-"`
				Login []string `query:"-"`
			}{},
			Response: BasicResponse[api.User]{},
		},
		// Videos
		{
			// https://dev.twitch.tv/docs/api/reference/#get-videos
			Resource: VideosResource,
			Name:     "Videos",
			Method:   http.MethodGet,
			Path:     api.EndpointVideos,
			DocsURL:  "#get-videos",
			Comments: []string{
				"Gets information about one or more published videos.", "You may get videos by ID, by user, or by game/category.", "",
				"# Authorization", "", "Requires an app access token or user access token.",
			},
			Params: struct {
				ID       []string `query:"-"`
				UserID   []string `query:"-"`
				GameID   []string `query:"-"`
				Language string   `query:"-"`
				Period   string   `query:"-"`
				Sort     string   `query:"-"`
				Type     string   `query:"-"`
				After    string   `query:"-"`
				First    int      `query:"-"`
			}{},
			Response: struct {
				Data       api.Video
				Pagination api.Pagination
			}{},
		},
		{
			// https://dev.twitch.tv/docs/api/reference/#delete-videos
			Resource: VideosResource,
			Name:     "Videos",
			Method:   http.MethodDelete,
			Path:     api.EndpointVideos,
			DocsURL:  "#delete-videos",
			Comments: []string{
				"Deletes one or more videos.", "You may delete past broadcasts, highlights, or uploads.", "",
				"# Authorization", "", "Requires a user access token that includes the channel:manage:videos scope.",
			},
			Params: struct {
				ID string `query:"required"`
			}{},
		},
		// Whispers
		{
			// https://dev.twitch.tv/docs/api/reference/#send-whisper
			Resource: WhispersResource,
			Name:     "SendWhisper",
			Method:   http.MethodPost,
			Path:     api.EndpointWhispers,
			DocsURL:  "#send-whisper",
			Comments: []string{
				"Sends a whisper message to the specified user.", "",
				"# Rate Limits", "", "You may whisper to a maximum of 40 unique recipients per day.",
				"Within the per day limit, you may whisper a maximum of 3 whispers per second and a maximum of 100 whispers per minute.", "",
				"# Authorization", "", "The user sending the whisper must have a verified phone number (see the Phone Number setting in your Security and Privacy settings).", "",
				"Requires a user access token that includes the user:manage:whispers scope.",
			},
			Params: struct {
				FromUserID string `query:"required"`
				ToUserID   string `query:"required"`
				Message    string `body:"required"`
			}{},
		},
	}
)

var resourceMapping map[*TwitchAPIResource][]*TwitchAPIEndpoint

func NewTwitchAPIResource(name string, subresources ...*TwitchAPIResource) *TwitchAPIResource {
	return &TwitchAPIResource{Name: name, SubResources: subresources}
}

func init() {
	resourceMapping = make(map[*TwitchAPIResource][]*TwitchAPIEndpoint)
	for _, e := range Endpoints {
		resourceMapping[e.Resource] = append(resourceMapping[e.Resource], e)
	}
}

func main() {
	var resourceCount, endpointCount int
	for _, resource := range Resources {
		tmplOutput := GetFileName(resource)
		out, err := os.Create(tmplOutput)
		if err != nil {
			slog.Error("failed to create output file",
				slog.String("file", tmplOutput),
				slog.Any("error", err),
			)
			_ = out.Close()
			_ = os.Remove(tmplOutput)
			continue
		}

		endpoints := resourceMapping[resource]
		if err := Generate(out, TemplateAPIHeader, resource); err != nil {
			slog.Error("failed to create output file",
				slog.String("file", tmplOutput),
				slog.Any("error", err),
			)
			_ = out.Close()
			_ = os.Remove(tmplOutput)
			continue
		}

		numResources, numEndpoints, err := GenerateResource(out, resource, endpoints)
		if err != nil {
			slog.Error("failed to generate resource",
				slog.String("output", tmplOutput),
				slog.String("template", TemplateAPIResource),
				slog.Any("error", err),
			)
		}
		resourceCount += numResources
		endpointCount += numEndpoints
		_ = out.Close()
	}
	slog.Info("generated resources", slog.Int("resources", resourceCount), slog.Int("endpoints", endpointCount))
}

func GenerateResource(w io.Writer, resource *TwitchAPIResource, endpoints []*TwitchAPIEndpoint) (numResources, numEndpoints int, err error) {
	if err = Generate(w, TemplateAPIResource, resource); err != nil {
		return
	}

	for _, subresource := range resource.SubResources {
		subresource.Name = resource.Name + subresource.Name
		nr, ne, err := GenerateResource(w, subresource, resourceMapping[subresource])
		if err != nil {
			return numResources, numEndpoints, err
		}
		numResources += nr
		numEndpoints += ne
	}

	for _, endpoint := range endpoints {
		if strings.HasPrefix(endpoint.DocsURL, "#") {
			endpoint.DocsURL = "https://dev.twitch.tv/docs/api/reference/" + endpoint.DocsURL
		}

		if err = Generate(w, TemplateAPICall, endpoint); err != nil {
			return
		}
		numEndpoints++
	}
	numResources++
	return
}

func Generate(w io.Writer, tmplFile string, data any) error {
	path := strings.Split(tmplFile, "/")
	tmplName := path[len(path)-1]
	tmpl, err := template.New(tmplName).Funcs(templateFuncs).ParseFiles(tmplFile)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, tmplName, data)
}

func GetFileName(r *TwitchAPIResource) string {
	return "twitch_" + strings.ToLower(ToSnakeCase(strings.ReplaceAll(r.Name, " ", ""))) + ".go"
}

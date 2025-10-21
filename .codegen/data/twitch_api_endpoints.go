package data

import (
	"net/http"
	"time"

	"github.com/adeithe/go-twitch/api"
)

// BasicResponse is a generic wrapper for Twitch API responses.
type BasicResponse[T any] struct {
	Data T
}

// Endpoints is the list of all API endpoints to generate. Resource mapping is done automatically using the Resource field.
var Endpoints = []*TwitchAPIEndpoint{
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
			BroadcasterID string `body:"broadcaster_id,required"`
			Length        int    `body:"length,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			ExtensionID string    `query:"extension_id"`
			Type        string    `query:"type"`
			After       string    `query:"after"`
			First       int       `query:"first"`
			StartedAt   time.Time `query:"started_at"`
			EndedAt     time.Time `query:"ended_at"`
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
			GameID    string    `query:"game_id"`
			Type      string    `query:"type"`
			After     string    `query:"after"`
			First     int       `query:"first"`
			StartedAt time.Time `query:"started_at"`
			EndedAt   time.Time `query:"ended_at"`
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
			UserID    string    `query:"user_id"`
			Period    string    `query:"period"`
			Count     int       `query:"count"`
			StartedAt time.Time `query:"started_at"`
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
			BroadcasterID string `query:"broadcaster_id"`
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
			ID          []string `query:"id"`
			ExtensionID string   `query:"extension_id,required"`
			After       string   `query:"after"`
			First       int      `query:"first"`
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
			BroadcasterID []string `query:"broadcaster_id,required"`
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
			BroadcasterID               string                                  `query:"broadcaster_id,required"`
			GameID                      string                                  `body:"game_id"`
			BroadcasterLanguage         string                                  `body:"broadcaster_language"`
			Title                       string                                  `body:"title"`
			Tags                        []string                                `body:"tags"`
			ContentClassificationLabels []api.ChannelContentClassificationLabel `body:"content_classification_labels"`
			Delay                       int                                     `body:"delay"`
			IsBrandedContent            bool                                    `body:"is_branded_content"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			UserID        string `query:"user_id,required"`
			BroadcasterID string `query:"broadcaster_id"`
			After         string `query:"after"`
			First         int    `query:"first"`
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
			UserID        string `query:"user_id"`
			BroadcasterID string `query:"broadcaster_id,required"`
			After         string `query:"after"`
			First         int    `query:"first"`
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
			BroadcasterID                     string `query:"broadcaster_id,required"`
			Title                             string `body:"title,required"`
			Prompt                            string `body:"prompt"`
			BackgroundColor                   string `body:"prompt"`
			Cost                              int64  `body:"cost,required"`
			MaxPerStream                      int    `body:"max_per_stream"`
			MaxPerUserPerStream               int    `body:"max_per_user_per_stream"`
			GlobalCooldownSeconds             int    `body:"global_cooldown_seconds"`
			IsEnabled                         bool   `body:"is_enabled"`
			IsUserInputRequired               bool   `body:"is_user_input_required"`
			IsMaxPerStreamEnabled             bool   `body:"is_max_per_stream_enabled"`
			IsMaxPerUserPerStreamEnabled      bool   `body:"is_max_per_user_per_stream_enabled"`
			IsGlobalCooldownEnabled           bool   `body:"is_global_cooldown_enabled"`
			ShouldRedemptionsSkipRequestQueue bool   `body:"should_redemptions_skip_request_queue"`
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
			ID            string `query:"id,required"`
			BroadcasterID string `query:"broadcaster_id,required"`
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
			ID                    []string `query:"id"`
			BroadcasterID         string   `query:"broadcaster,required"`
			OnlyManageableRewards bool     `query:"only_manageable_rewards"`
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
			ID            string `query:"id"`
			BroadcasterID string `query:"broadcaster_id,required"`
			RewardID      string `query:"reward_id,required"`
			Status        string `query:"status,required"`
			Sort          string `query:"sort"`
			After         string `query:"after"`
			First         int    `query:"first"`
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
			ID                                string `query:"id,required"`
			BroadcasterID                     string `query:"broadcaster_id,required"`
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
			ID            []string `query:"id,required"`
			BroadcasterID string   `query:"broadcaster_id,required"`
			RewardID      string   `query:"reward_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
			After         string `query:"after"`
			First         int    `query:"first"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
			ModeratorID   string `query:"moderator_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			EmoteSetID []string `query:"emote_set_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
			ModeratorID   string `query:"moderator_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
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
			UserID        string `query:"user_id,required"`
			BroadcasterID string `query:"broadcaster_id"`
			After         string `query:"after"`
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
			BroadcasterID                 string `query:"broadcaster_id,required"`
			ModeratorID                   string `query:"moderator_id,required"`
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
			BroadcasterID string `query:"broadcaster_id,required"`
			ModeratorID   string `query:"moderator_id,required"`
			Message       string `body:"message,required"`
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
			FromBroadcasterID string `query:"from_broadcaster_id,required"`
			ToBroadcasterID   string `query:"to_broadcaster_id,required"`
			ModeratorID       string `query:"moderator_id,required"`
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
			BroadcasterID         string `body:"broadcaster_id,required"`
			SenderID              string `body:"sender_id,required"`
			Message               string `body:"message,required"`
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
			UserID []string `query:"user_id,required"`
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
			UserID string `query:"user_id,required"`
			Color  string `query:"color,required"`
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
			BroadcasterID string `query:"-,required"`
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
		Comments: []string{
			"Gets one or more video clips that were captured from streams.", "",
			"At least one of ID, GameID, and BroadcasterID are required. They are mutually exclusive.", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
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
			ClipID        string `query:"clip_id,required"`
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
			ShardCount int `body:"shard_count,required"`
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
			"Updates a conduit's shard count.", "To delete shards, update the count to a lower number, and the shards above the count will be deleted.",
			"For example, if the existing shard count is 100, by resetting shard count to 50, shards 50-99 are disabled.", "",
			"# Authorization", "", "Requires an app access token.",
		},
		Params: struct {
			ID         string `body:"id,required"`
			ShardCount int    `body:"shard_count,required"`
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
			ID string `query:"id,required"`
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
			ConduitID string `query:"conduit_id,required"`
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
			ConduitID string             `body:"-,required"`
			Shards    []api.ConduitShard `body:"-,required"`
		}{},
		Response: BasicResponse[api.ConduitShard]{},
	},
	// Content Classification Labels
	{
		// https://dev.twitch.tv/docs/api/reference/#get-content-classification-labels
		Resource: CCLResource,
		Name:     "ContentLabels",
		Method:   http.MethodGet,
		Path:     api.EndpointContentLabels,
		DocsURL:  "#get-content-classification-labels",
		Comments: []string{
			"Gets information about Twitch content classification labels.", "",
			"The locale parameter is a ISO 3166-1 alpha-2 which tells Twitch which translation to use for the response. (Default: en-US)", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
		Params: struct {
			Locale string `query:"-"`
		}{},
		Response: BasicResponse[api.ContentClassificationLabel]{},
	},
	// Entitlements
	{
		// https://dev.twitch.tv/docs/api/reference/#get-drops-entitlements
		Resource: EntitlementsDropsResource,
		Name:     "DropsEntitlements",
		Method:   http.MethodGet,
		Path:     api.EndpointEntitlementsDrops,
		DocsURL:  "#get-drops-entitlements",
		Comments: []string{
			"Gets an organization's list of entitlements that have been granted to a game, a user, or both.", "",
			"Entitlements returned in the response body data are not guaranteed to be sorted by any field returned by the API.",
			"To retrieve CLAIMED or FULFILLED entitlements, use the FulfillmentStatus method to filter results.",
			"To retrieve entitlements for a specific game, use the GameID method to filter results.", "",
			"# Authorization", "", "Requires an app access token or user access token.", "",
			"The associated Client ID for the access token must be owned by a user who is a member of the organization that holds ownership of the game.",
		},
		Params: struct {
			ID                []string `query:"-"`
			UserID            string   `query:"-"`
			GameID            string   `query:"-"`
			FulfillmentStatus string   `query:"-"`
			After             string   `query:"-"`
			First             int      `query:"-"`
		}{},
		Response: BasicResponse[api.DropEntitlement]{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#update-drops-entitlements
		Resource: EntitlementsDropsResource,
		Name:     "DropsEntitlements",
		Method:   http.MethodPatch,
		Path:     api.EndpointEntitlementsDrops,
		DocsURL:  "#update-drops-entitlements",
		Comments: []string{
			"Updates the Drop entitlement's fulfillment status.", "",
			"# Authorization", "", "Requires an app access token or user access token.", "",
			"The associated Client ID for the access token must be owned by a user who is a member of the organization that holds ownership of the game.",
		},
		Params: struct {
			EntitlementID     []string `body:"entitlement_ids"`
			FulfillmentStatus string   `body:"fulfillment_status"`
		}{},
		Response: BasicResponse[api.UpdatedDropEntitlement]{},
	},
	// Extensions
	{
		// https://dev.twitch.tv/docs/api/reference/#get-extensions
		Resource: ExtensionsResource,
		Name:     "Extensions",
		Method:   http.MethodGet,
		Path:     api.EndpointExtensions,
		DocsURL:  "#get-extensions",
		Comments: []string{
			"Gets information about an extension.", "",
			"# Authorization", "", "Requires a signed JSON Web Token (JWT) created by an Extension Backend Service (EBS).",
			"The signed JWT must include the \"role\" field (see JWT Schema), and the \"role\" field must be set to external.",
		},
		Params: struct {
			ExtensionID      string `query:"-,required"`
			ExtensionVersion string `query:"-"`
		}{},
		Response: BasicResponse[api.Extension]{},
	},
	// EventSub
	{
		// https://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions
		Resource: EventSubResource,
		Name:     "EventSubSubscriptions",
		Method:   http.MethodGet,
		Path:     api.EndpointEventSubSubscriptions,
		DocsURL:  "#get-eventsub-subscriptions",
		Comments: []string{
			"Gets a list of all EventSub subscriptions that the authenticated app has created.", "",
			"# Authorization", "", "Requires an app access token.",
		},
		Params: struct {
			SubscriptionID string `query:"-"`
			Status         string `query:"-"`
			Type           string `query:"-"`
			UserID         string `query:"-"`
			After          string `query:"-"`
		}{},
		Response: struct {
			TotalCost  int
			MaxCost    int
			Data       api.EventSubSubscription
			Pagination api.Pagination
		}{},
	},
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
			ID   []string `query:"-"`
			Name []string `query:"-"`
			IGDB []string `query:"igdb_id"`
		}{},
		Response: BasicResponse[api.Game]{},
	},
	// Goals
	{
		// https://dev.twitch.tv/docs/api/reference/#get-creator-goals
		Resource: GoalsResource,
		Name:     "Goals",
		Method:   http.MethodGet,
		Path:     api.EndpointGoals,
		DocsURL:  "#get-creator-goals",
		Comments: []string{
			"Gets the broadcaster's list of active goals. Use this endpoint to get the current progress of each goal.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:read:goals scope.",
		},
		Params: struct {
			BroadcasterID string `query:"broadcaster_id,required"`
		}{},
		Response: BasicResponse[api.CreatorGoal]{},
	},
	// Guest Star
	{
		// https://dev.twitch.tv/docs/api/reference/#get-guest-star-session
		Resource: GuestStarResource,
		Name:     "GuestStarSession",
		Method:   http.MethodGet,
		Path:     api.EndpointGuestStarSession,
		DocsURL:  "#get-guest-star-session",
		Comments: []string{
			"Gets information about an ongoing Guest Star session for a particular channel.", "",
			"# Authorization", "",
			"Requires OAuth Scope: channel:read:guest_star, channel:manage:guest_star, moderator:read:guest_star or moderator:manage:guest_star.", "",
			"Guests must be either invited or assigned a slot within the session.",
		},
		Params: struct {
			BroadcasterID string `query:"broadcaster_id,required"`
			ModeratorID   string `query:"moderator_id,required"`
		}{},
		Response: BasicResponse[api.GuestStarSession]{},
	},
	// Hype Train
	{
		// https://dev.twitch.tv/docs/api/reference/#get-hype-train-status
		Resource: HypeTrainResource,
		Name:     "HypeTrainStatus",
		Method:   http.MethodGet,
		Path:     api.EndpointHypeTrainGetStatus,
		DocsURL:  "#get-hype-train-status",
		Comments: []string{
			"Get the status of a Hype Train for the specified broadcaster.", "",
			"# Authorization", "", "Requires an user access token.", "",
			"Requires OAuth Scope: channel:read:hype_train.", "",
			"Requires that the user access token belongs to BroadcasterID.",
		},
		Params: struct {
			BroadcasterID string `query:"broadcaster_id,required"`
		}{},
		Response: BasicResponse[api.HypeTrainStatusInfo]{},
	},
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
			BroadcasterID string            `query:"-,required"`
			ModeratorID   string            `query:"-,required"`
			Ban           []api.OutboundBan `body:"data,required"`
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
			BroadcasterID string `query:"-,required"`
			ModeratorID   string `query:"-,required"`
			UserID        string `query:"-,required"`
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
			BroadcasterID string `query:"-,required"`
			ModeratorID   string `query:"-,required"`
			MessageID     string `query:"-"`
		}{},
	},
	// Polls
	{
		// https://dev.twitch.tv/docs/api/reference/#get-polls
		Resource: PollsResource,
		Name:     "Polls",
		Method:   http.MethodGet,
		Path:     api.EndpointPolls,
		DocsURL:  "#get-polls",
		Comments: []string{
			"Gets a list of polls that have been created in the past 90 days.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:read:polls or channel:manage:polls scope.",
		},
		Params: struct {
			ID            []string `query:"-"`
			BroadcasterID string   `query:"-,required"`
			After         string   `query:"-"`
			First         int      `query:"-"`
		}{},
		Response: struct {
			Data       api.Poll
			Pagination api.Pagination
		}{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#create-poll
		Resource: PollsResource,
		Name:     "Polls",
		Method:   http.MethodPost,
		Path:     api.EndpointPolls,
		DocsURL:  "#create-poll",
		Comments: []string{
			"Creates a poll that viewers in the broadcaster's channel can vote on.", "",
			"The poll begins as soon as it's created. You may run only one poll at a time.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:polls scope.",
		},
		Params: struct {
			BroadcasterID              string               `body:"broadcaster_id,required"`
			Title                      string               `body:"title,required"`
			Choices                    []api.OutboundChoice `body:"choices"`
			Duration                   int                  `body:"duration,required"`
			ChannelPointsPerVote       int                  `body:"channel_points_per_vote"`
			ChannelPointsVotingEnabled bool                 `body:"channel_points_voting_enabled"`
		}{},
		Response: BasicResponse[api.Poll]{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#end-poll
		Resource: PollsResource,
		Name:     "Polls",
		Method:   http.MethodPatch,
		Path:     api.EndpointPolls,
		DocsURL:  "#end-poll",
		Comments: []string{
			"Ends an active poll. You have the option to end it or end it and archive it.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:polls scope.",
		},
		Params: struct {
			BroadcasterID string `body:"broadcaster_id,required"`
			ID            string `body:"id,required"`
			Status        string `body:"status,required"`
		}{},
		Response: BasicResponse[api.Poll]{},
	},
	// Predictions
	{
		// https://dev.twitch.tv/docs/api/reference/#get-predictions
		Resource: PredictionsResource,
		Name:     "Predictions",
		Method:   http.MethodGet,
		Path:     api.EndpointPredictions,
		DocsURL:  "#get-predictions",
		Comments: []string{
			"Gets a list of Channel Points Predictions that the broadcaster created.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:read:predictions or channel:manage:predictions scope.",
		},
		Params: struct {
			BroadcasterID string   `query:"-,required"`
			ID            []string `query:"-"`
			After         string   `query:"-"`
			First         int      `query:"-"`
		}{},
		Response: struct {
			Data       api.Prediction
			Pagination api.Pagination
		}{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#create-prediction
		Resource: PredictionsResource,
		Name:     "Predictions",
		Method:   http.MethodPost,
		Path:     api.EndpointPredictions,
		DocsURL:  "#create-prediction",
		Comments: []string{
			"Creates a Channel Points Prediction.", "",
			"With a Channel Points Prediction, the broadcaster poses a question and viewers try to predict the outcome. The prediction runs as soon as it's created. The broadcaster may run only one prediction at a time.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:predictions scope.",
		},
		Params: struct {
			BroadcasterID    string               `body:"-,required"`
			Title            string               `body:"title,required"`
			Outcome          []api.OutboundChoice `body:"outcomes"`
			PredictionWindow int                  `body:"prediction_window,required"`
		}{},
		Response: BasicResponse[api.Prediction]{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#end-prediction
		Resource: PredictionsResource,
		Name:     "Predictions",
		Method:   http.MethodPatch,
		Path:     api.EndpointPredictions,
		DocsURL:  "#end-prediction",
		Comments: []string{
			"Locks, resolves, or cancels a Channel Points Prediction.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:predictions scope.",
		},
		Params: struct {
			BroadcasterID    string `body:"-,required"`
			ID               string `body:"-,required"`
			Status           string `body:"-,required"`
			WinningOutcomeID string `body:"-"`
		}{},
		Response: BasicResponse[api.Prediction]{},
	},
	// Raids
	{
		// https://dev.twitch.tv/docs/api/reference/#start-a-raid
		Resource: RaidsResource,
		Name:     "Raid",
		Method:   http.MethodPost,
		Path:     api.EndpointRaids,
		DocsURL:  "#start-a-raid",
		Comments: []string{
			"Raid another channel by sending the broadcaster's viewers to the targeted channel.", "",
			"When you call the API from a chat bot or extension, the Twitch UX pops up a window at the top of the chat room that identifies the number of viewers in the raid.",
			"The raid occurs when the broadcaster clicks Raid Now or after the 90-second countdown expires.", "",
			"# Rate Limit", "", "The limit is 10 requests within a 10-minute window.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:raids scope.",
		},
		Params: struct {
			FromBroadcasterID string `body:"from_broadcaster_id,required"`
			ToBroadcasterID   string `body:"to_broadcaster_id,required"`
		}{},
		Response: BasicResponse[api.InitializedRaid]{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#cancel-a-raid
		Resource: RaidsResource,
		Name:     "Raid",
		Method:   http.MethodDelete,
		Path:     api.EndpointRaids,
		DocsURL:  "#cancel-a-raid",
		Comments: []string{
			"Cancels a pending raid that was initiated by the broadcaster.", "",
			"You can cancel a raid at any point up until the broadcaster clicks Raid Now in the Twitch UX or the 90-second countdown expires.", "",
			"# Rate Limit", "", "The limit is 10 requests within a 10-minute window.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:raids scope.",
		},
		Params: struct {
			BroadcasterID string `query:"-,required"`
		}{},
	},
	// Schedule
	{
		// https://dev.twitch.tv/docs/api/reference/#get-channel-stream-schedule
		Resource: ScheduleResource,
		Name:     "ChannelStreamSchedule",
		Method:   http.MethodGet,
		Path:     api.EndpointScheduleGetChannelStreamSchedule,
		DocsURL:  "#get-channel-stream-schedule",
		Comments: []string{
			"Gets the broadcaster's streaming schedule. You can get the entire schedule or specific segments of the schedule.", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
		Params: struct {
			BroadcasterID string    `query:"-,required"`
			ID            string    `query:"-"`
			After         string    `query:"-"`
			First         int       `query:"-"`
			StartTime     time.Time `query:"-"`
		}{},
		Response: struct {
			Data       api.StreamSchedule
			Pagination api.Pagination
		}{},
	},
	// Search
	{
		// https://dev.twitch.tv/docs/api/reference/#search-categories
		Resource: SearchCategoriesResource,
		Name:     "SearchCategories",
		Method:   http.MethodGet,
		Path:     api.EndpointSearchCategories,
		DocsURL:  "#search-categories",
		Comments: []string{
			"Gets the games or categories that match the specified query.", "",
			"To match, the category's name must contain all parts of the query string.",
			"For example, if the query string is 42, the response includes any category name that contains 42 in the title.",
			"If the query string is a phrase like love computer, the response includes any category name that contains the words love and computer anywhere in the name.",
			"The comparison is case insensitive.", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
		Params: struct {
			Query string `query:"-,required"`
			After string `query:"-"`
			First int    `query:"-"`
		}{},
		Response: struct {
			Data       api.CategorySearchResult
			Pagination api.Pagination
		}{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#search-channels
		Resource: SearchChannelsResource,
		Name:     "SearchChannels",
		Method:   http.MethodGet,
		Path:     api.EndpointSearchChannels,
		DocsURL:  "#search-channels",
		Comments: []string{
			"Gets the channels that match the specified query and have streamed content within the past 6 months.", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
		Params: struct {
			Query    string `query:"-,required"`
			After    string `query:"-"`
			First    int    `query:"-"`
			LiveOnly bool   `query:"-"`
		}{},
		Response: struct {
			Data       api.ChannelSearchResult
			Pagination api.Pagination
		}{},
	},
	// Streams
	{
		// https://dev.twitch.tv/docs/api/reference/#get-stream-key
		Resource: StreamKeyResource,
		Name:     "StreamKey",
		Method:   http.MethodGet,
		Path:     api.EndpointStreamsGetKey,
		DocsURL:  "#get-stream-key",
		Comments: []string{
			"Gets the channel's stream key.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:read:stream_key scope.",
		},
		Params: struct {
			BroadcasterID string `query:"broadcaster_id,required"`
		}{},
		Response: BasicResponse[api.StreamKey]{},
	},
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
	{
		// https://dev.twitch.tv/docs/api/reference/#get-followed-streams
		Resource: StreamsFollowedResource,
		Name:     "StreamsFollowed",
		Method:   http.MethodGet,
		Path:     api.EndpointStreamsFollowed,
		DocsURL:  "#get-followed-streams",
		Comments: []string{
			"Gets the list of broadcasters that the user follows and that are streaming live.", "",
			"# Authorization", "", "Requires a user access token that includes the user:read:follows scope.",
		},
		Params: struct {
			UserID string `query:"user_id,required"`
			After  string `query:"-"`
			First  int    `query:"-"`
		}{},
		Response: struct {
			Data       api.Stream
			Pagination api.Pagination
		}{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#create-stream-marker
		Resource: StreamsMarkersResource,
		Name:     "StreamMarker",
		Method:   http.MethodPost,
		Path:     api.EndpointStreamsMarkers,
		DocsURL:  "#create-stream-marker",
		Comments: []string{
			"Creates a marker for a stream that is currently live.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:manage:broadcast scope.",
		},
		Params: struct {
			BroadcasterID string `body:"user_id,required"`
			Description   string `body:"description"`
		}{},
		Response: BasicResponse[api.StreamMarkerData]{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#get-stream-markers
		Resource: StreamsMarkersResource,
		Name:     "StreamMarker",
		Method:   http.MethodGet,
		Path:     api.EndpointStreamsMarkers,
		DocsURL:  "#get-stream-markers",
		Comments: []string{
			"Gets a list of markers from the user's most recent stream or from the specified VOD/video.", "",
			"A marker is an arbitrary point in a live stream that the broadcaster or editor marked, so they can return to that spot later to create video highlights (see Video Producer, Highlights in the Twitch UX).", "",
			"# Authorization", "", "Requires a user access token that includes the user:read:broadcast or channel:manage:broadcast scope.",
		},
		Params: struct {
			BroadcasterID string `query:"user_id,required"`
			VideoID       string `query:"-,required"`
			Before        string `query:"-"`
			After         string `query:"-"`
			First         int    `query:"-"`
		}{},
		Response: BasicResponse[api.StreamMarker]{},
	},
	// Subscriptions
	{
		// https://dev.twitch.tv/docs/api/reference/#get-broadcaster-subscriptions
		Resource: SubscriptionsResource,
		Name:     "BroadcasterSubscriptions",
		Method:   http.MethodGet,
		Path:     api.EndpointSubscriptionsGetBroadcasterSubscriptions,
		DocsURL:  "#get-broadcaster-subscriptions",
		Comments: []string{
			"Gets a list of users that subscribe to the specified broadcaster.", "",
			"# Authorization", "", "Requires a user access token that includes the channel:read:subscriptions scope.", "",
			"A Twitch extensions may use an app access token if the broadcaster has granted the channel:read:subscriptions scope from within the Twitch Extensions manager.",
		},
		Params: struct {
			UserID        []string `query:"-"`
			BroadcasterID string   `query:"broadcaster_id,required"`
			Before        string   `query:"-"`
			After         string   `query:"-"`
			First         int      `query:"-"`
		}{},
		Response: struct {
			Total      int
			Points     int
			Data       api.ChannelSubscription
			Pagination api.Pagination
		}{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#check-user-subscription
		Resource: SubscriptionsUserResource,
		Name:     "UserSubscription",
		Method:   http.MethodGet,
		Path:     api.EndpointSubscriptionsCheckUserSubscription,
		DocsURL:  "#check-user-subscription",
		Comments: []string{
			"Checks whether the user subscribes to the broadcaster's channel.", "",
			"# Authorization", "", "Requires a user access token that includes the user:read:subscriptions scope.", "",
			"A Twitch extensions may use an app access token if the broadcaster has granted the user:read:subscriptions scope from within the Twitch Extensions manager.",
		},
		Params: struct {
			BroadcasterID string `query:"broadcaster_id,required"`
			UserID        string `query:"user_id,required"`
		}{},
		Response: BasicResponse[api.UserSubscriptionStatus]{},
	},
	// Teams
	{
		// https://dev.twitch.tv/docs/api/reference/#get-channel-teams
		Resource: TeamsChannelsResource,
		Name:     "ChannelTeams",
		Method:   http.MethodGet,
		Path:     api.EndpointTeamsGetChannelTeams,
		DocsURL:  "#get-channel-teams",
		Comments: []string{
			"Gets the list of Twitch teams that the broadcaster is a member of.", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
		Params: struct {
			BroadcasterID string `query:"broadcaster_id,required"`
		}{},
		Response: BasicResponse[api.ChannelTeam]{},
	},
	{
		// https://dev.twitch.tv/docs/api/reference/#get-teams
		Resource: TeamsResource,
		Name:     "Teams",
		Method:   http.MethodGet,
		Path:     api.EndpointTeams,
		DocsURL:  "#get-teams",
		Comments: []string{
			"Gets information about the specified Twitch team.", "",
			"You are required to specify Name or ID as they are mutually exclusive.", "",
			"# Authorization", "", "Requires an app access token or user access token.",
		},
		Params: struct {
			Name string `query:"-"`
			ID   string `query:"-"`
		}{},
		Response: BasicResponse[api.Team]{},
	},
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
			ID string `query:"-,required"`
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
			FromUserID string `query:"-,required"`
			ToUserID   string `query:"-,required"`
			Message    string `body:"-,required"`
		}{},
	},
}

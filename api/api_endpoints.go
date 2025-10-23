package api

const (
	// BaseURL is the base URL for the Twitch API.
	BaseURL string = "https://api.twitch.tv"
	// TwitchAPIVersionHelix is the base path for the Helix API.
	TwitchAPIVersionHelix string = "/helix"
)

const (
	// EndpointAdsStartCommercial is the endpoint for starting a commercial.
	EndpointAdsStartCommercial = TwitchAPIVersionHelix + "/channels/commercial"
	// EndpointAdsGetAdsSchedule is the endpoint for getting ad schedule information.
	EndpointAdsGetAdsSchedule = TwitchAPIVersionHelix + "/channels/ads"
	// EndpointAdsSnoozeNextAd is the endpoint for snoozing the next ad.
	EndpointAdsSnoozeNextAd = TwitchAPIVersionHelix + "/channels/ads/schedule/snooze"
	// EndpointAnalyticsGetExtensionAnalytics is the endpoint for getting extension analytics.
	EndpointAnalyticsGetExtensionAnalytics = TwitchAPIVersionHelix + "/analytics/extensions"
	// EndpointAnalyticsGetGameAnalytics is the endpoint for getting game analytics.
	EndpointAnalyticsGetGameAnalytics = TwitchAPIVersionHelix + "/analytics/games"
	// EndpointBitsGetLeaderboard is the endpoint for getting the Bits leaderboard.
	EndpointBitsGetLeaderboard = TwitchAPIVersionHelix + "/bits/leaderboard"
	// EndpointBitsGetCheermotes is the endpoint for getting cheermotes.
	EndpointBitsGetCheermotes = TwitchAPIVersionHelix + "/bits/cheermotes"
	// EndpointBitsGetExtensionTransactions is the endpoint for getting extension transactions.
	EndpointBitsGetExtensionTransactions = TwitchAPIVersionHelix + "/extensions/transactions"
	// EndpointChannelPointsCreateCustomRewards is the endpoint for creating custom rewards.
	EndpointChannelPointsCreateCustomRewards = TwitchAPIVersionHelix + "/channel_points/custom_rewards"
	// EndpointChannelPointsDeleteCustomReward is the endpoint for deleting a custom reward.
	EndpointChannelPointsDeleteCustomReward = TwitchAPIVersionHelix + "/channel_points/custom_rewards"
	// EndpointChannelPointsGetCustomRewards is the endpoint for getting custom rewards.
	EndpointChannelPointsGetCustomRewards = TwitchAPIVersionHelix + "/channel_points/custom_rewards"
	// EndpointChannelPointsGetCustomRewardRedemptions is the endpoint for getting custom reward redemptions.
	EndpointChannelPointsGetCustomRewardRedemptions = TwitchAPIVersionHelix + "/channel_points/custom_rewards/redemptions"
	// EndpointChannelPointsUpdateCustomReward is the endpoint for updating a custom reward.
	EndpointChannelPointsUpdateCustomReward = TwitchAPIVersionHelix + "/channel_points/custom_rewards"
	// EndpointChannelPointsUpdateRedemptionStatus is the endpoint for updating the status of a reward redemption.
	EndpointChannelPointsUpdateRedemptionStatus = TwitchAPIVersionHelix + "/channel_points/custom_rewards/redemptions"
	// EndpointChannelsGetInformation is the endpoint for getting channel information.
	EndpointChannelsGetInformation = TwitchAPIVersionHelix + "/channels"
	// EndpointChannelsModifyInformation is the endpoint for modifying channel information.
	EndpointChannelsModifyInformation = TwitchAPIVersionHelix + "/channels"
	// EndpointChannelsGetEditors is the endpoint for getting channel editors.
	EndpointChannelsGetEditors = TwitchAPIVersionHelix + "/channels/editors"
	// EndpointChannelsGetFollowedChannels is the endpoint for getting followed channels.
	EndpointChannelsGetFollowedChannels = TwitchAPIVersionHelix + "/channels/followed"
	// EndpointChannelsGetFollowers is the endpoint for getting channel followers.
	EndpointChannelsGetFollowers = TwitchAPIVersionHelix + "/channels/followers"
	// EndpointCharityGetCampaign is the endpoint for getting charity campaigns.
	EndpointCharityGetCampaign = TwitchAPIVersionHelix + "/charity/campaigns"
	// EndpointCharityGetCampaignDonations is the endpoint for getting charity campaign donations.
	EndpointCharityGetCampaignDonations = TwitchAPIVersionHelix + "/charity/donations"
	// EndpointChatGetChatters is the endpoint for getting chatters in a channel.
	EndpointChatGetChatters = TwitchAPIVersionHelix + "/chat/chatters"
	// EndpointChatGetChannelEmotes is the endpoint for getting channel emotes.
	EndpointChatGetChannelEmotes = TwitchAPIVersionHelix + "/chat/emotes"
	// EndpointChatGetGlobalEmotes is the endpoint for getting global emotes.
	EndpointChatGetGlobalEmotes = TwitchAPIVersionHelix + "/chat/emotes/global"
	// EndpointChatGetEmoteSets is the endpoint for getting emote sets.
	EndpointChatGetEmoteSets = TwitchAPIVersionHelix + "/chat/emotes/set"
	// EndpointChatGetChannelBadges is the endpoint for getting a channels chat badges.
	EndpointChatGetChannelBadges = TwitchAPIVersionHelix + "/chat/badges"
	// EndpointChatGetGlobalBadges is the endpoint for getting global chat badges.
	EndpointChatGetGlobalBadges = TwitchAPIVersionHelix + "/chat/badges/global"
	// EndpointChatGetSettings is the endpoint for getting chat settings.
	EndpointChatGetSettings = TwitchAPIVersionHelix + "/chat/settings"
	// EndpointChatGetSharedChatSession is the endpoint for getting a shared chat session.
	EndpointChatGetSharedChatSession = TwitchAPIVersionHelix + "/shared_chat/session"
	// EndpointChatGetUserEmotes is the endpoint for getting user emotes.
	EndpointChatGetUserEmotes = TwitchAPIVersionHelix + "/chat/emotes/user"
	// EndpointChatUpdateSettings is the endpoint for updating a chatrooms settings.
	EndpointChatUpdateSettings = TwitchAPIVersionHelix + "/chat/settings"
	// EndpointChatSendAnnouncement is the endpoint for sending an announcement to a chatroom.
	EndpointChatSendAnnouncement = TwitchAPIVersionHelix + "/chat/announcements"
	// EndpointChatSendShoutout is the endpoint for sending a shoutout in chat.
	EndpointChatSendShoutout = TwitchAPIVersionHelix + "/chat/shoutouts"
	// EndpointChatSendMessage is the endpoint for sending a message in chat.
	EndpointChatSendMessage = TwitchAPIVersionHelix + "/chat/messages"
	// EndpointChatGetUserColor is the endpoint for getting a user's chat color.
	EndpointChatGetUserColor = TwitchAPIVersionHelix + "/chat/color"
	// EndpointChatUpdateUserColor is the endpoint for updating a user's chat color.
	EndpointChatUpdateUserColor = TwitchAPIVersionHelix + "/chat/color"
	// EndpointClips is the endpoint for the Twitch Clips API.
	EndpointClips = TwitchAPIVersionHelix + "/clips"
	// EndpointClipsGetClipsDownload is the endpoint for getting clip download URLs.
	EndpointClipsGetClipsDownload = TwitchAPIVersionHelix + "/clips/downloads"
	// EndpointConduitsShards is the endpoint for the Twitch Eventsub Conduit Shards API.
	EndpointConduitsShards = TwitchAPIVersionHelix + "/eventsub/conduits/shards"
	// EndpointConduits is the endpoint for the Twitch Eventsub Conduits API.
	EndpointConduits = TwitchAPIVersionHelix + "/eventsub/conduits"
	// EndpointContentLabels is the endpoint for getting information about content classification labels.
	EndpointContentLabels = TwitchAPIVersionHelix + "/content_classification_labels"
	// EndpointEntitlementsDrops is the endpoint for the Twitch Entitlements Drops API.
	EndpointEntitlementsDrops = TwitchAPIVersionHelix + "/entitlements/drops"
	// EndpointEventSubSubscriptions is the endpoint for managing EventSub subscriptions.
	EndpointEventSubSubscriptions = TwitchAPIVersionHelix + "/eventsub/subscriptions"
	// EndpointExtensionsConfiguration is the endpoint for managing extension configurations.
	EndpointExtensionsConfiguration = TwitchAPIVersionHelix + "/extensions/configurations"
	// EndpointExtensionsRequiredConfiguration is the endpoint for managing required extension configurations.
	EndpointExtensionsRequiredConfiguration = TwitchAPIVersionHelix + "/extensions/required_configuration"
	// EndpointExtensionsPubSub is the endpoint for sending extension messages via PubSub.
	EndpointExtensionsPubSub = TwitchAPIVersionHelix + "/extensions/pubsub"
	// EndpointExtensionsLiveChannels is the endpoint for fetching live channels by extension.
	EndpointExtensionsLiveChannels = TwitchAPIVersionHelix + "/extensions/live"
	// EndpointExtensionsSecrets is the endpoint for managing extension secrets.
	EndpointExtensionsSecrets = TwitchAPIVersionHelix + "/extensions/jwt/secrets"
	// EndpointExtensionsChat is the endpoint for sending chat messages via extensions.
	EndpointExtensionsChat = TwitchAPIVersionHelix + "/extensions/chat"
	// EndpointExtensions is the endpoint for getting information about an extension.
	EndpointExtensions = TwitchAPIVersionHelix + "/extensions"
	// EndpointExtensionsReleased is the endpoint for getting information about released extensions.
	EndpointExtensionsReleased = TwitchAPIVersionHelix + "/extensions/released"
	// EndpointExtensionsBitsProducts is the endpoint for managing an extensions Bits product.
	EndpointExtensionsBitsProducts = TwitchAPIVersionHelix + "/bits/extensions"
	// EndpointGames is the endpoint for the Twitch Games API.
	EndpointGames = TwitchAPIVersionHelix + "/games"
	// EndpointGamesTop is the endpoint for the Twitch Top Games API.
	EndpointGamesTop = TwitchAPIVersionHelix + "/games/top"
	// EndpointGoals is the endpoint for managing creator goals.
	EndpointGoals = TwitchAPIVersionHelix + "/goals"
	// EndpointGuestStarChannelSettings is the endpoint for managing guest star channel settings.
	EndpointGuestStarChannelSettings = TwitchAPIVersionHelix + "/guest_star/channel_settings"
	// EndpointGuestStarSession is the endpoint for managing a guest star session.
	EndpointGuestStarSession = TwitchAPIVersionHelix + "/guest_star/session"
	// EndpointGuestStarInvites is the endpoint for managing guest star invites.
	EndpointGuestStarInvites = TwitchAPIVersionHelix + "/guest_star/invites"
	// EndpointGuestStarSlot is the endpoint for managing a user that was invited to the guest star session.
	EndpointGuestStarSlot = TwitchAPIVersionHelix + "/guest_star/slot"
	// EndpointGuestStarSlotSettings is the endpoint for managing an invited users settings in a guest star session.
	EndpointGuestStarSlotSettings = TwitchAPIVersionHelix + "/guest_star/slot_settings"
	// EndpointHypeTrainGetEvents is the endpoint for getting Hype Train events.
	//
	// Deprecated: Scheduled for removal on December 4, 2025. Use "Get Hype Train Status" instead.
	EndpointHypeTrainGetEvents = TwitchAPIVersionHelix + "/hypetrain/events"
	// EndpointHypeTrainGetStatus is the endpoint for getting the status of a Hype Train.
	EndpointHypeTrainGetStatus = TwitchAPIVersionHelix + "/hypetrain/status"
	// EndpointModerationCheckAutoModStatus is the endpoint for checking if a message would be held by AutoMod.
	EndpointModerationCheckAutoModStatus = TwitchAPIVersionHelix + "/moderation/enforcements/status"
	// EndpointModerationManageHeldAutoModMessages is the endpoint for managing a message held by AutoMod.
	EndpointModerationManageHeldAutoModMessages = TwitchAPIVersionHelix + "/moderation/automod/message"
	// EndpointModerationAutoModSettings is the endpoint for getting and maanaging AutoMod settings for a channel.
	EndpointModerationAutoModSettings = TwitchAPIVersionHelix + "/moderation/automod/settings"
	// EndpointModerationGetBannedUsers is the endpoint for getting banned users in a channel.
	EndpointModerationGetBannedUsers = TwitchAPIVersionHelix + "/moderation/banned"
	// EndpointModerationBans is the endpoint for managing banned users in a channel.
	EndpointModerationBans = TwitchAPIVersionHelix + "/moderation/bans"
	// EndpointModerationUnbanRequests is the endpoint for managing unban requests for a channel.
	EndpointModerationUnbanRequests = TwitchAPIVersionHelix + "/moderation/unban_requests"
	// EndpointModerationBlockedTerms is the endpoint for managing a channels blocked terms.
	EndpointModerationBlockedTerms = TwitchAPIVersionHelix + "/moderation/blocked_terms"
	// EndpointModerationDeleteChatMessages is the endpoint for deleting chat messages in a channel.
	EndpointModerationDeleteChatMessages = TwitchAPIVersionHelix + "/moderation/chat"
	// EndpointModerationModeratedChannels is the endpoint for checking which channels the authenticated user moderates.
	EndpointModerationModeratedChannels = TwitchAPIVersionHelix + "/moderation/channels"
	// EndpointModerationModerators is the endpoint for managing moderators in a channel.
	EndpointModerationModerators = TwitchAPIVersionHelix + "/moderation/moderators"
	// EndpointModerationVIPs is the endpoint for managing VIPs in a channel.
	EndpointModerationVIPs = TwitchAPIVersionHelix + "/channels/vips"
	// EndpointModerationShieldMode is the endpoint for managing a channels Shield Mode settings.
	EndpointModerationShieldMode = TwitchAPIVersionHelix + "/moderation/shield_mode"
	// EndpointModerationWarnChatUser is the endpoint for warning a user in chat.
	EndpointModerationWarnChatUser = TwitchAPIVersionHelix + "/moderation/warnings"
	// EndpointPolls is the endpoint for the Twitch Polls API.
	EndpointPolls = TwitchAPIVersionHelix + "/polls"
	// EndpointPredictions is the endpoint for managing channel predictions.
	EndpointPredictions = TwitchAPIVersionHelix + "/predictions"
	// EndpointRaids is the endpoint for managing raids.
	EndpointRaids = TwitchAPIVersionHelix + "/raids"
	// EndpointScheduleGetChannelStreamSchedule is the endpoint for getting the stream schedule for a channel.
	EndpointScheduleGetChannelStreamSchedule = TwitchAPIVersionHelix + "/schedule"
	// EndpointScheduleGetChannelCalendar is the endpoint for getting the iCalendar schedule for a channel.
	EndpointScheduleGetChannelCalendar = TwitchAPIVersionHelix + "/schedule/icalendar"
	// EndpointScheduleChannelSettings is the endpoint for getting or updating channel schedule settings.
	EndpointScheduleChannelSettings = TwitchAPIVersionHelix + "/schedule/settings"
	// EndpointScheduleCreateChannelSegment is the endpoint for managing stream schedule segments for a channel.
	EndpointScheduleCreateChannelSegment = TwitchAPIVersionHelix + "/schedule/segment"
	// EndpointSearchCategories is the endpoint for searching categories.
	EndpointSearchCategories = TwitchAPIVersionHelix + "/search/categories"
	// EndpointSearchChannels is the endpoint for searching channels.
	EndpointSearchChannels = TwitchAPIVersionHelix + "/search/channels"
	// EndpointStreamsGetKey is the endpoint for getting a stream key.
	EndpointStreamsGetKey = TwitchAPIVersionHelix + "/streams/key"
	// EndpointStreams is the endpoint for getting stream information.
	EndpointStreams = TwitchAPIVersionHelix + "/streams"
	// EndpointStreamsFollowed is the endpoint for getting followed streams.
	EndpointStreamsFollowed = TwitchAPIVersionHelix + "/streams/followed"
	// EndpointStreamsMarkers is the endpoint for managing a markers for a stream.
	EndpointStreamsMarkers = TwitchAPIVersionHelix + "/streams/markers"
	// EndpointSubscriptionsGetBroadcasterSubscriptions is the endpoint for getting broadcaster's subscriptions.
	EndpointSubscriptionsGetBroadcasterSubscriptions = TwitchAPIVersionHelix + "/subscriptions"
	// EndpointSubscriptionsCheckUserSubscription is the endpoint for getting a user's subscription.
	EndpointSubscriptionsCheckUserSubscription = TwitchAPIVersionHelix + "/subscriptions/user"
	// EndpointTagsGetAllStreamTags is the endpoint for getting all stream tags.
	EndpointTagsGetAllStreamTags = TwitchAPIVersionHelix + "/tags/streams"
	// EndpointTagsGetStreamTags is the endpoint for getting tags for a stream.
	EndpointTagsGetStreamTags = TwitchAPIVersionHelix + "/streams/tags"
	// EndpointTeams is the endpoint for getting teams.
	EndpointTeams = TwitchAPIVersionHelix + "/teams"
	// EndpointTeamsGetChannelTeams is the endpoint for getting a team by name.
	EndpointTeamsGetChannelTeams = TwitchAPIVersionHelix + "/teams/channel"
	// EndpointUsers is the endpoint for getting users.
	EndpointUsers = TwitchAPIVersionHelix + "/users"
	// EndpointUsersGetAuthorizationByUser is the endpoint for getting authorizations by user.
	EndpointUsersGetAuthorizationByUser = TwitchAPIVersionHelix + "/authorization/users"
	// EndpointUsersBlocks is the endpoint for managing blocked users.
	EndpointUsersBlocks = TwitchAPIVersionHelix + "/users/blocks"
	// EndpointUsersAllExtensions is the endpoint for getting all extensions a user has installed.
	EndpointUsersAllExtensions = TwitchAPIVersionHelix + "/users/extensions/list"
	// EndpointUsersActiveExtensions is the endpoint for getting active extensions for a user.
	EndpointUsersActiveExtensions = TwitchAPIVersionHelix + "/users/extensions"
	// EndpointVideos is the endpoint for the Twitch Videos API.
	EndpointVideos = TwitchAPIVersionHelix + "/videos"
	// EndpointWhispers is the endpoint for sending whispers.
	EndpointWhispers = TwitchAPIVersionHelix + "/whispers"
)

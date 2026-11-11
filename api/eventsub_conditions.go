package api

// EventSub condition helpers for https://dev.twitch.tv/docs/eventsub/eventsub-reference/#conditions

func eventSubBroadcasterOnly(broadcasterUserID string) EventSubCondition {
	return EventSubCondition{BroadcasterUserID: broadcasterUserID}
}

func eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return EventSubCondition{
		BroadcasterUserID: broadcasterUserID,
		ModeratorUserID:   moderatorUserID,
	}
}

func eventSubBroadcasterUser(broadcasterUserID, userID string) EventSubCondition {
	return EventSubCondition{
		BroadcasterUserID: broadcasterUserID,
		UserID:            userID,
	}
}

// WithAutomodMessageHoldCondition creates a condition for the automod.message.hold subscription.
func WithAutomodMessageHoldCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithAutomodMessageUpdateCondition creates a condition for the automod.message.update subscription.
func WithAutomodMessageUpdateCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithAutomodSettingsUpdateCondition creates a condition for the automod.settings.update subscription.
func WithAutomodSettingsUpdateCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithAutomodTermsUpdateCondition creates a condition for the automod.terms.update subscription.
func WithAutomodTermsUpdateCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelAdBreakBeginCondition creates a condition for the channel.ad_break.begin subscription.
func WithChannelAdBreakBeginCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelBanCondition creates a condition for the channel.ban subscription.
func WithChannelBanCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelBitsUseCondition creates a condition for the channel.bits.use subscription.
func WithChannelBitsUseCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelChatClearCondition creates a condition for the channel.chat.clear subscription.
func WithChannelChatClearCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatClearUserMessagesCondition creates a condition for the channel.chat.clear_user_messages subscription.
func WithChannelChatClearUserMessagesCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatMessageCondition creates a condition for the channel.chat.message subscription.
func WithChannelChatMessageCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatMessageDeleteCondition creates a condition for the channel.chat.message_delete subscription.
func WithChannelChatMessageDeleteCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatNotificationCondition creates a condition for the channel.chat.notification subscription.
func WithChannelChatNotificationCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatSettingsUpdateCondition creates a condition for the channel.chat_settings.update subscription.
func WithChannelChatSettingsUpdateCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatUserMessageHoldCondition creates a condition for the channel.chat.user_message_hold subscription.
func WithChannelChatUserMessageHoldCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelChatUserMessageUpdateCondition creates a condition for the channel.chat.user_message_update subscription.
func WithChannelChatUserMessageUpdateCondition(broadcasterUserID, userID string) EventSubCondition {
	return eventSubBroadcasterUser(broadcasterUserID, userID)
}

// WithChannelSubscribeCondition creates a condition for the channel.subscribe subscription.
func WithChannelSubscribeCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSubscriptionEndCondition creates a condition for the channel.subscription.end subscription.
func WithChannelSubscriptionEndCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSubscriptionGiftCondition creates a condition for the channel.subscription.gift subscription.
func WithChannelSubscriptionGiftCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSubscriptionMessageCondition creates a condition for the channel.subscription.message subscription.
func WithChannelSubscriptionMessageCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelCheerCondition creates a condition for the channel.cheer subscription.
func WithChannelCheerCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelUpdateCondition creates a condition for the channel.update subscription.
func WithChannelUpdateCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelFollowCondition creates a condition for the channel.follow subscription.
func WithChannelFollowCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelUnbanCondition creates a condition for the channel.unban subscription.
func WithChannelUnbanCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelUnbanRequestCreateCondition creates a condition for the channel.unban_request.create subscription.
func WithChannelUnbanRequestCreateCondition(moderatorUserID, broadcasterUserID string) EventSubCondition {
	return EventSubCondition{
		ModeratorUserID:   moderatorUserID,
		BroadcasterUserID: broadcasterUserID,
	}
}

// WithChannelUnbanRequestResolveCondition creates a condition for the channel.unban_request.resolve subscription.
func WithChannelUnbanRequestResolveCondition(moderatorUserID, broadcasterUserID string) EventSubCondition {
	return EventSubCondition{
		ModeratorUserID:   moderatorUserID,
		BroadcasterUserID: broadcasterUserID,
	}
}

// WithChannelRaidFromBroadcasterCondition creates a condition for the channel.raid subscription when monitoring outgoing raids.
func WithChannelRaidFromBroadcasterCondition(fromBroadcasterUserID string) EventSubCondition {
	return EventSubCondition{FromBroadcasterUserID: fromBroadcasterUserID}
}

// WithChannelRaidToBroadcasterCondition creates a condition for the channel.raid subscription when monitoring incoming raids.
func WithChannelRaidToBroadcasterCondition(toBroadcasterUserID string) EventSubCondition {
	return EventSubCondition{ToBroadcasterUserID: toBroadcasterUserID}
}

// WithChannelModerateCondition creates a condition for the channel.moderate subscription.
func WithChannelModerateCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelModerateV2Condition creates a condition for the channel.moderate v2 subscription.
func WithChannelModerateV2Condition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelModeratorAddCondition creates a condition for the channel.moderator.add subscription.
func WithChannelModeratorAddCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelModeratorRemoveCondition creates a condition for the channel.moderator.remove subscription.
func WithChannelModeratorRemoveCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelGuestStarSessionBeginCondition creates a condition for the channel.guest_star_session.begin subscription.
func WithChannelGuestStarSessionBeginCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelGuestStarSessionEndCondition creates a condition for the channel.guest_star_session.end subscription.
func WithChannelGuestStarSessionEndCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelGuestStarGuestUpdateCondition creates a condition for the channel.guest_star_guest.update subscription.
func WithChannelGuestStarGuestUpdateCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelGuestStarSettingsUpdateCondition creates a condition for the channel.guest_star_settings.update subscription.
func WithChannelGuestStarSettingsUpdateCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelPointsAutomaticRewardRedemptionAddCondition creates a condition for the channel.channel_points_automatic_reward_redemption.add subscription.
func WithChannelPointsAutomaticRewardRedemptionAddCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPointsAutomaticRewardRedemptionAddV2Condition creates a condition for the channel.channel_points_automatic_reward_redemption.add v2 subscription.
func WithChannelPointsAutomaticRewardRedemptionAddV2Condition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPointsCustomRewardAddCondition creates a condition for the channel.channel_points_custom_reward.add subscription.
func WithChannelPointsCustomRewardAddCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPointsCustomRewardUpdateCondition creates a condition for the channel.channel_points_custom_reward.update subscription.
func WithChannelPointsCustomRewardUpdateCondition(broadcasterUserID string, rewardID ...string) EventSubCondition {
	condition := eventSubBroadcasterOnly(broadcasterUserID)
	if len(rewardID) > 0 {
		condition.RewardID = rewardID[0]
	}
	return condition
}

// WithChannelPointsCustomRewardRemoveCondition creates a condition for the channel.channel_points_custom_reward.remove subscription.
func WithChannelPointsCustomRewardRemoveCondition(broadcasterUserID string, rewardID ...string) EventSubCondition {
	condition := eventSubBroadcasterOnly(broadcasterUserID)
	if len(rewardID) > 0 {
		condition.RewardID = rewardID[0]
	}
	return condition
}

// WithChannelPointsCustomRewardRedemptionAddCondition creates a condition for the channel.channel_points_custom_reward_redemption.add subscription.
func WithChannelPointsCustomRewardRedemptionAddCondition(broadcasterUserID string, rewardID ...string) EventSubCondition {
	condition := eventSubBroadcasterOnly(broadcasterUserID)
	if len(rewardID) > 0 {
		condition.RewardID = rewardID[0]
	}
	return condition
}

// WithChannelPointsCustomRewardRedemptionUpdateCondition creates a condition for the channel.channel_points_custom_reward_redemption.update subscription.
func WithChannelPointsCustomRewardRedemptionUpdateCondition(broadcasterUserID string, rewardID ...string) EventSubCondition {
	condition := eventSubBroadcasterOnly(broadcasterUserID)
	if len(rewardID) > 0 {
		condition.RewardID = rewardID[0]
	}
	return condition
}

// WithChannelCustomPowerUpRedemptionAddCondition creates a condition for the channel.channel_custom_power_up_redemption.add subscription.
func WithChannelCustomPowerUpRedemptionAddCondition(broadcasterUserID string, rewardID ...string) EventSubCondition {
	condition := eventSubBroadcasterOnly(broadcasterUserID)
	if len(rewardID) > 0 {
		condition.RewardID = rewardID[0]
	}
	return condition
}

// WithChannelPollBeginCondition creates a condition for the channel.poll.begin subscription.
func WithChannelPollBeginCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPollProgressCondition creates a condition for the channel.poll.progress subscription.
func WithChannelPollProgressCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPollEndCondition creates a condition for the channel.poll.end subscription.
func WithChannelPollEndCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPredictionBeginCondition creates a condition for the channel.prediction.begin subscription.
func WithChannelPredictionBeginCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPredictionProgressCondition creates a condition for the channel.prediction.progress subscription.
func WithChannelPredictionProgressCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPredictionLockCondition creates a condition for the channel.prediction.lock subscription.
func WithChannelPredictionLockCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelPredictionEndCondition creates a condition for the channel.prediction.end subscription.
func WithChannelPredictionEndCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSharedChatSessionBeginCondition creates a condition for Channel Shared Chat Session Begin events.
func WithChannelSharedChatSessionBeginCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSharedChatSessionUpdateCondition creates a condition for Channel Shared Chat Session Update events.
func WithChannelSharedChatSessionUpdateCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSharedChatSessionEndCondition creates a condition for Channel Shared Chat Session End events.
func WithChannelSharedChatSessionEndCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelSuspiciousUserMessageCondition creates a condition for the channel.suspicious_user.message subscription.
func WithChannelSuspiciousUserMessageCondition(moderatorUserID, broadcasterUserID string) EventSubCondition {
	return EventSubCondition{
		ModeratorUserID:   moderatorUserID,
		BroadcasterUserID: broadcasterUserID,
	}
}

// WithChannelSuspiciousUserUpdateCondition creates a condition for the channel.suspicious_user.update subscription.
func WithChannelSuspiciousUserUpdateCondition(moderatorUserID, broadcasterUserID string) EventSubCondition {
	return EventSubCondition{
		ModeratorUserID:   moderatorUserID,
		BroadcasterUserID: broadcasterUserID,
	}
}

// WithChannelVIPAddCondition creates a condition for the channel.vip.add subscription.
func WithChannelVIPAddCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelVIPRemoveCondition creates a condition for the channel.vip.remove subscription.
func WithChannelVIPRemoveCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithChannelWarningAcknowledgeCondition creates a condition for the channel.warning.acknowledge subscription.
func WithChannelWarningAcknowledgeCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithChannelWarningSendCondition creates a condition for the channel.warning.send subscription.
func WithChannelWarningSendCondition(broadcasterUserID, moderatorUserID string) EventSubCondition {
	return eventSubBroadcasterModerator(broadcasterUserID, moderatorUserID)
}

// WithConduitShardDisabledCondition creates a condition for the conduit.shard.disabled subscription.
func WithConduitShardDisabledCondition(clientID string, conduitID ...string) EventSubCondition {
	condition := EventSubCondition{ClientID: clientID}
	if len(conduitID) > 0 {
		condition.ConduitID = conduitID[0]
	}
	return condition
}

// WithDropEntitlementGrantCondition creates a condition for the drop.entitlement.grant subscription.
func WithDropEntitlementGrantCondition(organizationID, categoryID, campaignID string) EventSubCondition {
	return EventSubCondition{
		OrganizationID: organizationID,
		CategoryID:     categoryID,
		CampaignID:     campaignID,
	}
}

// WithExtensionBitsTransactionCreateCondition creates a condition for the extension.bits_transaction.create subscription.
func WithExtensionBitsTransactionCreateCondition(extensionClientID string) EventSubCondition {
	return EventSubCondition{ExtensionClientID: extensionClientID}
}

// WithGoalsCondition creates a condition for the channel.goal progress subscription types.
func WithGoalsCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithHypeTrainBeginCondition creates a condition for the channel.hype_train.begin subscription.
func WithHypeTrainBeginCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithHypeTrainProgressCondition creates a condition for the channel.hype_train.progress subscription.
func WithHypeTrainProgressCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithHypeTrainEndCondition creates a condition for the channel.hype_train.end subscription.
func WithHypeTrainEndCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithStreamOnlineCondition creates a condition for the stream.online subscription.
func WithStreamOnlineCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithStreamOfflineCondition creates a condition for the stream.offline subscription.
func WithStreamOfflineCondition(broadcasterUserID string) EventSubCondition {
	return eventSubBroadcasterOnly(broadcasterUserID)
}

// WithUserAuthorizationGrantCondition creates a condition for the user.authorization.grant subscription.
func WithUserAuthorizationGrantCondition(clientID string) EventSubCondition {
	return EventSubCondition{ClientID: clientID}
}

// WithUserAuthorizationRevokeCondition creates a condition for the user.authorization.revoke subscription.
func WithUserAuthorizationRevokeCondition(clientID string) EventSubCondition {
	return EventSubCondition{ClientID: clientID}
}

// WithUserUpdateCondition creates a condition for the user.update subscription.
func WithUserUpdateCondition(userID string) EventSubCondition {
	return EventSubCondition{UserID: userID}
}

// WithWhisperReceivedCondition creates a condition for Whisper Received events.
func WithWhisperReceivedCondition(userID string) EventSubCondition {
	return EventSubCondition{UserID: userID}
}

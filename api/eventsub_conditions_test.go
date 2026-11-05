package api_test

import (
	"encoding/json"
	"testing"

	"github.com/adeithe/go-twitch/api"
	"github.com/stretchr/testify/require"
)

const (
	testBroadcasterUserID = "1337"
	testModeratorUserID   = "9001"
	testUserID            = "42"
	testClientID          = "client-id"
	testConduitID         = "conduit-id"
	testExtensionClientID = "extension-client-id"
	testOrganizationID    = "org-id"
	testCategoryID        = "category-id"
	testCampaignID        = "campaign-id"
	testRewardID          = "reward-id"
)

func assertEventSubConditionJSON(t *testing.T, cond api.EventSubCondition, expected string) {
	t.Helper()
	bs, err := json.Marshal(cond)
	require.NoError(t, err)
	require.JSONEq(t, expected, string(bs))
}

func TestEventSubConditionHelpers(t *testing.T) {
	t.Parallel()

	broadcasterOnly := `{"broadcaster_user_id":"1337"}`
	broadcasterModerator := `{"broadcaster_user_id":"1337","moderator_user_id":"9001"}`
	broadcasterUser := `{"broadcaster_user_id":"1337","user_id":"42"}`
	broadcasterReward := `{"broadcaster_user_id":"1337","reward_id":"reward-id"}`

	for _, tc := range []struct {
		name     string
		cond     api.EventSubCondition
		expected string
	}{
		{"WithAutomodMessageHoldCondition", api.WithAutomodMessageHoldCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithAutomodMessageUpdateCondition", api.WithAutomodMessageUpdateCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithAutomodSettingsUpdateCondition", api.WithAutomodSettingsUpdateCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithAutomodTermsUpdateCondition", api.WithAutomodTermsUpdateCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelAdBreakBeginCondition", api.WithChannelAdBreakBeginCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelBanCondition", api.WithChannelBanCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelBitsUseCondition", api.WithChannelBitsUseCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelChatClearCondition", api.WithChannelChatClearCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatClearUserMessagesCondition", api.WithChannelChatClearUserMessagesCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatMessageCondition", api.WithChannelChatMessageCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatMessageDeleteCondition", api.WithChannelChatMessageDeleteCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatNotificationCondition", api.WithChannelChatNotificationCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatSettingsUpdateCondition", api.WithChannelChatSettingsUpdateCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatUserMessageHoldCondition", api.WithChannelChatUserMessageHoldCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelChatUserMessageUpdateCondition", api.WithChannelChatUserMessageUpdateCondition(testBroadcasterUserID, testUserID), broadcasterUser},
		{"WithChannelSubscribeCondition", api.WithChannelSubscribeCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSubscriptionEndCondition", api.WithChannelSubscriptionEndCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSubscriptionGiftCondition", api.WithChannelSubscriptionGiftCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSubscriptionMessageCondition", api.WithChannelSubscriptionMessageCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelCheerCondition", api.WithChannelCheerCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelUpdateCondition", api.WithChannelUpdateCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelFollowCondition", api.WithChannelFollowCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelUnbanCondition", api.WithChannelUnbanCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelUnbanRequestCreateCondition", api.WithChannelUnbanRequestCreateCondition(testModeratorUserID, testBroadcasterUserID), broadcasterModerator},
		{"WithChannelUnbanRequestResolveCondition", api.WithChannelUnbanRequestResolveCondition(testModeratorUserID, testBroadcasterUserID), broadcasterModerator},
		{"WithChannelRaidFromBroadcasterCondition", api.WithChannelRaidFromBroadcasterCondition(testBroadcasterUserID), `{"from_broadcaster_user_id":"1337"}`},
		{"WithChannelRaidToBroadcasterCondition", api.WithChannelRaidToBroadcasterCondition(testBroadcasterUserID), `{"to_broadcaster_user_id":"1337"}`},
		{"WithChannelModerateCondition", api.WithChannelModerateCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelModerateV2Condition", api.WithChannelModerateV2Condition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelModeratorAddCondition", api.WithChannelModeratorAddCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelModeratorRemoveCondition", api.WithChannelModeratorRemoveCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelGuestStarSessionBeginCondition", api.WithChannelGuestStarSessionBeginCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelGuestStarSessionEndCondition", api.WithChannelGuestStarSessionEndCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelGuestStarGuestUpdateCondition", api.WithChannelGuestStarGuestUpdateCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelGuestStarSettingsUpdateCondition", api.WithChannelGuestStarSettingsUpdateCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelPointsAutomaticRewardRedemptionAddCondition", api.WithChannelPointsAutomaticRewardRedemptionAddCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsAutomaticRewardRedemptionAddV2Condition", api.WithChannelPointsAutomaticRewardRedemptionAddV2Condition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsCustomRewardAddCondition", api.WithChannelPointsCustomRewardAddCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsCustomRewardUpdateCondition", api.WithChannelPointsCustomRewardUpdateCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsCustomRewardUpdateCondition/reward_id", api.WithChannelPointsCustomRewardUpdateCondition(testBroadcasterUserID, testRewardID), broadcasterReward},
		{"WithChannelPointsCustomRewardRemoveCondition", api.WithChannelPointsCustomRewardRemoveCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsCustomRewardRemoveCondition/reward_id", api.WithChannelPointsCustomRewardRemoveCondition(testBroadcasterUserID, testRewardID), broadcasterReward},
		{"WithChannelPointsCustomRewardRedemptionAddCondition", api.WithChannelPointsCustomRewardRedemptionAddCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsCustomRewardRedemptionAddCondition/reward_id", api.WithChannelPointsCustomRewardRedemptionAddCondition(testBroadcasterUserID, testRewardID), broadcasterReward},
		{"WithChannelPointsCustomRewardRedemptionUpdateCondition", api.WithChannelPointsCustomRewardRedemptionUpdateCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPointsCustomRewardRedemptionUpdateCondition/reward_id", api.WithChannelPointsCustomRewardRedemptionUpdateCondition(testBroadcasterUserID, testRewardID), broadcasterReward},
		{"WithChannelCustomPowerUpRedemptionAddCondition", api.WithChannelCustomPowerUpRedemptionAddCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelCustomPowerUpRedemptionAddCondition/reward_id", api.WithChannelCustomPowerUpRedemptionAddCondition(testBroadcasterUserID, testRewardID), broadcasterReward},
		{"WithChannelPollBeginCondition", api.WithChannelPollBeginCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPollProgressCondition", api.WithChannelPollProgressCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPollEndCondition", api.WithChannelPollEndCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPredictionBeginCondition", api.WithChannelPredictionBeginCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPredictionProgressCondition", api.WithChannelPredictionProgressCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPredictionLockCondition", api.WithChannelPredictionLockCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelPredictionEndCondition", api.WithChannelPredictionEndCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSharedChatSessionBeginCondition", api.WithChannelSharedChatSessionBeginCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSharedChatSessionUpdateCondition", api.WithChannelSharedChatSessionUpdateCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSharedChatSessionEndCondition", api.WithChannelSharedChatSessionEndCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelSuspiciousUserMessageCondition", api.WithChannelSuspiciousUserMessageCondition(testModeratorUserID, testBroadcasterUserID), broadcasterModerator},
		{"WithChannelSuspiciousUserUpdateCondition", api.WithChannelSuspiciousUserUpdateCondition(testModeratorUserID, testBroadcasterUserID), broadcasterModerator},
		{"WithChannelVIPAddCondition", api.WithChannelVIPAddCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelVIPRemoveCondition", api.WithChannelVIPRemoveCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithChannelWarningAcknowledgeCondition", api.WithChannelWarningAcknowledgeCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithChannelWarningSendCondition", api.WithChannelWarningSendCondition(testBroadcasterUserID, testModeratorUserID), broadcasterModerator},
		{"WithConduitShardDisabledCondition", api.WithConduitShardDisabledCondition(testClientID), `{"client_id":"client-id"}`},
		{"WithConduitShardDisabledCondition/conduit_id", api.WithConduitShardDisabledCondition(testClientID, testConduitID), `{"client_id":"client-id","conduit_id":"conduit-id"}`},
		{"WithDropEntitlementGrantCondition", api.WithDropEntitlementGrantCondition(testOrganizationID, testCategoryID, testCampaignID), `{"organization_id":"org-id","category_id":"category-id","campaign_id":"campaign-id"}`},
		{"WithExtensionBitsTransactionCreateCondition", api.WithExtensionBitsTransactionCreateCondition(testExtensionClientID), `{"extension_client_id":"extension-client-id"}`},
		{"WithGoalsCondition", api.WithGoalsCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithHypeTrainBeginCondition", api.WithHypeTrainBeginCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithHypeTrainProgressCondition", api.WithHypeTrainProgressCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithHypeTrainEndCondition", api.WithHypeTrainEndCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithStreamOnlineCondition", api.WithStreamOnlineCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithStreamOfflineCondition", api.WithStreamOfflineCondition(testBroadcasterUserID), broadcasterOnly},
		{"WithUserAuthorizationGrantCondition", api.WithUserAuthorizationGrantCondition(testClientID), `{"client_id":"client-id"}`},
		{"WithUserAuthorizationRevokeCondition", api.WithUserAuthorizationRevokeCondition(testClientID), `{"client_id":"client-id"}`},
		{"WithUserUpdateCondition", api.WithUserUpdateCondition(testUserID), `{"user_id":"42"}`},
		{"WithWhisperReceivedCondition", api.WithWhisperReceivedCondition(testUserID), `{"user_id":"42"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertEventSubConditionJSON(t, tc.cond, tc.expected)
		})
	}
}

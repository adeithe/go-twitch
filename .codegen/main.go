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
	BitsResource                     = NewTwitchAPIResource("Bits", BitsLeaderboardResource, BitsCheermotesResource, BitsExtensionsResource)
	BitsLeaderboardResource          = NewTwitchAPIResource("Leaderboard")
	BitsCheermotesResource           = NewTwitchAPIResource("Cheermotes")
	BitsExtensionsResource           = NewTwitchAPIResource("Extensions")
	ChannelsResource                 = NewTwitchAPIResource("Channels", ChannelEditorsResource, ChannelFollowedResource, ChannelFollowersResource)
	ChannelEditorsResource           = NewTwitchAPIResource("Editors")
	ChannelFollowedResource          = NewTwitchAPIResource("Followed")
	ChannelFollowersResource         = NewTwitchAPIResource("Followers")
	ChannelPointsResource            = NewTwitchAPIResource("ChannelPoints", ChannelPointsRewardsResource, ChannelPointsRedemptionsResource)
	ChannelPointsRewardsResource     = NewTwitchAPIResource("Rewards")
	ChannelPointsRedemptionsResource = NewTwitchAPIResource("Redemptions")
	CharityResource                  = NewTwitchAPIResource("Charity", CharityCampaignResource, CharityDonationsResource)
	CharityCampaignResource          = NewTwitchAPIResource("Campaign")
	CharityDonationsResource         = NewTwitchAPIResource("Donations")
	ChatResource                     = NewTwitchAPIResource("Chat", ChatChattersResource)
	ChatChattersResource             = NewTwitchAPIResource("Chatters")

	// Resources is the list of top-level API resources to generate. Subresources are included automatically.
	Resources = []*TwitchAPIResource{
		AdsResource, AnalyticsResource, BitsResource,
	}

	// Endpoints is the list of all API endpoints to generate. Resource mapping is done automatically using the Resource field.
	Endpoints = []*TwitchAPIEndpoint{
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
			Name:     "CharityCampaign",
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
				Data       api.CharityCampaign
				Pagination api.Pagination
			}{},
		},
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
	}
)

func NewTwitchAPIResource(name string, subresources ...*TwitchAPIResource) *TwitchAPIResource {
	return &TwitchAPIResource{Name: name, SubResources: subresources}
}

func main() {
	resourceMapping := make(map[*TwitchAPIResource][]*TwitchAPIEndpoint)
	for _, e := range Endpoints {
		resourceMapping[e.Resource] = append(resourceMapping[e.Resource], e)
	}

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

		count, err := GenerateResource(out, resource, endpoints)
		if err != nil {
			slog.Error("failed to generate resource",
				slog.String("output", tmplOutput),
				slog.String("template", TemplateAPIResource),
				slog.Any("error", err),
			)
		}

		for _, subresource := range resource.SubResources {
			subresource.Name = resource.Name + subresource.Name
			n, err := GenerateResource(out, subresource, resourceMapping[subresource])
			if err != nil {
				slog.Error("failed to generate subresource",
					slog.String("output", tmplOutput),
					slog.String("template", TemplateAPIResource),
					slog.Any("error", err),
				)
			}
			count += n
			resourceCount++
		}
		endpointCount += count
		_ = out.Close()
	}
	slog.Info("generated resources", slog.Int("endpoints", endpointCount))
}

func GenerateResource(w io.Writer, resource *TwitchAPIResource, endpoints []*TwitchAPIEndpoint) (n int, err error) {
	if err = Generate(w, TemplateAPIResource, resource); err != nil {
		return
	}

	for _, endpoint := range endpoints {
		if strings.HasPrefix(endpoint.DocsURL, "#") {
			endpoint.DocsURL = "https://dev.twitch.tv/docs/api/reference/" + endpoint.DocsURL
		}

		if err = Generate(w, TemplateAPICall, endpoint); err != nil {
			return
		}
		n++
	}
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

package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
)

type EndpointTestCase struct {
	name     string
	endpoint func(*apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint
	fetch    func(*api.Client, ...api.RequestOption) (func(t *testing.T), error)
}

func TestAPI_Ads(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Start Commercial",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/channels/commercial", &api.ResponseData[api.Commercial]{
					Data: []api.Commercial{{Length: 60, RetryAfter: 480}},
				}, apitest.RequireBodyParam("broadcaster_id"), apitest.RequireBodyParam("length"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Ads.Insert("1234", 60).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, 60, res.Data[0].Length)
					require.Equal(t, 480, res.Data[0].RetryAfter)
				}, err
			},
		},
		{
			"Get Ad Schedule",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				timestamp := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T23:08:18+00:00"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/channels/ads", &api.ResponseData[api.AdSchedule]{
					Data: []api.AdSchedule{{
						DurationSeconds: 60,
						PrerollFreeTime: 90,
						SnoozeCount:     1,
						SnoozeRefreshAt: timestamp,
						NextAdAt:        timestamp,
						LastAdAt:        timestamp,
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				timestamp := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T23:08:18+00:00"))
				res, err := api.Ads.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, 60, res.Data[0].DurationSeconds)
					require.Equal(t, 90, res.Data[0].PrerollFreeTime)
					require.Equal(t, 1, res.Data[0].SnoozeCount)
					require.Equal(t, timestamp.Unix(), res.Data[0].SnoozeRefreshAt.Unix())
					require.Equal(t, timestamp.Unix(), res.Data[0].NextAdAt.Unix())
					require.Equal(t, timestamp.Unix(), res.Data[0].LastAdAt.Unix())
				}, err
			},
		},
		{
			"Snooze Next Ad",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				timestamp := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T23:08:18+00:00"))
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/channels/ads/schedule/snooze", &api.ResponseData[api.AdsSnoozed]{
					Data: []api.AdsSnoozed{{
						SnoozeCount:     1,
						SnoozeRefreshAt: timestamp,
						NextAdAt:        timestamp,
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				timestamp := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T23:08:18+00:00"))
				res, err := api.Ads.Snooze.Insert("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, 1, res.Data[0].SnoozeCount)
					require.Equal(t, timestamp.Unix(), res.Data[0].SnoozeRefreshAt.Unix())
					require.Equal(t, timestamp.Unix(), res.Data[0].NextAdAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Analytics(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Extension Analytics",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-03-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-06-01T00:00:00Z"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/analytics/extensions", &api.ResponseData[api.ExtensionAnalyticsReport]{
					Data: []api.ExtensionAnalyticsReport{{
						ExtensionID: "ext123",
						Type:        "overview_v2",
						URL:         "https://twitch-piper-reports.s3-us-west-2.amazonaws.com/dynamic/LoL%20ADC...",
						DateRange:   api.DateRange{start, end},
					}},
					Pagination: api.Pagination{
						Cursor: "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6NX19",
					},
				}, apitest.QueryParamEquals("extension_id", "ext123"), apitest.QueryParamEquals("type", "overview_v2"), apitest.QueryParamEquals("started_at", "2018-03-01T00:00:00Z"), apitest.QueryParamEquals("ended_at", "2018-06-01T00:00:00Z"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-03-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-06-01T00:00:00Z"))
				res, err := api.Analytics.Extensions.List().ExtensionID("ext123").Type("overview_v2").StartedAt(start).EndedAt(end).After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "ext123", res.Data[0].ExtensionID)
					require.Equal(t, "overview_v2", res.Data[0].Type)
					require.Equal(t, "https://twitch-piper-reports.s3-us-west-2.amazonaws.com/dynamic/LoL%20ADC...", res.Data[0].URL)
					require.Equal(t, start.Unix(), res.Data[0].DateRange.StartedAt.Unix())
					require.Equal(t, end.Unix(), res.Data[0].DateRange.EndedAt.Unix())
					require.Equal(t, "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6NX19", res.Pagination.Cursor)
				}, err
			},
		},
		{
			"Get Game Analytics",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-03-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-06-01T00:00:00Z"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/analytics/games", &api.ResponseData[api.GameAnalyticsReport]{
					Data: []api.GameAnalyticsReport{{
						GameID:    "game123",
						Type:      "overview_v2",
						URL:       "https://twitch-piper-reports.s3-us-west-2.amazonaws.com/games/66170/overview/15183...",
						DateRange: api.DateRange{start, end},
					}},
				}, apitest.QueryParamEquals("game_id", "game123"), apitest.QueryParamEquals("type", "overview_v2"), apitest.QueryParamEquals("started_at", "2018-03-01T00:00:00Z"), apitest.QueryParamEquals("ended_at", "2018-06-01T00:00:00Z"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-03-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2018-06-01T00:00:00Z"))
				res, err := api.Analytics.Games.List().GameID("game123").Type("overview_v2").StartedAt(start).EndedAt(end).After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "game123", res.Data[0].GameID)
					require.Equal(t, "overview_v2", res.Data[0].Type)
					require.Equal(t, "https://twitch-piper-reports.s3-us-west-2.amazonaws.com/games/66170/overview/15183...", res.Data[0].URL)
					require.Equal(t, start.Unix(), res.Data[0].DateRange.StartedAt.Unix())
					require.Equal(t, end.Unix(), res.Data[0].DateRange.EndedAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Bits(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Bits Leaderboard",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-31T23:59:59Z"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/bits/leaderboard", &api.ResponseData[api.BitsLeaderboardEntry]{
					Total: 1,
					Data: []api.BitsLeaderboardEntry{{
						UserID:   "user123",
						UserName: "CoolViewer",
						Rank:     1,
						Score:    5000,
					}},
					DateRange: api.DateRange{start, end},
				}, apitest.QueryParamEquals("user_id", "1234"), apitest.QueryParamEquals("count", "1"), apitest.QueryParamEquals("period", "month"), apitest.QueryParamEquals("started_at", "2023-08-01T00:00:00Z"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-31T23:59:59Z"))
				res, err := api.Bits.Leaderboard.List().UserID("1234").Count(1).Period("month").StartedAt(start).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, res.Total)
					require.Equal(t, "user123", res.Data[0].UserID)
					require.Equal(t, "CoolViewer", res.Data[0].UserName)
					require.Equal(t, 1, res.Data[0].Rank)
					require.Equal(t, 5000, res.Data[0].Score)
					require.Equal(t, start.Unix(), res.DateRange.StartedAt.Unix())
					require.Equal(t, end.Unix(), res.DateRange.EndedAt.Unix())
				}, err
			},
		},
		{
			"Get Bits Cheermotes",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				updated := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/bits/cheermotes", &api.ResponseData[api.Cheermote]{
					Data: []api.Cheermote{{
						Prefix: "xd",
						Type:   "global_first_party",
						Order:  3,
						Tiers: []api.CheermoteTier{{
							ID:             "1",
							MinBits:        1,
							Color:          "Blue",
							CanCheer:       true,
							ShowInBitsCard: true,
						}},
						IsCharitable: true,
						LastUpdated:  updated,
					}},
				}, apitest.QueryParamEquals("broadcaster_id", "1234"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				updated := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				res, err := api.Bits.Cheermotes.List().BroadcasterID("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "xd", res.Data[0].Prefix)
					require.Equal(t, "global_first_party", res.Data[0].Type)
					require.Equal(t, 3, res.Data[0].Order)
					require.Len(t, res.Data[0].Tiers, 1)
					require.Equal(t, "1", res.Data[0].Tiers[0].ID)
					require.Equal(t, 1, res.Data[0].Tiers[0].MinBits)
					require.Equal(t, "Blue", res.Data[0].Tiers[0].Color)
					require.True(t, res.Data[0].Tiers[0].CanCheer)
					require.True(t, res.Data[0].Tiers[0].ShowInBitsCard)
					require.True(t, res.Data[0].IsCharitable)
					require.Equal(t, updated.Unix(), res.Data[0].LastUpdated.Unix())
				}, err
			},
		},
		{
			"Get Extension Transactions",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-31T23:59:59Z"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/extensions/transactions", &api.ResponseData[api.ExtensionTransaction]{
					Total: 1,
					Data: []api.ExtensionTransaction{{
						ID:               "1234",
						BroadcasterID:    "broad123",
						BroadcasterLogin: "twitchdev",
						BroadcasterName:  "TwitchDev",
						UserID:           "user123",
						UserLogin:        "coolviewer",
						UserName:         "CoolViewer",
						ProductType:      "bits_in_extension",
						Product: api.ExtensionProduct{
							SKU:    "sku123",
							Domain: "extensions.twitch.tv",
							Cost: api.ExtensionProductCost{
								Amount: 100,
								Type:   "bits",
							},
							DisplayName:   "100 Bits Pack",
							InDevelopment: true,
							Broadcast:     false,
							Expiration:    end,
						},
						Timestamp: start,
					}},
					DateRange: api.DateRange{start, end},
				}, apitest.QueryParamEquals("extension_id", "ext123"), apitest.QueryParamEquals("id", "1234"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-31T23:59:59Z"))
				res, err := api.Bits.Extensions.List("ext123").ID("1234").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, "1234", res.Data[0].ID)
					require.Equal(t, "broad123", res.Data[0].BroadcasterID)
					require.Equal(t, "twitchdev", res.Data[0].BroadcasterLogin)
					require.Equal(t, "TwitchDev", res.Data[0].BroadcasterName)
					require.Equal(t, "user123", res.Data[0].UserID)
					require.Equal(t, "coolviewer", res.Data[0].UserLogin)
					require.Equal(t, "CoolViewer", res.Data[0].UserName)
					require.Equal(t, "bits_in_extension", res.Data[0].ProductType)
					require.Equal(t, "sku123", res.Data[0].Product.SKU)
					require.Equal(t, "extensions.twitch.tv", res.Data[0].Product.Domain)
					require.Equal(t, 100, res.Data[0].Product.Cost.Amount)
					require.Equal(t, "bits", res.Data[0].Product.Cost.Type)
					require.Equal(t, "100 Bits Pack", res.Data[0].Product.DisplayName)
					require.True(t, res.Data[0].Product.InDevelopment)
					require.False(t, res.Data[0].Product.Broadcast)
					require.Equal(t, end.Unix(), res.Data[0].Product.Expiration.Unix())
					require.Equal(t, start.Unix(), res.Data[0].Timestamp.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Channels(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Channel Information",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/channels", &api.ResponseData[api.Channel]{
					Data: []api.Channel{{
						GameID:              "game123",
						GameName:            "Awesome Game",
						BroadcasterID:       "1234",
						BroadcasterLogin:    "twitchdev",
						BroadcasterName:     "TwitchDev",
						BroadcasterLanguage: "en",
						Title:               "Playing Awesome Game!",
						Delay:               10,
					}},
				}, apitest.QueryParamEquals("broadcaster_id", "5678"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Channels.List("5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "game123", res.Data[0].GameID)
					require.Equal(t, "Awesome Game", res.Data[0].GameName)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "twitchdev", res.Data[0].BroadcasterLogin)
					require.Equal(t, "TwitchDev", res.Data[0].BroadcasterName)
					require.Equal(t, "en", res.Data[0].BroadcasterLanguage)
					require.Equal(t, "Playing Awesome Game!", res.Data[0].Title)
					require.Equal(t, 10, res.Data[0].Delay)
				}, err
			},
		},
		{
			"Modify Channel Information",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodPatch, "/helix/channels",
					apitest.QueryParamEquals("broadcaster_id", "5678"),
					apitest.RequireBodyParam("title"),
					apitest.RequireBodyParam("game_id"),
					apitest.RequireBodyParam("broadcaster_language"),
					apitest.RequireBodyParam("delay"),
					apitest.RequireBodyParam("is_branded_content"),
					apitest.RequireBodyParam("content_classification_labels"),
					apitest.RequireBodyParam("tags"),
				)
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := client.Channels.Modify("5678").GameID("game456").BroadcasterLanguage("en").Title("This should be fun!").Delay(15).Tags("fun").ContentClassificationLabels(api.ChannelContentClassificationLabel{ID: "mature", IsEnabled: true}).IsBrandedContent(false).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
		{
			"Get Channel Editors",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/channels/editors", &api.ResponseData[api.ChannelEditor]{
					Data: []api.ChannelEditor{{
						UserID:   "1234",
						UserName: "Cool_Editor",
						AddedAt:  start,
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z"))
				res, err := api.Channels.Editors.List("5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, "Cool_Editor", res.Data[0].UserName)
					require.Equal(t, start.Unix(), res.Data[0].AddedAt.Unix())
				}, err
			},
		},
		{
			"Get Followed Channels",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/channels/followed", &api.ResponseData[api.Followed]{
					Total: 1,
					Data: []api.Followed{{
						BroadcasterID:    "1234",
						BroadcasterLogin: "twitchdev",
						BroadcasterName:  "TwitchDev",
						FollowedAt:       Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z")),
					}},
					Pagination: api.Pagination{
						Cursor: "xyz987",
					},
				}, apitest.RequireQueryParam("user_id"), apitest.QueryParamEquals("broadcaster_id", "1234"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Channels.Followed.List("5678").BroadcasterID("1234").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, res.Total)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "twitchdev", res.Data[0].BroadcasterLogin)
					require.Equal(t, "TwitchDev", res.Data[0].BroadcasterName)
					require.Equal(t, "2023-08-01T00:00:00Z", res.Data[0].FollowedAt.Format(time.RFC3339))
				}, err
			},
		},
		{
			"Get Channel Followers",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/channels/followers", &api.ResponseData[api.Follower]{
					Total: 1,
					Data: []api.Follower{{
						UserID:     "1234",
						UserLogin:  "twitchdev",
						UserName:   "TwitchDev",
						FollowedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T00:00:00Z")),
					}},
					Pagination: api.Pagination{
						Cursor: "xyz987",
					},
				}, apitest.QueryParamEquals("user_id", "1234"), apitest.QueryParamEquals("broadcaster_id", "5678"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Channels.Followers.List("5678").UserID("1234").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, res.Total)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, "twitchdev", res.Data[0].UserLogin)
					require.Equal(t, "TwitchDev", res.Data[0].UserName)
					require.Equal(t, "2023-08-01T00:00:00Z", res.Data[0].FollowedAt.Format(time.RFC3339))
				}, err
			},
		},
	})
}

func TestAPI_ChannelPoints(t *testing.T) {}

func TestAPI_Charity(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Charity Campaign",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/charity/campaigns", &api.ResponseData[api.CharityCampaign]{
					Data: []api.CharityCampaign{{
						ID:                 "camp123",
						BroadcasterID:      "1234",
						BroadcasterName:    "TwitchDev",
						CharityName:        "Save the Whales",
						CharityDescription: "Help us save the whales!",
						CharityLogo:        "https://example.com/logo.png",
						TargetAmount:       api.CharityCampaignAmount{Value: 10000, Decimal: 2, Currency: "USD"},
						CurrentAmount:      api.CharityCampaignAmount{Value: 5000, Decimal: 2, Currency: "USD"},
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Charity.Campaign.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "TwitchDev", res.Data[0].BroadcasterName)
					require.Equal(t, "Save the Whales", res.Data[0].CharityName)
					require.Equal(t, "Help us save the whales!", res.Data[0].CharityDescription)
					require.Equal(t, "https://example.com/logo.png", res.Data[0].CharityLogo)
					require.Equal(t, 10000, res.Data[0].TargetAmount.Value)
					require.Equal(t, 2, res.Data[0].TargetAmount.Decimal)
					require.Equal(t, "USD", res.Data[0].TargetAmount.Currency)
					require.Equal(t, 5000, res.Data[0].CurrentAmount.Value)
					require.Equal(t, 2, res.Data[0].CurrentAmount.Decimal)
					require.Equal(t, "USD", res.Data[0].CurrentAmount.Currency)
				}, err
			},
		},
		{
			"Get Charity Campaign Donations",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/charity/donations", &api.ResponseData[api.CharityCampaignDonation]{
					Data: []api.CharityCampaignDonation{{
						ID:         "don123",
						CampaignID: "camp123",
						UserID:     "1234",
						UserLogin:  "generousviewer",
						UserName:   "GenerousViewer",
						Amount:     api.CharityCampaignAmount{Value: 500, Decimal: 2, Currency: "USD"},
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Charity.Donations.List("1234").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "don123", res.Data[0].ID)
					require.Equal(t, "camp123", res.Data[0].CampaignID)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, "generousviewer", res.Data[0].UserLogin)
					require.Equal(t, "GenerousViewer", res.Data[0].UserName)
					require.Equal(t, 500, res.Data[0].Amount.Value)
					require.Equal(t, 2, res.Data[0].Amount.Decimal)
					require.Equal(t, "USD", res.Data[0].Amount.Currency)
				}, err
			},
		},
	})
}

func TestAPI_Chat(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Global Chat Badges",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/badges/global", &api.ResponseData[api.ChatBadge]{
					Data: []api.ChatBadge{{
						SetID: "1",
						Versions: []api.ChatBadgeVersion{{
							ID:          "1",
							ImageSize1X: "https://static-cdn.jtvnw.net/badges/v1/1/1",
							ImageSize2X: "https://static-cdn.jtvnw.net/badges/v1/1/2",
							ImageSize4X: "https://static-cdn.jtvnw.net/badges/v1/1/4",
						}},
					}},
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Badges.Global.List().Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1", res.Data[0].SetID)
					require.Len(t, res.Data[0].Versions, 1)
					require.Equal(t, "1", res.Data[0].Versions[0].ID)
					require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/1/1", res.Data[0].Versions[0].ImageSize1X)
					require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/1/2", res.Data[0].Versions[0].ImageSize2X)
					require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/1/4", res.Data[0].Versions[0].ImageSize4X)
				}, err
			},
		},
		{
			"Get Chat Badges",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/badges", &api.ResponseData[api.ChatBadge]{
					Data: []api.ChatBadge{{
						SetID: "1",
						Versions: []api.ChatBadgeVersion{{
							ID:          "1",
							ImageSize1X: "https://static-cdn.jtvnw.net/badges/v1/1/1",
							ImageSize2X: "https://static-cdn.jtvnw.net/badges/v1/1/2",
							ImageSize4X: "https://static-cdn.jtvnw.net/badges/v1/1/4",
						}},
					}},
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Badges.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1", res.Data[0].SetID)
					require.Len(t, res.Data[0].Versions, 1)
					require.Equal(t, "1", res.Data[0].Versions[0].ID)
					require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/1/1", res.Data[0].Versions[0].ImageSize1X)
					require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/1/2", res.Data[0].Versions[0].ImageSize2X)
					require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/1/4", res.Data[0].Versions[0].ImageSize4X)
				}, err
			},
		},
		{
			"Get Chat User Color",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/color", &api.ResponseData[api.UserChatColor]{
					Data: []api.UserChatColor{{
						UserID: "1234",
						Color:  "#FF0000",
					}},
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.User.Color.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, "#FF0000", res.Data[0].Color)
				}, err
			},
		},
		{
			"Update User Chat Color",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodPut, "/helix/chat/color", apitest.QueryParamEquals("user_id", "1234"), apitest.QueryParamEquals("color", "#FF0000"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.User.Color.Update("1234", "#FF0000").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
		{
			"Get Global Emotes",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/emotes/global", &api.ResponseData[api.Emote]{
					Data: []api.Emote{{
						ID:   "emote123",
						Name: "CoolEmote",
						Images: api.SizedImage{
							Size1x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0",
							Size2x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0",
							Size4x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0",
						},
						Format:     []string{"static", "animated"},
						Scale:      []string{"1.0", "2.0", "3.0"},
						ThemeMode:  []string{"light", "dark"},
						EmoteType:  "subscriptions",
						EmoteSetID: "set123",
					}},
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Emotes.Global.List().Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "emote123", res.Data[0].ID)
					require.Equal(t, "CoolEmote", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0", res.Data[0].Images.Size1x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0", res.Data[0].Images.Size2x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0", res.Data[0].Images.Size4x)
					require.Equal(t, []string{"static", "animated"}, res.Data[0].Format)
					require.Equal(t, []string{"1.0", "2.0", "3.0"}, res.Data[0].Scale)
					require.Equal(t, []string{"light", "dark"}, res.Data[0].ThemeMode)
					require.Equal(t, "subscriptions", res.Data[0].EmoteType)
					require.Equal(t, "set123", res.Data[0].EmoteSetID)
				}, err
			},
		},
		{
			"Get Channel Emotes",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/emotes", &api.ResponseData[api.Emote]{
					Data: []api.Emote{{
						ID:   "emote123",
						Name: "CoolEmote",
						Images: api.SizedImage{
							Size1x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0",
							Size2x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0",
							Size4x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0",
						},
						Format:     []string{"static", "animated"},
						Scale:      []string{"1.0", "2.0", "3.0"},
						ThemeMode:  []string{"light", "dark"},
						EmoteType:  "subscriptions",
						EmoteSetID: "set123",
					}},
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Emotes.Channel.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "emote123", res.Data[0].ID)
					require.Equal(t, "CoolEmote", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0", res.Data[0].Images.Size1x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0", res.Data[0].Images.Size2x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0", res.Data[0].Images.Size4x)
					require.Equal(t, []string{"static", "animated"}, res.Data[0].Format)
					require.Equal(t, []string{"1.0", "2.0", "3.0"}, res.Data[0].Scale)
					require.Equal(t, []string{"light", "dark"}, res.Data[0].ThemeMode)
					require.Equal(t, "subscriptions", res.Data[0].EmoteType)
					require.Equal(t, "set123", res.Data[0].EmoteSetID)
				}, err
			},
		},
		{
			"Get User Emotes",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/emotes/user", &api.ResponseData[api.Emote]{
					Data: []api.Emote{{
						ID:   "emote123",
						Name: "CoolEmote",
						Images: api.SizedImage{
							Size1x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0",
							Size2x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0",
							Size4x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0",
						},
						Format:     []string{"static", "animated"},
						Scale:      []string{"1.0", "2.0", "3.0"},
						ThemeMode:  []string{"light", "dark"},
						EmoteType:  "subscriptions",
						EmoteSetID: "set123",
					}},
				}, apitest.RequireQueryParam("user_id"), apitest.QueryParamEquals("broadcaster_id", "4567"), apitest.QueryParamEquals("after", "abc123"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Emotes.User.List("1234").BroadcasterID("4567").After("abc123").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "emote123", res.Data[0].ID)
					require.Equal(t, "CoolEmote", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0", res.Data[0].Images.Size1x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0", res.Data[0].Images.Size2x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0", res.Data[0].Images.Size4x)
					require.Equal(t, []string{"static", "animated"}, res.Data[0].Format)
					require.Equal(t, []string{"1.0", "2.0", "3.0"}, res.Data[0].Scale)
					require.Equal(t, []string{"light", "dark"}, res.Data[0].ThemeMode)
					require.Equal(t, "subscriptions", res.Data[0].EmoteType)
					require.Equal(t, "set123", res.Data[0].EmoteSetID)
				}, err
			},
		},
		{
			"Get Emote Sets",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/emotes/set", &api.ResponseData[api.Emote]{
					Data: []api.Emote{{
						ID:   "emote123",
						Name: "CoolEmote",
						Images: api.SizedImage{
							Size1x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0",
							Size2x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0",
							Size4x: "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0",
						},
						Format:     []string{"static", "animated"},
						Scale:      []string{"1.0", "2.0", "3.0"},
						ThemeMode:  []string{"light", "dark"},
						EmoteType:  "subscriptions",
						EmoteSetID: "set123",
					}},
				}, apitest.QueryParamEquals("emote_set_id", "1234"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.EmoteSets.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "emote123", res.Data[0].ID)
					require.Equal(t, "CoolEmote", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/1.0", res.Data[0].Images.Size1x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/2.0", res.Data[0].Images.Size2x)
					require.Equal(t, "https://static-cdn.jtvnw.net/emoticons/v1/emote123/3.0", res.Data[0].Images.Size4x)
					require.Equal(t, []string{"static", "animated"}, res.Data[0].Format)
					require.Equal(t, []string{"1.0", "2.0", "3.0"}, res.Data[0].Scale)
					require.Equal(t, []string{"light", "dark"}, res.Data[0].ThemeMode)
					require.Equal(t, "subscriptions", res.Data[0].EmoteType)
					require.Equal(t, "set123", res.Data[0].EmoteSetID)
				}, err
			},
		},
		{
			"Get Chatters",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/chatters", &api.ResponseData[api.UserInfo]{
					Data: []api.UserInfo{{
						UserID:    "1234",
						UserLogin: "coolviewer",
						UserName:  "CoolViewer",
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, "coolviewer", res.Data[0].UserLogin)
					require.Equal(t, "CoolViewer", res.Data[0].UserName)
				}, err
			},
		},
		{
			"Get Chat Settings",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/chat/settings", &api.ResponseData[api.ChatSettings]{
					Data: []api.ChatSettings{{
						BroadcasterID:                "1234",
						SlowMode:                     true,
						SlowModeWaitTimeSeconds:      10,
						FollowMode:                   true,
						FollowModeDurationMinutes:    5,
						EmoteMode:                    false,
						NonModeratorChatDelay:        true,
						NonModeratorChatDelaySeconds: 15,
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Settings.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.True(t, res.Data[0].SlowMode)
					require.Equal(t, 10, res.Data[0].SlowModeWaitTimeSeconds)
					require.True(t, res.Data[0].FollowMode)
					require.Equal(t, 5, res.Data[0].FollowModeDurationMinutes)
					require.False(t, res.Data[0].EmoteMode)
					require.True(t, res.Data[0].NonModeratorChatDelay)
					require.Equal(t, 15, res.Data[0].NonModeratorChatDelaySeconds)
				}, err
			},
		},
		{
			"Update Chat Settings",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodPatch, "/helix/chat/settings",
					apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"),
					apitest.RequireBodyParam("slow_mode"), apitest.RequireBodyParam("slow_mode_wait_time"),
					apitest.RequireBodyParam("non_moderator_chat_delay"), apitest.RequireBodyParam("non_moderator_chat_delay_duration"),
					apitest.RequireBodyParam("follower_mode"), apitest.RequireBodyParam("follower_mode_duration"),
					apitest.RequireBodyParam("emote_mode"), apitest.RequireBodyParam("unique_chat_mode"),
				)
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Settings.Modify("1234", "5678").
					EmoteMode(true).UniqueChatMode(true).SubscriberMode(true).
					NonModeratorChatDelay(true).NonModeratorChatDelayDuration(15).
					FollowerMode(true).FollowerModeDuration(5).
					SlowMode(true).SlowModeWaitTime(10).
					Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
		{
			"Get Shared Chat Session",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/shared_chat/session", &api.ResponseData[api.SharedChatSession]{
					Data: []api.SharedChatSession{{
						SessionID:         "session123",
						HostBroadcasterID: "1234",
						Participants: []api.SharedChatSessionParticipant{{
							BroadcasterID: "4567",
						}},
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
						UpdatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.Chat.Shared.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "session123", res.Data[0].SessionID)
					require.Equal(t, "1234", res.Data[0].HostBroadcasterID)
					require.Len(t, res.Data[0].Participants, 1)
					require.Equal(t, "4567", res.Data[0].Participants[0].BroadcasterID)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
		{
			"Send Chat Announcement",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodPost, "/helix/chat/announcements", apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"), apitest.RequireBodyParam("message"), apitest.RequireBodyParam("color"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Announcement.Insert("1234", "5678", "9876").Color("purple").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
		{
			"Send A Shoutout",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodPost, "/helix/chat/shoutouts",
					apitest.RequireQueryParam("from_broadcaster_id"), apitest.RequireQueryParam("to_broadcaster_id"),
					apitest.RequireQueryParam("moderator_id"),
				)
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Shoutout.Insert("1234", "5678", "9876").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
		{
			"Send Chat Message",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodPost, "/helix/chat/messages",
					apitest.RequireBodyParam("broadcaster_id"), apitest.RequireBodyParam("sender_id"),
					apitest.RequireBodyParam("message"), apitest.RequireBodyParam("replay_parent_message_id"),
					apitest.RequireBodyParam("for_source_only"),
				)
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Insert("1234", "5678", "Hello, chat!").ReplayParentMessageID("abc123").ForSourceOnly(true).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
	})
}

func TestAPI_Clips(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Create Clip",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/clips", &api.ResponseData[api.EditableClip]{
					Data: []api.EditableClip{{
						ID:      "FiveWordsForClipSlug",
						EditURL: "http://clips.twitch.tv/FiveWordsForClipSlug/edit",
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Clips.Insert("1234").HasDelay(false).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "FiveWordsForClipSlug", res.Data[0].ID)
					require.Equal(t, "http://clips.twitch.tv/FiveWordsForClipSlug/edit", res.Data[0].EditURL)
				}, err
			},
		},
		{
			"Get Clips",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/clips", &api.ResponseData[api.Clip]{
					Data: []api.Clip{{
						ID:              "FiveWordsForClipSlug",
						URL:             "http://clips.twitch.tv/FiveWordsForClipSlug",
						Title:           "Amazing Clip Title",
						EmbedURL:        "http://clips.twitch.tv/embed?clip=FiveWordsForClipSlug",
						BroadcasterID:   "1234",
						BroadcasterName: "CoolStreamer",
						CreatorID:       "5678",
						CreatorName:     "AwesomeViewer",
						VideoID:         "abcd1234efgh",
						GameID:          "9876",
						Language:        "en",
						ViewCount:       500,
						Duration:        60,
						VODOffset:       480,
						Featured:        false,
						ThumbnailURL:    "http://clips-media-assets2.twitch.tv/12345-preview-480x272.jpg",
						CreatedAt:       Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:34:56Z")),
					}},
				}, apitest.QueryParamEquals("id", "FiveWordsForClipSlug"), apitest.QueryParamEquals("broadcaster_id", "1234"), apitest.QueryParamEquals("game_id", "9876"), apitest.QueryParamEquals("after", "cba321"), apitest.QueryParamEquals("first", "1"), apitest.QueryParamEquals("started_at", "2023-08-01T12:00:00Z"), apitest.QueryParamEquals("ended_at", "2023-08-02T12:00:00Z"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				startedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				endedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-02T12:00:00Z"))
				res, err := api.Clips.List().ID("FiveWordsForClipSlug").BroadcasterID("1234").GameID("9876").IsFeatured(false).Before("abc123").After("cba321").First(1).StartedAt(startedAt).EndedAt(endedAt).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "FiveWordsForClipSlug", res.Data[0].ID)
					require.Equal(t, "http://clips.twitch.tv/FiveWordsForClipSlug", res.Data[0].URL)
					require.Equal(t, "Amazing Clip Title", res.Data[0].Title)
					require.Equal(t, "http://clips.twitch.tv/embed?clip=FiveWordsForClipSlug", res.Data[0].EmbedURL)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "CoolStreamer", res.Data[0].BroadcasterName)
					require.Equal(t, "5678", res.Data[0].CreatorID)
					require.Equal(t, "AwesomeViewer", res.Data[0].CreatorName)
					require.Equal(t, "abcd1234efgh", res.Data[0].VideoID)
					require.Equal(t, "9876", res.Data[0].GameID)
					require.Equal(t, "en", res.Data[0].Language)
					require.Equal(t, 500, res.Data[0].ViewCount)
					require.InDelta(t, 60, res.Data[0].Duration, 0.00)
					require.Equal(t, 480, res.Data[0].VODOffset)
					require.False(t, res.Data[0].Featured)
					require.Equal(t, "http://clips-media-assets2.twitch.tv/12345-preview-480x272.jpg", res.Data[0].ThumbnailURL)
					createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:34:56Z"))
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
		{
			"Get Clips Download",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/clips/downloads", &api.ResponseData[api.DownloadableClip]{
					Data: []api.DownloadableClip{
						{
							ClipID:               "InexpensiveDistinctFoxChefFrank",
							LandscapeDownloadURL: "https://production.assets.clips.twitchcdn.net/yFZG...",
						},
						{
							ClipID:               "SpinelessCloudyLeopardMcaT",
							LandscapeDownloadURL: "https://production.assets.clips.twitchcdn.net/542j...",
						},
					},
				}, apitest.RequireQueryParam("clip_id"), apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("editor_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Clips.Download.List().ClipID("InexpensiveDistinctFoxChefFrank").ClipID("SpinelessCloudyLeopardMcaT").BroadcasterID("141981764").EditorID("141981764").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 2)
					require.Equal(t, "InexpensiveDistinctFoxChefFrank", res.Data[0].ClipID)
					require.Equal(t, "https://production.assets.clips.twitchcdn.net/yFZG...", res.Data[0].LandscapeDownloadURL)
					require.Equal(t, "SpinelessCloudyLeopardMcaT", res.Data[1].ClipID)
					require.Equal(t, "https://production.assets.clips.twitchcdn.net/542j...", res.Data[1].LandscapeDownloadURL)
				}, err
			},
		},
	})
}

func TestAPI_Conduits(t *testing.T) {}

func TestAPI_ContentLabels(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Content Classification Labels",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/content_classification_labels", &api.ResponseData[api.ContentClassificationLabel]{
					Data: []api.ContentClassificationLabel{{
						ID:          "label123",
						Name:        "Violence",
						Description: "Content that depicts violence",
					}},
				}, apitest.QueryParamEquals("locale", "en-US"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.ContentLabels.List().Locale("en-US").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "label123", res.Data[0].ID)
					require.Equal(t, "Violence", res.Data[0].Name)
					require.Equal(t, "Content that depicts violence", res.Data[0].Description)
				}, err
			},
		},
	})
}

func TestAPI_Entitlements(t *testing.T) {}

func TestAPI_Extensions(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Extensions",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/extensions", &api.ResponseData[api.Extension]{
					Data: []api.Extension{{
						ID:            "ext123",
						Name:          "Cool Extension",
						Description:   "An awesome Twitch extension",
						AuthorName:    "twitchdev",
						State:         "Release",
						Version:       "0.0.9",
						ViewerSummary: "Test ALL the extensions features!",
					}},
				}, apitest.QueryParamEquals("extension_id", "ext123"), apitest.QueryParamEquals("extension_version", "0.0.9"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Extensions.List("ext123").ExtensionVersion("0.0.9").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "ext123", res.Data[0].ID)
					require.Equal(t, "Cool Extension", res.Data[0].Name)
					require.Equal(t, "An awesome Twitch extension", res.Data[0].Description)
					require.Equal(t, "twitchdev", res.Data[0].AuthorName)
					require.Equal(t, "Release", res.Data[0].State)
					require.Equal(t, "0.0.9", res.Data[0].Version)
					require.Equal(t, "Test ALL the extensions features!", res.Data[0].ViewerSummary)
				}, err
			},
		},
	})
}

func TestAPI_EventSub(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get EventSub Subscriptions",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/eventsub/subscriptions", &api.ResponseData[api.EventSubSubscription]{
					Data: []api.EventSubSubscription{{
						ID:        "sub123",
						Status:    "enabled",
						Type:      "channel.follow",
						Version:   "1",
						Condition: map[string]any{"broadcaster_user_id": "1234"},
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.QueryParamEquals("subscription_id", "sub123"), apitest.QueryParamEquals("status", "enabled"), apitest.QueryParamEquals("type", "channel.follow"), apitest.QueryParamEquals("user_id", "1234"), apitest.QueryParamEquals("after", "abc123"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.EventSub.List().SubscriptionID("sub123").Status("enabled").Type("channel.follow").UserID("1234").After("abc123").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "sub123", res.Data[0].ID)
					require.Equal(t, "enabled", res.Data[0].Status)
					require.Equal(t, "channel.follow", res.Data[0].Type)
					require.Equal(t, "1", res.Data[0].Version)
					require.Equal(t, "1234", res.Data[0].Condition["broadcaster_user_id"])
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Games(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Top Games",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/games/top", &api.ResponseData[api.Game]{
					Data: []api.Game{{
						ID:        "21779",
						Name:      "Fortnite",
						BoxArtURL: "https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg",
						IGDB:      "124024",
					}},
				}, apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("before", "def456"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Games.Top.List().After("abc123").Before("def456").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "21779", res.Data[0].ID)
					require.Equal(t, "Fortnite", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg", res.Data[0].BoxArtURL)
					require.Equal(t, "124024", res.Data[0].IGDB)
				}, err
			},
		},
		{
			"Get Games",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/games", &api.ResponseData[api.Game]{
					Data: []api.Game{{
						ID:        "21779",
						Name:      "Fortnite",
						BoxArtURL: "https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg",
						IGDB:      "124024",
					}},
				}, apitest.QueryParamEquals("id", "21779"), apitest.QueryParamEquals("name", "Fortnite"), apitest.QueryParamEquals("igdb_id", "124024"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Games.List().ID("21779").Name("Fortnite").IGDB("124024").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "21779", res.Data[0].ID)
					require.Equal(t, "Fortnite", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg", res.Data[0].BoxArtURL)
					require.Equal(t, "124024", res.Data[0].IGDB)
				}, err
			},
		},
	})
}

func TestAPI_Goals(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Creator Goals",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/goals", &api.ResponseData[api.CreatorGoal]{
					Data: []api.CreatorGoal{{
						Type:             "follower",
						ID:               "goal123",
						BroadcasterID:    "1234",
						BroadcasterLogin: "twitchdev",
						BroadcasterName:  "TwitchDev",
						Description:      "Reach 100 followers!",
						CurrentAmount:    75,
						TargetAmount:     100,
						CreatedAt:        Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.Goals.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "goal123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "TwitchDev", res.Data[0].BroadcasterName)
					require.Equal(t, "Reach 100 followers!", res.Data[0].Description)
					require.Equal(t, 75, res.Data[0].CurrentAmount)
					require.Equal(t, 100, res.Data[0].TargetAmount)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_GuestStar(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Guest Star Sessions",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/guest_star/session", &api.ResponseData[api.GuestStarSession]{
					Data: []api.GuestStarSession{{
						ID: "session123",
						Guests: []api.GuestStarGuest{{
							SlotID:    "slot1",
							IsLive:    false,
							UserID:    "9876",
							UserLogin: "coolviewer",
							UserName:  "CoolViewer",
							Volume:    100,
							AudioSettings: api.MediaSettings{
								IsHostEnabled:  true,
								IsGuestEnabled: true,
								IsAvailable:    true,
							},
							VideoSettings: api.MediaSettings{
								IsHostEnabled:  false,
								IsGuestEnabled: false,
								IsAvailable:    false,
							},
							AssignedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T11:50:00Z")),
							JoinedAt:   Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
						}},
					}},
				}, apitest.QueryParamEquals("broadcaster_id", "1234"), apitest.QueryParamEquals("moderator_id", "4567"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.GuestStar.List("1234", "4567").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "session123", res.Data[0].ID)
					require.Len(t, res.Data[0].Guests, 1)
					require.Equal(t, "slot1", res.Data[0].Guests[0].SlotID)
					require.False(t, res.Data[0].Guests[0].IsLive)
					require.Equal(t, "9876", res.Data[0].Guests[0].UserID)
					require.Equal(t, "coolviewer", res.Data[0].Guests[0].UserLogin)
					require.Equal(t, "CoolViewer", res.Data[0].Guests[0].UserName)
					require.Equal(t, 100, res.Data[0].Guests[0].Volume)
					require.True(t, res.Data[0].Guests[0].AudioSettings.IsHostEnabled)
					require.True(t, res.Data[0].Guests[0].AudioSettings.IsGuestEnabled)
					require.True(t, res.Data[0].Guests[0].AudioSettings.IsAvailable)
					require.False(t, res.Data[0].Guests[0].VideoSettings.IsHostEnabled)
					require.False(t, res.Data[0].Guests[0].VideoSettings.IsGuestEnabled)
					require.False(t, res.Data[0].Guests[0].VideoSettings.IsAvailable)
					assignedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T11:50:00Z"))
					joinedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
					require.Equal(t, assignedAt.Unix(), res.Data[0].Guests[0].AssignedAt.Unix())
					require.Equal(t, joinedAt.Unix(), res.Data[0].Guests[0].JoinedAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_HypeTrain(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Hype Train Status",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/hypetrain/status", &api.ResponseData[api.HypeTrainStatusInfo]{
					Data: []api.HypeTrainStatusInfo{{
						Current: &api.HypeTrainStatus{
							ID:               "train123",
							BroadcasterID:    "1234",
							BroadcasterLogin: "twitchdev",
							BroadcasterName:  "TwitchDev",
							Level:            3,
							Total:            1500,
							Progress:         4,
							Goal:             1600,
							StartedAt:        Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
							ExpiresAt:        Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:10:00Z")),
						},
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.HypeTrain.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.NotNil(t, res.Data[0].Current)
					require.Equal(t, "train123", res.Data[0].Current.ID)
					require.Equal(t, "1234", res.Data[0].Current.BroadcasterID)
					require.Equal(t, "twitchdev", res.Data[0].Current.BroadcasterLogin)
					require.Equal(t, "TwitchDev", res.Data[0].Current.BroadcasterName)
					require.Equal(t, 3, res.Data[0].Current.Level)
					require.Equal(t, 1500, res.Data[0].Current.Total)
					require.Equal(t, 4, res.Data[0].Current.Progress)
					require.Equal(t, 1600, res.Data[0].Current.Goal)
					startedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
					expiresAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:10:00Z"))
					require.Equal(t, startedAt.Unix(), res.Data[0].Current.StartedAt.Unix())
					require.Equal(t, expiresAt.Unix(), res.Data[0].Current.ExpiresAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Moderation(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Ban Users",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				expiresAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:01Z"))
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/moderation/bans", &api.ResponseData[api.IssuedBan]{
					Data: []api.IssuedBan{{
						BroadcasterID: "9876",
						ModeratorID:   "5678",
						UserID:        "1234",
						CreatedAt:     createdAt,
						ExpiresAt:     &expiresAt,
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"), apitest.RequireBodyParam("data"))
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				expiresAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:01Z"))
				res, err := client.Moderation.Bans.Insert("1234", "5678", api.WithTimeout("1234", "", 1)).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "9876", res.Data[0].BroadcasterID)
					require.Equal(t, "5678", res.Data[0].ModeratorID)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
					require.Equal(t, expiresAt.Unix(), res.Data[0].ExpiresAt.Unix())
				}, err
			},
		},
		{
			"Unban User",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodDelete, "/helix/moderation/bans",
					apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"),
					apitest.QueryParamEquals("user_id", "9876"),
				)
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := client.Moderation.Bans.Delete("1234", "5678", "9876").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
		{
			"Delete Chat Messages",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodDelete, "/helix/moderation/chat",
					apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"),
					apitest.QueryParamEquals("message_id", "chat123"),
				)
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := client.Moderation.ClearChat.Delete("1234", "5678").MessageID("chat123").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
	})
}

func TestAPI_Polls(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Polls",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/polls", &api.ResponseData[api.Poll]{
					Data: []api.Poll{{
						ID:            "poll123",
						BroadcasterID: "1234",
						Title:         "Which game should I play?",
						Status:        "ACTIVE",
						Choices: []api.PollChoice{{
							ID:    "choice1",
							Title: "Game A",
							Votes: 10,
						}, {
							ID:    "choice2",
							Title: "Game B",
							Votes: 15,
						}},
						StartedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.QueryParamEquals("id", "poll123"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.Polls.List("1234").ID("poll123").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "poll123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "Which game should I play?", res.Data[0].Title)
					require.Equal(t, "ACTIVE", res.Data[0].Status)
					require.Len(t, res.Data[0].Choices, 2)
					require.Equal(t, "choice1", res.Data[0].Choices[0].ID)
					require.Equal(t, "Game A", res.Data[0].Choices[0].Title)
					require.Equal(t, 10, res.Data[0].Choices[0].Votes)
					require.Equal(t, "choice2", res.Data[0].Choices[1].ID)
					require.Equal(t, "Game B", res.Data[0].Choices[1].Title)
					require.Equal(t, 15, res.Data[0].Choices[1].Votes)
					require.Equal(t, createdAt.Unix(), res.Data[0].StartedAt.Unix())
				}, err
			},
		},
		{
			"Create Poll",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/polls", &api.ResponseData[api.Poll]{
					Data: []api.Poll{{
						ID:            "poll123",
						BroadcasterID: "1234",
						Title:         "Which game should I play?",
						Status:        "ACTIVE",
						Choices: []api.PollChoice{
							{ID: "choice1", Title: "Game A"},
							{ID: "choice2", Title: "Game B"},
						},
						StartedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireBodyParam("broadcaster_id"), apitest.RequireBodyParam("title"), apitest.RequireBodyParam("choices"), apitest.RequireBodyParam("duration"), apitest.RequireBodyParam("channel_points_voting_enabled"), apitest.RequireBodyParam("channel_points_per_vote"))
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := client.Polls.Insert("1234", "Which game should I play?", 90).Choices(api.WithChoice("Game A"), api.WithChoice("Game B")).ChannelPointsVotingEnabled(true).ChannelPointsPerVote(500).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "poll123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "Which game should I play?", res.Data[0].Title)
					require.Equal(t, "ACTIVE", res.Data[0].Status)
					require.Len(t, res.Data[0].Choices, 2)
					require.Equal(t, "choice1", res.Data[0].Choices[0].ID)
					require.Equal(t, "Game A", res.Data[0].Choices[0].Title)
					require.Equal(t, "choice2", res.Data[0].Choices[1].ID)
					require.Equal(t, "Game B", res.Data[0].Choices[1].Title)
					require.Equal(t, createdAt.Unix(), res.Data[0].StartedAt.Unix())
				}, err
			},
		},
		{
			"End Poll",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPatch, "/helix/polls", &api.ResponseData[api.Poll]{
					Data: []api.Poll{{
						ID:            "poll123",
						BroadcasterID: "1234",
						Title:         "Which game should I play?",
						Status:        "COMPLETED",
						Choices: []api.PollChoice{
							{ID: "choice1", Title: "Game A", Votes: 10},
							{ID: "choice2", Title: "Game B", Votes: 15},
						},
						StartedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireBodyParam("broadcaster_id"), apitest.RequireBodyParam("id"), apitest.RequireBodyParam("status"))
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := client.Polls.Modify("1234", "poll123", "TERMINATED").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "poll123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "Which game should I play?", res.Data[0].Title)
					require.Equal(t, "COMPLETED", res.Data[0].Status)
					require.Len(t, res.Data[0].Choices, 2)
					require.Equal(t, "choice1", res.Data[0].Choices[0].ID)
					require.Equal(t, "Game A", res.Data[0].Choices[0].Title)
					require.Equal(t, 10, res.Data[0].Choices[0].Votes)
					require.Equal(t, "choice2", res.Data[0].Choices[1].ID)
					require.Equal(t, "Game B", res.Data[0].Choices[1].Title)
					require.Equal(t, 15, res.Data[0].Choices[1].Votes)
					require.Equal(t, createdAt.Unix(), res.Data[0].StartedAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Predictions(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Predictions",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/predictions", &api.ResponseData[api.Prediction]{
					Data: []api.Prediction{{
						ID:            "pred123",
						BroadcasterID: "1234",
						Title:         "Who will win?",
						Status:        "ACTIVE",
						Outcomes: []api.PredictionOutcome{{
							ID:            "outcome1",
							Title:         "Team A",
							Users:         10,
							ChannelPoints: 5000,
						}, {
							ID:            "outcome2",
							Title:         "Team B",
							Users:         15,
							ChannelPoints: 7500,
						}},
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.QueryParamEquals("id", "pred123"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.Predictions.List("1234").ID("pred123").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "pred123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "Who will win?", res.Data[0].Title)
					require.Equal(t, "ACTIVE", res.Data[0].Status)
					require.Len(t, res.Data[0].Outcomes, 2)
					require.Equal(t, "outcome1", res.Data[0].Outcomes[0].ID)
					require.Equal(t, "Team A", res.Data[0].Outcomes[0].Title)
					require.Equal(t, 10, res.Data[0].Outcomes[0].Users)
					require.Equal(t, 5000, res.Data[0].Outcomes[0].ChannelPoints)
					require.Equal(t, "outcome2", res.Data[0].Outcomes[1].ID)
					require.Equal(t, "Team B", res.Data[0].Outcomes[1].Title)
					require.Equal(t, 15, res.Data[0].Outcomes[1].Users)
					require.Equal(t, 7500, res.Data[0].Outcomes[1].ChannelPoints)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
		{
			"Create Prediction",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/predictions", &api.ResponseData[api.Prediction]{
					Data: []api.Prediction{{
						ID:            "pred123",
						BroadcasterID: "1234",
						Title:         "Who will win?",
						Status:        "ACTIVE",
						Outcomes: []api.PredictionOutcome{
							{ID: "outcome1", Title: "Team A"},
							{ID: "outcome2", Title: "Team B"},
						},
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireBodyParam("broadcaster_id"), apitest.RequireBodyParam("title"), apitest.RequireBodyParam("outcomes"), apitest.RequireBodyParam("prediction_window"))
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := client.Predictions.Insert("1234", "Who will win?", 60).Outcome(api.WithChoice("Team A"), api.WithChoice("Team B")).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "pred123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "Who will win?", res.Data[0].Title)
					require.Equal(t, "ACTIVE", res.Data[0].Status)
					require.Len(t, res.Data[0].Outcomes, 2)
					require.Equal(t, "outcome1", res.Data[0].Outcomes[0].ID)
					require.Equal(t, "Team A", res.Data[0].Outcomes[0].Title)
					require.Equal(t, "outcome2", res.Data[0].Outcomes[1].ID)
					require.Equal(t, "Team B", res.Data[0].Outcomes[1].Title)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
		{
			"End Prediction",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPatch, "/helix/predictions", &api.ResponseData[api.Prediction]{
					Data: []api.Prediction{{
						ID:            "pred123",
						BroadcasterID: "1234",
						Title:         "Who will win?",
						Status:        "RESOLVED",
						Outcomes: []api.PredictionOutcome{
							{ID: "outcome1", Title: "Team A"},
							{ID: "outcome2", Title: "Team B"},
						},
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireBodyParam("broadcaster_id"), apitest.RequireBodyParam("id"), apitest.RequireBodyParam("status"), apitest.RequireBodyParam("winning_outcome_id"))
			},
			func(client *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := client.Predictions.Modify("1234", "pred123", "LOCKED").WinningOutcomeID("outcome1").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "pred123", res.Data[0].ID)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "Who will win?", res.Data[0].Title)
					require.Equal(t, "RESOLVED", res.Data[0].Status)
					require.Len(t, res.Data[0].Outcomes, 2)
					require.Equal(t, "outcome1", res.Data[0].Outcomes[0].ID)
					require.Equal(t, "Team A", res.Data[0].Outcomes[0].Title)
					require.Equal(t, "outcome2", res.Data[0].Outcomes[1].ID)
					require.Equal(t, "Team B", res.Data[0].Outcomes[1].Title)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Raids(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Start A Raid",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/raids", &api.ResponseData[api.InitializedRaid]{
					Data: []api.InitializedRaid{{
						Mature:    true,
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
					}},
				}, apitest.RequireBodyParam("from_broadcaster_id"), apitest.RequireBodyParam("to_broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.Raids.Insert("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.True(t, res.Data[0].Mature)
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
				}, err
			},
		},
		{
			"Cancel A Raid",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockValidator(mock, http.MethodDelete, "/helix/raids", apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Raids.Delete("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, http.StatusOK, res.StatusCode)
				}, err
			},
		},
	})
}

func TestAPI_Schedule(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Channel Stream Schedule",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/schedule", &api.ResponseData[api.StreamSchedule]{
					Data: []api.StreamSchedule{{
						BroadcasterID:   "1234",
						BroadcasterName: "CoolStreamer",
						Segments: []api.StreamScheduleSegment{{
							ID:       "segment123",
							Title:    "Morning Stream",
							StartsAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T10:00:00Z")),
							EndsAt:   Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z")),
						}},
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.QueryParamEquals("id", "4567"), apitest.QueryParamEquals("start_time", "2023-08-01T10:00:00Z"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				start := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T10:00:00Z"))
				end := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"))
				res, err := api.Schedule.List("1234").ID("4567").StartTime(start).After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "CoolStreamer", res.Data[0].BroadcasterName)
					require.Len(t, res.Data[0].Segments, 1)
					require.Equal(t, "segment123", res.Data[0].Segments[0].ID)
					require.Equal(t, "Morning Stream", res.Data[0].Segments[0].Title)
					require.Equal(t, start.Unix(), res.Data[0].Segments[0].StartsAt.Unix())
					require.Equal(t, end.Unix(), res.Data[0].Segments[0].EndsAt.Unix())
				}, err
			},
		},
	})
}

func TestAPI_Search(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Search Categories",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/search/categories", &api.ResponseData[api.CategorySearchResult]{
					Data: []api.CategorySearchResult{{
						ID:        "21779",
						Name:      "League of Legends",
						BoxArtURL: "https://static-cdn.jtvnw.net/ttv-boxart/League%20of%20Legends-{width}x{height}.jpg",
					}},
				}, apitest.RequireQueryParam("query"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Search.Categories.List("League").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "21779", res.Data[0].ID)
					require.Equal(t, "League of Legends", res.Data[0].Name)
					require.Equal(t, "https://static-cdn.jtvnw.net/ttv-boxart/League%20of%20Legends-{width}x{height}.jpg", res.Data[0].BoxArtURL)
				}, err
			},
		},
		{
			"Search Channels",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/search/channels", &api.ResponseData[api.ChannelSearchResult]{
					Data: []api.ChannelSearchResult{{
						BroadcasterID:       "141981764",
						BroadcasterLogin:    "twitchdev",
						BroadcasterName:     "TwitchDev",
						BroadcasterLanguage: "en",
						GameID:              "509670",
						GameName:            "Science & Technology",
						Title:               "Welcome to Twitch Dev",
						Live:                true,
						ThumbnailURL:        "https://static-cdn.jtvnw.net/jtv_user_pictures/twitchdev-profile_image-...",
					}},
				}, apitest.RequireQueryParam("query"), apitest.QueryParamEquals("live_only", "true"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Search.Channels.List("TwitchDev").LiveOnly(true).After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "141981764", res.Data[0].BroadcasterID)
					require.Equal(t, "twitchdev", res.Data[0].BroadcasterLogin)
					require.Equal(t, "TwitchDev", res.Data[0].BroadcasterName)
					require.Equal(t, "en", res.Data[0].BroadcasterLanguage)
					require.Equal(t, "509670", res.Data[0].GameID)
					require.Equal(t, "Science & Technology", res.Data[0].GameName)
					require.Equal(t, "Welcome to Twitch Dev", res.Data[0].Title)
					require.True(t, res.Data[0].Live)
					require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/twitchdev-profile_image-...", res.Data[0].ThumbnailURL)
				}, err
			},
		},
	})
}

func TestAPI_Streams(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Stream Key",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/streams/key", &api.ResponseData[api.StreamKey]{
					Data: []api.StreamKey{{
						Key: "abcd-efgh-ijkl-mnop",
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Streams.StreamKey.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "abcd-efgh-ijkl-mnop", res.Data[0].Key)
				}, err
			},
		},
		{
			"Get Streams",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/streams", &api.ResponseData[api.Stream]{
					Data: []api.Stream{{
						ID:          "123456789",
						UserID:      "987654321",
						UserLogin:   "cool_streamer",
						UserName:    "Cool_Streamer",
						GameID:      "21779",
						GameName:    "League of Legends",
						Type:        "live",
						Title:       "Playing some League!",
						ViewerCount: 1500,
						StartedAt:   Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T20:00:00Z")),
					}},
				}, apitest.RequireQueryParam("user_id"), apitest.RequireQueryParam("game_id"), apitest.QueryParamEquals("language", "en"), apitest.QueryParamEquals("period", "all"), apitest.QueryParamEquals("type", "live"), apitest.QueryParamEquals("before", "cba321"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Streams.List().UserID("987654321").GameID("1234").Language("en").Period("all").Type("live").Before("cba321").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "123456789", res.Data[0].ID)
					require.Equal(t, "987654321", res.Data[0].UserID)
					require.Equal(t, "cool_streamer", res.Data[0].UserLogin)
					require.Equal(t, "Cool_Streamer", res.Data[0].UserName)
					require.Equal(t, "21779", res.Data[0].GameID)
					require.Equal(t, "League of Legends", res.Data[0].GameName)
					require.Equal(t, "live", res.Data[0].Type)
					require.Equal(t, "Playing some League!", res.Data[0].Title)
					require.Equal(t, 1500, res.Data[0].ViewerCount)
					startedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T20:00:00Z"))
					require.Equal(t, startedAt.Unix(), res.Data[0].StartedAt.Unix())
				}, err
			},
		},
		{
			"Get Followed Streams",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/streams/followed", &api.ResponseData[api.Stream]{
					Data: []api.Stream{{
						ID:          "123456789",
						UserID:      "987654321",
						UserLogin:   "cool_streamer",
						UserName:    "Cool_Streamer",
						GameID:      "21779",
						GameName:    "League of Legends",
						Type:        "live",
						Title:       "Playing some League!",
						ViewerCount: 1500,
						StartedAt:   Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T20:00:00Z")),
					}},
				}, apitest.RequireQueryParam("user_id"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Streams.Followed.List("1234").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "123456789", res.Data[0].ID)
					require.Equal(t, "987654321", res.Data[0].UserID)
					require.Equal(t, "cool_streamer", res.Data[0].UserLogin)
					require.Equal(t, "Cool_Streamer", res.Data[0].UserName)
					require.Equal(t, "21779", res.Data[0].GameID)
					require.Equal(t, "League of Legends", res.Data[0].GameName)
					require.Equal(t, "live", res.Data[0].Type)
					require.Equal(t, "Playing some League!", res.Data[0].Title)
					require.Equal(t, 1500, res.Data[0].ViewerCount)
					startedAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T20:00:00Z"))
					require.Equal(t, startedAt.Unix(), res.Data[0].StartedAt.Unix())
				}, err
			},
		},
		{
			"Create Stream Marker",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/streams/markers", &api.ResponseData[api.StreamMarkerData]{
					Data: []api.StreamMarkerData{{
						ID:              "marker123",
						CreatedAt:       Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T21:00:00Z")),
						Description:     "My first marker",
						PositionSeconds: 3600,
					}},
				}, apitest.RequireBodyParam("user_id"), apitest.RequireBodyParam("description"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Streams.Markers.Insert("1234").Description("test").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "marker123", res.Data[0].ID)
					createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T21:00:00Z"))
					require.Equal(t, createdAt.Unix(), res.Data[0].CreatedAt.Unix())
					require.Equal(t, "My first marker", res.Data[0].Description)
					require.Equal(t, 3600, res.Data[0].PositionSeconds)
				}, err
			},
		},
		{
			"Get Stream Markers",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/streams/markers", &api.ResponseData[api.StreamMarker]{
					Data: []api.StreamMarker{
						{
							UserID:    "1234",
							UserLogin: "twitchdev",
							UserName:  "TwitchDev",
							Videos: []api.StreamMarkerVideo{
								{
									VideoID: "video123",
									Markers: []api.StreamMarkerData{
										{
											ID:              "marker123",
											CreatedAt:       Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T21:00:00Z")),
											Description:     "My first marker",
											PositionSeconds: 3600,
										},
									},
								},
							},
						},
					},
				}, apitest.RequireQueryParam("user_id"), apitest.RequireQueryParam("video_id"), apitest.QueryParamEquals("before", "123cba"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Streams.Markers.List("1234", "video123").Before("123cba").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].UserID)
					require.Equal(t, "twitchdev", res.Data[0].UserLogin)
					require.Equal(t, "TwitchDev", res.Data[0].UserName)
					require.Len(t, res.Data[0].Videos, 1)
					require.Equal(t, "video123", res.Data[0].Videos[0].VideoID)
					require.Len(t, res.Data[0].Videos[0].Markers, 1)
					require.Equal(t, "marker123", res.Data[0].Videos[0].Markers[0].ID)
					createdAt := Must[time.Time](t)(time.Parse(time.RFC3339, "2023-08-01T21:00:00Z"))
					require.Equal(t, createdAt.Unix(), res.Data[0].Videos[0].Markers[0].CreatedAt.Unix())
					require.Equal(t, "My first marker", res.Data[0].Videos[0].Markers[0].Description)
					require.Equal(t, 3600, res.Data[0].Videos[0].Markers[0].PositionSeconds)
				}, err
			},
		},
	})
}

func TestAPI_Subscriptions(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Broadcaster Subscriptions",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/subscriptions", &api.ResponseData[api.ChannelSubscription]{
					Data: []api.ChannelSubscription{{
						BroadcasterID:   "1234",
						BroadcasterName: "CoolStreamer",
						IsGift:          false,
						PlanName:        "Premium Plan",
						Tier:            "1000",
						UserID:          "5678",
						UserName:        "CoolViewer",
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.QueryParamEquals("user_id", "4567"), apitest.QueryParamEquals("before", "cba321"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Subscriptions.List("1234").UserID("4567").Before("cba321").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "CoolStreamer", res.Data[0].BroadcasterName)
					require.False(t, res.Data[0].IsGift)
					require.Equal(t, "Premium Plan", res.Data[0].PlanName)
					require.Equal(t, "1000", res.Data[0].Tier)
					require.Equal(t, "5678", res.Data[0].UserID)
					require.Equal(t, "CoolViewer", res.Data[0].UserName)
				}, err
			},
		},
		{
			"Check User Subscription",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/subscriptions/user", &api.ResponseData[api.UserSubscriptionStatus]{
					Data: []api.UserSubscriptionStatus{{
						BroadcasterID:   "1234",
						BroadcasterName: "CoolStreamer",
						IsGift:          false,
						Tier:            "1000",
					}},
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("user_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Subscriptions.Subscribed.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, res.Data, 1)
					require.Equal(t, "1234", res.Data[0].BroadcasterID)
					require.Equal(t, "CoolStreamer", res.Data[0].BroadcasterName)
					require.False(t, res.Data[0].IsGift)
					require.Equal(t, "1000", res.Data[0].Tier)
				}, err
			},
		},
	})
}

func TestAPI_Teams(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Channel Teams",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/teams/channel", &api.ResponseData[api.Team]{
					Data: []api.Team{{
						ID:        "1234",
						TeamName:  "AwesomeTeam",
						Info:      "We are an awesome team!",
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2015-09-15T17:16:03Z")),
						UpdatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2016-09-15T17:16:03Z")),
					}},
				}, apitest.RequireQueryParam("broadcaster_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Teams.Channels.List("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
				}, err
			},
		},
		{
			"Get Teams",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/teams", &api.ResponseData[api.Team]{
					Data: []api.Team{{
						ID:        "1234",
						TeamName:  "AwesomeTeam",
						Info:      "We are an awesome team!",
						CreatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2015-09-15T17:16:03Z")),
						UpdatedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2016-09-15T17:16:03Z")),
					}},
				}, apitest.RequireQueryParam("name"), apitest.QueryParamEquals("id", "1234"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Teams.List().Name("AwesomeTeam").ID("1234").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
				}, err
			},
		},
	})
}

func TestAPI_Users(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Users",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/users", &api.ResponseData[api.User]{
					Data: []api.User{{
						UserID:          "141981764",
						UserLogin:       "twitchdev",
						UserName:        "TwitchDev",
						Type:            "staff",
						BroadcasterType: "partner",
						Description:     "Just a Twitch Dev.",
						ProfileImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/twitchdev-profile_image-...",
						OfflineImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/twitchdev-channel_offline_image-...",
						CreatedAt:       Must[time.Time](t)(time.Parse(time.RFC3339, "2016-04-21T22:48:16Z")),
					}},
				}, apitest.QueryParamEquals("id", "141981764"), apitest.QueryParamEquals("login", "twitchdev"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Users.List().ID("141981764").Login("twitchdev").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
				}, err
			},
		},
		{
			"Get Authorization By User",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/authorization/users", &api.ResponseData[api.UserAuthorization]{
					Data: []api.UserAuthorization{{
						UserID:    "141981764",
						UserLogin: "twitchdev",
						UserName:  "TwitchDev",
						Scopes:    []string{"user:read:email", "channel:manage:broadcast"},
					}},
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Users.Authorization.List("141981764").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
					require.Len(t, res.Data, 1)
					require.Equal(t, "141981764", res.Data[0].UserID)
					require.Equal(t, "twitchdev", res.Data[0].UserLogin)
					require.Equal(t, "TwitchDev", res.Data[0].UserName)
					require.ElementsMatch(t, []string{"user:read:email", "channel:manage:broadcast"}, res.Data[0].Scopes)
				}, err
			},
		},
	})
}

func TestAPI_Videos(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Videos",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, "/helix/videos", &api.ResponseData[api.Video]{
					Data: []api.Video{{
						ID:               "335921245",
						BroadcasterID:    "141981764",
						BroadcasterLogin: "twitchdev",
						BroadcasterName:  "TwitchDev",
						Title:            "Twitch Developers 101",
						Description:      "Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.",
						URL:              "https://www.twitch.tv/videos/335921245",
						ThumbnailURL:     "https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg",
						Viewable:         "public",
						ViewCount:        1863062,
						Language:         "en",
						Type:             "upload",
						Duration:         api.VideoDuration(time.Second * 201),
						MutedSegments: []api.VideoMutedSegment{{
							Duration: 30,
							Offset:   120,
						}},
						CreatedAt:   Must[time.Time](t)(time.Parse(time.RFC3339, "2018-10-02T17:21:19Z")),
						PublishedAt: Must[time.Time](t)(time.Parse(time.RFC3339, "2018-10-02T17:21:19Z")),
					}},
				}, apitest.RequireQueryParam("user_id"), apitest.QueryParamEquals("id", "9876"), apitest.QueryParamEquals("game_id", "4567"), apitest.QueryParamEquals("language", "en-US"), apitest.QueryParamEquals("period", "all"), apitest.QueryParamEquals("type", "upload"), apitest.QueryParamEquals("sort", "descending"), apitest.QueryParamEquals("after", "abc123"), apitest.QueryParamEquals("first", "1"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Videos.List().UserID("1234").ID("9876").GameID("4567").Language("en-US").Type("upload").Sort("descending").Period("all").After("abc123").First(1).Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
				}, err
			},
		},
		{
			"Delete Videos",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodDelete, "/helix/videos", &api.ResponseData[string]{
					Data: []string{"1234", "9876"},
				}, apitest.RequireQueryParam("id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Videos.Delete("1234").ID("9876").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
				}, err
			},
		},
	})
}

func TestAPI_Whispers(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Send Whisper",
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodPost, "/helix/whispers", &api.ResponseData[any]{},
					apitest.RequireQueryParam("from_user_id"), apitest.RequireQueryParam("to_user_id"), apitest.BodyParamEquals("message", "Hello!"),
				)
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Whispers.Insert("1234", "5678", "Hello!").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Equal(t, 200, res.StatusCode)
				}, err
			},
		},
	})
}

func RunEndpointTestCases(t *testing.T, tests []EndpointTestCase) {
	t.Helper()
	mock := apitest.NewMockAPI(t, apitest.WithTLS(), apitest.EnableHTTP2())
	clientID, clientSecret, err := mock.RegisterApplication()
	authorization := api.AppAccess(clientID, clientSecret)
	oauthEndpoint := mock.OAuthTokenEndpoint()
	require.NoError(t, err)
	require.NotEmpty(t, clientID)
	require.NotEmpty(t, clientSecret)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(authorization))
	for _, tt := range tests {
		endpoint := tt.endpoint(mock)
		t.Run(fmt.Sprintf("AppToken_%s", tt.name), func(t *testing.T) {
			check, err := tt.fetch(client)
			require.NoError(t, err)
			require.Exactly(t, 1, oauthEndpoint.TimesCalled)
			require.Exactly(t, 1, oauthEndpoint.Successes)
			require.Exactly(t, 0, oauthEndpoint.Failures)
			require.Exactly(t, 1, endpoint.TimesCalled)
			require.Exactly(t, 1, endpoint.Successes)
			require.Exactly(t, 0, endpoint.Failures)
			check(t)
		})

		t.Run(fmt.Sprintf("UserToken_%s", tt.name), func(t *testing.T) {
			token, err := mock.NewBearerToken(clientID)
			require.NoError(t, err)

			check, err := tt.fetch(client, api.WithBearerToken(token.AccessToken))
			require.NoError(t, err)
			require.Exactly(t, 1, oauthEndpoint.TimesCalled)
			require.Exactly(t, 1, oauthEndpoint.Successes)
			require.Exactly(t, 0, oauthEndpoint.Failures)
			require.Exactly(t, 2, endpoint.TimesCalled)
			require.Exactly(t, 2, endpoint.Successes)
			require.Exactly(t, 0, endpoint.Failures)
			check(t)
		})
	}
}

func Must[T any](t *testing.T) func(T, error) T {
	return func(v T, err error) T {
		require.NoError(t, err)
		return v
	}
}

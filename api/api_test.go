package api_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
)

type EndpointTestCase struct {
	name     string
	opts     []apitest.MockTwitchAPIOption
	endpoint func(*apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint
	fetch    func(*api.Client, ...api.RequestOption) (func(t *testing.T), error)
}

func TestAPI_Ads(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Start Commercial",
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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

func TestAPI_Channels(t *testing.T) {}

func TestAPI_ChannelPoints(t *testing.T) {}

func TestAPI_Charity(t *testing.T) {}

func TestAPI_Chat(t *testing.T) {}

func TestAPI_Clips(t *testing.T) {}

func TestAPI_Conduits(t *testing.T) {}

func TestAPI_ContentLabels(t *testing.T) {}

func TestAPI_Entitlements(t *testing.T) {}

func TestAPI_Extensions(t *testing.T) {}

func TestAPI_EventSub(t *testing.T) {}

func TestAPI_Games(t *testing.T) {}

func TestAPI_GuestStar(t *testing.T) {}

func TestAPI_HypeTrain(t *testing.T) {}

func TestAPI_Moderation(t *testing.T) {}

func TestAPI_Polls(t *testing.T) {}

func TestAPI_Predictions(t *testing.T) {}

func TestAPI_Raids(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Start A Raid",
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
	})
}

func TestAPI_Videos(t *testing.T) {
	RunEndpointTestCases(t, []EndpointTestCase{
		{
			"Get Videos",
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
			[]apitest.MockTwitchAPIOption{},
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := apitest.NewMockAPI(t, tt.opts...)
			endpoint := tt.endpoint(mock)

			clientID, secret, err := mock.RegisterApplication()
			require.NoError(t, err)
			require.NotEmpty(t, clientID)
			require.NotEmpty(t, secret)

			token, err := mock.NewBearerToken(clientID)
			require.NoError(t, err)
			require.NotEmpty(t, token)

			client := api.New(clientID, api.WithHTTPClient(mock.Client()))
			check, err := tt.fetch(client, api.WithBearerToken(token))
			require.NoError(t, err)
			require.Exactly(t, 1, endpoint.TimesCalled)
			require.Exactly(t, 1, endpoint.Successes)
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

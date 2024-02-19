package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	clientID     string
	clientSecret string
	bearerToken  string
	httpClient   HTTPClient

	Ads           *AdsResource
	Analytics     *AnalyticsResource
	Bits          *BitsResource
	Channels      *ChannelsResource
	ChannelPoints *ChannelPointsResource
	Charity       *CharityResource
	Chat          *ChatResource
	Clips         *ClipsResource
	Entitlements  *EntitlementsResource
	Extensions    *ExtensionsResource
	EventSub      *EventSubResource
	Games         *GamesResource
	Goals         *GoalsResource
	GuestStar     *GuestStarResource
	HypeTrain     *HypeTrainResource
	Moderation    *ModerationResource
	Polls         *PollsResource
	Predictions   *PredictionsResource
	Raids         *RaidsResource
	Schedule      *ScheduleResource
	Search        *SearchResource
	Streams       *StreamsResource
	Subscriptions *SubscriptionsResource
	Tags          *TagsResource
	Teams         *TeamsResource
	Users         *UsersResource
	Videos        *VideosResource
	Whispers      *WhispersResource
}

const BaseURL = "https://api.twitch.tv/helix"

// New creates a new API client for Twitch.
func New(clientID string, opts ...ClientOption) *Client {
	defaultOpts := []ClientOption{
		WithHTTPClient(http.DefaultClient),
	}

	client := &Client{clientID: clientID}
	for _, opt := range append(defaultOpts, opts...) {
		opt(client)
	}

	client.Ads = NewAdsResource(client)
	client.Analytics = NewAnalyticsResource(client)
	client.Bits = NewBitsResource(client)
	client.Channels = NewChannelsResource(client)
	client.ChannelPoints = NewChannelPointsResource(client)
	client.Charity = NewCharityResource(client)
	client.Chat = NewChatResource(client)
	client.Clips = NewClipsResource(client)
	client.Entitlements = NewEntitlementsResource(client)
	client.Extensions = NewExtensionsResource(client)
	client.EventSub = NewEventSubResource(client)
	client.Games = NewGamesResource(client)
	client.Goals = NewGoalsResource(client)
	client.GuestStar = NewGuestStarResource(client)
	client.HypeTrain = NewHypeTrainResource(client)
	client.Moderation = NewModerationResource(client)
	client.Polls = NewPollsResource(client)
	client.Predictions = NewPredictionsResource(client)
	client.Raids = NewRaidsResource(client)
	client.Schedule = NewScheduleResource(client)
	client.Search = NewSearchResource(client)
	client.Streams = NewStreamsResource(client)
	client.Subscriptions = NewSubscriptionsResource(client)
	client.Tags = NewTagsResource(client)
	client.Teams = NewTeamsResource(client)
	client.Users = NewUsersResource(client)
	client.Videos = NewVideosResource(client)
	client.Whispers = NewWhispersResource(client)
	return client
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, opts ...RequestOption) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", BaseURL, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-ID", c.clientID)
	for _, opt := range opts {
		opt(req)
	}

	if c.bearerToken != "" {
		if req.Header.Get("Authorization") == "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
		}
	}

	return c.httpClient.Do(req)
}

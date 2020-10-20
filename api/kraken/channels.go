package kraken

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ChannelsData stores a list of channels returned by the Twitch Channels API.
type ChannelsData struct {
	Total int       `json:"_total"`
	Data  []Channel `json:"channels"`
}

// Channel stores data about a Twitch channel.
type Channel struct {
	ID                           int    `json:"_id"`
	Login                        string `json:"name"`
	DisplayName                  string `json:"display_name"`
	Title                        string `json:"status"`
	Description                  string `json:"description"`
	Game                         string `json:"game"`
	URL                          string `json:"url"`
	Logo                         string `json:"logo"`
	VideoBanner                  string `json:"video_banner"`
	ProfileBanner                string `json:"profile_banner"`
	ProfileBannerBackgroundColor string `json:"profile_banner_background_color"`
	Views                        int    `json:"views"`
	Followers                    int    `json:"followers"`
	IsMature                     bool   `json:"mature"`
	IsPartner                    bool   `json:"partner"`
	IsPrivateVideo               bool   `json:"private_video"`
	IsPrivacyOptionsEnabled      bool   `json:"privacy_options_enabled"`
	StreamKey                    string `json:"stream_key"`
	Email                        string `json:"email"`
	BroadcastLanguage            string `json:"broadcaster_language"`
	BroadcastSoftware            string `json:"broadcaster_software"`
	Language                     string `json:"language"`
	Type                         string `json:"broadcaster_type"`
	CreatedAt                    string `json:"created_at"`
	UpdatedAt                    string `json:"updated_at"`
}

// GetOwnChannel retrieves a channel object based on the clients OAuth token.
//
// See: https://dev.twitch.tv/docs/v5/reference/channels#get-channel
func (client *Client) GetOwnChannel() (*Channel, error) {
	if len(client.token) < 1 {
		return nil, errors.New("no oauth token is assigned to the client")
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("channel?_t=%d", time.Now().UTC().Unix()), nil)
	if err != nil {
		return nil, err
	}
	channel := &Channel{}
	if err := json.Unmarshal(res.Body, &channel); err != nil {
		return nil, err
	}
	client.Self.Channel = *channel
	return channel, nil
}

// GetChannelsByID retrieves a specified channel object.
//
// See: https://dev.twitch.tv/docs/v5/reference/channels#get-channel-by-id
func (client *Client) GetChannelsByID(ids ...string) (*ChannelsData, error) {
	if len(ids) < 1 {
		return nil, errors.New("you must provide at least 1 channel id")
	}
	if len(ids) > 100 {
		return nil, errors.New("you can only request a total of 100 channels per request")
	}
	params := fmt.Sprintf("&id=%s", strings.Join(ids, ","))
	res, err := client.Request(http.MethodGet, fmt.Sprintf("channels?_t=%d%s", time.Now().UTC().Unix(), params), nil)
	if err != nil {
		return nil, err
	}
	channels := &ChannelsData{}
	if err := json.Unmarshal(res.Body, &channels); err != nil {
		return nil, err
	}
	return channels, nil
}

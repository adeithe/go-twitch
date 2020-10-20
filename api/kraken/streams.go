package kraken

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// StreamType parameter used to filter live streams.
type StreamType string

const (
	// Live on Twitch.
	Live StreamType = "live"
	// Playlist appears to always return nothing.
	Playlist StreamType = "playlist"
	// WatchParty streams watching something via Amazon Prime Video.
	WatchParty StreamType = "watch_party"
	// All streams on Twitch, no filters.
	All StreamType = "all"
)

// StreamOpts stores options for requests to the Twitch Streams API.
type StreamOpts struct {
	ChannelIDs []string
	Limit      int
	Offset     int
	Game       string
	Language   string
	Type       StreamType
}

// StreamSummary stores summary data about a category on Twitch.
type StreamSummary struct {
	Channels int `json:"channels"`
	Viewers  int `json:"viewers"`
}

// StreamsData stores a list of streams returned by the Twitch Streams API.
type StreamsData struct {
	Data []Stream `json:"streams"`
}

// Stream stores data about a livestream on Twitch.
type Stream struct {
	ID           int      `json:"_id"`
	Game         string   `json:"game"`
	CommunityID  string   `json:"community_id"`
	CommunityIDs []string `json:"community_ids"`
	Viewers      int      `json:"viewers"`
	VideoHeight  int      `json:"video_height"`
	Delay        int      `json:"delay"`
	AverageFPS   int      `json:"average_fps"`
	Preview      ImageURL `json:"preview"`
	Channel      Channel  `json:"channel"`
	IsPlaylist   bool     `json:"is_playlist"`
	Platform     string   `json:"broadcast_platform"`
	Type         string   `json:"stream_type"`
	CreatedAt    string   `json:"created_at"`
}

// GetStreamSummary retrieves a stream object for the specified game.
//
// See: https://dev.twitch.tv/docs/v5/reference/streams#get-streams-summary
func (client *Client) GetStreamSummary(game string) (*StreamSummary, error) {
	game = url.QueryEscape(strings.ToLower(game))
	res, err := client.Request(http.MethodGet, fmt.Sprintf("streams/summary?_t=%d&game=%s", time.Now().UTC().Unix(), game), nil)
	if err != nil {
		return nil, err
	}
	summary := &StreamSummary{}
	if err := json.Unmarshal(res.Body, &summary); err != nil {
		return nil, err
	}
	return summary, nil
}

// GetStreams retrieves a list of stream objects based on the specified StreamOpts.
//
// See: https://dev.twitch.tv/docs/v5/reference/streams#get-live-streams
func (client *Client) GetStreams(opts StreamOpts) (*StreamsData, error) {
	if len(opts.ChannelIDs) > 100 {
		return nil, errors.New("you can only request a total of 100 streams at a time")
	}
	params := ""
	if len(opts.ChannelIDs) > 0 {
		params += fmt.Sprintf("&channel=%s", strings.Join(opts.ChannelIDs, ","))
	}
	if opts.Limit > 0 {
		params += fmt.Sprintf("&limit=%d", opts.Limit)
	}
	if opts.Offset > 0 {
		params += fmt.Sprintf("&offset=%d", opts.Offset)
	}
	if len(opts.Game) > 0 {
		params += fmt.Sprintf("&game=%s", url.QueryEscape(strings.ToLower(opts.Game)))
	}
	if len(opts.Language) > 0 {
		params += fmt.Sprintf("&language=%s", opts.Language)
	}
	if len(opts.Type) > 0 {
		params += fmt.Sprintf("&stream_type=%s", opts.Type)
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("streams?_t=%d%s", time.Now().UTC().Unix(), params), nil)
	if err != nil {
		return nil, err
	}
	streams := &StreamsData{}
	if err := json.Unmarshal(res.Body, &streams); err != nil {
		return nil, err
	}
	return streams, nil
}

package helix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// StreamOpts stores options for requests to the Twitch Streams API.
type StreamOpts struct {
	First      int
	UserIDs    []string
	UserLogins []string
	GameIDs    []string
	Languages  []string
	After      Pagination
	Before     Pagination
}

// StreamsData stores a list of streams returned by the Twitch Streams API.
type StreamsData struct {
	Data       []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Stream stores data about a livestream on Twitch.
type Stream struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	Login        string   `json:"user_name"`
	GameID       string   `json:"game_id"`
	Title        string   `json:"title"`
	ThumbnailURL string   `json:"thumbnail_url"`
	TagIDs       []string `json:"tag_ids"`
	Language     string   `json:"language"`
	ViewerCount  int      `json:"viewer_count"`
	Type         string   `json:"type"`
	StartedAt    string   `json:"started_at"`
}

// GetStreams retrieves a list of stream objects based on the specified StreamOpts.
func (client *Client) GetStreams(opts StreamOpts) (*StreamsData, error) {
	if len(opts.UserIDs)+len(opts.UserLogins) > 100 {
		return nil, errors.New("you can only request a total of 100 streams at a time")
	}
	if len(opts.GameIDs) > 100 {
		return nil, errors.New("you can only request a total of 100 game ids at a time")
	}
	if len(opts.Languages) > 100 {
		return nil, errors.New("you can only request a total of 100 languages at a time")
	}
	params := ""
	if opts.First > 0 {
		params += fmt.Sprintf("&first=%d", opts.First)
	}
	if len(opts.UserIDs) > 0 {
		params += fmt.Sprintf("&user_id=%s", strings.Join(opts.UserIDs, "&user_id="))
	}
	if len(opts.UserLogins) > 0 {
		params += fmt.Sprintf("&user_login=%s", strings.Join(opts.UserLogins, "&user_login="))
	}
	if len(opts.GameIDs) > 0 {
		params += fmt.Sprintf("&game_id=%s", strings.Join(opts.GameIDs, "&game_id="))
	}
	if len(opts.Languages) > 0 {
		params += fmt.Sprintf("&language=%s", strings.Join(opts.Languages, "&language="))
	}
	if len(opts.After.Cursor) > 0 {
		params += fmt.Sprintf("&after=%s", opts.After.Cursor)
	}
	if len(opts.Before.Cursor) > 0 {
		params += fmt.Sprintf("&before=%s", opts.Before.Cursor)
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

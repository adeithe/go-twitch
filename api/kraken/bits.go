package kraken

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BitOpts stores options for requests to the Twitch Bits API.
type BitOpts struct {
	ChannelID string
}

// CheermotesData stores cheer actions returned by the Twitch Bits API.
type CheermotesData struct {
	Data []Cheermote `json:"actions"`
}

// Cheermote stores data about an available cheer action.
type Cheermote struct {
	Prefix string   `json:"prefix"`
	Scales []string `json:"scales"`
	Tiers  []struct {
		ID       string                       `json:"id"`
		MinBits  int                          `json:"min_bits"`
		Color    string                       `json:"color"`
		Images   map[string]map[string]string `json:"images"`
		CanCheer bool                         `json:"can_cheer"`
	} `json:"tiers"`
	Backgrounds []string `json:"backgrounds"`
	States      []string `json:"states"`
	Priority    int      `json:"priority"`
	Type        string   `json:"type"`
	UpdatedAt   string   `json:"updated_at"`
}

// GetCheermotes retrieves the list of available cheermotes, animated emotes to which viewers can assign Bits, to cheer in chat.
//
// See: https://dev.twitch.tv/docs/v5/reference/bits#get-cheermotes
func (client *Client) GetCheermotes(opts BitOpts) (*CheermotesData, error) {
	params := ""
	if len(opts.ChannelID) > 0 {
		params += fmt.Sprint("&channel_id=", opts.ChannelID)
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("bits/actions?_t=%d%s", time.Now().UTC().Unix(), params), nil)
	if err != nil {
		return nil, err
	}
	cheermotes := &CheermotesData{}
	if err := json.Unmarshal(res.Body, &cheermotes); err != nil {
		return nil, err
	}
	return cheermotes, nil
}

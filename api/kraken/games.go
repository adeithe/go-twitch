package kraken

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GameOpts stores options for requests to the Twitch Games API.
type GameOpts struct {
	Limit  int
	Offset int
}

// TopGames stores a list of top categories returned by the Twitch Games API.
type TopGames struct {
	Total int `json:"_total"`
	Data  []struct {
		Channels int  `json:"channels"`
		Viewers  int  `json:"viewers"`
		Game     Game `json:"game"`
	} `json:"top"`
}

// Game stores data about a category on Twitch.
type Game struct {
	ID            int      `json:"_id"`
	GiantbombID   int      `json:"giantbomb_id"`
	Name          string   `json:"name"`
	Box           ImageURL `json:"box"`
	Logo          ImageURL `json:"logo"`
	LocalizedName string   `json:"localized_name"`
	Locale        string   `json:"locale"`
}

// GetTopGames retrieves games sorted by number of current viewers on Twitch, most popular first.
//
// See: https://dev.twitch.tv/docs/v5/reference/games#get-top-games
func (client *Client) GetTopGames(opts GameOpts) (*TopGames, error) {
	params := ""
	if opts.Limit > 0 {
		params += fmt.Sprintf("&limit=%d", opts.Limit)
	}
	if opts.Offset > 0 {
		params += fmt.Sprintf("&offset=%d", opts.Offset)
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("games/top?_t=%d%s", time.Now().UTC().Unix(), params), nil)
	if err != nil {
		return nil, err
	}
	top := &TopGames{}
	if err := json.Unmarshal(res.Body, &top); err != nil {
		return nil, err
	}
	return top, nil
}

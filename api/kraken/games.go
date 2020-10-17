package kraken

// GameOpts stores data about a request to the Twitch Games API.
type GameOpts struct {
	Limit  int
	Offset int
}

// GamesData stores data about a response from the Twitch Games API.
type GamesData struct {
	Total int `json:"_total"`
	Top   []struct {
		Channels int  `json:"channels"`
		Viewers  int  `json:"viewers"`
		Data     Game `json:"game"`
	} `json:"top"`
}

// Game stores data about a game returned by the Twitch Games API.
type Game struct {
	ID            int      `json:"_id"`
	GiantbombID   int      `json:"giantbomb_id"`
	Name          string   `json:"name"`
	Box           ImageURL `json:"box"`
	Logo          ImageURL `json:"logo"`
	LocalizedName string   `json:"localized_name"`
	Locale        string   `json:"locale"`
}

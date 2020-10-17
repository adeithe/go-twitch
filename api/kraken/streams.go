package kraken

// StreamOpts stores data about a request to the Twitch Stream API.
type StreamOpts struct {
	ChannelIDs []string
	Game       string
	Language   string
	Limit      int
	Offset     int
	Type       string
}

// StreamSummary stores the number of lives channels and active viewers for a category.
type StreamSummary struct {
	Channels int `json:"channels"`
	Viewers  int `json:"viewers"`
}

// StreamData stores data about a response from the Twitch Stream API.
type StreamData struct {
	Total int      `json:"_total"`
	Data  []Stream `json:"streams"`
}

// Stream stores data about a stream returned by the Twitch Stream API.
type Stream struct {
	ID           int      `json:"_id"`
	Game         string   `json:"game"`
	Platform     string   `json:"broadcast_platform"`
	CommunityID  string   `json:"community_id"`
	CommunityIDs []string `json:"community_ids"`
	Viewers      int      `json:"viewers"`
	VideoHeight  int      `json:"video_height"`
	FPS          int      `json:"average_fps"`
	Delay        int      `json:"delay"`
	IsPlaylist   bool     `json:"is_playlist"`
	Channel      Channel  `json:"channel"`
	Thumbnail    ImageURL `json:"preview"`
	Type         string   `json:"stream_type"`
	CreatedAt    string   `json:"created_at"`
}

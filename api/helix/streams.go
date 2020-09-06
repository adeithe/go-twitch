package helix

// StreamOpts stores data about a request to the Twitch Stream API.
type StreamOpts struct {
	First      int
	GameIDs    []string
	Language   []string
	UserIDs    []string
	UserLogins []string
	Before     string
	After      string
}

// StreamData stores data about a response from the Twitch Stream API.
type StreamData struct {
	Data       []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Stream stores data about a stream returned by the Twitch Stream API.
type Stream struct {
	ID        string   `json:"id"`
	GameID    string   `json:"game_id"`
	UserID    string   `json:"user_id"`
	UserName  string   `json:"user_name"`
	Title     string   `json:"title"`
	Viewers   int      `json:"viewer_count"`
	Language  string   `json:"language"`
	TagIDs    []string `json:"tag_ids"`
	Thumbnail string   `json:"thumbnail_url"`
	Type      string   `json:"type"`
	CreatedAt string   `json:"started_at"`
}

package kraken

// Error stores data about a error that occurred with the Twitch API.
type Error struct {
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

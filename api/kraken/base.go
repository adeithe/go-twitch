package kraken

// Error stores data about a error that occurred with the Twitch API.
type Error struct {
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// ImageURL stores links for images, thumbnails and artwork.
type ImageURL struct {
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Template string `json:"template"`
}

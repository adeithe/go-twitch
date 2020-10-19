package kraken

// ImageURL stores image URLs as provided by Twitch.
type ImageURL struct {
	Large    string `json:"large"`
	Medium   string `json:"medium"`
	Small    string `json:"small"`
	Template string `json:"template"`
}

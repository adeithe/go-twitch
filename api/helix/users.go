package helix

// UserOpts stores data about a request to the Twitch User API.
type UserOpts struct {
	GetSelf bool
	IDs     []string
	Logins  []string
}

// User stores data about a user returned by the Twitch User API.
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Email           string `json:"email,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	Description     string `json:"description,omitempty"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ChannelViews    int    `json:"view_count"`
	BroadcasterType string `json:"broadcaster_type"`
	Type            string `json:"type"`
}

package helix

type UserOpts struct {
	GetSelf bool
	IDs     []string
	Logins  []string
}

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

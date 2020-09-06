package kraken

// UserOpts stores data about a request to the Twitch User API.
type UserOpts struct {
	IDs    []string
	Logins []string
}

// UserData stores data about a response from the Twitch User API.
type UserData struct {
	Total int    `json:"_total"`
	Data  []User `json:"users"`
}

// User stores data about a specific user returned by the Twitch User API.
type User struct {
	ID              string `json:"_id"`
	Login           string `json:"name"`
	DisplayName     string `json:"display_name,omitempty"`
	Description     string `json:"bio,omitempty"`
	ProfileImageURL string `json:"logo,omitempty"`
	Type            string `json:"type"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

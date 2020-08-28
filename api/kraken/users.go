package kraken

type UserOpts struct {
	IDs    []string
	Logins []string
}

type UserData struct {
	Total int    `json:"_total"`
	Data  []User `json:"users"`
}

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

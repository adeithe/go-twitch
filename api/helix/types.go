package helix

// Pagination stores a cursor for the Twitch Helix API.
type Pagination struct {
	Cursor string `json:"cursor"`
}

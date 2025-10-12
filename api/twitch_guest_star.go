package api

// GuestStarResource represents the Twitch Guest Star API.
type GuestStarResource struct {
	client *Client
}

// NewGuestStarResource creates a new GuestStarResource.
func NewGuestStarResource(client *Client) *GuestStarResource {
	return &GuestStarResource{client}
}

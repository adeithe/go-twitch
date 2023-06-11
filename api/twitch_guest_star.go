package api

type GuestStarResource struct {
	client *Client
}

func NewGuestStarResource(client *Client) *GuestStarResource {
	return &GuestStarResource{client}
}

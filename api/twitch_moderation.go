package api

type ModerationResource struct {
	client *Client
}

func NewModerationResource(client *Client) *ModerationResource {
	return &ModerationResource{client}
}

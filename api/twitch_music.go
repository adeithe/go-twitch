package api

type MusicResource struct {
	client *Client
}

func NewMusicResource(client *Client) *MusicResource {
	return &MusicResource{client}
}

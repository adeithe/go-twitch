package api

type ClipsResource struct {
	client *Client
}

func NewClipsResource(client *Client) *ClipsResource {
	return &ClipsResource{client}
}

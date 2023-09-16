package api

type VideosResource struct {
	client *Client
}

func NewVideosResource(client *Client) *VideosResource {
	return &VideosResource{client}
}

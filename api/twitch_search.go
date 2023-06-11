package api

type SearchResource struct {
	client *Client
}

func NewSearchResource(client *Client) *SearchResource {
	return &SearchResource{client}
}

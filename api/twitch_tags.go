package api

type TagsResource struct {
	client *Client
}

func NewTagsResource(client *Client) *TagsResource {
	return &TagsResource{client}
}

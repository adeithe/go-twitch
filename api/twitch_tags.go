package api

// TagsResource represents the Twitch Tags API.
type TagsResource struct {
	client *Client
}

// NewTagsResource creates a new TagsResource.
func NewTagsResource(client *Client) *TagsResource {
	return &TagsResource{client}
}

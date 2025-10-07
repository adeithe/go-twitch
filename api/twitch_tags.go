package api

// TagsResource represents the Twitch Tags API.
type TagsResource struct {
	client *Client
}

const (
	// EndpointTagsGetAllStreamTags is the endpoint for getting all stream tags.
	EndpointTagsGetAllStreamTags = TwitchAPIVersionHelix + "/tags/streams"
	// EndpointTagsGetStreamTags is the endpoint for getting tags for a stream.
	EndpointTagsGetStreamTags = TwitchAPIVersionHelix + "/streams/tags"
)

// NewTagsResource creates a new TagsResource.
func NewTagsResource(client *Client) *TagsResource {
	return &TagsResource{client}
}

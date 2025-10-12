package api

// ContentLabelsResource represents the Twitch Content Classification Labels API.
type ContentLabelsResource struct {
	client *Client
}

// NewContentLabelsResource creates a new ContentLabelsResource.
func NewContentLabelsResource(client *Client) *ContentLabelsResource {
	return &ContentLabelsResource{client}
}

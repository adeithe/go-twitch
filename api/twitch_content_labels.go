package api

// ContentLabelsResource represents the Twitch Content Classification Labels API.
type ContentLabelsResource struct {
	client *Client
}

// EndpointContentLabels is the endpoint for getting information about content classification labels.
const EndpointContentLabels = TwitchAPIVersionHelix + "/content_classification_labels"

// NewContentLabelsResource creates a new ContentLabelsResource.
func NewContentLabelsResource(client *Client) *ContentLabelsResource {
	return &ContentLabelsResource{client}
}

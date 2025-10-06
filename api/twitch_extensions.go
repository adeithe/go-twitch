package api

// ExtensionsResource represents the Twitch Extensions API.
type ExtensionsResource struct {
	client *Client
}

// NewExtensionsResource creates a new ExtensionsResource.
func NewExtensionsResource(client *Client) *ExtensionsResource {
	return &ExtensionsResource{client}
}

package api

type ExtensionsResource struct {
	client *Client
}

func NewExtensionsResource(client *Client) *ExtensionsResource {
	return &ExtensionsResource{client}
}

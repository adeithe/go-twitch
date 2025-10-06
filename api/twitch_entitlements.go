package api

// EntitlementsResource represends the Twitch Entitlements API.
type EntitlementsResource struct {
	client *Client
}

// NewEntitlementsResource creates a new EntitlementsResource.
func NewEntitlementsResource(client *Client) *EntitlementsResource {
	return &EntitlementsResource{client}
}

package api

// EntitlementsResource represends the Twitch Entitlements API.
type EntitlementsResource struct {
	client *Client
}

// EndpointEntitlements is the endpoint for the Twitch Entitlements API.
const EndpointEntitlements = TwitchAPIVersionHelix + "/entitlements/drops"

// NewEntitlementsResource creates a new EntitlementsResource.
func NewEntitlementsResource(client *Client) *EntitlementsResource {
	return &EntitlementsResource{client}
}

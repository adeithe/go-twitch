package api

// ExtensionsResource represents the Twitch Extensions API.
type ExtensionsResource struct {
	client *Client
}

const (
	// EndpointExtensionsConfiguration is the endpoint for managing extension configurations.
	EndpointExtensionsConfiguration = TwitchAPIVersionHelix + "/extensions/configurations"
	// EndpointExtensionsRequiredConfiguration is the endpoint for managing required extension configurations.
	EndpointExtensionsRequiredConfiguration = TwitchAPIVersionHelix + "/extensions/required_configuration"
	// EndpointExtensionsPubSub is the endpoint for sending extension messages via PubSub.
	EndpointExtensionsPubSub = TwitchAPIVersionHelix + "/extensions/pubsub"
	// EndpointExtensionsLiveChannels is the endpoint for fetching live channels by extension.
	EndpointExtensionsLiveChannels = TwitchAPIVersionHelix + "/extensions/live"
	// EndpointExtensionsSecrets is the endpoint for managing extension secrets.
	EndpointExtensionsSecrets = TwitchAPIVersionHelix + "/extensions/jwt/secrets"
	// EndpointExtensionsChat is the endpoint for sending chat messages via extensions.
	EndpointExtensionsChat = TwitchAPIVersionHelix + "/extensions/chat"
	// EndpointExtensions is the endpoint for getting information about an extension.
	EndpointExtensions = TwitchAPIVersionHelix + "/extensions"
	// EndpointExtensionsReleased is the endpoint for getting information about released extensions.
	EndpointExtensionsReleased = TwitchAPIVersionHelix + "/extensions/released"
	// EndpointExtensionsBitsProducts is the endpoint for managing an extensions Bits product.
	EndpointExtensionsBitsProducts = TwitchAPIVersionHelix + "/bits/extensions"
)

// NewExtensionsResource creates a new ExtensionsResource.
func NewExtensionsResource(client *Client) *ExtensionsResource {
	return &ExtensionsResource{client}
}

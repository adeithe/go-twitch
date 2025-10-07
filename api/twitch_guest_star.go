package api

// GuestStarResource represents the Twitch Guest Star API.
type GuestStarResource struct {
	client *Client
}

const (
	// EndpointGuestStarChannelSettings is the endpoint for managing guest star channel settings.
	EndpointGuestStarChannelSettings = TwitchAPIVersionHelix + "/guest_star/channel_settings"
	// EndpointGuestStarSession is the endpoint for managing a guest star session.
	EndpointGuestStarSession = TwitchAPIVersionHelix + "/guest_star/session"
	// EndpointGuestStarInvites is the endpoint for managing guest star invites.
	EndpointGuestStarInvites = TwitchAPIVersionHelix + "/guest_star/invites"
	// EndpointGuestStarSlot is the endpoint for managing a user that was invited to the guest star session.
	EndpointGuestStarSlot = TwitchAPIVersionHelix + "/guest_star/slot"
	// EndpointGuestStarSlotSettings is the endpoint for managing an invited users settings in a guest star session.
	EndpointGuestStarSlotSettings = TwitchAPIVersionHelix + "/guest_star/slot_settings"
)

// NewGuestStarResource creates a new GuestStarResource.
func NewGuestStarResource(client *Client) *GuestStarResource {
	return &GuestStarResource{client}
}

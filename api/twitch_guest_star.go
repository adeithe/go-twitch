package api

import (
	"context"
	"net/http"
)

// GuestStarResource represents the Twitch GuestStar API.
type GuestStarResource struct {
	client *Client

	// Session provides access to the Twitch Session API.
	Session *GuestStarSessionResource
}

// NewGuestStarResource creates a new GuestStarResource.
func NewGuestStarResource(client *Client) *GuestStarResource {
	r := &GuestStarResource{client: client}
	r.Session = NewGuestStarSessionResource(client)
	return r
}

// GuestStarSessionResource represents the Twitch GuestStarSession API.
type GuestStarSessionResource struct {
	client *Client
}

// NewGuestStarSessionResource creates a new GuestStarSessionResource.
func NewGuestStarSessionResource(client *Client) *GuestStarSessionResource {
	return &GuestStarSessionResource{client}
}

// GuestStarSessionListCall represents a GET call to a Twitch GuestStar API endpoint.
type GuestStarSessionListCall struct {
	resource *GuestStarResource
	opts     []RequestOption
}

// GuestStarSessionListResponse represents the response from a GET request to /helix/guest_star/session.
type GuestStarSessionListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the GuestStarSession data returned by the Twitch API.
	Data []GuestStarSession
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/guest_star/session.
//
// Gets information about an ongoing Guest Star session for a particular channel.
//
// # Authorization
//
// Requires OAuth Scope: channel:read:guest_star, channel:manage:guest_star, moderator:read:guest_star or moderator:manage:guest_star.
//
// Guests must be either invited or assigned a slot within the session.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-guest-star-session
func (r *GuestStarResource) List(broadcasterID string, moderatorID string) *GuestStarSessionListCall {
	c := &GuestStarSessionListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		ModeratorID(moderatorID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *GuestStarSessionListCall) BroadcasterID(broadcasterID string) *GuestStarSessionListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ModeratorID sets the ModeratorID query parameter.
func (api *GuestStarSessionListCall) ModeratorID(moderatorID string) *GuestStarSessionListCall {
	api.opts = append(api.opts, SetQueryParameter("moderator_id", moderatorID))
	return api
}

// Do executes the request.
func (api *GuestStarSessionListCall) Do(ctx context.Context, opts ...RequestOption) (*GuestStarSessionListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/guest_star/session", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[GuestStarSession](res)
	if err != nil {
		return nil, err
	}

	return &GuestStarSessionListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

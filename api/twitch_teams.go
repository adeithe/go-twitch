package api

import (
	"context"
	"net/http"
)

// TeamsResource represents the Twitch Teams API.
type TeamsResource struct {
	client *Client

	// Channels provides access to the Twitch Channels API.
	Channels *TeamsChannelsResource
}

// NewTeamsResource creates a new TeamsResource.
func NewTeamsResource(client *Client) *TeamsResource {
	r := &TeamsResource{client: client}
	r.Channels = NewTeamsChannelsResource(client)
	return r
}

// TeamsChannelsResource represents the Twitch TeamsChannels API.
type TeamsChannelsResource struct {
	client *Client
}

// NewTeamsChannelsResource creates a new TeamsChannelsResource.
func NewTeamsChannelsResource(client *Client) *TeamsChannelsResource {
	return &TeamsChannelsResource{client}
}

// ChannelTeamsListCall represents a GET call to a Twitch TeamsChannels API endpoint.
type ChannelTeamsListCall struct {
	resource *TeamsChannelsResource
	opts     []RequestOption
}

// ChannelTeamsListResponse represents the response from a GET request to /helix/teams/channel.
type ChannelTeamsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ChannelTeam data returned by the Twitch API.
	Data []ChannelTeam
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/teams/channel.
//
// Gets the list of Twitch teams that the broadcaster is a member of.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-teams
func (r *TeamsChannelsResource) List(broadcasterID string) *ChannelTeamsListCall {
	c := &ChannelTeamsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelTeamsListCall) BroadcasterID(broadcasterID string) *ChannelTeamsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *ChannelTeamsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelTeamsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/teams/channel", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChannelTeam](res)
	if err != nil {
		return nil, err
	}

	return &ChannelTeamsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// TeamsListCall represents a GET call to a Twitch Teams API endpoint.
type TeamsListCall struct {
	resource *TeamsResource
	opts     []RequestOption
}

// TeamsListResponse represents the response from a GET request to /helix/teams.
type TeamsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Team data returned by the Twitch API.
	Data []Team
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/teams.
//
// Gets information about the specified Twitch team.
//
// You are required to specify Name or ID as they are mutually exclusive.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-teams
func (r *TeamsResource) List() *TeamsListCall {
	return &TeamsListCall{resource: r}
}

// Name sets the Name query parameter.
func (api *TeamsListCall) Name(name string) *TeamsListCall {
	api.opts = append(api.opts, SetQueryParameter("name", name))
	return api
}

// ID sets the ID query parameter.
func (api *TeamsListCall) ID(id string) *TeamsListCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// Do executes the request.
func (api *TeamsListCall) Do(ctx context.Context, opts ...RequestOption) (*TeamsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/teams", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Team](res)
	if err != nil {
		return nil, err
	}

	return &TeamsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

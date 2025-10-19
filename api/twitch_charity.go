package api

import (
	"context"
	"net/http"
)

// CharityResource represents the Twitch Charity API.
type CharityResource struct {
	client *Client

	// Campaign provides access to the Twitch Campaign API.
	Campaign *CharityCampaignResource
	// Donations provides access to the Twitch Donations API.
	Donations *CharityDonationsResource
}

// NewCharityResource creates a new CharityResource.
func NewCharityResource(client *Client) *CharityResource {
	r := &CharityResource{client: client}
	r.Campaign = NewCharityCampaignResource(client)
	r.Donations = NewCharityDonationsResource(client)
	return r
}

// CharityCampaignResource represents the Twitch CharityCampaign API.
type CharityCampaignResource struct {
	client *Client
}

// NewCharityCampaignResource creates a new CharityCampaignResource.
func NewCharityCampaignResource(client *Client) *CharityCampaignResource {
	return &CharityCampaignResource{client}
}

// CharityCampaignListCall represents a GET call to a Twitch CharityCampaign API endpoint.
type CharityCampaignListCall struct {
	resource *CharityCampaignResource
	opts     []RequestOption
}

// CharityCampaignListResponse represents the response from a GET request to /helix/charity/campaigns.
type CharityCampaignListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CharityCampaign data returned by the Twitch API.
	Data []CharityCampaign
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/charity/campaigns.
//
// Gets information about the charity campaign that a broadcaster is running.
// For example, the  fundraising goal for a campaign and the current amount of donations.
//
// To receive events when progress is made towards the campaign goal or the broadcaster changes the fundraising goal, subscribe to the channel.charity_campaign.progress subscription type.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:charity scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-charity-campaign
func (r *CharityCampaignResource) List(broadcasterID string) *CharityCampaignListCall {
	c := &CharityCampaignListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *CharityCampaignListCall) BroadcasterID(broadcasterID string) *CharityCampaignListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *CharityCampaignListCall) Do(ctx context.Context, opts ...RequestOption) (*CharityCampaignListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/charity/campaigns", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CharityCampaign](res)
	if err != nil {
		return nil, err
	}

	return &CharityCampaignListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// CharityDonationsResource represents the Twitch CharityDonations API.
type CharityDonationsResource struct {
	client *Client
}

// NewCharityDonationsResource creates a new CharityDonationsResource.
func NewCharityDonationsResource(client *Client) *CharityDonationsResource {
	return &CharityDonationsResource{client}
}

// CharityDonationsListCall represents a GET call to a Twitch CharityDonations API endpoint.
type CharityDonationsListCall struct {
	resource *CharityDonationsResource
	opts     []RequestOption
}

// CharityDonationsListResponse represents the response from a GET request to /helix/charity/campaigns.
type CharityDonationsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the CharityCampaignDonation data returned by the Twitch API.
	Data []CharityCampaignDonation
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/charity/campaigns.
//
// Gets the list of donations that users have made to the broadcasters active charity campaign.
//
// To receive events as donations occur, subscribe to the channel.charity_campaign.donate subscription type.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:charity scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-charity-campaign-donations
func (r *CharityDonationsResource) List(broadcasterID string) *CharityDonationsListCall {
	c := &CharityDonationsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *CharityDonationsListCall) BroadcasterID(broadcasterID string) *CharityDonationsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// After sets the After query parameter.
func (api *CharityDonationsListCall) After(after string) *CharityDonationsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *CharityDonationsListCall) First(first int) *CharityDonationsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *CharityDonationsListCall) Do(ctx context.Context, opts ...RequestOption) (*CharityDonationsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/charity/campaigns", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[CharityCampaignDonation](res)
	if err != nil {
		return nil, err
	}

	return &CharityDonationsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

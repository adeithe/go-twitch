package api

// CharityResource handles charity related API calls.
type CharityResource struct {
	client *Client
}

const (
	// EndpointCharityGetCampaign is the endpoint for getting charity campaigns.
	EndpointCharityGetCampaign = TwitchAPIVersionHelix + "/charity/campaigns"
	// EndpointCharityGetCampaignDonations is the endpoint for getting charity campaign donations.
	EndpointCharityGetCampaignDonations = TwitchAPIVersionHelix + "/charity/donations"
)

// NewCharityResource creates a new CharityResource.
func NewCharityResource(client *Client) *CharityResource {
	return &CharityResource{client}
}

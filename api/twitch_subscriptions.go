package api

import (
	"context"
	"net/http"
)

// SubscriptionsResource represents the Twitch Subscriptions API.
type SubscriptionsResource struct {
	client *Client

	// Subscribed provides access to the Twitch Subscribed API.
	Subscribed *SubscriptionsSubscribedResource
}

// NewSubscriptionsResource creates a new SubscriptionsResource.
func NewSubscriptionsResource(client *Client) *SubscriptionsResource {
	r := &SubscriptionsResource{client: client}
	r.Subscribed = NewSubscriptionsSubscribedResource(client)
	return r
}

// SubscriptionsSubscribedResource represents the Twitch SubscriptionsSubscribed API.
type SubscriptionsSubscribedResource struct {
	client *Client
}

// NewSubscriptionsSubscribedResource creates a new SubscriptionsSubscribedResource.
func NewSubscriptionsSubscribedResource(client *Client) *SubscriptionsSubscribedResource {
	return &SubscriptionsSubscribedResource{client}
}

// UserSubscriptionListCall represents a GET call to a Twitch SubscriptionsSubscribed API endpoint.
type UserSubscriptionListCall struct {
	resource *SubscriptionsSubscribedResource
	opts     []RequestOption
}

// UserSubscriptionListResponse represents the response from a GET request to /helix/subscriptions/user.
type UserSubscriptionListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the UserSubscriptionStatus data returned by the Twitch API.
	Data []UserSubscriptionStatus
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/subscriptions/user.
//
// Checks whether the user subscribes to the broadcaster's channel.
//
// # Authorization
//
// Requires a user access token that includes the user:read:subscriptions scope.
//
// A Twitch extensions may use an app access token if the broadcaster has granted the user:read:subscriptions scope from within the Twitch Extensions manager.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#check-user-subscription
func (r *SubscriptionsSubscribedResource) List(broadcasterID string, userID string) *UserSubscriptionListCall {
	c := &UserSubscriptionListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		UserID(userID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *UserSubscriptionListCall) BroadcasterID(broadcasterID string) *UserSubscriptionListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// UserID sets the UserID query parameter.
func (api *UserSubscriptionListCall) UserID(userID string) *UserSubscriptionListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// Do executes the request.
func (api *UserSubscriptionListCall) Do(ctx context.Context, opts ...RequestOption) (*UserSubscriptionListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/subscriptions/user", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[UserSubscriptionStatus](res)
	if err != nil {
		return nil, err
	}

	return &UserSubscriptionListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// BroadcasterSubscriptionsListCall represents a GET call to a Twitch Subscriptions API endpoint.
type BroadcasterSubscriptionsListCall struct {
	resource *SubscriptionsResource
	opts     []RequestOption
}

// BroadcasterSubscriptionsListResponse represents the response from a GET request to /helix/subscriptions.
type BroadcasterSubscriptionsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Total is the int data returned by the Twitch API.
	Total int
	// Points is the int data returned by the Twitch API.
	Points int
	// Data is the ChannelSubscription data returned by the Twitch API.
	Data []ChannelSubscription
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/subscriptions.
//
// Gets a list of users that subscribe to the specified broadcaster.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:subscriptions scope.
//
// A Twitch extensions may use an app access token if the broadcaster has granted the channel:read:subscriptions scope from within the Twitch Extensions manager.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-broadcaster-subscriptions
func (r *SubscriptionsResource) List(broadcasterID string) *BroadcasterSubscriptionsListCall {
	c := &BroadcasterSubscriptionsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// UserID adds to the UserID query parameter.
func (api *BroadcasterSubscriptionsListCall) UserID(userIDs ...string) *BroadcasterSubscriptionsListCall {
	for _, userID := range userIDs {
		api.opts = append(api.opts, AddQueryParameter("user_id", userID))
	}
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *BroadcasterSubscriptionsListCall) BroadcasterID(broadcasterID string) *BroadcasterSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Before sets the Before query parameter.
func (api *BroadcasterSubscriptionsListCall) Before(before string) *BroadcasterSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("before", before))
	return api
}

// After sets the After query parameter.
func (api *BroadcasterSubscriptionsListCall) After(after string) *BroadcasterSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *BroadcasterSubscriptionsListCall) First(first int) *BroadcasterSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *BroadcasterSubscriptionsListCall) Do(ctx context.Context, opts ...RequestOption) (*BroadcasterSubscriptionsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/subscriptions", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ChannelSubscription](res)
	if err != nil {
		return nil, err
	}

	return &BroadcasterSubscriptionsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Total:      data.Total,
		Points:     data.Points,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

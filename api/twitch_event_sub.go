package api

import (
	"context"
	"net/http"
)

// EventSubResource represents the Twitch EventSub API.
type EventSubResource struct {
	client *Client
}

// NewEventSubResource creates a new EventSubResource.
func NewEventSubResource(client *Client) *EventSubResource {
	return &EventSubResource{client}
}

// EventSubSubscriptionsListCall represents a GET call to a Twitch EventSub API endpoint.
type EventSubSubscriptionsListCall struct {
	resource *EventSubResource
	opts     []RequestOption
}

// EventSubSubscriptionsListResponse represents the response from a GET request to /helix/eventsub/subscriptions.
type EventSubSubscriptionsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// TotalCost is the int data returned by the Twitch API.
	TotalCost int
	// MaxCost is the int data returned by the Twitch API.
	MaxCost int
	// Data is the EventSubSubscription data returned by the Twitch API.
	Data []EventSubSubscription
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/eventsub/subscriptions.
//
// Gets a list of all EventSub subscriptions that the authenticated app has created.
//
// # Authorization
//
// Requires an app access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions
func (r *EventSubResource) List() *EventSubSubscriptionsListCall {
	return &EventSubSubscriptionsListCall{resource: r}
}

// SubscriptionID sets the SubscriptionID query parameter.
func (api *EventSubSubscriptionsListCall) SubscriptionID(subscriptionID string) *EventSubSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("subscription_id", subscriptionID))
	return api
}

// Status sets the Status query parameter.
func (api *EventSubSubscriptionsListCall) Status(status string) *EventSubSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("status", status))
	return api
}

// Type sets the Type query parameter.
func (api *EventSubSubscriptionsListCall) Type(t string) *EventSubSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("type", t))
	return api
}

// UserID sets the UserID query parameter.
func (api *EventSubSubscriptionsListCall) UserID(userID string) *EventSubSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// After sets the After query parameter.
func (api *EventSubSubscriptionsListCall) After(after string) *EventSubSubscriptionsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// Do executes the request.
func (api *EventSubSubscriptionsListCall) Do(ctx context.Context, opts ...RequestOption) (*EventSubSubscriptionsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/eventsub/subscriptions", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[EventSubSubscription](res)
	if err != nil {
		return nil, err
	}

	return &EventSubSubscriptionsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		TotalCost:  data.TotalCost,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

package api

import (
	"bytes"
	"context"
	"encoding/json"
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

// EventSubSubscriptionsInsertCall represents a POST call to a Twitch EventSub API endpoint.
type EventSubSubscriptionsInsertCall struct {
	resource *EventSubResource
	body     map[string]any
}

// EventSubSubscriptionsInsertResponse represents the response from a POST request to /helix/eventsub/subscriptions.
type EventSubSubscriptionsInsertResponse struct {
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
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/eventsub/subscriptions.
//
// Create a EventSub subscription.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-eventsub-subscriptions
func (r *EventSubResource) Insert() *EventSubSubscriptionsInsertCall {
	return &EventSubSubscriptionsInsertCall{resource: r, body: make(map[string]any)}
}

// Type sets the Type body parameter.
func (api *EventSubSubscriptionsInsertCall) Type(t string) *EventSubSubscriptionsInsertCall {
	api.body["type"] = t
	return api
}

// Version sets the Version body parameter.
func (api *EventSubSubscriptionsInsertCall) Version(version string) *EventSubSubscriptionsInsertCall {
	api.body["version"] = version
	return api
}

// Condition sets the Condition body parameter.
func (api *EventSubSubscriptionsInsertCall) Condition(condition EventSubCondition) *EventSubSubscriptionsInsertCall {
	api.body["condition"] = condition
	return api
}

// Transport sets the Transport body parameter.
func (api *EventSubSubscriptionsInsertCall) Transport(transport Transport) *EventSubSubscriptionsInsertCall {
	api.body["transport"] = transport
	return api
}

// SessionID sets the SessionID body parameter.
func (api *EventSubSubscriptionsInsertCall) SessionID(sessionID string) *EventSubSubscriptionsInsertCall {
	api.body["session_id"] = sessionID
	return api
}

// ConduitID sets the ConduitID body parameter.
func (api *EventSubSubscriptionsInsertCall) ConduitID(conduitID string) *EventSubSubscriptionsInsertCall {
	api.body["conduit_id"] = conduitID
	return api
}

// Do executes the request.
func (api *EventSubSubscriptionsInsertCall) Do(ctx context.Context, opts ...RequestOption) (*EventSubSubscriptionsInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/eventsub/subscriptions", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[EventSubSubscription](res)
	if err != nil {
		return nil, err
	}

	return &EventSubSubscriptionsInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		TotalCost:  data.TotalCost,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// EventSubSubscriptionsDeleteCall represents a DELETE call to a Twitch EventSub API endpoint.
type EventSubSubscriptionsDeleteCall struct {
	resource *EventSubResource
	opts     []RequestOption
}

// EventSubSubscriptionsDeleteResponse represents the response from a DELETE request to /helix/eventsub/subscriptions.
type EventSubSubscriptionsDeleteResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Delete creates a new DELETE request to /helix/eventsub/subscriptions.
//
// Create a EventSub subscription.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#delete-eventsub-subscriptions
func (r *EventSubResource) Delete(id string) *EventSubSubscriptionsDeleteCall {
	c := &EventSubSubscriptionsDeleteCall{resource: r}
	return c.
		ID(id)
}

// ID sets the ID query parameter.
func (api *EventSubSubscriptionsDeleteCall) ID(id string) *EventSubSubscriptionsDeleteCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// Do executes the request.
func (api *EventSubSubscriptionsDeleteCall) Do(ctx context.Context, opts ...RequestOption) (*EventSubSubscriptionsDeleteResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "DELETE", "/helix/eventsub/subscriptions", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &EventSubSubscriptionsDeleteResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
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
// Requires an app access token or user access token.
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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// PollsResource represents the Twitch Polls API.
type PollsResource struct {
	client *Client
}

// NewPollsResource creates a new PollsResource.
func NewPollsResource(client *Client) *PollsResource {
	return &PollsResource{client}
}

// PollsListCall represents a GET call to a Twitch Polls API endpoint.
type PollsListCall struct {
	resource *PollsResource
	opts     []RequestOption
}

// PollsListResponse represents the response from a GET request to /helix/polls.
type PollsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Poll data returned by the Twitch API.
	Data []Poll
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/polls.
//
// Gets a list of polls that have been created in the past 90 days.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:polls or channel:manage:polls scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-polls
func (r *PollsResource) List(broadcasterID string) *PollsListCall {
	c := &PollsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// ID adds to the ID query parameter.
func (api *PollsListCall) ID(ids ...string) *PollsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *PollsListCall) BroadcasterID(broadcasterID string) *PollsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// After sets the After query parameter.
func (api *PollsListCall) After(after string) *PollsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *PollsListCall) First(first int) *PollsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *PollsListCall) Do(ctx context.Context, opts ...RequestOption) (*PollsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/polls", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Poll](res)
	if err != nil {
		return nil, err
	}

	return &PollsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// PollsInsertCall represents a POST call to a Twitch Polls API endpoint.
type PollsInsertCall struct {
	resource *PollsResource
	body     map[string]any
}

// PollsInsertResponse represents the response from a POST request to /helix/polls.
type PollsInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Poll data returned by the Twitch API.
	Data []Poll
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/polls.
//
// Creates a poll that viewers in the broadcaster's channel can vote on.
//
// The poll begins as soon as it's created. You may run only one poll at a time.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:polls scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-poll
func (r *PollsResource) Insert(broadcasterID string, title string, duration int) *PollsInsertCall {
	c := &PollsInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		Title(title).
		Duration(duration)
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *PollsInsertCall) BroadcasterID(broadcasterID string) *PollsInsertCall {
	api.body["broadcaster_id"] = broadcasterID
	return api
}

// Title sets the Title body parameter.
func (api *PollsInsertCall) Title(title string) *PollsInsertCall {
	api.body["title"] = title
	return api
}

// Choices sets the Choices body parameter.
func (api *PollsInsertCall) Choices(choicess ...OutboundChoice) *PollsInsertCall {
	api.body["choices"] = choicess
	return api
}

// Duration sets the Duration body parameter.
func (api *PollsInsertCall) Duration(duration int) *PollsInsertCall {
	api.body["duration"] = duration
	return api
}

// ChannelPointsPerVote sets the ChannelPointsPerVote body parameter.
func (api *PollsInsertCall) ChannelPointsPerVote(channelPointsPerVote int) *PollsInsertCall {
	api.body["channel_points_per_vote"] = channelPointsPerVote
	return api
}

// ChannelPointsVotingEnabled sets the ChannelPointsVotingEnabled body parameter.
func (api *PollsInsertCall) ChannelPointsVotingEnabled(channelPointsVotingEnabled bool) *PollsInsertCall {
	api.body["channel_points_voting_enabled"] = channelPointsVotingEnabled
	return api
}

// Do executes the request.
func (api *PollsInsertCall) Do(ctx context.Context, opts ...RequestOption) (*PollsInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/polls", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Poll](res)
	if err != nil {
		return nil, err
	}

	return &PollsInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// PollsModifyCall represents a PATCH call to a Twitch Polls API endpoint.
type PollsModifyCall struct {
	resource *PollsResource
	body     map[string]any
}

// PollsModifyResponse represents the response from a PATCH request to /helix/polls.
type PollsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Poll data returned by the Twitch API.
	Data []Poll
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/polls.
//
// Ends an active poll. You have the option to end it or end it and archive it.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:polls scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#end-poll
func (r *PollsResource) Modify(broadcasterID string, id string, status string) *PollsModifyCall {
	c := &PollsModifyCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		ID(id).
		Status(status)
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *PollsModifyCall) BroadcasterID(broadcasterID string) *PollsModifyCall {
	api.body["broadcaster_id"] = broadcasterID
	return api
}

// ID sets the ID body parameter.
func (api *PollsModifyCall) ID(id string) *PollsModifyCall {
	api.body["id"] = id
	return api
}

// Status sets the Status body parameter.
func (api *PollsModifyCall) Status(status string) *PollsModifyCall {
	api.body["status"] = status
	return api
}

// Do executes the request.
func (api *PollsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*PollsModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/polls", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Poll](res)
	if err != nil {
		return nil, err
	}

	return &PollsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

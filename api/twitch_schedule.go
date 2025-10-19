package api

import (
	"context"
	"net/http"
	"time"
)

// ScheduleResource represents the Twitch Schedule API.
type ScheduleResource struct {
	client *Client
}

// NewScheduleResource creates a new ScheduleResource.
func NewScheduleResource(client *Client) *ScheduleResource {
	return &ScheduleResource{client}
}

// ChannelStreamScheduleListCall represents a GET call to a Twitch Schedule API endpoint.
type ChannelStreamScheduleListCall struct {
	resource *ScheduleResource
	opts     []RequestOption
}

// ChannelStreamScheduleListResponse represents the response from a GET request to /helix/schedule/icalendar.
type ChannelStreamScheduleListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the StreamSchedule data returned by the Twitch API.
	Data []StreamSchedule
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/schedule/icalendar.
//
// Gets the broadcaster's streaming schedule. You can get the entire schedule or specific segments of the schedule.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-channel-stream-schedule
func (r *ScheduleResource) List(broadcasterID string) *ChannelStreamScheduleListCall {
	c := &ChannelStreamScheduleListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ChannelStreamScheduleListCall) BroadcasterID(broadcasterID string) *ChannelStreamScheduleListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ID sets the ID query parameter.
func (api *ChannelStreamScheduleListCall) ID(id string) *ChannelStreamScheduleListCall {
	api.opts = append(api.opts, SetQueryParameter("id", id))
	return api
}

// After sets the After query parameter.
func (api *ChannelStreamScheduleListCall) After(after string) *ChannelStreamScheduleListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ChannelStreamScheduleListCall) First(first int) *ChannelStreamScheduleListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// StartTime sets the StartTime query parameter.
func (api *ChannelStreamScheduleListCall) StartTime(startTime time.Time) *ChannelStreamScheduleListCall {
	api.opts = append(api.opts, SetQueryParameter("start_time", startTime.Format(time.RFC3339)))
	return api
}

// Do executes the request.
func (api *ChannelStreamScheduleListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelStreamScheduleListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/schedule/icalendar", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[StreamSchedule](res)
	if err != nil {
		return nil, err
	}

	return &ChannelStreamScheduleListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

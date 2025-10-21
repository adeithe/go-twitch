package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// StreamsResource represents the Twitch Streams API.
type StreamsResource struct {
	client *Client

	// StreamKey provides access to the Twitch StreamKey API.
	StreamKey *StreamsStreamKeyResource
	// Followed provides access to the Twitch Followed API.
	Followed *StreamsFollowedResource
	// Markers provides access to the Twitch Markers API.
	Markers *StreamsMarkersResource
}

// NewStreamsResource creates a new StreamsResource.
func NewStreamsResource(client *Client) *StreamsResource {
	r := &StreamsResource{client: client}
	r.StreamKey = NewStreamsStreamKeyResource(client)
	r.Followed = NewStreamsFollowedResource(client)
	r.Markers = NewStreamsMarkersResource(client)
	return r
}

// StreamsStreamKeyResource represents the Twitch StreamsStreamKey API.
type StreamsStreamKeyResource struct {
	client *Client
}

// NewStreamsStreamKeyResource creates a new StreamsStreamKeyResource.
func NewStreamsStreamKeyResource(client *Client) *StreamsStreamKeyResource {
	return &StreamsStreamKeyResource{client}
}

// StreamKeyListCall represents a GET call to a Twitch StreamsStreamKey API endpoint.
type StreamKeyListCall struct {
	resource *StreamsStreamKeyResource
	opts     []RequestOption
}

// StreamKeyListResponse represents the response from a GET request to /helix/streams/key.
type StreamKeyListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the StreamKey data returned by the Twitch API.
	Data []StreamKey
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/streams/key.
//
// Gets the channel's stream key.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:stream_key scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-stream-key
func (r *StreamsStreamKeyResource) List(broadcasterID string) *StreamKeyListCall {
	c := &StreamKeyListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *StreamKeyListCall) BroadcasterID(broadcasterID string) *StreamKeyListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *StreamKeyListCall) Do(ctx context.Context, opts ...RequestOption) (*StreamKeyListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/streams/key", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[StreamKey](res)
	if err != nil {
		return nil, err
	}

	return &StreamKeyListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// StreamsFollowedResource represents the Twitch StreamsFollowed API.
type StreamsFollowedResource struct {
	client *Client
}

// NewStreamsFollowedResource creates a new StreamsFollowedResource.
func NewStreamsFollowedResource(client *Client) *StreamsFollowedResource {
	return &StreamsFollowedResource{client}
}

// StreamsFollowedListCall represents a GET call to a Twitch StreamsFollowed API endpoint.
type StreamsFollowedListCall struct {
	resource *StreamsFollowedResource
	opts     []RequestOption
}

// StreamsFollowedListResponse represents the response from a GET request to /helix/streams/followed.
type StreamsFollowedListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Stream data returned by the Twitch API.
	Data []Stream
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/streams/followed.
//
// Gets the list of broadcasters that the user follows and that are streaming live.
//
// # Authorization
//
// Requires a user access token that includes the user:read:follows scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-followed-streams
func (r *StreamsFollowedResource) List(userID string) *StreamsFollowedListCall {
	c := &StreamsFollowedListCall{resource: r}
	return c.
		UserID(userID)
}

// UserID sets the UserID query parameter.
func (api *StreamsFollowedListCall) UserID(userID string) *StreamsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", userID))
	return api
}

// After sets the After query parameter.
func (api *StreamsFollowedListCall) After(after string) *StreamsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *StreamsFollowedListCall) First(first int) *StreamsFollowedListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *StreamsFollowedListCall) Do(ctx context.Context, opts ...RequestOption) (*StreamsFollowedListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/streams/followed", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Stream](res)
	if err != nil {
		return nil, err
	}

	return &StreamsFollowedListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// StreamsMarkersResource represents the Twitch StreamsMarkers API.
type StreamsMarkersResource struct {
	client *Client
}

// NewStreamsMarkersResource creates a new StreamsMarkersResource.
func NewStreamsMarkersResource(client *Client) *StreamsMarkersResource {
	return &StreamsMarkersResource{client}
}

// StreamMarkerInsertCall represents a POST call to a Twitch StreamsMarkers API endpoint.
type StreamMarkerInsertCall struct {
	resource *StreamsMarkersResource
	body     map[string]any
}

// StreamMarkerInsertResponse represents the response from a POST request to /helix/streams/markers.
type StreamMarkerInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the StreamMarkerData data returned by the Twitch API.
	Data []StreamMarkerData
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/streams/markers.
//
// Creates a marker for a stream that is currently live.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:broadcast scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-stream-marker
func (r *StreamsMarkersResource) Insert(broadcasterID string) *StreamMarkerInsertCall {
	c := &StreamMarkerInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *StreamMarkerInsertCall) BroadcasterID(broadcasterID string) *StreamMarkerInsertCall {
	api.body["user_id"] = broadcasterID
	return api
}

// Description sets the Description body parameter.
func (api *StreamMarkerInsertCall) Description(description string) *StreamMarkerInsertCall {
	api.body["description"] = description
	return api
}

// Do executes the request.
func (api *StreamMarkerInsertCall) Do(ctx context.Context, opts ...RequestOption) (*StreamMarkerInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/streams/markers", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[StreamMarkerData](res)
	if err != nil {
		return nil, err
	}

	return &StreamMarkerInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// StreamMarkerListCall represents a GET call to a Twitch StreamsMarkers API endpoint.
type StreamMarkerListCall struct {
	resource *StreamsMarkersResource
	opts     []RequestOption
}

// StreamMarkerListResponse represents the response from a GET request to /helix/streams/markers.
type StreamMarkerListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the StreamMarker data returned by the Twitch API.
	Data []StreamMarker
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/streams/markers.
//
// Gets a list of markers from the user's most recent stream or from the specified VOD/video.
//
// A marker is an arbitrary point in a live stream that the broadcaster or editor marked, so they can return to that spot later to create video highlights (see Video Producer, Highlights in the Twitch UX).
//
// # Authorization
//
// Requires a user access token that includes the user:read:broadcast or channel:manage:broadcast scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-stream-markers
func (r *StreamsMarkersResource) List(broadcasterID string, videoID string) *StreamMarkerListCall {
	c := &StreamMarkerListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID).
		VideoID(videoID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *StreamMarkerListCall) BroadcasterID(broadcasterID string) *StreamMarkerListCall {
	api.opts = append(api.opts, SetQueryParameter("user_id", broadcasterID))
	return api
}

// VideoID sets the VideoID query parameter.
func (api *StreamMarkerListCall) VideoID(videoID string) *StreamMarkerListCall {
	api.opts = append(api.opts, SetQueryParameter("video_id", videoID))
	return api
}

// Before sets the Before query parameter.
func (api *StreamMarkerListCall) Before(before string) *StreamMarkerListCall {
	api.opts = append(api.opts, SetQueryParameter("before", before))
	return api
}

// After sets the After query parameter.
func (api *StreamMarkerListCall) After(after string) *StreamMarkerListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *StreamMarkerListCall) First(first int) *StreamMarkerListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *StreamMarkerListCall) Do(ctx context.Context, opts ...RequestOption) (*StreamMarkerListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/streams/markers", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[StreamMarker](res)
	if err != nil {
		return nil, err
	}

	return &StreamMarkerListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// StreamsListCall represents a GET call to a Twitch Streams API endpoint.
type StreamsListCall struct {
	resource *StreamsResource
	opts     []RequestOption
}

// StreamsListResponse represents the response from a GET request to /helix/streams.
type StreamsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Stream data returned by the Twitch API.
	Data []Stream
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/streams.
//
// Gets a list of all streams.
// The list is in descending order by the number of viewers watching the stream.
// Because viewers come and go during a stream, it's possible to find duplicate or missing streams in the list as you page through the results.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-streams
func (r *StreamsResource) List() *StreamsListCall {
	return &StreamsListCall{resource: r}
}

// ID adds to the ID query parameter.
func (api *StreamsListCall) ID(ids ...string) *StreamsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// UserID adds to the UserID query parameter.
func (api *StreamsListCall) UserID(userIDs ...string) *StreamsListCall {
	for _, userID := range userIDs {
		api.opts = append(api.opts, AddQueryParameter("user_id", userID))
	}
	return api
}

// GameID adds to the GameID query parameter.
func (api *StreamsListCall) GameID(gameIDs ...string) *StreamsListCall {
	for _, gameID := range gameIDs {
		api.opts = append(api.opts, AddQueryParameter("game_id", gameID))
	}
	return api
}

// Language sets the Language query parameter.
func (api *StreamsListCall) Language(language string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("language", language))
	return api
}

// Period sets the Period query parameter.
func (api *StreamsListCall) Period(period string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("period", period))
	return api
}

// Type sets the Type query parameter.
func (api *StreamsListCall) Type(t string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("type", t))
	return api
}

// Before sets the Before query parameter.
func (api *StreamsListCall) Before(before string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("before", before))
	return api
}

// After sets the After query parameter.
func (api *StreamsListCall) After(after string) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *StreamsListCall) First(first int) *StreamsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *StreamsListCall) Do(ctx context.Context, opts ...RequestOption) (*StreamsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/streams", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Stream](res)
	if err != nil {
		return nil, err
	}

	return &StreamsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

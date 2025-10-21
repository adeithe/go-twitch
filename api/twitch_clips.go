package api

import (
	"context"
	"net/http"
	"time"
)

// ClipsResource represents the Twitch Clips API.
type ClipsResource struct {
	client *Client

	// Download provides access to the Twitch Download API.
	Download *ClipsDownloadResource
}

// NewClipsResource creates a new ClipsResource.
func NewClipsResource(client *Client) *ClipsResource {
	r := &ClipsResource{client: client}
	r.Download = NewClipsDownloadResource(client)
	return r
}

// ClipsDownloadResource represents the Twitch ClipsDownload API.
type ClipsDownloadResource struct {
	client *Client
}

// NewClipsDownloadResource creates a new ClipsDownloadResource.
func NewClipsDownloadResource(client *Client) *ClipsDownloadResource {
	return &ClipsDownloadResource{client}
}

// ClipsDownloadListCall represents a GET call to a Twitch ClipsDownload API endpoint.
type ClipsDownloadListCall struct {
	resource *ClipsDownloadResource
	opts     []RequestOption
}

// ClipsDownloadListResponse represents the response from a GET request to /helix/clips/downloads.
type ClipsDownloadListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the DownloadableClip data returned by the Twitch API.
	Data []DownloadableClip
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/clips/downloads.
//
// Provides URLs to download the video file for the specified clips.
//
// # Rate Limits
//
// Limited to 100 requests per minute.
//
// # Authorization
//
// Requires an app access token or user access token that includes the editor:manage:clips or channel:manage:clips scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-clips-download
func (r *ClipsDownloadResource) List(clipID string) *ClipsDownloadListCall {
	c := &ClipsDownloadListCall{resource: r}
	return c.
		ClipID(clipID)
}

// ClipID sets the ClipID query parameter.
func (api *ClipsDownloadListCall) ClipID(clipID string) *ClipsDownloadListCall {
	api.opts = append(api.opts, SetQueryParameter("clip_id", clipID))
	return api
}

// EditorID sets the EditorID query parameter.
func (api *ClipsDownloadListCall) EditorID(editorID string) *ClipsDownloadListCall {
	api.opts = append(api.opts, SetQueryParameter("editor_id", editorID))
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ClipsDownloadListCall) BroadcasterID(broadcasterID string) *ClipsDownloadListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// Do executes the request.
func (api *ClipsDownloadListCall) Do(ctx context.Context, opts ...RequestOption) (*ClipsDownloadListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/clips/downloads", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[DownloadableClip](res)
	if err != nil {
		return nil, err
	}

	return &ClipsDownloadListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// CreateClipInsertCall represents a POST call to a Twitch Clips API endpoint.
type CreateClipInsertCall struct {
	resource *ClipsResource
	opts     []RequestOption
}

// CreateClipInsertResponse represents the response from a POST request to /helix/clips.
type CreateClipInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Clip data returned by the Twitch API.
	Data []Clip
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/clips.
//
// Creates a clip for a stream.
//
// # Authorization
//
// Requires a user access token that includes the clips:edit scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-clip
func (r *ClipsResource) Insert(broadcasterID string) *CreateClipInsertCall {
	c := &CreateClipInsertCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *CreateClipInsertCall) BroadcasterID(broadcasterID string) *CreateClipInsertCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// HasDelay sets the HasDelay query parameter.
func (api *CreateClipInsertCall) HasDelay(hasDelay bool) *CreateClipInsertCall {
	api.opts = append(api.opts, SetQueryParameter("has_delay", hasDelay))
	return api
}

// Do executes the request.
func (api *CreateClipInsertCall) Do(ctx context.Context, opts ...RequestOption) (*CreateClipInsertResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/clips", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Clip](res)
	if err != nil {
		return nil, err
	}

	return &CreateClipInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// ClipsListCall represents a GET call to a Twitch Clips API endpoint.
type ClipsListCall struct {
	resource *ClipsResource
	opts     []RequestOption
}

// ClipsListResponse represents the response from a GET request to /helix/clips.
type ClipsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Clip data returned by the Twitch API.
	Data []Clip
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/clips.
//
// Gets one or more video clips that were captured from streams.
//
// At least one of ID, GameID, and BroadcasterID are required. They are mutually exclusive.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-clips
func (r *ClipsResource) List() *ClipsListCall {
	return &ClipsListCall{resource: r}
}

// ID adds to the ID query parameter.
func (api *ClipsListCall) ID(ids ...string) *ClipsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *ClipsListCall) BroadcasterID(broadcasterID string) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// GameID sets the GameID query parameter.
func (api *ClipsListCall) GameID(gameID string) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("game_id", gameID))
	return api
}

// Before sets the Before query parameter.
func (api *ClipsListCall) Before(before string) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("before", before))
	return api
}

// After sets the After query parameter.
func (api *ClipsListCall) After(after string) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *ClipsListCall) First(first int) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// IsFeatured sets the IsFeatured query parameter.
func (api *ClipsListCall) IsFeatured(isFeatured bool) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("is_featured", isFeatured))
	return api
}

// StartedAt sets the StartedAt query parameter.
func (api *ClipsListCall) StartedAt(startedAt time.Time) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("started_at", startedAt.Format(time.RFC3339)))
	return api
}

// EndedAt sets the EndedAt query parameter.
func (api *ClipsListCall) EndedAt(endedAt time.Time) *ClipsListCall {
	api.opts = append(api.opts, SetQueryParameter("ended_at", endedAt.Format(time.RFC3339)))
	return api
}

// Do executes the request.
func (api *ClipsListCall) Do(ctx context.Context, opts ...RequestOption) (*ClipsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/clips", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Clip](res)
	if err != nil {
		return nil, err
	}

	return &ClipsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

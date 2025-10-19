package api

import (
	"context"
	"net/http"
)

// ContentLabelsResource represents the Twitch ContentLabels API.
type ContentLabelsResource struct {
	client *Client
}

// NewContentLabelsResource creates a new ContentLabelsResource.
func NewContentLabelsResource(client *Client) *ContentLabelsResource {
	return &ContentLabelsResource{client}
}

// ContentLabelsListCall represents a GET call to a Twitch ContentLabels API endpoint.
type ContentLabelsListCall struct {
	resource *ContentLabelsResource
	opts     []RequestOption
}

// ContentLabelsListResponse represents the response from a GET request to /helix/content_classification_labels.
type ContentLabelsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the ContentClassificationLabel data returned by the Twitch API.
	Data []ContentClassificationLabel
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/content_classification_labels.
//
// Gets information about Twitch content classification labels.
//
// The locale parameter is a ISO 3166-1 alpha-2 which tells Twitch which translation to use for the response. (Default: en-US)
//
// # Authorization
//
// Requires an app access token or user access token.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-content-classification-labels
func (r *ContentLabelsResource) List() *ContentLabelsListCall {
	return &ContentLabelsListCall{resource: r}
}

// Locale sets the Locale query parameter.
func (api *ContentLabelsListCall) Locale(locale string) *ContentLabelsListCall {
	api.opts = append(api.opts, SetQueryParameter("locale", locale))
	return api
}

// Do executes the request.
func (api *ContentLabelsListCall) Do(ctx context.Context, opts ...RequestOption) (*ContentLabelsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/content_classification_labels", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[ContentClassificationLabel](res)
	if err != nil {
		return nil, err
	}

	return &ContentLabelsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

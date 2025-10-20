package api

import (
	"context"
	"net/http"
)

// ExtensionsResource represents the Twitch Extensions API.
type ExtensionsResource struct {
	client *Client
}

// NewExtensionsResource creates a new ExtensionsResource.
func NewExtensionsResource(client *Client) *ExtensionsResource {
	return &ExtensionsResource{client}
}

// ExtensionsListCall represents a GET call to a Twitch Extensions API endpoint.
type ExtensionsListCall struct {
	resource *ExtensionsResource
	opts     []RequestOption
}

// ExtensionsListResponse represents the response from a GET request to /helix/extensions.
type ExtensionsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Extension data returned by the Twitch API.
	Data []Extension
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/extensions.
//
// Gets information about an extension.
//
// # Authorization
//
// Requires a signed JSON Web Token (JWT) created by an Extension Backend Service (EBS).
// The signed JWT must include the "role" field (see JWT Schema), and the "role" field must be set to external.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-extensions
func (r *ExtensionsResource) List(extensionID string) *ExtensionsListCall {
	c := &ExtensionsListCall{resource: r}
	return c.
		ExtensionID(extensionID)
}

// ExtensionID sets the ExtensionID query parameter.
func (api *ExtensionsListCall) ExtensionID(extensionID string) *ExtensionsListCall {
	api.opts = append(api.opts, SetQueryParameter("extension_id", extensionID))
	return api
}

// ExtensionVersion sets the ExtensionVersion query parameter.
func (api *ExtensionsListCall) ExtensionVersion(extensionVersion string) *ExtensionsListCall {
	api.opts = append(api.opts, SetQueryParameter("extension_version", extensionVersion))
	return api
}

// Do executes the request.
func (api *ExtensionsListCall) Do(ctx context.Context, opts ...RequestOption) (*ExtensionsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/extensions", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Extension](res)
	if err != nil {
		return nil, err
	}

	return &ExtensionsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

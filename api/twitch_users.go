package api

import (
	"context"
	"net/http"
)

// UsersResource represents the Twitch Users API.
type UsersResource struct {
	client *Client
}

// NewUsersResource creates a new UsersResource.
func NewUsersResource(client *Client) *UsersResource {
	return &UsersResource{client}
}

// UsersListCall represents a GET call to a Twitch Users API endpoint.
type UsersListCall struct {
	resource *UsersResource
	opts     []RequestOption
}

// UsersListResponse represents the response from a GET request to /helix/users.
type UsersListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the User data returned by the Twitch API.
	Data []User
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/users.
//
// Gets information about one or more users.
// You may specify users by ID or by login name.
//
// You may look up users using their user ID, login name, or both but the sum total of the number of users you may look up is 100.
// For example, you may specify 50 IDs and 50 names or 100 IDs or names, but you cannot specify 100 IDs and 100 names.
//
// If you don't specify IDs or login names, the request returns information about the user in the access token if you specify a user access token.
//
// # Authorization
//
// Requires an app access token or user access token.
//
// To include the user's verified email address in the response, you must use a user access token that includes the user:read:email scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-users
func (r *UsersResource) List() *UsersListCall {
	return &UsersListCall{resource: r}
}

// ID adds to the ID query parameter.
func (api *UsersListCall) ID(iDs ...string) *UsersListCall {
	for _, iD := range iDs {
		api.opts = append(api.opts, AddQueryParameter("id", iD))
	}
	return api
}

// Login adds to the Login query parameter.
func (api *UsersListCall) Login(logins ...string) *UsersListCall {
	for _, login := range logins {
		api.opts = append(api.opts, AddQueryParameter("login", login))
	}
	return api
}

// Do executes the request.
func (api *UsersListCall) Do(ctx context.Context, opts ...RequestOption) (*UsersListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/users", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[User](res)
	if err != nil {
		return nil, err
	}

	return &UsersListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

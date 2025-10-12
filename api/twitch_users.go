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

// UsersListCall is a call to the Users List endpoint.
type UsersListCall struct {
	resource *UsersResource
	opts     []RequestOption
}

// UsersListResponse is the response from the Users List endpoint.
type UsersListResponse struct {
	Header http.Header
	Data   []User
}

// List creates a request to list users based on the specified criteria.
//
// The email field will be empty unless the access token has the user:read:email scope.
func (r *UsersResource) List() *UsersListCall {
	return &UsersListCall{resource: r}
}

// ID filters the results to the specified user IDs.
func (c *UsersListCall) ID(ids []string) *UsersListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("id", id))
	}
	return c
}

// Login filters the results to the specified usernames.
func (c *UsersListCall) Login(logins []string) *UsersListCall {
	for _, login := range logins {
		c.opts = append(c.opts, AddQueryParameter("login", login))
	}
	return c
}

// Do executes the request.
func (c *UsersListCall) Do(ctx context.Context, opts ...RequestOption) (*UsersListResponse, error) {
	res, err := c.resource.client.DoRequest(ctx, http.MethodGet, EndpointUsers, nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[User](res)
	if err != nil {
		return nil, err
	}

	return &UsersListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

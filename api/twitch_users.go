package api

import (
	"context"
	"net/http"
	"time"
)

type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	Email           string    `json:"email,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type UsersResource struct {
	client *Client
}

func NewUsersResource(client *Client) *UsersResource {
	return &UsersResource{client}
}

type UsersListCall struct {
	resource *UsersResource
	opts     []RequestOption
}

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
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/users", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[User](res)
	if err != nil {
		return nil, err
	}

	return &UsersListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

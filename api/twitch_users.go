package api

type UsersResource struct {
	client *Client
}

func NewUsersResource(client *Client) *UsersResource {
	return &UsersResource{client}
}

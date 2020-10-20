package helix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// UserOpts stores options for requests to the Twitch Users API.
type UserOpts struct {
	IDs    []string
	Logins []string
}

// UsersData stores a list of users returned by the Twitch Users API.
type UsersData struct {
	Data []User `json:"data"`
}

// User stores data about a user on Twitch.
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	Email           string `json:"email"`
	ViewCount       int    `json:"view_count"`
	BroadcasterType string `json:"broadcaster_type"`
	Type            string `json:"type"`
}

// GetOwnUser retrieves a Twitch User object based on the clients OAuth Token.
func (client *Client) GetOwnUser() (*User, error) {
	if len(client.bearer) < 1 {
		return nil, errors.New("a bearer token is required to use this endpoint")
	}
	users, err := client.GetUsers(UserOpts{})
	if err != nil {
		return nil, err
	}
	if len(users.Data) < 1 {
		return nil, errors.New("unable to get user")
	}
	user := users.Data[0]
	client.Self = user
	return &user, nil
}

// GetUsers retrieves a list of users based on the specified UserOpts.
func (client *Client) GetUsers(opts UserOpts) (*UsersData, error) {
	if len(opts.IDs)+len(opts.Logins) > 100 {
		return nil, errors.New("you can only request a total of 100 users at a time")
	}
	params := ""
	if len(opts.IDs) > 0 {
		params += fmt.Sprintf("&id=%s", strings.Join(opts.IDs, "&id="))
	}
	if len(opts.Logins) > 0 {
		params += fmt.Sprintf("&login=%s", strings.Join(opts.Logins, "&login="))
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("users?_t=%d%s", time.Now().UTC().Unix(), params), nil)
	if err != nil {
		return nil, err
	}
	users := &UsersData{}
	if err := json.Unmarshal(res.Body, &users); err != nil {
		return nil, err
	}
	return users, nil
}

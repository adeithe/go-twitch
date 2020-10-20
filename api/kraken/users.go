package kraken

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
	Total int    `json:"_total"`
	Data  []User `json:"users"`
}

// User stores data about a user on Twitch.
type User struct {
	ID                 string `json:"_id"`
	Login              string `json:"name"`
	DisplayName        string `json:"display_name"`
	Bio                string `json:"bio"`
	Logo               string `json:"logo"`
	Email              string `json:"email"`
	IsPartner          bool   `json:"partnered"`
	IsEmailVerified    bool   `json:"email_verified"`
	IsTwitterConnected bool   `json:"twitter_connected"`
	Notifications      struct {
		Push  bool `json:"push"`
		Email bool `json:"email"`
	} `json:"notifications"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetOwnUser retrieves a user object based on the clients OAuth token.
//
// See: https://dev.twitch.tv/docs/v5/reference/users#get-user
func (client *Client) GetOwnUser() (*User, error) {
	if len(client.token) < 1 {
		return nil, errors.New("no oauth token is assigned to the client")
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("user?_t=%d", time.Now().UTC().Unix()), nil)
	if err != nil {
		return nil, err
	}
	user := &User{}
	if err := json.Unmarshal(res.Body, &user); err != nil {
		return nil, err
	}
	client.Self.User = *user
	return user, nil
}

// GetUsers retrieves a list of users based on the specified UserOpts.
//
// See: https://dev.twitch.tv/docs/v5/reference/users#get-users
func (client *Client) GetUsers(opts UserOpts) (*UsersData, error) {
	if len(opts.IDs)+len(opts.Logins) > 100 {
		return nil, errors.New("you can only request a total of 100 users at a time")
	}
	params := ""
	if len(opts.IDs) > 0 {
		params += fmt.Sprint("&id=", strings.Join(opts.IDs, ","))
	}
	if len(opts.Logins) > 0 {
		params += fmt.Sprint("&login=", strings.Join(opts.Logins, ","))
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

package api

import (
	"github.com/Adeithe/go-twitch/api/helix"
	"github.com/Adeithe/go-twitch/api/kraken"
)

// Client used to store data about a Twitch application and/or user.
type Client struct {
	ID     string
	bearer string
}

// IClient interface containing methods for the API Client.
type IClient interface {
	NewBearer(string) *Client
	Kraken() *kraken.Client
	Helix() *helix.Client

	Login(string, string) (*TwitchLogin, error)
}

// BaseURL for the Twitch API.
const BaseURL = "https://api.twitch.tv"

var _ IClient = &Client{}

// Official Twitch API Client - This should be used sparingly or not at all.
var Official *Client = &Client{ID: "kimne78kx3ncx6brgo4mv6wki5h1ko"}

// New API Client
func New(clientID string) *Client {
	return &Client{ID: clientID}
}

// NewBearer will create a copy of the Twitch API client using the provided Bearer token.
func (client Client) NewBearer(token string) *Client {
	return &Client{
		ID:     client.ID,
		bearer: token,
	}
}

// Kraken provides an interface for Twitch Kraken API endpoints.
//
// Deprecated: Twitch API v5 (Kraken) is deprecated. You should use the New Twitch API (Helix) instead.
func (client Client) Kraken() *kraken.Client {
	return kraken.New(client.ID, client.bearer)
}

// Helix provides an interface for Twitch Helix API endpoints.
func (client Client) Helix() *helix.Client {
	return helix.New(client.ID, client.bearer)
}

// Login will attempt to login via Twitch using a username/password combination.
// This attempts to mimic a user submitting the login form.
// It isn't perfect and has a good chance of failing entirely due to CAPTCHA.
func (client Client) Login(username string, password string) (*TwitchLogin, error) {
	if client.ID != Official.ID {
		return Official.Login(username, password)
	}
	login := &TwitchLogin{
		Username:  username,
		password:  password,
		ErrorCode: -1,
	}
	err := login.Verify("")
	return login, err
}

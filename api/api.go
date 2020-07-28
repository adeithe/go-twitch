package api

type IClient interface {
	NewBearer(string) *Client
	Login(string, string) (*TwitchLogin, error)
}

type Client struct {
	ClientID string
	bearer   string
}

// BaseURL for the Twitch API.
const BaseURL = "https://api.twitch.tv"

var _ IClient = &Client{}

// Official Twitch API Client - This should be used sparingly or not at all.
var Official *Client = &Client{ClientID: "kimne78kx3ncx6brgo4mv6wki5h1ko"}

// New API Client
func New(clientID string) *Client {
	return &Client{ClientID: clientID}
}

// NewBearer will create a copy of the Twitch API client using the provided Bearer token.
func (api Client) NewBearer(token string) *Client {
	return &Client{
		ClientID: api.ClientID,
		bearer:   token,
	}
}

// Login will attempt to login via Twitch using a username/password combination.
// This attempts to mimic a user submitting the login form.
// It isn't perfect and has a good chance of failing entirely due to CAPTCHA.
func (api Client) Login(username string, password string) (*TwitchLogin, error) {
	if api.ClientID != Official.ClientID {
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

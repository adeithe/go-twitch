package kraken

import (
	"encoding/json"
	"fmt"

	"github.com/Adeithe/go-twitch/api/request"
)

// Client stores an API ClientID and OAuth Token
type Client struct {
	ID    string
	token string

	Self struct {
		User    User
		Channel Channel
	}
}

// IKraken contains all methods available to the Kraken API Client.
type IKraken interface {
	Request(string, string, interface{}) (request.HTTPResponse, error)

	GetCheermotes(BitOpts) (*CheermotesData, error)

	GetOwnChannel() (*Channel, error)
	GetChannelsByID(...string) (*ChannelsData, error)

	GetTopGames(GameOpts) (*TopGames, error)

	GetIngestServers() (*IngestsData, error)

	GetStreamSummary(string) (*StreamSummary, error)
	GetStreams(StreamOpts) (*StreamsData, error)

	GetAllTeams(TeamOpts) (*TeamsData, error)
	GetTeam(string) (*TeamsData, error)

	GetOwnUser() (*User, error)
	GetUsers(UserOpts) (*UsersData, error)
}

// BaseURL is the API path that will never change.
const BaseURL = "https://api.twitch.tv/kraken"

var _ IKraken = &Client{}

// New Kraken API Client.
//
// Deprecated: Twitch API v5 (Kraken) is deprecated. You should use the New Twitch API (Helix) instead.
func New(id, token string) *Client {
	client := &Client{ID: id, token: token}
	if len(client.token) > 0 {
		client.GetOwnUser()
		client.GetOwnChannel()
	}
	return client
}

// Request Twitch Kraken Endpoints and get an HTTP response back.
func (client *Client) Request(method, path string, body interface{}) (request.HTTPResponse, error) {
	req := request.New(method, BaseURL, path)
	req.Headers["Accept"] = "application/vnd.twitchtv.v5+json"
	req.Headers["Client-ID"] = client.ID
	if len(client.token) > 0 {
		req.Headers["Authorization"] = fmt.Sprintf("OAuth %s", client.token)
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return request.HTTPResponse{}, err
	}
	req.Body = bytes
	return req.Do()
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch/api/helix"
)

// Helix client aka the New Twitch API
type Helix struct {
	Client *Client
}

// IHelix interface for the Helix Client
type IHelix interface {
	Request(string, string, interface{}) (HTTPResponse, error)

	GetUsers(helix.UserOpts) ([]helix.User, error)
	GetStreams(helix.StreamOpts) (helix.StreamData, error)
}

var _ IHelix = &Helix{}

// GetUsers retrieves data about user logins and ids as requested. Limit 100 total users.
func (api Helix) GetUsers(opts helix.UserOpts) ([]helix.User, error) {
	if (len(opts.IDs) + len(opts.Logins)) > 100 {
		return []helix.User{}, fmt.Errorf("you can only request a total of 100 IDs and Logins at one time")
	}
	params := ""
	if !opts.GetSelf {
		if len(opts.IDs) > 0 {
			params += fmt.Sprintf("&id=%s", strings.Join(opts.IDs, "&id="))
		}
		if len(opts.Logins) > 0 {
			params += fmt.Sprintf("&login=%s", strings.Join(opts.Logins, "&login="))
		}
	}
	res, err := api.Request(http.MethodGet, fmt.Sprintf("users?_t=%d%s", time.Now().Unix(), params), nil)
	if err != nil {
		return []helix.User{}, err
	}
	data := helix.Response{}
	users := []helix.User{}
	if err := json.Unmarshal(res.Body, &data); err != nil {
		return users, err
	}
	if data.Status > 0 && len(data.Error) > 0 && len(data.Message) > 0 {
		return users, fmt.Errorf("%s: %s (status: %d)", data.Error, data.Message, data.Status)
	}
	if err := json.Unmarshal(data.Data, &users); err != nil {
		return users, err
	}
	return users, err
}

// GetStreams retrieves data about streams as requested. Limit 100 channel ids.
func (api Helix) GetStreams(opts helix.StreamOpts) (helix.StreamData, error) {
	if len(opts.GameIDs) > 100 {
		return helix.StreamData{}, fmt.Errorf("you may only specify a max of 100 game ids")
	}
	if len(opts.Language) > 100 {
		return helix.StreamData{}, fmt.Errorf("you may only specify a max of 100 languages")
	}
	if len(opts.UserIDs) > 100 {
		return helix.StreamData{}, fmt.Errorf("you may only specify a max of 100 user ids")
	}
	if len(opts.UserLogins) > 100 {
		return helix.StreamData{}, fmt.Errorf("you may only specify a max of 100 user logs")
	}
	params := ""
	if opts.First > 0 {
		params += fmt.Sprintf("&first=%d", opts.First)
	}
	if len(opts.GameIDs) > 0 {
		params += fmt.Sprintf("&game_id=%s", strings.Join(opts.GameIDs, "&game_id="))
	}
	if len(opts.Language) > 0 {
		params += fmt.Sprintf("&language=%s", strings.Join(opts.Language, "&language="))
	}
	if len(opts.UserIDs) > 0 {
		params += fmt.Sprintf("&user_id=%s", strings.Join(opts.UserIDs, "&user_id="))
	}
	if len(opts.UserLogins) > 0 {
		params += fmt.Sprintf("&user_login=%s", strings.Join(opts.UserLogins, "&user_login="))
	}
	res, err := api.Request(http.MethodGet, fmt.Sprintf("streams?_t=%d%s", time.Now().Unix(), params), nil)
	if err != nil {
		return helix.StreamData{}, err
	}
	data := helix.Response{}
	streams := helix.StreamData{}
	if err := json.Unmarshal(res.Body, &data); err != nil {
		return streams, err
	}
	if data.Status > 0 && len(data.Error) > 0 && len(data.Message) > 0 {
		return streams, fmt.Errorf("%s: %s (status: %d)", data.Error, data.Message, data.Status)
	}
	if err := json.Unmarshal(res.Body, &streams); err != nil {
		return streams, err
	}
	return streams, nil
}

// Request sends an API request to the Twitch Helix API endpoint.
func (api Helix) Request(method string, path string, body interface{}) (HTTPResponse, error) {
	req := NewRequest(method, fmt.Sprintf("%s/%s", BaseURL, "helix"), path)
	req.Headers["Content-Type"] = "application/json"
	req.Headers["Client-ID"] = api.Client.ID
	if len(api.Client.bearer) > 0 {
		req.Headers["Authorization"] = fmt.Sprintf("Bearer %s", api.Client.bearer)
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return HTTPResponse{}, err
	}
	req.Body = bytes
	return req.Do()
}

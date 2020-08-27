package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
}

var _ IHelix = &Helix{}

// GetUsers retrieves data about user logins and ids as requested. Limit 100 total users.
func (api Helix) GetUsers(opts helix.UserOpts) ([]helix.User, error) {
	path := "users"
	if !opts.GetSelf {
		if (len(opts.IDs) + len(opts.Logins)) > 100 {
			return []helix.User{}, fmt.Errorf("you can only request a total of 100 IDs and Logins at one time")
		}
		path += "?"
		if len(opts.IDs) > 0 {
			path += fmt.Sprintf("id=%s", strings.Join(opts.IDs, "&id="))
		}
		if len(opts.Logins) > 0 {
			if !strings.HasSuffix(path, "?") {
				path += "&"
			}
			path += fmt.Sprintf("login=%s", strings.Join(opts.Logins, "&login="))
		}
	}
	res, err := api.Request(http.MethodGet, path, nil)
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

// Request sends an API request to the Twitch Helix API endpoint.
func (api Helix) Request(method string, path string, body interface{}) (HTTPResponse, error) {
	req := NewRequest(method, fmt.Sprintf("%s/%s", BaseURL, "helix"), path)
	req.Headers["Content-Type"] = "application/json"
	req.Headers["Client-ID"] = api.Client.ID
	if len(api.Client.bearer) > 0 {
		req.Headers["Authorization"] = fmt.Sprintf("Bearer %s", api.Client.bearer)
	}
	body, err := json.Marshal(body)
	if err != nil {
		return HTTPResponse{}, err
	}
	return req.Do()
}

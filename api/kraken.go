package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch/api/kraken"
)

// Kraken client aka Twitch API v5
type Kraken struct {
	Client *Client
}

// IKraken interface for the Kraken Client
type IKraken interface {
	Request(string, string, interface{}) (HTTPResponse, error)
	IsError([]byte) (kraken.Error, bool)

	GetUsers(kraken.UserOpts) (kraken.UserData, error)
	GetStreams(kraken.StreamOpts) (kraken.StreamData, error)
}

var _ IKraken = &Kraken{}

// GetUsers retrieves data about user logins and ids as requested. Limit 100 total users.
func (api Kraken) GetUsers(opts kraken.UserOpts) (kraken.UserData, error) {
	if (len(opts.IDs) + len(opts.Logins)) > 100 {
		return kraken.UserData{}, fmt.Errorf("you can only request a total of 100 IDs and Logins at one time")
	}
	params := ""
	if len(opts.IDs) > 0 {
		params += fmt.Sprintf("&id=%s", strings.Join(opts.IDs, ","))
	}
	if len(opts.Logins) > 0 {
		params += fmt.Sprintf("&login=%s", strings.Join(opts.Logins, ","))
	}
	res, err := api.Request(http.MethodGet, fmt.Sprintf("users?_t=%d%s", time.Now().Unix(), params), nil)
	if err != nil {
		return kraken.UserData{}, err
	}
	users := kraken.UserData{}
	if data, ok := api.IsError(res.Body); !ok {
		return users, fmt.Errorf("%s: %s (status: %d)", data.Error, data.Message, data.Status)
	}
	if err := json.Unmarshal(res.Body, &users); err != nil {
		return users, err
	}
	return users, nil
}

// GetStreams retrieves data about streams as requested. Limit 100 channel ids.
func (api Kraken) GetStreams(opts kraken.StreamOpts) (kraken.StreamData, error) {
	if len(opts.ChannelIDs) > 100 {
		return kraken.StreamData{}, fmt.Errorf("you can only request a total of 100 channel ids at one time")
	}
	params := ""
	if len(opts.ChannelIDs) > 0 {
		params += fmt.Sprintf("&channel=%s", strings.Join(opts.ChannelIDs, ","))
	}
	if len(opts.Game) > 0 {
		params += fmt.Sprintf("&game=%s", opts.Game)
	}
	if len(opts.Language) > 0 {
		params += fmt.Sprintf("&language=%s", opts.Language)
	}
	if len(opts.Type) > 0 {
		params += fmt.Sprintf("&stream_type=%s", opts.Type)
	}
	if opts.Limit > 0 {
		params += fmt.Sprintf("&limit=%d", opts.Limit)
	}
	if opts.Offset > 0 {
		params += fmt.Sprintf("&offset=%d", opts.Offset)
	}
	res, err := api.Request(http.MethodGet, fmt.Sprintf("streams?_t=%d%s", time.Now().Unix(), params), nil)
	if err != nil {
		return kraken.StreamData{}, err
	}
	streams := kraken.StreamData{}
	if data, ok := api.IsError(res.Body); !ok {
		return streams, fmt.Errorf("%s: %s (status: %d)", data.Error, data.Message, data.Status)
	}
	if err := json.Unmarshal(res.Body, &streams); err != nil {
		return streams, err
	}
	streams.Total = len(streams.Data)
	return streams, nil
}

// IsError returns error data for a request if available.
func (api Kraken) IsError(bytes []byte) (kraken.Error, bool) {
	var data kraken.Error
	if err := json.Unmarshal(bytes, &data); err != nil {
		data = kraken.Error{
			Error:   "internal",
			Status:  500,
			Message: "failed to parse response data",
		}
		return data, false
	}
	if data.Status > 0 && len(data.Error) > 0 && len(data.Message) > 0 {
		return data, false
	}
	return data, true
}

// Request sends an API request to the Twitch Kraken API endpoint.
func (api Kraken) Request(method string, path string, body interface{}) (HTTPResponse, error) {
	req := NewRequest(method, fmt.Sprintf("%s/%s", BaseURL, "kraken"), path)
	req.Headers["Accept"] = "application/vnd.twitchtv.v5+json"
	req.Headers["Client-ID"] = api.Client.ID
	if len(api.Client.bearer) > 0 {
		req.Headers["Authorization"] = fmt.Sprintf("OAuth %s", api.Client.bearer)
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return HTTPResponse{}, err
	}
	req.Body = bytes
	return req.Do()
}

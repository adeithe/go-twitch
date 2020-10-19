package helix

import (
	"encoding/json"
	"fmt"

	"github.com/Adeithe/go-twitch/api/request"
)

type Client struct {
	ID     string
	bearer string
}

type IHelix interface {
	Request(string, string, interface{}) (request.HTTPResponse, error)
}

// BaseURL is the API path that will never change.
const BaseURL = "https://api.twitch.tv/helix"

var _ IHelix = &Client{}

// New Helix API Client.
func New(id, bearer string) *Client {
	return &Client{ID: id, bearer: bearer}
}

// Request Twitch Helix Endpoints and get an HTTP response back.
func (client Client) Request(method, path string, body interface{}) (request.HTTPResponse, error) {
	req := request.New(method, BaseURL, path)
	req.Headers["Client-ID"] = client.ID
	if len(client.bearer) > 0 {
		req.Headers["Authorization"] = fmt.Sprintf("Bearer %s", client.bearer)
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return request.HTTPResponse{}, err
	}
	req.Body = bytes
	return req.Do()
}

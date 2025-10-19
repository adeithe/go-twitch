package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// WhispersResource represents the Twitch Whispers API.
type WhispersResource struct {
	client *Client
}

// NewWhispersResource creates a new WhispersResource.
func NewWhispersResource(client *Client) *WhispersResource {
	return &WhispersResource{client}
}

// SendWhisperInsertCall represents a POST call to a Twitch Whispers API endpoint.
type SendWhisperInsertCall struct {
	resource *WhispersResource
	body     map[string]any
	opts     []RequestOption
}

// SendWhisperInsertResponse represents the response from a POST request to /helix/whispers.
type SendWhisperInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/whispers.
//
// Sends a whisper message to the specified user.
//
// # Rate Limits
//
// You may whisper to a maximum of 40 unique recipients per day.
// Within the per day limit, you may whisper a maximum of 3 whispers per second and a maximum of 100 whispers per minute.
//
// # Authorization
//
// The user sending the whisper must have a verified phone number (see the Phone Number setting in your Security and Privacy settings).
//
// Requires a user access token that includes the user:manage:whispers scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#send-whisper
func (r *WhispersResource) Insert(fromUserID string, toUserID string, message string) *SendWhisperInsertCall {
	c := &SendWhisperInsertCall{resource: r, body: make(map[string]any)}
	return c.
		FromUserID(fromUserID).
		ToUserID(toUserID).
		Message(message)
}

// FromUserID sets the FromUserID query parameter.
func (api *SendWhisperInsertCall) FromUserID(fromUserID string) *SendWhisperInsertCall {
	api.opts = append(api.opts, SetQueryParameter("from_user_id", fromUserID))
	return api
}

// ToUserID sets the ToUserID query parameter.
func (api *SendWhisperInsertCall) ToUserID(toUserID string) *SendWhisperInsertCall {
	api.opts = append(api.opts, SetQueryParameter("to_user_id", toUserID))
	return api
}

// Message sets the Message body parameter.
func (api *SendWhisperInsertCall) Message(message string) *SendWhisperInsertCall {
	api.body["message"] = message
	return api
}

// Do executes the request.
func (api *SendWhisperInsertCall) Do(ctx context.Context, opts ...RequestOption) (*SendWhisperInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/whispers", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	if err != nil {
		return nil, err
	}

	return &SendWhisperInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Request:    res.Request,
	}, nil
}

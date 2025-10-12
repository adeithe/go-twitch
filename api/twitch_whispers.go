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

// WhispersInsertCall is a call to the Whispers Insert endpoint.
type WhispersInsertCall struct {
	resource *WhispersResource
	opts     []RequestOption
	body     map[string]any
}

// Insert creates a new WhispersInsertCall to send a whisper from senderID to recipientID.
func (r *WhispersResource) Insert(senderID, recipientID string) *WhispersInsertCall {
	c := &WhispersInsertCall{resource: r, body: make(map[string]any)}
	c.opts = append(c.opts, SetQueryParameter("from_user_id", senderID))
	c.opts = append(c.opts, SetQueryParameter("to_user_id", recipientID))
	return c
}

// Message the whisper message to send. The message must not be empty.
//
// The maximum message lengths are:
//   - 500 characters if the user you're sending the message to hasn't whispered you before.
//   - 10,000 characters if the user you're sending the message to has whispered you before.
//
// Messages that exceed the maximum length are truncated.
func (c *WhispersInsertCall) Message(message string) *WhispersInsertCall {
	c.body["message"] = message
	return c
}

// Do executes the request.
//
//	req := client.Whispers.SendWhisper("123", "456").Message("Hello")
//	data, err := req.Do(ctx, api.WithBearerToken("kpvy3cjboyptmdkiacwr0c19hotn5s")
func (c *WhispersInsertCall) Do(ctx context.Context, opts ...RequestOption) error {
	bs, err := json.Marshal(c.body)
	if err != nil {
		return err
	}

	res, err := c.resource.client.DoRequest(ctx, http.MethodPost, EndpointWhispers, bytes.NewReader(bs), opts...)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	_, err = decodeResponse[any](res)
	return err
}

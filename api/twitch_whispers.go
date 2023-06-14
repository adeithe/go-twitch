package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type WhispersResource struct {
	client *Client
}

func NewWhispersResource(client *Client) *WhispersResource {
	return &WhispersResource{client}
}

type SendWhisperRequest struct {
	resource   *WhispersResource
	fromUserId string
	toUserId   string
	message    string
}

func (r *WhispersResource) SendWhisper(senderId, recipientId string) *SendWhisperRequest {
	return &SendWhisperRequest{r, senderId, recipientId, ""}
}

// SenderID sets the sender ID for the request.
//
// This ID must match the user ID in the user access token and the account must have a verified phone number.
func (c *SendWhisperRequest) SenderID(senderId string) *SendWhisperRequest {
	c.fromUserId = senderId
	return c
}

// RecipientID the ID of the user to receive the whisper.
func (c *SendWhisperRequest) RecipientID(recipientId string) *SendWhisperRequest {
	c.toUserId = recipientId
	return c
}

// Message the whisper message to send. The message must not be empty.
//
// The maximum message lengths are:
//   - 500 characters if the user you're sending the message to hasn't whispered you before.
//   - 10,000 characters if the user you're sending the message to has whispered you before.
//
// Messages that exceed the maximum length are truncated.
func (c *SendWhisperRequest) Message(message string) *SendWhisperRequest {
	c.message = message
	return c
}

// Do executes the request.
//
//	req := client.Whispers.SendWhisper("123", "456").Message("Hello")
//	data, err := req.Do(ctx, api.WithBearerToken("kpvy3cjboyptmdkiacwr0c19hotn5s")
func (c *SendWhisperRequest) Do(ctx context.Context, opts ...RequestOption) error {
	bs, err := json.Marshal(map[string]any{"message": c.message})
	if err != nil {
		return err
	}

	query := url.Values{}
	query.Set("from_user_id", c.fromUserId)
	query.Set("to_user_id", c.toUserId)
	res, err := c.resource.client.doRequest(ctx, http.MethodPost, fmt.Sprintf("/whispers?%s", query.Encode()), bytes.NewReader(bs), opts...)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = decodeResponse[any](res)
	return err
}

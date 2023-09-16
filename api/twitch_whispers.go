package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type WhispersResource struct {
	client *Client
}

func NewWhispersResource(client *Client) *WhispersResource {
	return &WhispersResource{client}
}

type WhispersInsertCall struct {
	resource *WhispersResource
	opts     []RequestOption
	message  string
}

func (r *WhispersResource) Insert(senderId, recipientId string) *WhispersInsertCall {
	c := &WhispersInsertCall{resource: r}
	c.opts = append(c.opts, SetQueryParameter("from_user_id", senderId))
	c.opts = append(c.opts, SetQueryParameter("to_user_id", recipientId))
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
	c.message = message
	return c
}

// Do executes the request.
//
//	req := client.Whispers.SendWhisper("123", "456").Message("Hello")
//	data, err := req.Do(ctx, api.WithBearerToken("kpvy3cjboyptmdkiacwr0c19hotn5s")
func (c *WhispersInsertCall) Do(ctx context.Context, opts ...RequestOption) error {
	bs, err := json.Marshal(map[string]any{"message": c.message})
	if err != nil {
		return err
	}

	res, err := c.resource.client.doRequest(ctx, http.MethodPost, "/whispers", bytes.NewReader(bs), opts...)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = decodeResponse[any](res)
	return err
}

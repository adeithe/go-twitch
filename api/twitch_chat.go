package api

type ChatResource struct {
	client *Client
}

func NewChatResource(client *Client) *ChatResource {
	return &ChatResource{client}
}

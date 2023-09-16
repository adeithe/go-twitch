package api

type EntitlementsResource struct {
	client *Client
}

func NewEntitlementsResource(client *Client) *EntitlementsResource {
	return &EntitlementsResource{client}
}

package api

type BitsResource struct {
	client *Client
}

func NewBitsResource(client *Client) *BitsResource {
	return &BitsResource{client}
}

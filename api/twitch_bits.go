package api

type BitsResource struct {
	client *Client

	Cheermotes *CheermotesResource
}

func NewBitsResource(client *Client) *BitsResource {
	r := &BitsResource{client: client}
	r.Cheermotes = NewCheermotesResource(client)
	return r
}

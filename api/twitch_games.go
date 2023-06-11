package api

type GamesResource struct {
	client *Client
}

func NewGamesResource(client *Client) *GamesResource {
	return &GamesResource{client}
}

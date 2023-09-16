package api

type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
	IGDB      string `json:"igdb_id"`
}

type GamesResource struct {
	client *Client

	Top *TopGamesResource
}

func NewGamesResource(client *Client) *GamesResource {
	c := &GamesResource{client: client}
	c.Top = NewTopGamesResource(client)
	return c
}

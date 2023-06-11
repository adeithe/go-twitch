package api

type TeamsResource struct {
	client *Client
}

func NewTeamsResource(client *Client) *TeamsResource {
	return &TeamsResource{client}
}

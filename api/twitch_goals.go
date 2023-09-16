package api

type GoalsResource struct {
	client *Client
}

func NewGoalsResource(client *Client) *GoalsResource {
	return &GoalsResource{client}
}

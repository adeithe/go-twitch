package graphql

import (
	"context"
	"net/http"

	"github.com/Adeithe/go-twitch/api"
	"github.com/shurcooL/graphql"
)

// Client stores data about a GraphQL client
type Client struct {
	ID     string
	bearer string

	graphql *graphql.Client
}

// URL is the address for the GraphQL server
const URL = "https://gql.twitch.tv/gql"

// New Twitch GraphQL Client
//
// This uses the official Twitch client by default and therefore should be used sparingly or not at all.
func New() (client *Client) {
	client = &Client{ID: api.Official.ID}
	client.graphql = graphql.NewClient(URL, &http.Client{Transport: httpTransport{client, http.DefaultTransport}})
	return
}

// SetBearer sets the token sent with GraphQL requests
func (client *Client) SetBearer(token string) {
	client.bearer = token
}

// CustomQuery executes a query on the GraphQL server
//
// See: https://github.com/shurcooL/graphql
func (client Client) CustomQuery(query interface{}, vars map[string]interface{}) error {
	return client.graphql.Query(context.Background(), query, vars)
}

// CustomMutation executes a mutation on the GraphQL server
//
// See: https://github.com/shurcooL/graphql
func (client Client) CustomMutation(mutation interface{}, vars map[string]interface{}) error {
	return client.graphql.Mutate(context.Background(), mutation, vars)
}

// GetUsersByID retrieves an array of users from Twitch based on their User IDs
func (client *Client) GetUsersByID(ids ...string) ([]User, error) {
	if len(ids) > 100 {
		return []User{}, ErrTooManyArguments
	}
	users := QueryUserIDs{}
	vars := map[string]interface{}{"ids": toIDs(ids...)}
	if err := client.CustomQuery(&users, vars); err != nil {
		return []User{}, err
	}
	return users.Data, nil
}

// GetUsersByLogin retrieves an array of users from Twitch based on their usernames
func (client *Client) GetUsersByLogin(logins ...string) ([]User, error) {
	if len(logins) > 100 {
		return []User{}, ErrTooManyArguments
	}
	users := QueryUserLogins{}
	vars := map[string]interface{}{"logins": toStrings(logins...)}
	if err := client.CustomQuery(&users, vars); err != nil {
		return []User{}, err
	}
	return users.Data, nil
}

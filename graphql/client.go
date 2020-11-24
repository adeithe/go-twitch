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

// IsUsernameAvailable returns true if the provided username is not taken on Twitch
func (client *Client) IsUsernameAvailable(username string) (bool, error) {
	user := GQLUsernameAvailabilityQuery{}
	vars := map[string]interface{}{"username": graphql.String(username)}
	if err := client.CustomQuery(&user, vars); err != nil {
		return false, err
	}
	return user.IsAvailable, nil
}

// GetCurrentUser retrieves the current user based on the clients authentication token
func (client Client) GetCurrentUser() (*User, error) {
	if len(client.bearer) < 1 {
		return nil, ErrTokenNotSet
	}
	user := GQLCurrentUserQuery{}
	if err := client.CustomQuery(&user, nil); err != nil {
		return nil, err
	}
	return user.Data, nil
}

// GetUsersByID retrieves an array of users from Twitch based on their User IDs
func (client Client) GetUsersByID(ids ...string) ([]User, error) {
	if len(ids) > 100 {
		return []User{}, ErrTooManyArguments
	}
	users := GQLUserIDsQuery{}
	vars := map[string]interface{}{"ids": toIDs(ids...)}
	if err := client.CustomQuery(&users, vars); err != nil {
		return []User{}, err
	}
	return users.Data, nil
}

// GetUsersByLogin retrieves an array of users from Twitch based on their usernames
func (client Client) GetUsersByLogin(logins ...string) ([]User, error) {
	if len(logins) > 100 {
		return []User{}, ErrTooManyArguments
	}
	users := GQLUserLoginsQuery{}
	vars := map[string]interface{}{"logins": toStrings(logins...)}
	if err := client.CustomQuery(&users, vars); err != nil {
		return []User{}, err
	}
	return users.Data, nil
}

// GetChannelsByID retrieves an array of channels from Twitch based on their IDs
func (client Client) GetChannelsByID(ids ...string) ([]Channel, error) {
	if len(ids) > 100 {
		return []Channel{}, ErrTooManyArguments
	}
	channels := GQLChannelIDsQuery{}
	vars := map[string]interface{}{"ids": toIDs(ids...)}
	if err := client.CustomQuery(&channels, vars); err != nil {
		return []Channel{}, err
	}
	return channels.Data, nil
}

// GetChannelsByName retrieves an array of channels from Twitch based on their names
func (client Client) GetChannelsByName(names ...string) ([]Channel, error) {
	if len(names) > 100 {
		return []Channel{}, ErrTooManyArguments
	}
	channels := GQLChannelNamesQuery{}
	vars := map[string]interface{}{"names": toStrings(names...)}
	if err := client.CustomQuery(&channels, vars); err != nil {
		return []Channel{}, err
	}
	return channels.Data, nil
}

// GetStreams retrieves data about streams available on Twitch
func (client Client) GetStreams(opts StreamQueryOpts) (*StreamsQuery, error) {
	if opts.First < 1 || opts.First > 100 {
		opts.First = 25
	}
	streams := GQLStreamsQuery{}
	vars := map[string]interface{}{
		"first":   graphql.Int(opts.First),
		"after":   opts.After,
		"options": opts.Options,
	}
	if err := client.CustomQuery(&streams, vars); err != nil {
		return nil, err
	}
	return streams.Data, nil
}

// GetGames retrieves data about games available on Twitch
func (client Client) GetGames(opts GameQueryOpts) (*GamesQuery, error) {
	if opts.First < 1 || opts.First > 100 {
		opts.First = 25
	}
	games := GQLGamesQuery{}
	vars := map[string]interface{}{
		"first":   graphql.Int(opts.First),
		"after":   opts.After,
		"options": opts.Options,
	}
	if err := client.CustomQuery(&games, vars); err != nil {
		return nil, err
	}
	return games.Data, nil
}

package graphql

import "github.com/shurcooL/graphql"

// QueryUserIDs GraphQL query for 1-100 Users by their IDs
type QueryUserIDs struct {
	Data []User `graphql:"users(ids: $ids)"`
}

// QueryUserLogins GraphQL query for 1-100 Users by their usernames
type QueryUserLogins struct {
	Data []User `graphql:"users(logins: $logins)"`
}

func toIDs(ids ...string) (gids []graphql.ID) {
	for _, id := range ids {
		gids = append(gids, id)
	}
	return
}

func toStrings(strings ...string) (strs []graphql.String) {
	for _, str := range strings {
		strs = append(strs, graphql.String(str))
	}
	return
}

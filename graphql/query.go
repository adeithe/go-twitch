package graphql

import (
	"time"

	"github.com/shurcooL/graphql"
)

// GQLUsernameAvailabilityQuery GraphQL query to check if a username is available on Twitch
type GQLUsernameAvailabilityQuery struct {
	IsAvailable bool `graphql:"isUsernameAvailable(username: $username)"`
}

// GQLCurrentUserQuery GraphQL query for current authenticated user
type GQLCurrentUserQuery struct {
	Data *User `graphql:"currentUser"`
}

// GQLUserIDsQuery GraphQL query for 1-100 Users by their IDs
type GQLUserIDsQuery struct {
	Data []User `graphql:"users(ids: $ids)"`
}

// GQLUserLoginsQuery GraphQL query for 1-100 Users by their usernames
type GQLUserLoginsQuery struct {
	Data []User `graphql:"users(logins: $logins)"`
}

// GQLChannelIDsQuery GraphQL query for 1-100 Channels by their IDs
type GQLChannelIDsQuery struct {
	Data []Channel `graphql:"users(ids: $ids)"`
}

// GQLChannelNamesQuery GraphQL query for 1-100 Channels by name
type GQLChannelNamesQuery struct {
	Data []Channel `graphql:"users(logins: $names)"`
}

// GQLStreamsQuery GraphQL query for streams on Twitch
type GQLStreamsQuery struct {
	Data *StreamsQuery `graphql:"streams(first: $first, after: $after, options: $options)"`
}

// GQLUserVideosQuery GraphQL query for a users videos on Twitch
type GQLUserVideosQuery struct {
	Data *UserVideosQuery `graphql:"user(id: $id)"`
}

// GQLClipQuery GraphQL query for a clip on Twitch
type GQLClipQuery struct {
	Data *Clip `graphql:"clip(slug: $slug)"`
}

// GQLGamesQuery GraphQL query for games on Twitch
type GQLGamesQuery struct {
	Data *GamesQuery `graphql:"games(first: $first, after: $after, options: $options)"`
}

// GQLFollowersQuery GraphQL query for a users followers on Twitch
type GQLFollowersQuery struct {
	Data *FollowersQuery `graphql:"user(id: $id)"`
}

// StreamsQuery stores data returned from GQLStreamsQuery
type StreamsQuery struct {
	ResponseID   graphql.ID `graphql:"responseID"`
	GenerationID graphql.ID `graphql:"generationID"`
	Data         []struct {
		TrackingID graphql.ID `graphql:"trackingID"`
		Stream     Stream     `graphql:"node"`
		Cursor     Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
}

// FollowersQuery stores data returned from GQLFollowersQuery
type FollowersQuery struct {
	Followers struct {
		TotalCount int32
		Data       []struct {
			User       User `graphql:"node"`
			FollowedAt time.Time
			Cursor     Cursor
		} `graphql:"edges"`
		PageInfo PageInfo
	} `graphql:"followers(first: $first, after: $after)"`
}

// UserVideosQuery stores data returned from GQLUserVideosQuery
type UserVideosQuery struct {
	Videos struct {
		TotalCount int32
		Data       []struct {
			Video  Video `graphql:"node"`
			Cursor Cursor
		} `graphql:"edges"`
	} `graphql:"videos(first: $first, after: $after, sort: TIME)"`
}

// GamesQuery stores data returned from GQLGamesQuery
type GamesQuery struct {
	Data []struct {
		TrackingID graphql.ID `graphql:"trackingID"`
		Game       Game       `graphql:"node"`
		Cursor     Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
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

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

// GQLVideosQuery GraphQL query for videos on Twitch
type GQLVideosQuery struct {
	Data *VideosQuery `graphql:"videos(first: $first, after: $after)"`
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
	Data *struct {
		Followers *FollowersQuery `graphql:"followers(first: $first, after: $after)"`
	} `graphql:"user(id: $id)"`
}

// GQLModsQuery GraphQL query for getting mods for a user on Twitch
type GQLModsQuery struct {
	Data *struct {
		Mods *ModsQuery `graphql:"mods(first: $first, after: $after)"`
	} `graphql:"user(id: $id)"`
}

// GQLVIPsQuery GraphQL query for getting VIPs for a user on Twitch
type GQLVIPsQuery struct {
	Data *struct {
		VIPs *VIPsQuery `graphql:"vips(first: $first, after: $after)"`
	} `graphql:"user(id: $id)"`
}

// GQLUserVideosQuery GraphQL query for a users videos on Twitch
type GQLUserVideosQuery struct {
	Data *struct {
		Videos *VideosQuery `graphql:"videos(first: $first, after: $after, sort: TIME)"`
	} `graphql:"user(id: $id)"`
}

// StreamsQuery stores data returned from GQLStreamsQuery
type StreamsQuery struct {
	ResponseID   graphql.ID `graphql:"responseID"`
	GenerationID graphql.ID `graphql:"generationID"`
	Streams      []struct {
		TrackingID graphql.ID `graphql:"trackingID"`
		Stream     Stream     `graphql:"node"`
		Cursor     Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
}

// ModsQuery stores data returned from GQLUserModsQuery and GQLChannelModsQuery
type ModsQuery struct {
	Mods []struct {
		User      User `graphql:"node"`
		IsActive  bool
		GrantedAt time.Time
		Cursor    Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
}

// VIPsQuery stores data returned from GQLUserVIPsQuery and GQLChannelVIPsQuery
type VIPsQuery struct {
	VIPs []struct {
		User      User `graphql:"node"`
		GrantedAt time.Time
		Cursor    Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
}

// FollowersQuery stores data returned from GQLFollowersQuery
type FollowersQuery struct {
	TotalCount int32
	Followers  []struct {
		User       User `graphql:"node"`
		FollowedAt time.Time
		Cursor     Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
}

// VideosQuery stores data returned from GQLUserVideosQuery
type VideosQuery struct {
	TotalCount int32
	Videos     []struct {
		Video  Video `graphql:"node"`
		Cursor Cursor
	} `graphql:"edges"`
	PageInfo PageInfo
}

// GamesQuery stores data returned from GQLGamesQuery
type GamesQuery struct {
	Games []struct {
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

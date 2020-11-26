package graphql

// Cursor pagaination string
type Cursor string

// StreamQueryOpts stores various options for querying streams on Twitch
type StreamQueryOpts struct {
	First   int32
	After   Cursor
	Options StreamOptions
}

// StreamOptions stores various options for querying streams on Twitch
type StreamOptions struct {
	Locale string
	Tags   []string
}

// GameQueryOpts stores various options for querying games on Twitch
type GameQueryOpts struct {
	First   int32
	After   Cursor
	Options GameOptions
}

// GameOptions stores various options for querying games on Twitch
type GameOptions struct {
	Locale string
	Tags   []string
}

// FollowOpts stores various options for querying followers on Twitch
type FollowOpts struct {
	First int32
	After Cursor
}

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

// VideoQueryOpts stores various options for querying videos on Twitch
type VideoQueryOpts struct {
	First int32
	After Cursor
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

// FollowQueryOpts stores various options for querying followers on Twitch
type FollowQueryOpts struct {
	First int32
	After Cursor
}

// ModsQueryOpts stores various options for querying mods on Twitch
type ModsQueryOpts struct {
	First int32
	After Cursor
}

// VIPsQueryOpts stores various options for querying vips on Twitch
type VIPsQueryOpts struct {
	First int32
	After Cursor
}

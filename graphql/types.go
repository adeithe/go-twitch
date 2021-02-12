package graphql

import (
	"errors"
	"time"

	"github.com/shurcooL/graphql"
)

var (
	// ErrInvalidArgument returned when an argument is invalid
	ErrInvalidArgument = errors.New("one or more arguments are invalid")
	// ErrTokenNotSet returned when a method requires an authorization token but no token is set
	ErrTokenNotSet = errors.New("missing authorization token")
	// ErrTooManyArguments returned when a method receives more arguments than allowed by the GraphQL server
	ErrTooManyArguments = errors.New("too many arguments provided")
)

// PageInfo stores data about available uses for a Cursor
type PageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
}

// User stores data about a user on Twitch
type User struct {
	ID               graphql.ID
	Login            string
	DisplayName      string
	ChannelURL       string `graphql:"profileURL"`
	BannerImageURL   string `graphql:"bannerImageURL"`
	OfflineImageURL  string `graphql:"offlineImageURL"`
	ChatColor        string
	Description      string
	ProfileViewCount int32
	Stream           *Stream
	Hosting          *struct {
		Channel
		Stream Stream
	}
	Roles struct {
		IsAffiliate           bool
		IsPartner             bool
		IsExtensionsDeveloper bool
		IsGlobalMod           bool
		IsStaff               bool
		IsSiteAdmin           bool
	}
	HasPrime  bool
	HasTurbo  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Stream stores data about a livestream on Twitch
type Stream struct {
	ID                  graphql.ID
	ViewersCount        int32
	Channel             *Channel `graphql:"broadcaster"`
	Game                *Game
	BroadcasterSoftware string
	AverageFPS          float64 `graphql:"averageFPS"`
	Bitrate             float64
	Codec               string
	Type                string
	CreatedAt           time.Time
	UpdatedAt           time.Time `graphql:"lastUpdatedAt"`
}

// Channel stores data about a users channel on Twitch
type Channel struct {
	ID                graphql.ID
	Name              string `graphql:"login"`
	DisplayName       string
	URL               string `graphql:"profileURL"`
	BannerImageURL    string `graphql:"bannerImageURL"`
	OfflineImageURL   string `graphql:"offlineImageURL"`
	ChatColor         string
	Description       string
	ProfileViewCount  int32
	BroadcastSettings struct {
		Title            string
		Game             *Game
		Language         string
		IsMature         bool
		LiveNotification *struct {
			IsDefault bool
			Text      string `graphql:"liveUpNotification"`
		} `graphql:"liveUpNotificationInfo"`
	}
	ChatSettings struct {
		Rules                          []string
		IsFastSubsModeEnabled          bool
		IsLinksBlocked                 bool  `graphql:"blockLinks"`
		IsVerifiedAccountRequired      bool  `graphql:"requireVerifiedAccount"`
		IsSubOnly                      bool  `graphql:"isSubscribersOnlyModeEnabled"`
		IsEmoteOnly                    bool  `graphql:"isEmoteOnlyModeEnabled"`
		IsUniqueModeEnabled            bool  `graphql:"isUniqueChatModeEnabled"`
		SlowModeDurationInSeconds      int32 `graphql:"slowModeDurationSeconds"`
		FollowersOnlyDurationInMinutes int32 `graphql:"followersOnlyDurationMinutes"`
		ChatDelayInMilliseconds        int32 `graphql:"chatDelayMs"`
	}
	Roles struct {
		IsAffiliate           bool
		IsPartner             bool
		IsExtensionsDeveloper bool
		IsGlobalMod           bool
		IsStaff               bool
		IsSiteAdmin           bool
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Video stores data about a streams VOD on Twitch.
type Video struct {
	ID                 string
	Title              string
	Description        string
	ViewCount          int32
	AnimatedPreviewURL string `graphql:"animatedPreviewURL"`
	SeekPreviewsURL    string `graphql:"seekPreviewsURL"`
	Owner              User
	Game               Game
	Tags               []string
	OffsetSeconds      int32
	LengthInSeconds    int32 `graphql:"lengthSeconds"`
	PublishedAt        time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// Clip stores data about a section of a Twitch stream.
type Clip struct {
	ID                 string
	Slug               string
	Title              string
	ViewCount          int32
	Broadcaster        User
	Author             User `graphql:"curator"`
	Game               Game
	Video              *Video
	URL                string
	DurationSeconds    int32
	VideoOffsetSeconds int32
	CreatedAt          time.Time
}

// Game stores data about a category on Twitch
type Game struct {
	ID                graphql.ID
	Name              string
	DisplayName       string
	BroadcastersCount int32
	ViewersCount      int32
	FollowersCount    int32
	PopularityScore   int32
	GiantBombID       graphql.ID `graphql:"giantBombID"`
	PrestoID          graphql.ID `graphql:"prestoID"`
}

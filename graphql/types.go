package graphql

import (
	"errors"
	"time"

	"github.com/shurcooL/graphql"
)

var (
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
	BannerImageURL   string `graphql:"bannerImageURL"`
	OfflineImageURL  string `graphql:"offlineImageURL"`
	ChatColor        string
	Description      string
	ProfileViewCount int32
	Channel          *Channel
	Stream           *Stream
	Hosting          *struct {
		ID          graphql.ID
		Login       string
		DisplayName string
		Stream      Stream
		IsPartner   bool
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	HasPrime bool
	HasTurbo bool
	Roles    struct {
		IsAffiliate           bool
		IsPartner             bool
		IsExtensionsDeveloper bool
		IsExtensionsApprover  bool
		IsGlobalMod           bool
		IsStaff               bool
		IsSiteAdmin           bool
	}
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Stream stores data about a livestream on Twitch
type Stream struct {
	ID                  graphql.ID
	Title               string
	ViewersCount        int32
	Channel             *Channel
	Game                *Game
	ClipCount           int32
	BroadcasterSoftware string
	DelayLengthSeconds  int32
	AverageFPS          float64 `graphql:"averageFPS"`
	Bitrate             float64
	Codec               string
	Type                string
	CreatedAt           time.Time
	UpdatedAt           time.Time `graphql:"lastUpdatedAt"`
}

// Channel stores data about a users channel on Twitch
type Channel struct {
	ID          graphql.ID
	Name        string
	DisplayName string
	Game        string
	URL         string
	Settings    struct {
		Broadcast struct {
			Title            string
			Game             *Game
			Language         string
			IsMature         bool
			LiveNotification *struct {
				IsDefault bool
				Text      string `graphql:"liveUpNotification"`
			} `graphql:"liveUpNotificationInfo"`
		} `graphql:"broadcastSettings"`
		Chat struct {
			Rules                             []string
			AutoModLevel                      int32
			IsFastSubsModeEnabled             bool
			IsOptedOutOfGlobalBannedWordsList bool
			IsBroadcasterLanguageModeEnabled  bool
			IsLinksBlocked                    bool  `graphql:"blockLinks"`
			IsVerifiedAccountRequired         bool  `graphql:"requireVerifiedAccount"`
			IsSubOnly                         bool  `graphql:"isSubscribersOnlyModeEnabled"`
			IsEmoteOnly                       bool  `graphql:"isEmoteOnlyModeEnabled"`
			IsUniqueModeEnabled               bool  `graphql:"isUniqueChatModeEnabled"`
			SlowModeDurationInSeconds         int32 `graphql:"slowModeDurationSeconds"`
			FollowersOnlyDurationInMinutes    int32 `graphql:"followersOnlyDurationMinutes"`
			ChatDelayInMilliseconds           int32 `graphql:"chatDelayMs"`
		} `graphql:"chatSettings"`
	} `graphql:"owner"`
	IsClipsEnabled    bool
	IsOptedOutOfDrops bool
}

// Game stores data about a category on Twitch
type Game struct {
	ID              graphql.ID
	Name            string
	DisplayName     string
	ChannelsCount   int32
	ViewersCount    int32
	FollowersCount  int32
	PopularityScore int32
	GiantBombID     graphql.ID `graphql:"giantBombID"`
	PrestoID        graphql.ID `graphql:"prestoID"`
}

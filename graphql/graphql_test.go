package graphql

import (
	"fmt"
	"testing"
)

const (
	testUserID    = "44322889"
	testUserLogin = "dallas"
)

func TestQueryUsersByID(t *testing.T) {
	gql := New()
	users, err := gql.GetUsersByID(testUserID)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user got %d", len(users))
	}
	user := users[0]
	if user.Login != testUserLogin {
		t.Fatalf("expected %s got %s", testUserLogin, user.Login)
	}
	t.Logf("got %d users", len(users))
}

func TestQueryUsersByLogin(t *testing.T) {
	gql := New()
	users, err := gql.GetUsersByLogin(testUserLogin)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user got %d", len(users))
	}
	user := users[0]
	if fmt.Sprint(user.ID) != testUserID {
		t.Fatalf("expected %s got %s", testUserID, user.ID)
	}
	t.Logf("got %d users", len(users))
}

func TestQueryChannelsByID(t *testing.T) {
	gql := New()
	channels, err := gql.GetChannelsByID(testUserID)
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 {
		t.Fatalf("expected 1 channel got %d", len(channels))
	}
	channel := channels[0]
	if channel.Name != testUserLogin {
		t.Fatalf("expected %s got %s", testUserLogin, channel.Name)
	}
	t.Logf("got %d channels", len(channels))
}

func TestQueryChannelsByName(t *testing.T) {
	gql := New()
	channels, err := gql.GetChannelsByName(testUserLogin)
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 {
		t.Fatalf("expected 1 channel got %d", len(channels))
	}
	channel := channels[0]
	if fmt.Sprint(channel.ID) != testUserID {
		t.Fatalf("expected %s got %s", testUserID, channel.ID)
	}
	t.Logf("got %d channels", len(channels))
}

func TestQueryStreams(t *testing.T) {
	gql := New()
	streams, err := gql.GetStreams(StreamQueryOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if streams == nil {
		t.Fatal("expected at least 1 stream got 0")
	}
	if len(streams.Data) < 1 {
		t.Fatalf("expected at least 1 stream got %d", len(streams.Data))
	}
	t.Logf("got %d streams", len(streams.Data))
}

func TestQueryGames(t *testing.T) {
	gql := New()
	games, err := gql.GetGames(GameQueryOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if games == nil {
		t.Fatal("expected at least 1 game got 0")
	}
	if len(games.Data) < 1 {
		t.Fatalf("expected at least 1 game got %d", len(games.Data))
	}
	t.Logf("got %d games", len(games.Data))
}

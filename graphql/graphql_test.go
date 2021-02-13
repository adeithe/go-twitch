package graphql

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	envUsername = os.Getenv("TWITCH_USERNAME")
	envToken    = os.Getenv("TWITCH_TOKEN")
)

func TestClient(t *testing.T) {
	gql := New()
	gql.ID = ""
	gql.SetBearer("abcd123")
	if gql.bearer != "abcd123" {
		t.Fatal("bearer token was not set successfully")
	}
	if _, err := gql.GetStreams(StreamQueryOpts{First: 1}); err == nil {
		t.Fatal("expected error did not occur")
	}
}

func TestErrors(t *testing.T) {
	gql := New()
	var args []string
	for i := 0; i < 101; i++ {
		args = append(args, fmt.Sprint(i))
	}
	if _, err := gql.GetUsersByID(args...); err == nil {
		t.Fatal("GetUsersByID didnt return an error")
	}
	if _, err := gql.GetUsersByLogin(args...); err == nil {
		t.Fatal("GetUsersByLogin didnt return an error")
	}
	if _, err := gql.GetChannelsByID(args...); err == nil {
		t.Fatal("GetChannelsByID didnt return an error")
	}
	if _, err := gql.GetChannelsByName(args...); err == nil {
		t.Fatal("GetChannelsByName didnt return an error")
	}
}

func TestIsUsernameAvailable(t *testing.T) {
	var available bool
	var err error
	gql := New()
	tries := 5
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < tries; i++ {
		chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789_")
		username := make([]rune, 20)
		for i := range username {
			username[i] = chars[rand.Intn(len(chars))]
		}
		available, _ = gql.IsUsernameAvailable(string(username))
		if available {
			tries = i
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	if !available {
		t.Fatalf("failed to find an available username after %d tries", tries)
	}
	t.Logf("found an available username after %d tries", tries+1)
}

func TestQueryUsers(t *testing.T) {
	gql := New()
	if _, err := gql.GetFollowersForUser(User{}, FollowQueryOpts{}); err == nil {
		t.Fatalf("expected error did not occur")
	}
	users, err := gql.GetUsersByID("44322889")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("expected: 1 got: %d", len(users))
	}
	users, err = gql.GetUsersByLogin(users[0].Login)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("expected: 1 got: %d", len(users))
	}
	if _, err := gql.GetFollowersForUser(users[0], FollowQueryOpts{}); err != nil {
		t.Fatal(err)
	}
}

func TestQueryChannels(t *testing.T) {
	gql := New()
	if _, err := gql.GetFollowersForChannel(Channel{}, FollowQueryOpts{}); err == nil {
		t.Fatalf("expected error did not occur")
	}
	channels, err := gql.GetChannelsByID("44322889")
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 {
		t.Fatalf("expected: 1 got: %d", len(channels))
	}
	channels, err = gql.GetChannelsByName(channels[0].Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 {
		t.Fatalf("expected: 1 got: %d", len(channels))
	}
	if _, err := gql.GetFollowersForChannel(channels[0], FollowQueryOpts{}); err != nil {
		t.Fatal(err)
	}
}

func TestQueryStreams(t *testing.T) {
	gql := New()
	data, err := gql.GetStreams(StreamQueryOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("streams query returned nil")
	}
	if len(data.Streams) < 1 {
		t.Fatalf("expected at least 1 stream got %d", len(data.Streams))
	}
	t.Logf("got %d streams", len(data.Streams))
}

func TestQueryGames(t *testing.T) {
	gql := New()
	data, err := gql.GetGames(GameQueryOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("games query returned nil")
	}
	if len(data.Games) < 1 {
		t.Fatalf("expected at least 1 game got %d", len(data.Games))
	}
	t.Logf("got %d games", len(data.Games))
}

func TestAuthenticated(t *testing.T) {
	if len(envUsername) < 1 {
		t.Skipf("TWITCH_USERNAME is not set. Skipping...")
	}
	if len(envToken) < 1 {
		t.Skipf("TWITCH_TOKEN is not set. Skipping...")
	}
	gql := New()
	gql.SetBearer(envToken)
	user, err := gql.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}
	if user.Login != strings.ToLower(envUsername) {
		t.Fatal("returned user was invalid")
	}
}

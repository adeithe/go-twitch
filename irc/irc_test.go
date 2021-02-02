package irc

import (
	"encoding/hex"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	envUsername = os.Getenv("TWITCH_USERNAME")
	envToken    = os.Getenv("TWITCH_TOKEN")

	testChannel = "dallas"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		in   string
		want error
	}{
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;user-id=12345;user-type=", ErrPartialMessage},
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;user-id=12345;user-type= :tmi.twitch.tv", ErrNoCommand},
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;user-id=12345;user-type= :tmi.twitch.tv GLOBALUSERSTATE", nil},
		{":testuser!testuser@testuser.tmi.twitch.tv JOIN #channel", nil},
		{"@emote-only=0;followers-only=-1;r9k=0;rituals=0;room-id=123;slow=0;subs-only=1 :tmi.twitch.tv ROOMSTATE #channel", nil},
		{"@badge-info=;badges=;color=;display-name=TestUser;emotes=;flags=;id=abcd123-0123-4abc-defg-1234567890;mod=0;room-id=123;subscriber=0;tmi-sent-ts=1612256273447;turbo=0;user-id=12345;user-type= :testuser!testuser@testuser.tmi.twitch.tv PRIVMSG #channel :this is a message", nil},
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;mod=0;subscriber=0;user-type= :tmi.twitch.tv USERSTATE #channel", nil},
		{"@login=testuser;room-id=;target-msg-id=abcd123-0123-4abc-defg-1234567890;tmi-sent-ts=1612256431213 :tmi.twitch.tv CLEARMSG #sodapoppin :this is a message", nil},
		{"@ban-duration=60;room-id=123;target-user-id=12345;tmi-sent-ts=1612256572313 :tmi.twitch.tv CLEARCHAT #channel :testuser", nil},
		{"@room-id=123;target-user-id=12345;tmi-sent-ts=1612256572313 :tmi.twitch.tv CLEARCHAT #channel :testuser", nil},
		{"@msg-id=msg_banned :tmi.twitch.tv NOTICE #channel :You are permanently banned from talking in channel.", nil},
		{":testuser@testuser.tmi.twitch.tv PART #channel", nil},
		{":tmi.twitch.tv RECONNECT", nil},
		{":tmi.twitch.tv PING", nil},
	}
	conn := &Conn{}
	for i, test := range tests {
		msg, err := NewParsedMessage(test.in)
		if err != nil {
			if test.want != nil {
				if test.want.Error() != err.Error() {
					t.Fatalf("Simulated line #%d failed, got: %s, want: %s", i, err, test.want)
				}
				continue
			}
			t.Fatalf("Simulated line #%d failed, got: %s, want: <nil>", i, err)
		}
		conn.handle(msg)
	}
}

func TestShardedConnection(t *testing.T) {
	client := New()
	client.SetMaxChannelsPerShard(0)
	if err := client.Join(testChannel); err != nil {
		t.Fatal(err)
	}
	if _, ok := client.GetChannel(testChannel); ok {
		if err := client.Leave(testChannel); err != nil {
			t.Fatal(err)
		}
	}
	if err := client.Join(testChannel); err != nil {
		t.Fatal(err)
	}
	if err := client.Leave(testChannel); err != nil {
		t.Fatal(err)
	}
	client.Close()
}

func TestSingleConnection(t *testing.T) {
	conn := &Conn{}
	if err := conn.Join(testChannel); err != nil {
		t.Fatal(err)
	}
	if err := conn.Reconnect(); err != nil {
		t.Fatal(err)
	}
	if !conn.IsConnected() {
		t.Fatal("reconnect failed")
	}
	latency, err := conn.Ping()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("set username to %s", conn.Username)
	t.Logf("latency of %dms", latency.Milliseconds())
	if err := conn.Leave(testChannel); err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestAuthenticatedConnection(t *testing.T) {
	if len(envUsername) < 1 {
		t.Skipf("TWITCH_USERNAME is not set. Skipping...")
	}
	if len(envToken) < 1 {
		t.Skipf("TWITCH_TOKEN is not set. Skipping...")
	}

	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		t.Fatal(err)
	}
	message := hex.EncodeToString(bytes)
	r := make(chan bool, 1)
	c := make(chan bool, 1)

	reader := New()
	reader.OnShardChannelUpdate(func(shardID int, msg RoomState) {
		r <- true
	})
	reader.OnShardMessage(func(shardID int, msg ChatMessage) {
		c <- msg.Text == message
	})
	if err := reader.Join(envUsername); err != nil {
		t.Fatal(err)
	}

	select {
	case <-r:
	case <-time.After(time.Second * 30):
		t.Fatal("failed to prepare chatroom reader")
	}

	writer := Conn{}
	writer.SetLogin(envUsername, envToken)
	if err := writer.Connect(); err != nil {
		t.Fatal(err)
	}
	if err := writer.Sayf(envUsername, "%s", message); err != nil {
		t.Fatal(err)
	}

	var success bool
	select {
	case success = <-c:
		if success {
			t.Log("verified message sent by writer")
		}
	case <-time.After(time.Second * 5):
	}
	if !success {
		t.Fatal("failed to verify that message was sent")
	}

	reader.Close()
	writer.Close()
}

package irc

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
	"time"
)

const (
	testChannel = "dallas"
)

func TestShardedConnection(t *testing.T) {
	client := New()
	client.SetMaxChannelsPerShard(0)
	if err := client.Join(testChannel); err != nil {
		t.Fatal(err)
	}
	if client.IsInChannel(testChannel) {
		if err := client.Leave(testChannel); err != nil {
			t.Fatal(err)
		}
	}
	if err := client.Join(testChannel); err != nil {
		t.Fatal(err)
	}
	client.Close()
}

func TestSingleConnection(t *testing.T) {
	conn := Conn{}
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
	t.Logf("latency of %dms", latency.Milliseconds())
	conn.Close()
}

func TestAuthenticatedConnection(t *testing.T) {
	username := os.Getenv("TWITCH_USERNAME")
	if len(username) < 1 {
		t.Skipf("TWITCH_USERNAME is not set. Skipping...")
	}
	token := os.Getenv("TWITCH_TOKEN")
	if len(token) < 1 {
		t.Skipf("TWITCH_TOKEN is not set. Skipping...")
	}

	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		t.Fatal(err)
	}
	message := hex.EncodeToString(bytes)
	c := make(chan bool, 1)

	reader := New()
	reader.OnShardMessage(func(shardID int, msg ChatMessage) {
		c <- msg.Message == message
	})
	if err := reader.Join(username); err != nil {
		t.Fatal(err)
	}

	writer := Conn{}
	writer.SetLogin(username, token)
	if err := writer.Say(username, message); err != nil {
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

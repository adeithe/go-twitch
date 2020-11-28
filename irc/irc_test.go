package irc

import "testing"

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

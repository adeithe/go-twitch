package irc

import (
	"context"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := &Conn{}
	if err := conn.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	latency, err := conn.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for msg := range conn.RawMessage {
			t.Log(msg.Raw)
		}
	}()

	t.Logf("set username to %s", conn.Username)
	t.Logf("latency of %dms", latency.Milliseconds())
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

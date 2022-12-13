package irc_test

import (
	"context"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestIRC_Events_OnRawMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ch := make(chan *irc.RawMessage, 100)
	conn, err := irc.New(irc.OnRawMessageChan(ch))
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
		for msg := range ch {
			if msg.Command == irc.CMDReady {
				break
			}
		}
		assert.NoError(t, conn.Close())
	}
}

func TestIRC_Events_OnDisconnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ch := make(chan *irc.Conn, 1)
	conn, err := irc.New(irc.OnDisconnectChan(ch))
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	for i := 0; i < 2; i++ {
		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.NoError(t, conn.Close())
			select {
			case <-ctx.Done():
				t.Fail()
				return
			case <-ch:
			}
		}
	}
}

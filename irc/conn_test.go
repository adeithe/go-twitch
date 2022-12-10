package irc_test

import (
	"context"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestIRC_Ping(t *testing.T) {
	t.Run("Connected", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.False(t, conn.IsConnected())

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.True(t, conn.IsConnected())

			latency, err := conn.Ping(ctx)
			assert.NoError(t, err)
			assert.NotZero(t, latency.Milliseconds())

			assert.NoError(t, conn.Close())
		}
	})

	t.Run("ErrNotConnected", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.False(t, conn.IsConnected())

		_, err = conn.Ping(ctx)
		assert.ErrorIs(t, err, irc.ErrNotConnected)
	})

	t.Run("Context Cancelled", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.False(t, conn.IsConnected())

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			pingCtx, pingCancel := context.WithCancel(ctx)
			pingCancel()

			_, err := conn.Ping(pingCtx)
			assert.ErrorIs(t, err, context.Canceled)
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("Context Deadline Exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.False(t, conn.IsConnected())

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			pingCtx, pingCancel := context.WithTimeout(ctx, time.Nanosecond)
			defer pingCancel()
			<-pingCtx.Done()

			_, err := conn.Ping(pingCtx)
			assert.ErrorIs(t, err, context.DeadlineExceeded)
			assert.NoError(t, conn.Close())
		}
	})
}

func TestIRC_Connect_WithAuth(t *testing.T) {
	t.Run("Connected", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New(irc.WithAuth("justinfan1234", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"))
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.False(t, conn.IsConnected())

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.NoError(t, conn.Connect(ctx)) // Already connected.
			assert.True(t, conn.IsConnected())
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("ErrLoginFailed", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New(irc.WithAuth("username", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"))
		assert.NoError(t, err)
		assert.NotNil(t, conn)

		assert.ErrorIs(t, irc.EnsureConnection(ctx, conn), irc.ErrLoginFailed)
	})
}

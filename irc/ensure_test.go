package irc_test

import (
	"context"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestIRC_EnsureConnection(t *testing.T) {
	t.Run("TLS", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotPanics(t, func() { conn.Close() })

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("NoTLS", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New(irc.WithAddress(irc.DefaultHostname, 6667), irc.WithoutTLS())
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotPanics(t, func() { conn.Close() })

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("Insecure", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := irc.New(irc.WithInsecure())
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotPanics(t, func() { conn.Close() })

		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("Context Canceled", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		cancel()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotPanics(t, func() { conn.Close() })

		assert.ErrorIs(t, irc.EnsureConnection(ctx, conn), context.Canceled)
	})

	t.Run("Context Deadline Exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		<-ctx.Done()

		conn, err := irc.New()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotPanics(t, func() { conn.Close() })

		assert.ErrorIs(t, irc.EnsureConnection(ctx, conn), context.DeadlineExceeded)
	})

	t.Run("ErrNilConnection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		assert.ErrorIs(t, irc.EnsureConnection(ctx, nil), irc.ErrNilConnection)
	})
}

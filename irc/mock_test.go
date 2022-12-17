package irc_test

import (
	"context"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestIRC_Mock_RunT(t *testing.T) {
	t.Run("Authenticated", func(t *testing.T) {
		mock := irc.RunT(t)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mock.AddUsers("username", "jtv")
		t.Logf("Mock IRC Server: %s", mock.Addr().String())
		conn, err := irc.New(irc.WithoutTLS(), mock.WithAddress(), irc.WithAuth("username", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"))
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.True(t, conn.IsConnected())
			assert.NoError(t, mock.Write("@emote-only=0;followers-only=-1;r9k=0;room-id=84316241;slow=0;subs-only=0 :tmi.twitch.tv ROOMSTATE #unknown"))

			channel, err := conn.JoinChannel("jtv")
			assert.NoError(t, err)
			assert.NotNil(t, channel)
			assert.NoError(t, channel.Sayf("foobar"))

			_, err = conn.JoinChannel("fakeuser")
			assert.ErrorIs(t, err, irc.ErrJoinFailed)
			assert.Contains(t, err.Error(), "This channel does not exist or has been suspended")

			channel, isJoined := conn.GetChannel("unknown")
			assert.True(t, isJoined)
			assert.NotNil(t, channel)
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("Anonymous", func(t *testing.T) {
		mock := irc.RunT(t)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mock.AddUsers("username", "jtv")
		t.Logf("Mock IRC Server: %s", mock.Addr().String())
		conn, err := irc.New(irc.WithoutTLS(), mock.WithAddress())
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.True(t, conn.IsConnected())

			channel, err := conn.JoinChannel("jtv")
			assert.NoError(t, err)
			assert.NotNil(t, channel)
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("Ping", func(t *testing.T) {
		mock := irc.RunT(t)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mock.AddUsers("username")
		t.Logf("Mock IRC Server: %s", mock.Addr().String())
		conn, err := irc.New(irc.WithoutTLS(), mock.WithAddress(), irc.WithAuth("username", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"))
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.True(t, conn.IsConnected())
			assert.NoError(t, mock.Write("PING :tmi.twitch.tv"))

			_, err := conn.Ping(ctx)
			assert.NoError(t, err)
			assert.NoError(t, conn.Close())
		}
	})

	t.Run("Login Failed", func(t *testing.T) {
		mock := irc.RunT(t)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		t.Logf("Mock IRC Server: %s", mock.Addr().String())
		conn, err := irc.New(irc.WithoutTLS(), mock.WithAddress(), irc.WithAuth("username", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"))
		assert.NoError(t, err)
		assert.NotNil(t, conn)

		err = irc.EnsureConnection(ctx, conn)
		assert.ErrorIs(t, err, irc.ErrLoginFailed)
		assert.Contains(t, err.Error(), "Login authentication failed")
	})

	t.Run("Invalid Message", func(t *testing.T) {
		mock := irc.RunT(t)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mock.AddUsers("username")
		t.Logf("Mock IRC Server: %s", mock.Addr().String())
		conn, err := irc.New(irc.WithoutTLS(), mock.WithAddress(), irc.WithAuth("username", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"))
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
			assert.True(t, conn.IsConnected())
			assert.NoError(t, conn.SendRaw(": JOIN #jtv invalid message"))
			assert.NoError(t, mock.Write(": JOIN #jtv invalid message"))
			assert.NoError(t, conn.Close())
		}
	})
}

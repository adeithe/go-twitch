package irc_test

import (
	"context"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestIRC_Channel_Join(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := irc.New()
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	if assert.NoError(t, irc.EnsureConnection(ctx, conn)) {
		channel, err := conn.JoinChannel("jtv")
		assert.NoError(t, err)
		assert.NotNil(t, channel)

		channel2, ok := conn.GetChannel("jtv")
		assert.True(t, ok)
		if assert.NotNil(t, channel2) {
			assert.True(t, channel.IsJoined())
			assert.True(t, channel.IsSubOnly())
			assert.False(t, channel.IsEmoteOnly())
			assert.False(t, channel.IsR9KMode())
			assert.Equal(t, channel.Username(), channel2.Username())
			assert.Equal(t, channel.RoomID(), channel2.RoomID())
			assert.NoError(t, conn.PartChannel(channel.Username()))
		}
		assert.NoError(t, conn.Close())
	}
}

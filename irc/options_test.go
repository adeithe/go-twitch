package irc_test

import (
	"testing"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestIRC_Options_WithAuth(t *testing.T) {
	t.Run("Username", func(t *testing.T) {
		_, err := irc.New(irc.WithAuth("", ""))
		assert.ErrorIs(t, err, irc.ErrInvalidField)
		assert.Contains(t, err.Error(), "username")
	})

	t.Run("Token", func(t *testing.T) {
		_, err := irc.New(irc.WithAuth("justinfan123", ""))
		assert.ErrorIs(t, err, irc.ErrInvalidField)
		assert.Contains(t, err.Error(), "token")
	})
}

func TestIRC_Options_WithAddress(t *testing.T) {
	t.Run("Host", func(t *testing.T) {
		_, err := irc.New(irc.WithAddress("", 6667))
		assert.ErrorIs(t, err, irc.ErrInvalidField)
		assert.Contains(t, err.Error(), "host")
	})

	t.Run("Port", func(t *testing.T) {
		_, err := irc.New(irc.WithAddress("irc.chat.twitch.tv", 0))
		assert.ErrorIs(t, err, irc.ErrInvalidField)
		assert.Contains(t, err.Error(), "port")
	})
}

func TestIRC_Options_WithHostname(t *testing.T) {
	t.Run("ErrInvalidField", func(t *testing.T) {
		_, err := irc.New(irc.WithHostname(""))
		assert.ErrorIs(t, err, irc.ErrInvalidField)
		assert.Contains(t, err.Error(), "hostname")
	})
}

func TestIRC_Options_WithBufferSize(t *testing.T) {
	t.Run("ErrInvalidField", func(t *testing.T) {
		_, err := irc.New(irc.WithBufferSize(0))
		assert.ErrorIs(t, err, irc.ErrInvalidField)
		assert.Contains(t, err.Error(), "bufferSize")
	})
}

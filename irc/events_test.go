package irc_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestEventHandlers(t *testing.T) {
	tests := []struct {
		raw    string
		events *irc.Events
		f      func(t *testing.T, events *irc.Events)
	}{
		{
			":tmi.twitch.tv 001 justinfan15365 :Welcome, GLHF!",
			&irc.Events{
				RawMessage: make(chan *irc.Message, 1),
			},
			func(t *testing.T, events *irc.Events) {
				assert.Equal(t, 1, len(events.RawMessage))
			},
		},
		{
			":tmi.twitch.tv 376 justinfan16432 :>",
			&irc.Events{
				Ready: make(chan struct{}, 1),
			},
			func(t *testing.T, events *irc.Events) {
				assert.Equal(t, 1, len(events.Ready))
			},
		},
		{
			"PING :tmi.twitch.tv",
			&irc.Events{
				Latency: make(chan time.Duration, 1),
			},
			func(t *testing.T, events *irc.Events) {
				assert.Equal(t, 1, len(events.Latency))
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			mock, err := NewMockServer(t, func(c net.Conn, m *irc.Message) {
				_, _ = c.Write([]byte(test.raw + "\r\n"))
			})
			_ = mock.Conn(test.events).Connect(ctx)

			assert.NoError(t, err)
			test.f(t, test.events)
		})
	}
}

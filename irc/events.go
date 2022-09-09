package irc

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type Events struct {
	Ready      chan struct{}
	Latency    chan time.Duration
	RawMessage chan *Message
}

func emit[T any](c chan T, v T) {
	if c != nil {
		if cap(c) > 0 && len(c) == cap(c) {
			log.Warn().Msg("buffered channel is full, dropping oldest event")
			<-c
		}
		c <- v
	}
}

func (c *Conn) handleMessage(msg *Message) {
	switch msg.Command {
	case CMDReady:
		emit(c.events.Ready, struct{}{})
	case CMDPing:
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, _ = c.Ping(ctx)
		}()
	case CMDPong:
		close(c.pingC)

	case CMDRoomState:
	case CMDJoin:
	case CMDPart:

	case CMDGlobalUserState:
	case CMDUserState:

	case CMDHostTarget:
	case CMDUserNotice:
	case CMDClearChat:
	case CMDClearMessage:
	case CMDNotice:
	case CMDPrivMessage:
	}
}

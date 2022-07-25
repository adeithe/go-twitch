package irc

import (
	"context"
	"time"
)

type Events struct {
	Ready      chan struct{}
	Latency    chan time.Duration
	RawMessage chan *Message
}

func emit[T any](c chan T, v T) {
	if c != nil {
		c <- v
	}
}

func (c *Conn) handleMessage(msg *Message) {
	switch msg.Command {
	case CMDReady:
		emit(c.events.Ready, struct{}{})
	case CMDPing:
		_, _ = c.Ping(context.Background())
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

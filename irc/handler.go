package irc

import (
	"context"
	"time"
)

func (conn *Conn) handle(msg Message) {
	switch msg.Command {
	case CMDReady:
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn.Ping(ctx)
		}()
	case CMDReconnect:

	case CMDPing:
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn.Ping(ctx)
		}()
	case CMDPong:
		close(conn.pingC)

	case CMDRoomState:
	case CMDJoin:
	case CMDPart:

	case CMDGlobalUserState:
		conn.UserState = NewGlobalUserState(msg)
	case CMDUserState:

	case CMDHostTarget:

	case CMDUserNotice:
		//notice := NewUserNotice(msg)

	case CMDClearChat:
		//ban := NewChatBan(msg)
	case CMDClearMessage:
		//delete := NewChatMessageDelete(msg)

	case CMDNotice:
		notice := NewServerNotice(msg)
		conn.handleServerNotice(notice)

	case CMDPrivMessage:
		//msg := NewChatMessage(msg)
	}
}

func (conn *Conn) handleServerNotice(notice ServerNotice) {
	if conn.ServerNotice == nil {
		conn.ServerNotice = make(chan ServerNotice, 1024)
	}
	if cap(conn.ServerNotice) > 0 && len(conn.ServerNotice) == cap(conn.ServerNotice) {
		<-conn.ServerNotice
	}
	conn.ServerNotice <- notice
}

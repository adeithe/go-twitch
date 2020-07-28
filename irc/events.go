package irc

import "time"

// OnReady event called when the IRC client is ready and logged in.
func (irc *Client) OnReady(handler func()) {
	irc.onReady = append(irc.onReady, handler)
}

// OnDisconnect event called when the IRC client disconnects from the server.
func (irc *Client) OnDisconnect(handler func()) {
	irc.onDisconnect = append(irc.onDisconnect, handler)
}

// OnPing event called when the IRC client sends a Ping.
func (irc *Client) OnPing(handler func()) {
	irc.onPing = append(irc.onPing, handler)
}

// OnPong event called when the IRC client receives a Ping response.
func (irc *Client) OnPong(handler func(latency time.Duration)) {
	irc.onPong = append(irc.onPong, handler)
}

// OnMessage event called when the IRC client receives a chat message.
func (irc *Client) OnMessage(handler func(message ChatMessage)) {
	irc.onMessage = append(irc.onMessage, handler)
}

// OnRawMessage event called when the IRC client receives a message.
func (irc *Client) OnRawMessage(handler func(message IRCMessage)) {
	irc.onRawMessage = append(irc.onRawMessage, handler)
}

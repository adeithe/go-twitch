package irc

type Events struct {
	OnServerNotice     []func(*Notice)
	OnChannelJoined    []func(*Channel)
	OnChannelNotice    []func(*Notice)
	OnChatMessage      []func(*ChatMessage)
	OnRawMessage       []func(*RawMessage)
	OnReconnectRequest []func(*Conn)
	OnDisconnect       []func(*Conn)
}

// OnServerNotice registers a function to be called when the client receives a server notice.
func OnServerNotice(f func(*Notice)) ConnOption {
	return func(c *Conn) error {
		c.events.OnServerNotice = append(c.events.OnServerNotice, f)
		return nil
	}
}

// OnServerNoticeChan registers a channel to be sent new server notices received by the client.
func OnServerNoticeChan(ch chan<- *Notice) ConnOption {
	return OnServerNotice(func(c *Notice) { ch <- c })
}

// OnChannelJoined registers a function to be called when the client successfully joins a channel.
func OnChannelJoined(f func(*Channel)) ConnOption {
	return func(c *Conn) error {
		c.events.OnChannelJoined = append(c.events.OnChannelJoined, f)
		return nil
	}
}

// OnChannelJoinedChan registers a channel to be sent new chatrooms joined by the client.
func OnChannelJoinedChan(ch chan<- *Channel) ConnOption {
	return OnChannelJoined(func(c *Channel) { ch <- c })
}

// OnChannelNotice registers a function to be called when the client receives a channel notice.
func OnChannelNotice(f func(*Notice)) ConnOption {
	return func(c *Conn) error {
		c.events.OnChannelNotice = append(c.events.OnChannelNotice, f)
		return nil
	}
}

// OnChannelNoticeChan registers a channel to be sent new channel notices received by the client.
func OnChannelNoticeChan(ch chan<- *Notice) ConnOption {
	return OnChannelNotice(func(c *Notice) { ch <- c })
}

// OnChatMessage registers a function to be called when a message is received for a channel.
func OnChatMessage(f func(*ChatMessage)) ConnOption {
	return func(c *Conn) error {
		c.events.OnChatMessage = append(c.events.OnChatMessage, f)
		return nil
	}
}

// OnChatMessageChan registers a channel to be sent new messages received for a channel.
func OnChatMessageChan(ch chan<- *ChatMessage) ConnOption {
	return OnChatMessage(func(c *ChatMessage) { ch <- c })
}

// OnRawMessage registers a function to be called when a raw message is received from the server.
func OnRawMessage(f func(*RawMessage)) ConnOption {
	return func(c *Conn) error {
		c.events.OnRawMessage = append(c.events.OnRawMessage, f)
		return nil
	}
}

// OnRawMessageChan registers a channel to be sent raw messages received from the server.
func OnRawMessageChan(ch chan<- *RawMessage) ConnOption {
	return OnRawMessage(func(msg *RawMessage) { ch <- msg })
}

// OnReconnectRequest registers a function to be called when the server requests that the client reconnects.
func OnReconnectRequest(f func(*Conn)) ConnOption {
	return func(c *Conn) error {
		c.events.OnReconnectRequest = append(c.events.OnReconnectRequest, f)
		return nil
	}
}

// OnReconnectRequestChan registers a channel to be sent the connection when the server requests that the client reconnects.
func OnReconnectRequestChan(ch chan<- *Conn) ConnOption {
	return OnReconnectRequest(func(c *Conn) { ch <- c })
}

// OnDisconnect registers a function to be called when the connection to the server is lost.
func OnDisconnect(f func(*Conn)) ConnOption {
	return func(c *Conn) error {
		c.events.OnDisconnect = append(c.events.OnDisconnect, f)
		return nil
	}
}

// OnDisconnectChan registers a channel to be sent the connection when the connection to the server is lost.
func OnDisconnectChan(ch chan<- *Conn) ConnOption {
	return OnDisconnect(func(c *Conn) { ch <- c })
}

func doEvent[T any](v []func(T)) func(T) {
	return func(a T) {
		for _, fn := range v {
			fn(a)
		}
	}
}

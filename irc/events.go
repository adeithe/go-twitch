package irc

type Events struct {
	OnRawMessage []func(*RawMessage)
	OnDisconnect []func(*Conn)
}

// OnRawMessage registers a function to be called when a raw message is received from the IRC server.
func OnRawMessage(f func(*RawMessage)) ConnOption {
	return func(c *Conn) error {
		c.events.OnRawMessage = append(c.events.OnRawMessage, f)
		return nil
	}
}

// OnRawMessageChan registers a channel to be sent raw messages from the IRC server.
func OnRawMessageChan(ch chan<- *RawMessage) ConnOption {
	return OnRawMessage(func(msg *RawMessage) { ch <- msg })
}

// OnDisconnect registers a function to be called when the connection to the IRC server is lost.
func OnDisconnect(f func(*Conn)) ConnOption {
	return func(c *Conn) error {
		c.events.OnDisconnect = append(c.events.OnDisconnect, f)
		return nil
	}
}

// OnDisconnectChan registers a channel to be sent the connection when the connection to the IRC server is lost.
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

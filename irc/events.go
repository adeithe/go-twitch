package irc

import "sync"

type Events struct {
	OnRawMessage []func(*RawMessage)
	OnDisconnect []func(*Conn)
}

func OnRawMessage(f func(*RawMessage)) ConnOption {
	return func(c *Conn) error {
		c.events.OnRawMessage = append(c.events.OnRawMessage, f)
		return nil
	}
}

func OnRawMessageChan(ch chan<- *RawMessage) ConnOption {
	return OnRawMessage(func(msg *RawMessage) { ch <- msg })
}

func OnDisconnect(f func(*Conn)) ConnOption {
	return func(c *Conn) error {
		c.events.OnDisconnect = append(c.events.OnDisconnect, f)
		return nil
	}
}

func OnDisconnectChan(ch chan<- *Conn) ConnOption {
	return OnDisconnect(func(c *Conn) { ch <- c })
}

func doEvent[T any](v []func(T)) func(T) {
	var wg sync.WaitGroup
	return func(a T) {
		for _, fn := range v {
			wg.Add(1)
			go func(fn func(T)) {
				defer wg.Done()
				fn(a)
			}(fn)
		}
		wg.Wait()
	}
}

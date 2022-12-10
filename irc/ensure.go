package irc

import (
	"context"
	"errors"
	"time"
)

var ErrNilConnection = errors.New("irc: nil connection")

// EnsureConnection ensures that the connection is established before returning using an exponential backoff between failures.
//
// Returns an error if the context is cancelled or the deadline is exceeded before a connection is established or if authentication fails.
func EnsureConnection(ctx context.Context, conn *Conn) error {
	if conn == nil {
		return ErrNilConnection
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for backoff := time.Second; ; backoff *= 2 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := conn.Connect(ctx); err != nil {
				if errors.Is(err, ErrLoginFailed) {
					return err
				}
				ticker.Reset(backoff)
				continue
			}
			return nil
		}
	}
}

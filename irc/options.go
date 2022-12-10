package irc

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type ConnOption func(*Conn) error

var ErrInvalidField = errors.New("irc: invalid field")

func WithAuth(username, token string) ConnOption {
	return func(conn *Conn) error {
		if strings.HasPrefix(strings.ToLower(token), "oauth:") {
			token = token[6:]
		}

		if username == "" {
			return fmt.Errorf("%w - username", ErrInvalidField)
		}
		if token == "" {
			return fmt.Errorf("%w - token", ErrInvalidField)
		}
		conn.username = username
		conn.token = token
		return nil
	}
}

func WithAddress(host string, port uint16) ConnOption {
	return func(conn *Conn) (err error) {
		if host == "" {
			return fmt.Errorf("%w - host", ErrInvalidField)
		}
		if port < 1 {
			return fmt.Errorf("%w - port", ErrInvalidField)
		}
		conn.addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
		return
	}
}

func WithHostname(hostname string) ConnOption {
	return func(conn *Conn) error {
		if hostname == "" {
			return fmt.Errorf("%w - hostname", ErrInvalidField)
		}
		conn.hostname = hostname
		return nil
	}
}

func WithoutTLS() ConnOption {
	return func(conn *Conn) error {
		conn.tls = false
		return nil
	}
}

func WithInsecure() ConnOption {
	return func(conn *Conn) error {
		conn.insecure = true
		return nil
	}
}

func WithBufferSize(bufferSize int) ConnOption {
	return func(conn *Conn) error {
		if bufferSize < 1 {
			return fmt.Errorf("%w - bufferSize", ErrInvalidField)
		}
		conn.bufferSize = bufferSize
		return nil
	}
}

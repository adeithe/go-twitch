package irc

import (
	"net"
	"strings"
)

type ConnOption func(*Conn)

// WithTags enables additional metadata to the command and membership messages.
// When enabling the tags capability, the commands capability is also enabled automatically due to Twitch requirements.
//
// See: https://dev.twitch.tv/docs/irc/capabilities
func WithTags() ConnOption {
	return func(c *Conn) {
		WithCommands()(c)
		WithCapability(CapabilityTags)(c)
	}
}

// WithCommands allows the client to use Twitch chat commands in chat while connected to the server.
//
// See: https://dev.twitch.tv/docs/irc/capabilities
func WithCommands() ConnOption {
	return WithCapability(CapabilityCommands)
}

// WithMembership enables JOIN and PART messages when users join and leave the chat room.
//
// See: https://dev.twitch.tv/docs/irc/capabilities
func WithMembership() ConnOption {
	return WithCapability(CapabilityMembership)
}

// WithCapability adds a capability to the connection.
//
// See: https://dev.twitch.tv/docs/irc/capabilities
func WithCapability(capability Capability) ConnOption {
	return func(conn *Conn) {
		for _, current := range conn.capabilities {
			if capability == current {
				return
			}
		}
		conn.capabilities = append(conn.capabilities, capability)
	}
}

// WithAuth attempts to authenticate with the server using the given username and token.
func WithAuth(username, token string) ConnOption {
	return func(conn *Conn) {
		conn.username = username
		conn.token = strings.TrimPrefix(token, "oauth:")
	}
}

// WithAddr sets the
func WithAddr(addr net.Addr) ConnOption {
	return func(c *Conn) {
		c.addr = addr
	}
}

// WithHostname sets the hostname to validate the certificate with when connecting to the server.
func WithHostname(hostname string) ConnOption {
	return func(conn *Conn) {
		conn.hostname = hostname
	}
}

// WithoutTLS disables TLS support when connecting to the server.
func WithoutTLS() ConnOption {
	return func(conn *Conn) {
		conn.tls = false
	}
}

// WithInsecure disables TLS certificate verification when connecting to the server.
func WithInsecure() ConnOption {
	return func(conn *Conn) {
		conn.insecure = true
	}
}

// WithBufferSize sets the size of the read buffer.
func WithBufferSize(size int) ConnOption {
	return func(conn *Conn) {
		if size > 0 {
			conn.bufferSize = size
		}
	}
}

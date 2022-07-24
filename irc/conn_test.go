package irc_test

import (
	"bufio"
	"context"
	"net"
	"net/textproto"
	"strings"
	"syscall"
	"testing"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func TestCreateClient(t *testing.T) {
	client := irc.New(nil)
	assert.NotNil(t, client)
}

func TestConnectFailure(t *testing.T) {
	ctx, cancel := NewTestContext(t)
	defer cancel()

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6667")
	assert.NoError(t, err)

	client := irc.New(nil, irc.WithAddr(addr))
	assert.ErrorIs(t, client.Connect(ctx), syscall.ECONNREFUSED)
}

func TestConnPingError(t *testing.T) {
	client := irc.New(nil)
	_, err := client.Ping(context.Background())
	assert.ErrorIs(t, err, irc.ErrNotConnected)
}

func TestConnSendRawNotConnected(t *testing.T) {
	client := irc.New(nil)
	assert.ErrorIs(t, client.SendRaw("PING"), irc.ErrNotConnected)
}

func TestConn(t *testing.T) {
	ctx, cancel := NewTestContext(t)
	defer cancel()

	events := &irc.Events{
		Ready: make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-events.Ready:
				cancel()
			}
		}
	}()

	client := irc.New(events)
	assert.NoError(t, client.Connect(ctx))
}

func TestConnCapabilities(t *testing.T) {
	ctx, cancel := NewTestContext(t)
	defer cancel()

	mock, err := NewMockServer(t, func(conn net.Conn, msg *irc.Message) {
		cancel()
	})
	if !assert.NoError(t, err) {
		return
	}

	client := mock.Conn(nil, irc.WithTags(), irc.WithCommands(), irc.WithMembership())
	assert.NoError(t, client.Connect(ctx))
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
	assert.ElementsMatch(t, mock.Capabilities(), []irc.Capability{
		irc.CapabilityTags,
		irc.CapabilityCommands,
		irc.CapabilityMembership,
	})
}

func TestConnAuthentication(t *testing.T) {
	ctx, cancel := NewTestContext(t)
	defer cancel()

	expectedUsername := "justinfan123"
	expectedToken := "Kappa123"

	mock, err := NewMockServer(t, func(conn net.Conn, msg *irc.Message) {
		cancel()
	})
	if !assert.NoError(t, err) {
		return
	}

	client := mock.Conn(nil,
		irc.WithAuth(expectedUsername, expectedToken),
		irc.WithTags(), irc.WithCommands(), irc.WithMembership(),
	)
	assert.NoError(t, client.Connect(ctx))
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
	mock.AssertAuth(t, expectedUsername, expectedToken)
}

func TestConnReaderParseError(t *testing.T) {
	ctx, cancel := NewTestContext(t)
	defer cancel()

	mock, err := NewMockServer(t, func(conn net.Conn, msg *irc.Message) {
		_, _ = conn.Write([]byte("@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0\r\n"))
	})
	if !assert.NoError(t, err) {
		return
	}

	client := mock.Conn(nil, irc.WithTags(), irc.WithCommands(), irc.WithMembership())
	assert.ErrorIs(t, client.Connect(ctx), irc.ErrPartialMessage)
}

func NewTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	deadline, _ := t.Deadline()
	return context.WithDeadline(context.Background(), deadline)
}

type ServeIRC func(net.Conn, *irc.Message)

type MockServer struct {
	lis  net.Listener
	conn net.Conn

	username, token string
	capabilities    []irc.Capability
}

func NewMockServer(t *testing.T, f ServeIRC) (*MockServer, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { _ = lis.Close() })

	srv := &MockServer{lis: lis}
	go srv.listen(f)
	return srv, nil
}

func (srv *MockServer) Conn(events *irc.Events, opts ...irc.ConnOption) *irc.Conn {
	if srv.conn != nil {
		_ = srv.conn.Close()
	}

	connOpts := []irc.ConnOption{
		irc.WithAddr(srv.lis.Addr()),
		irc.WithInsecure(),
		irc.WithoutTLS(),
	}
	connOpts = append(connOpts, opts...)
	return irc.New(events, connOpts...)
}

func (srv *MockServer) AssertAuth(t *testing.T, username, token string) bool {
	return assert.Equal(t, srv.username, username) && assert.Equal(t, srv.token, token)
}

func (srv *MockServer) Capabilities() []irc.Capability {
	return srv.capabilities
}

func (srv *MockServer) listen(f ServeIRC) {
	for {
		conn, err := srv.lis.Accept()
		if err != nil {
			break
		}
		defer conn.Close()
		srv.conn = conn

		reader := textproto.NewReader(bufio.NewReader(conn))
		for {
			line, err := reader.ReadLine()
			if err != nil {
				break
			}

			msg, err := irc.ParseMessage(line)
			if err != nil {
				_ = conn.Close()
				break
			}

			srv.handleMessage(conn, msg)

			if msg.Command != "PING" {
				continue
			}

			if f != nil {
				f(conn, msg)
			}
		}
	}
}

func (srv *MockServer) handleMessage(conn net.Conn, msg *irc.Message) {
	switch msg.Command {
	case "CAP":
		if len(msg.Params) < 1 || msg.Params[0] != "REQ" {
			break
		}
		for _, capability := range strings.Split(msg.Text, " ") {
			srv.capabilities = append(srv.capabilities, irc.Capability(capability))
		}
	case "PASS":
		if len(msg.Params) < 1 {
			break
		}
		srv.token = strings.TrimPrefix(msg.Params[0], "oauth:")
	case "NICK":
		if len(msg.Params) < 1 {
			break
		}
		srv.username = msg.Params[0]
	case "PING":
		_, _ = conn.Write([]byte("PONG\r\n"))
	}
}

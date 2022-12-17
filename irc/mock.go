package irc

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
)

type MockServer struct {
	listener net.Listener
	clients  []*mockClient
	users    []string
}

type mockClient struct {
	ready           bool
	authenticated   bool
	username, token string

	conn   net.Conn
	server *MockServer
}

// RunT creates a mock server for testing Twitch IRC clients.
//
// If the mock server fails to start, the test will fail immediately.
func RunT(t *testing.T) *MockServer {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = lis.Close() })

	srv := &MockServer{listener: lis}
	go srv.run()
	return srv
}

// AddUsers adds the provided logins to the mock server.
// Required for authentication and joining a chat.
//
// Any oauth token will be accepted so long as it is properly formatted.
func (s *MockServer) AddUsers(logins ...string) {
	s.users = append(s.users, logins...)
}

// WithAddress returns a ConnOption that will connect to the mock server.
func (s *MockServer) WithAddress() ConnOption {
	parts := strings.Split(s.listener.Addr().String(), ":")
	port, _ := strconv.ParseUint(parts[len(parts)-1], 0, 16)
	return WithAddress(strings.Join(parts[:len(parts)-1], ":"), uint16(port))
}

// Write sends the given lines to all connected clients.
func (s *MockServer) Write(lines ...string) error {
	for _, client := range s.clients {
		_, _ = client.conn.Write([]byte(strings.Join(lines, "\r\n") + "\r\n"))
	}
	return nil
}

// Addr returns the address of the mock server.
func (s *MockServer) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *MockServer) run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		client := &mockClient{conn: conn, server: s}
		s.clients = append(s.clients, client)
		go client.handle()
	}
}

func (c *mockClient) Write(lines ...string) error {
	_, err := c.conn.Write([]byte(strings.Join(lines, "\r\n") + "\r\n"))
	return err
}

func (c *mockClient) handle() {
	reader := textproto.NewReader(bufio.NewReaderSize(c.conn, 4096))
	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}

		msg, err := ParseMessage(line)
		if err != nil {
			continue
		}

		if !c.ready {
			c.handleLoginMessage(msg)
			continue
		}
		c.handleMessage(msg)
	}
}

func (c *mockClient) handleMessage(msg *RawMessage) {
	switch msg.Command {
	case CMDPing:
		_ = c.Write(":tmi.twitch.tv PONG tmi.twitch.tv :" + msg.Text)
	case CMDJoin:
		if len(msg.Params) > 0 {
			channel := strings.ToLower(msg.Params[0])[1:]
			for _, user := range c.server.users {
				if strings.EqualFold(channel, user) {
					if c.authenticated {
						_ = c.Write(fmt.Sprintf("@badge-info=subscriber/3;badges=subscriber/1;color=;display-name=%s;emote-sets=0,33,50,237,793,2126,3517,4578,5569,9400,10337,12239;mod=0;subscriber=1;user-type= :tmi.twitch.tv USERSTATE #%s", c.username, channel))
					}
					_ = c.Write("@emote-only=0;followers-only=-1;r9k=0;rituals=0;room-id=12345678;slow=0;subs-only=0 :tmi.twitch.tv ROOMSTATE #" + channel)
					return
				}
			}
			_ = c.Write(fmt.Sprintf("@msg-id=msg_channel_suspended :tmi.twitch.tv NOTICE #%s :%s", channel, "This channel does not exist or has been suspended."))
		}
	case CMDPrivMessage:
		if len(msg.Params) > 0 {
			channel := strings.ToLower(msg.Params[0])[1:]
			if c.authenticated {
				_ = c.Write(fmt.Sprintf("@badge-info=subscriber/2;badges=subscriber/1;color=;display-name=%s;emote-sets=0,33,50,237,793,2126,3517,4578,5569,9400,10337,12239;id=aef784d3-1a4a-4dd3-909c-f02da16e8221;mod=0;subscriber=1;user-type= :tmi.twitch.tv USERSTATE #%s", c.username, channel))
				return
			}
		}
	}
}

func (c *mockClient) handleLoginMessage(msg *RawMessage) {
	switch msg.Command {
	case "CAP":
		if len(msg.Params) > 0 && msg.Params[0] == "REQ" {
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv CAP * ACK :%s", msg.Text))
		}
	case "PASS":
		if len(msg.Params) > 0 {
			c.token = msg.Params[0]
		}
	case "NICK":
		defer func() { c.token = "" }()
		if len(msg.Params) > 0 {
			c.username = strings.ToLower(msg.Params[0])
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 001 %s :Welcome, GLHF!", c.username))
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 002 %s :Your host is tmi.twitch.tv", c.username))
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 003 %s :This server is rather new", c.username))
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 004 %s :-", c.username))
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 375 %s :-", c.username))
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 372 %s :You are in a maze of twisty passages, all alike.", c.username))
			if !strings.HasPrefix(strings.ToLower(c.username), "justinfan") {
				if !strings.HasPrefix(strings.ToLower(c.token), "oauth:") {
					_ = c.Write(":tmi.twitch.tv NOTICE * :Improperly formatted auth")
					_ = c.conn.Close()
					return
				}
				for _, user := range c.server.users {
					if strings.EqualFold(c.username, user) {
						c.ready = true
						c.authenticated = true
						_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 376 %s :>", c.username))
						_ = c.Write(fmt.Sprintf("@badge-info=;badges=;color=;display-name=%s;emote-sets=0,33,50,237,793,2126,3517,4578,5569,9400,10337,12239;user-id=12345678;user-type= :tmi.twitch.tv GLOBALUSERSTATE", c.username))
						return
					}
				}
				_ = c.Write(":tmi.twitch.tv NOTICE * :Login authentication failed")
				_ = c.conn.Close()
				return
			}
			c.ready = true
			_ = c.Write(fmt.Sprintf(":tmi.twitch.tv 376 %s :>", c.username))
		}
	}
}

package irc

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

// Conn is a single connection to the Twitch IRC service
type Conn struct {
	UserState GlobalUserState
	Username  string
	token     string

	socket      net.Conn
	isConnected bool

	Latency      time.Duration
	ServerNotice chan ServerNotice
	RawMessage   chan Message

	pingC    chan bool
	pingMx   sync.Mutex
	writerMx sync.Mutex
}

// IConn is a generic IRC connection
type IConn interface {
	SetLogin(string, string) error

	Connect(context.Context) error
	SendRaw(context.Context, ...string) error
	Ping(context.Context) (time.Duration, error)
	Close() error
}

var _ IConn = &Conn{}

// IP for the IRC server
const IP = "irc.chat.twitch.tv"

var DefaultBuffer = 1024

func (conn *Conn) SetLogin(username, token string) error {
	if conn.isConnected {
		return ErrAlreadyConnected
	}
	conn.Username = username
	if strings.HasPrefix(strings.ToLower(token), "oauth:") {
		token = token[6:]
	}
	conn.token = token
	return nil
}

// Connect attempts to open a connection to the IRC server
func (conn *Conn) Connect(ctx context.Context) error {
	if conn.isConnected {
		return ErrAlreadyConnected
	}
	dialer := &net.Dialer{KeepAlive: time.Second * 10}
	socket, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", IP, 6697), nil)
	if err != nil {
		return err
	}
	if len(conn.Username) < 1 || len(conn.token) < 1 {
		conn.SetLogin(fmt.Sprintf("justinfan%d", rand.Intn(99899)+100), "Kappa123")
	}
	conn.socket = socket
	conn.isConnected = true
	go conn.reader()
	return conn.SendRaw(ctx,
		"CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands",

		fmt.Sprintf("PASS oauth:%s", conn.token),
		fmt.Sprintf("NICK %s", conn.Username),
	)
}

// SendRaw writes the provided messages to the IRC server
func (conn *Conn) SendRaw(ctx context.Context, raw ...string) error {
	if !conn.isConnected {
		return ErrNotConnected
	}
	c := make(chan error, 1)
	conn.writerMx.Lock()
	defer conn.writerMx.Unlock()
	go func() {
		for _, msg := range raw {
			select {
			case <-ctx.Done():
				c <- ctx.Err()
				return
			default:
				if _, err := conn.socket.Write([]byte(msg + "\r\n")); err != nil {
					c <- err
					return
				}
			}
		}
		close(c)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func (conn *Conn) Ping(ctx context.Context) (time.Duration, error) {
	conn.pingMx.Lock()
	defer conn.pingMx.Unlock()
	if !conn.isConnected {
		return 0, ErrNotConnected
	}
	conn.pingC = make(chan bool, 1)
	start := time.Now()
	if err := conn.SendRaw(ctx, "PING"); err != nil {
		return 0, err
	}
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-conn.pingC:
		conn.Latency = time.Since(start)
	}
	return conn.Latency, nil
}

func (conn *Conn) Close() error {
	return conn.socket.Close()
}

func (conn *Conn) reader() error {
	reader := textproto.NewReader(bufio.NewReader(conn.socket))
	for {
		line, err := reader.ReadLine()
		if err != nil {
			conn.isConnected = false
			return err
		}
		msg, err := NewParsedMessage(line)
		if err != nil {
			continue
		}
		conn.handle(msg)
		if conn.RawMessage == nil {
			conn.RawMessage = make(chan Message, DefaultBuffer)
		}
		if cap(conn.RawMessage) > 0 && len(conn.RawMessage) == cap(conn.RawMessage) {
			<-conn.RawMessage
		}
		conn.RawMessage <- msg
	}
}

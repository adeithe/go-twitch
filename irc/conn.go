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

// Conn a connection to a Twitch IRC server.
type Conn struct {
	events *Events

	username, token string
	capabilities    []Capability
	latency         time.Duration

	tls, insecure bool
	hostname      string
	addr          net.Addr
	bufferSize    int

	conn     net.Conn
	writerMx sync.Mutex
	close    chan error

	pingC  chan struct{}
	pingMx sync.Mutex
}

const DefaultHostname = "irc.chat.twitch.tv"

var (
	ErrNotConnected = fmt.Errorf("irc: not connected")
)

// New creates a new IRC connection.
func New(events *Events, opts ...ConnOption) *Conn {
	if events == nil {
		events = &Events{}
	}

	conn := &Conn{
		events: events,
		tls:    true,
		close:  make(chan error, 1),
	}

	tlsAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:6697", DefaultHostname))
	defaultOpts := []ConnOption{
		WithAuth(
			fmt.Sprintf("justinfan%d", rand.Intn(99899)+100),
			"Kappa123",
		),
		WithAddr(tlsAddr),
		WithHostname(DefaultHostname),
		WithBufferSize(4096),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(conn)
	}

	if !conn.tls && conn.addr == tlsAddr {
		conn.addr, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:6667", DefaultHostname))
	}
	return conn
}

// Connect attempts to connect to the IRC server.
//
// If successful, the goroutine is blocked until the connection is closed.
func (c *Conn) Connect(ctx context.Context) error {
	conn, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	c.conn = conn
	go c.reader()

	_ = c.requestCapabilities()
	_ = c.authenticate()
	_, _ = c.Ping(context.Background())

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-c.close:
	}
	c.conn = nil
	if err == context.Canceled {
		err = nil
	}
	return err
}

func (c *Conn) Ping(ctx context.Context) (time.Duration, error) {
	c.pingMx.Lock()
	defer c.pingMx.Unlock()

	start := time.Now()
	c.pingC = make(chan struct{})
	if err := c.SendRaw(string(CMDPing)); err != nil {
		return 0, err
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-c.pingC:
	}
	c.latency = time.Since(start)
	emit(c.events.Latency, c.latency)
	return c.latency, nil
}

// SendRaw sends a raw message to the IRC server.
func (c *Conn) SendRaw(raw ...string) error {
	c.writerMx.Lock()
	defer c.writerMx.Unlock()
	if c.conn == nil {
		return ErrNotConnected
	}

	var err error
	if len(raw) > 0 {
		_, err = c.conn.Write([]byte(strings.Join(raw, "\r\n") + "\r\n"))
	}
	return err
}

// Capabilities returns the list of capabilities that enabled on the connection.
func (c *Conn) Capabilities() []Capability {
	return c.capabilities
}

// Addr returns the address used to connect to the server.
func (c *Conn) Addr() net.Addr {
	return c.addr
}

func (c *Conn) dial() (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   time.Second * 10,
		KeepAlive: time.Second * 10,
	}
	if c.tls {
		config := &tls.Config{
			MinVersion:         tls.VersionTLS12,
			ServerName:         c.hostname,
			InsecureSkipVerify: c.insecure,
		}
		return tls.DialWithDialer(dialer, "tcp", c.addr.String(), config)
	}
	return dialer.Dial("tcp", c.addr.String())
}

func (c *Conn) reader() {
	reader := textproto.NewReader(bufio.NewReaderSize(c.conn, c.bufferSize))
	for {
		line, err := reader.ReadLine()
		if err != nil {
			c.close <- err
			break
		}

		msg, err := ParseMessage(line)
		if err != nil {
			c.close <- err
			break
		}

		emit(c.events.RawMessage, msg)
		c.handleMessage(msg)
	}
}

func (c *Conn) requestCapabilities() error {
	if len(c.capabilities) == 0 {
		return nil
	}

	var capabilities []string
	for _, capability := range c.Capabilities() {
		capabilities = append(capabilities, string(capability))
	}

	return c.SendRaw("CAP REQ :" + strings.Join(capabilities, " "))
}

func (c *Conn) authenticate() error {
	return c.SendRaw(fmt.Sprintf("PASS oauth:%s", c.token), fmt.Sprintf("NICK %s", c.username))
}

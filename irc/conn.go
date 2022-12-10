package irc

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

type Conn struct {
	username, token string

	tls, insecure bool
	hostname      string
	addr          net.Addr
	bufferSize    int

	ready        bool
	conn         net.Conn
	channels     map[string]*Channel
	writerMx     sync.Mutex
	channelsMx   sync.RWMutex
	connectionMx sync.Mutex
	lastMessage  time.Time

	wg     sync.WaitGroup
	readyC chan error
	pingC  chan struct{}
	pingMx sync.Mutex
}

const DefaultHostname = "irc.chat.twitch.tv"

var (
	ErrNotConnected = errors.New("irc: not connected")
	ErrLoginFailed  = errors.New("irc: authentication failed")
	ErrReadyTimeout = errors.New("irc: timed out while waiting for ready message")
)

func New(opts ...ConnOption) (*Conn, error) {
	defaultOpts := []ConnOption{
		WithAuth(fmt.Sprintf("justinfan%d", rand.Intn(99899)+100), "Kappa123"),
		WithAddress(DefaultHostname, 6697),
		WithHostname(DefaultHostname),
		WithBufferSize(4096),
	}

	conn := &Conn{
		channels: make(map[string]*Channel),

		tls:    true,
		readyC: make(chan error),
		pingC:  make(chan struct{}),
	}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(conn); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (c *Conn) IsConnected() bool {
	if c.conn == nil || !c.ready {
		return false
	}
	return time.Since(c.lastMessage) < time.Minute*5
}

func (c *Conn) GetChannel(channelName string) (*Channel, bool) {
	c.channelsMx.RLock()
	defer c.channelsMx.RUnlock()
	channel, ok := c.channels[channelName]
	return channel, ok
}

// JoinChannel joins the specified channel and returns a channel instance.
//
// NOTE: The room state will be updated asynchronously once acknowledged by the server.
func (c *Conn) JoinChannel(channelName string) (*Channel, error) {
	channelName = sanitizeUsername(channelName)

	c.channelsMx.RLock()
	channel, ok := c.channels[channelName]
	c.channelsMx.RUnlock()
	if !ok {
		c.channelsMx.Lock()
		channel = &Channel{conn: c, name: channelName}
		c.channels[channelName] = channel
		c.channelsMx.Unlock()
	}

	if err := c.SendRaw(fmt.Sprintf("JOIN #%s", channelName)); err != nil {
		return nil, err
	}
	return channel, nil
}

func (c *Conn) PartChannel(channelName string) error {
	channelName = sanitizeUsername(channelName)
	if _, ok := c.GetChannel(channelName); !ok {
		return nil
	}

	if err := c.SendRaw(fmt.Sprintf("PART #%s", channelName)); err != nil {
		return err
	}

	c.channelsMx.Lock()
	defer c.channelsMx.Unlock()
	delete(c.channels, channelName)
	return nil
}

func (c *Conn) Ping(ctx context.Context) (time.Duration, error) {
	c.pingMx.Lock()
	defer c.pingMx.Unlock()

	now := time.Now()
	if err := c.SendRaw("PING"); err != nil {
		return 0, err
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-c.pingC:
		return time.Since(now), nil
	}
}

func (c *Conn) SendRaw(lines ...string) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}
	return c.write(lines...)
}

func (c *Conn) Close() error {
	c.connectionMx.Lock()
	defer c.connectionMx.Unlock()
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return err
	}
	c.wg.Wait()
	return nil
}

func (c *Conn) Connect(ctx context.Context) error {
	c.connectionMx.Lock()
	defer c.connectionMx.Unlock()
	if c.IsConnected() {
		return nil
	}

	conn, err := c.dial(ctx)
	if err != nil {
		return err
	}
	c.conn = conn

	c.wg.Add(1)
	go c.reader()
	err = c.authenticate(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrReadyTimeout
	}
	return err
}

func (c *Conn) dial(ctx context.Context) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: time.Second * 10}
	if c.tls {
		tlsDialer := &tls.Dialer{
			NetDialer: dialer,
			Config: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				ServerName:         DefaultHostname,
				InsecureSkipVerify: c.insecure,
			},
		}
		return tlsDialer.DialContext(ctx, "tcp", c.addr.String())
	}
	return dialer.DialContext(ctx, "tcp", c.addr.String())
}

func (c *Conn) authenticate(ctx context.Context) error {
	lines := []string{
		"CAP REQ :twitch.tv/tags twitch.tv/commands twitch.tv/membership",

		fmt.Sprintf("PASS oauth:%s", c.token),
		fmt.Sprintf("NICK %s", c.username),
	}

	if err := c.write(lines...); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		_ = c.conn.Close()
		return ctx.Err()
	case err := <-c.readyC:
		return err
	}
}

func (c *Conn) write(lines ...string) error {
	c.writerMx.Lock()
	defer c.writerMx.Unlock()
	_, err := c.conn.Write([]byte(strings.Join(lines, "\r\n") + "\r\n"))
	return err
}

func (c *Conn) reader() {
	reader := textproto.NewReader(bufio.NewReaderSize(c.conn, c.bufferSize))
	defer c.wg.Done()
	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}

		msg, err := ParseMessage(line)
		if err != nil {
			continue
		}
		c.lastMessage = time.Now()

		if !c.ready {
			c.handleLogin(msg)
			continue
		}
		c.handleMessage(msg)
	}
	c.ready = false
}

func (c *Conn) handleLogin(msg *Message) {
	switch msg.Command {
	case CMDReady:
		c.ready = true
		c.readyC <- nil
	case CMDNotice:
		c.readyC <- fmt.Errorf("%w - %s", ErrLoginFailed, msg.Text)
	}
}

func (c *Conn) handleMessage(msg *Message) {
	switch msg.Command {
	case CMDPing:
		_ = c.write("PONG :" + msg.Text)
	case CMDPong:
		c.pingC <- struct{}{}

	case CMDRoomState:
		c.handleRoomState(msg)
	case CMDUserState:
		c.handleUserState(msg)
	}
}

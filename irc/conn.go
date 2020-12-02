package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

// Conn is a single connection to the Twitch IRC service
type Conn struct {
	Username string
	token    string
	state    UserState

	attempts    int
	socket      net.Conn
	isShard     bool
	isConnected bool
	latency     time.Duration
	channels    map[string]bool
	ping        chan bool

	onReconnect     []func()
	onDisconnect    []func()
	onLatencyUpdate []func(time.Duration)
	onChannelJoin   []func(string, string)
	onChannelLeave  []func(string, string)
	onMessage       []func(ChatMessage)
	onRawMessage    []func(Message)

	listeners sync.Mutex
	writer    sync.Mutex
}

// IConn is a generic IRC connection
type IConn interface {
	SetLogin(string, string) error
	IsConnected() bool
	IsInChannel(string) bool

	Connect() error
	SendRaw(...string) error
	Ping() (time.Duration, error)
	Join(...string) error
	Say(string, string) error
	Sayf(string, string, ...interface{}) error
	Leave(...string) error
	Reconnect() error
	Close()

	OnReady(func())
	OnDisconnect(func())
	OnLatencyUpdate(func(time.Duration))
	OnChannelJoin(func(string, string))
	OnChannelLeave(func(string, string))
	OnMessage(func(ChatMessage))
	OnRawMessage(func(Message))
}

// IP for the IRC server
const IP = "irc.chat.twitch.tv"

// SetLogin sets the username and authentication token for the connection
//
// Will return an error if the connection is already open
func (conn *Conn) SetLogin(username, token string) error {
	if conn.IsConnected() {
		return ErrAlreadyConnected
	}
	conn.Username = strings.ToLower(username)
	if strings.HasPrefix(strings.ToLower(token), "oauth:") {
		token = token[6:]
	}
	conn.token = token
	return nil
}

// IsConnected returns true if the connection is active
func (conn *Conn) IsConnected() bool {
	return conn.isConnected
}

// IsInChannel returns true if this connection is listening to the provided channel
func (conn *Conn) IsInChannel(channel string) bool {
	conn.listeners.Lock()
	defer conn.listeners.Unlock()
	for c := range conn.channels {
		if strings.ToLower(channel) == c {
			return true
		}
	}
	return false
}

// Connect attempts to open a connection to the IRC server
func (conn *Conn) Connect() error {
	if conn.isConnected {
		return ErrAlreadyConnected
	}
	dialer := &net.Dialer{KeepAlive: time.Second * 10}
	socket, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", IP, 6697), nil)
	if err != nil {
		return err
	}
	if len(conn.Username) < 1 || len(conn.token) < 1 {
		conn.SetLogin("justinfan123", "Kappa123")
	}
	conn.socket = socket
	conn.isConnected = true
	go conn.reader()
	return conn.SendRaw(
		"CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands",

		fmt.Sprintf("PASS oauth:%s", conn.token),
		fmt.Sprintf("NICK %s", conn.Username),
	)
}

// SendRaw writes the provided messages to the IRC server
//
// Will attempt to connect to the IRC server if it is not already connected
func (conn *Conn) SendRaw(raw ...string) error {
	if !conn.IsConnected() {
		if err := conn.Connect(); err != nil {
			return err
		}
	}
	conn.writer.Lock()
	defer conn.writer.Unlock()
	for _, msg := range raw {
		if _, err := conn.socket.Write([]byte(msg + "\r\n")); err != nil {
			return err
		}
	}
	return nil
}

// Ping sends a ping message to the IRC server
//
// This operation will block, giving the server up to 5 seconds to respond after correcting for latency before failing
func (conn *Conn) Ping() (time.Duration, error) {
	conn.listeners.Lock()
	defer conn.listeners.Unlock()
	if !conn.IsConnected() {
		return 0, ErrNotConnected
	}
	start := time.Now()
	conn.ping = make(chan bool, 1)
	if err := conn.SendRaw("PING"); err != nil {
		return conn.latency, err
	}
	timer := time.NewTimer(time.Second*5 + conn.latency)
	defer timer.Stop()
	select {
	case <-conn.ping:
	case <-timer.C:
		return conn.latency, ErrPingTimeout
	}
	conn.latency = time.Since(start)
	for _, f := range conn.onLatencyUpdate {
		go f(conn.latency)
	}
	return conn.latency, nil
}

// Join attempts to join a channel
func (conn *Conn) Join(channels ...string) error {
	for _, channel := range channels {
		if err := conn.SendRaw(fmt.Sprintf("JOIN #%s", strings.TrimPrefix(channel, "#"))); err != nil {
			return err
		}
	}
	return nil
}

// Say sends a message in the provided channel if authenticated
//
// If using a shards, you must create a single connection and use it as a writer
func (conn *Conn) Say(channel string, message string) error {
	if conn.isShard {
		return ErrShardedMessageSend
	}
	if strings.HasPrefix(strings.ToLower(conn.Username), "justinfan") {
		return ErrNotAuthenticated
	}
	return conn.SendRaw(fmt.Sprintf("PRIVMSG #%s :%s", strings.TrimPrefix(channel, "#"), message))
}

// Sayf sends a formatted message in the provided channel if authenticated
//
// If using a shards, you must create a single connection and use it as a writer
func (conn *Conn) Sayf(channel string, format string, a ...interface{}) error {
	return conn.Say(channel, fmt.Sprintf(format, a...))
}

// Leave attempts to leave a channel
func (conn *Conn) Leave(channels ...string) error {
	for _, channel := range channels {
		if err := conn.SendRaw(fmt.Sprintf("PART #%s", strings.TrimPrefix(channel, "#"))); err != nil {
			return err
		}
		conn.listeners.Lock()
		delete(conn.channels, channel)
		conn.listeners.Unlock()
	}
	return nil
}

// Reconnect closes and reopens the IRC connection
//
// Equivalent to Connect if the connection is not already open
func (conn *Conn) Reconnect() error {
	if conn.IsConnected() {
		conn.Close()
	}
	if err := conn.Connect(); err != nil {
		return err
	}
	for _, f := range conn.onReconnect {
		go f()
	}
	return nil
}

// Close disconnects from the IRC server
func (conn *Conn) Close() {
	if !conn.IsConnected() {
		return
	}
	conn.socket.Close()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
	return
}

// OnReconnect event called when the connection is reopened
func (conn *Conn) OnReconnect(f func()) {
	conn.onReconnect = append(conn.onReconnect, f)
}

// OnDisconnect event called when the connection was closed
func (conn *Conn) OnDisconnect(f func()) {
	conn.onDisconnect = append(conn.onDisconnect, f)
}

// OnLatencyUpdate event called after the latency to server has been updated
func (conn *Conn) OnLatencyUpdate(f func(time.Duration)) {
	conn.onLatencyUpdate = append(conn.onLatencyUpdate, f)
}

// OnChannelJoin event called after a user joins a chatroom
func (conn *Conn) OnChannelJoin(f func(string, string)) {
	conn.onChannelJoin = append(conn.onChannelJoin, f)
}

// OnChannelLeave event called after a user leeaves a chatroom
func (conn *Conn) OnChannelLeave(f func(string, string)) {
	conn.onChannelLeave = append(conn.onChannelLeave, f)
}

// OnMessage event called after a message is received
func (conn *Conn) OnMessage(f func(ChatMessage)) {
	conn.onMessage = append(conn.onMessage, f)
}

// OnRawMessage event called after a raw IRC message has been handled
func (conn *Conn) OnRawMessage(f func(Message)) {
	conn.onRawMessage = append(conn.onRawMessage, f)
}

func (conn *Conn) reader() {
	reader := textproto.NewReader(bufio.NewReader(conn.socket))
	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}
		msg, err := NewParsedMessage(line)
		if err != nil {
			continue
		}
		switch msg.Command {
		case CMDReady:
			go conn.Ping()
		case CMDReconnect:
			conn.Reconnect()
			return
		case CMDPing:
			go conn.Ping()
		case CMDPong:
			close(conn.ping)

		case CMDRoomState:
			conn.listeners.Lock()
			if conn.channels == nil {
				conn.channels = make(map[string]bool)
			}
			conn.channels[strings.TrimPrefix(msg.Params[0], "#")] = true
			conn.listeners.Unlock()
		case CMDJoin:
			for _, f := range conn.onChannelJoin {
				go f(strings.TrimPrefix(msg.Params[0], "#"), msg.Sender.Username)
			}
		case CMDPart:
			for _, f := range conn.onChannelLeave {
				go f(strings.TrimPrefix(msg.Params[0], "#"), msg.Sender.Username)
			}

		case CMDGlobalUserState:
			conn.state = NewUserState(msg)
		case CMDUserState:
			conn.state = NewUserState(msg)

		case CMDHostTarget:

		case CMDUserNotice:

		case CMDClearChat:
		case CMDClearMessage:

		case CMDNotice:

		case CMDPrivMessage:
			msg := NewChatMessage(msg)
			for _, f := range conn.onMessage {
				go f(msg)
			}
		}
		for _, f := range conn.onRawMessage {
			go f(msg)
		}
	}
	conn.isConnected = false
	for _, f := range conn.onDisconnect {
		go f()
	}
}

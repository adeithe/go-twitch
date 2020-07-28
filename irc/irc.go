package irc

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc/cmd"
)

type IClient interface {
	SetLogin(string, string)
	Join(...string)
	Leave(...string)
	Say(string, string)
	Sayf(string, string, ...interface{})
	Whisper(string, string)
	Whisperf(string, string, ...interface{})
	SendRaw(...string)

	Connect() error
	Ping() error
	Reconnect()
	Close()

	OnReady(func())
	OnDisconnect(func())
	OnPing(func())
	OnPong(func(time.Duration))
	OnRawMessage(func(IRCMessage))
}

type Client struct {
	Username  string
	token     string
	UserState UserState

	conn           net.Conn
	latency        time.Duration
	lastPing       time.Time
	isConnected    bool
	isReconnecting bool
	isReady        bool
	mu             sync.Mutex

	TLS           bool
	SkipTLSVerify bool

	onReady      []func()
	onDisconnect []func()
	onPing       []func()
	onPong       []func(time.Duration)
	onMessage    []func(ChatMessage)
	onRawMessage []func(IRCMessage)
}

// IP for the IRC server.
const IP = "irc.chat.twitch.tv"

var _ IClient = &Client{}

// New IRC client for Twitch chat.
func New() (irc *Client) {
	irc = &Client{TLS: true}
	irc.SetLogin("justinfan123", "Kappa123")
	return
}

// SetLogin will set the login details for IRC.
// If not called, read-only details will be used.
func (irc *Client) SetLogin(username string, oauth string) {
	irc.Username = username
	irc.token = strings.TrimPrefix(oauth, "oauth:")
}

// Join the chat for a Twitch channel.
func (irc *Client) Join(channels ...string) {
	for _, channel := range channels {
		irc.SendRaw(fmt.Sprintf("JOIN #%s", strings.TrimPrefix(channel, "#")))
	}
}

// Leave departs a Twitch channels chatroom.
func (irc *Client) Leave(channels ...string) {
	for _, channel := range channels {
		irc.SendRaw(fmt.Sprintf("PART #%s", strings.TrimPrefix(channel, "#")))
	}
}

// Say sends a message in the specified channel.
func (irc *Client) Say(channel string, message string) {
	irc.SendRaw(fmt.Sprintf("PRIVMSG #%s :%s", strings.TrimPrefix(channel, "#"), message))
}

// Sayf format and send a message to the specified channel.
func (irc *Client) Sayf(channel string, format string, args ...interface{}) {
	irc.Say(channel, fmt.Sprintf(format, args...))
}

// Whisper sends a message privately to the specified user.
func (irc *Client) Whisper(username string, message string) {
	irc.SendRaw(fmt.Sprintf("PRIVMSG #jtv :/whisper %s %s", strings.TrimPrefix(username, "#"), message))
}

// Whisperf format and send a whisper to the specified user.
func (irc *Client) Whisperf(username string, format string, args ...interface{}) {
	irc.Say(username, fmt.Sprintf(format, args...))
}

// SendRaw sends a raw message to the IRC server.
func (irc *Client) SendRaw(raw ...string) {
	irc.mu.Lock()
	defer irc.mu.Unlock()
	for _, msg := range raw {
		if !irc.isConnected {
			break
		}
		irc.conn.Write([]byte(msg + "\r\n"))
	}
}

// Connect to the IRC server.
func (irc *Client) Connect() (err error) {
	if irc.isConnected {
		return errors.New("client is already connected")
	}
	dialer := &net.Dialer{KeepAlive: time.Second * 10}
	if irc.TLS {
		conf := &tls.Config{InsecureSkipVerify: irc.SkipTLSVerify}
		irc.conn, err = tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", IP, 6697), conf)
	} else {
		irc.conn, err = dialer.Dial("tcp", fmt.Sprintf("%s:%d", IP, 6667))
	}
	if err != nil {
		return
	}
	irc.isReconnecting = false
	irc.isConnected = true
	go irc.reader()
	irc.SendRaw(
		"CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands",

		fmt.Sprintf("PASS oauth:%s", irc.token),
		fmt.Sprintf("NICK %s", irc.Username),
	)
	return
}

// Ping the IRC server.
func (irc *Client) Ping() error {
	if !irc.isConnected {
		return errors.New("client is not connected")
	}
	irc.lastPing = time.Now()
	irc.SendRaw("PING")
	for _, f := range irc.onPing {
		go f()
	}
	return nil
}

// Reconnect to the IRC server.
func (irc *Client) Reconnect() {
	if !irc.isConnected {
		return
	}
	irc.isReconnecting = true
	irc.Close()
}

// Close the IRC connection.
func (irc *Client) Close() {
	if !irc.isConnected {
		return
	}
	irc.conn.Close()
}

func (irc *Client) reader() {
	reader := textproto.NewReader(bufio.NewReader(irc.conn))
	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}
		messages := strings.Split(line, "\r\n")
		for _, msg := range messages {
			irc.handle(msg)
		}
	}
	irc.isReady = false
	irc.isConnected = false
	irc.lastPing = time.Time{}
	for _, f := range irc.onDisconnect {
		go f()
	}
	irc.conn.Close()
	if irc.isReconnecting {
		irc.Connect()
	}
}

func (irc *Client) handle(raw string) {
	fmt.Println(raw)
	if msg, err := NewParsedMessage(raw); err == nil {
		switch msg.Command {
		case cmd.Ready:
			irc.isReady = true
			for _, f := range irc.onReady {
				go f()
			}
		case cmd.Reconnect:
			irc.Reconnect()
		case cmd.Ping:
			irc.Ping()
		case cmd.Pong:
			irc.latency = time.Since(irc.lastPing)
			for _, f := range irc.onPong {
				go f(irc.latency)
			}

		case cmd.GlobalUserState:
			irc.UserState = NewUserState(msg)
		case cmd.UserState:
			irc.UserState = NewUserState(msg)

		case cmd.HostTarget:

		case cmd.UserNotice:

		case cmd.ClearChat:
		case cmd.ClearMessage:

		case cmd.Notice:

		case cmd.PrivMessage:
			message := NewChatMessage(msg)
			for _, f := range irc.onMessage {
				go f(message)
			}
		}
	}
}

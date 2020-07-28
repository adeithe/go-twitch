package pubsub

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/pubsub/nonce"
	"github.com/gorilla/websocket"
)

type IClient interface {
	UseToken(string)
	Listen(...string)
	Unlisten(...string)

	Connect() error
	Write(int, []byte) error
	Ping() error
	Reconnect()
	Close()

	GetLatency() time.Duration
	IsConnected() bool
	IsReconnecting() bool

	OnDisconnect(func())
	OnTopicListen(func(string))
	OnTopicUnlisten(func(string))
	OnTopicResponseError(func(string, string))
	OnMessage(func(string, json.RawMessage))
}

type Client struct {
	token string

	conn           *websocket.Conn
	topics         []string
	latency        time.Duration
	isConnected    bool
	isReconnecting bool
	mu             sync.Mutex

	pendingListens   map[string]string
	pendingUnlistens map[string]string
	lastPing         time.Time
	done             chan byte

	onDisconnect         []func()
	onTopicListen        []func(string)
	onTopicUnlisten      []func(string)
	onTopicResponseError []func(string, string)
	onMessage            []func(string, json.RawMessage)
}

// IP for the PubSub server
const IP = "pubsub-edge.twitch.tv"

var _ IClient = &Client{}

// New PubSub Client
func New() *Client {
	return &Client{
		pendingListens:   make(map[string]string),
		pendingUnlistens: make(map[string]string),
	}
}

// UseToken sets the token to provide when listening to new topics.
// If using different tokens, you may use UseToken before Listen as needed.
func (pubsub *Client) UseToken(token string) {
	pubsub.token = strings.TrimPrefix(token, "oauth:")
}

// Listen to new PubSub topics.
func (pubsub *Client) Listen(topics ...string) {
	for _, topic := range topics {
		nonce := nonce.New()
		data, err := json.Marshal(MessageData{Type: Listen, Nonce: nonce, Data: TopicListenData{[]string{topic}, pubsub.token}})
		if err != nil {
			continue
		}
		pubsub.pendingListens[nonce] = topic
		pubsub.Write(websocket.TextMessage, data)
	}
}

// Unlisten from subscribed PubSub topics.
func (pubsub *Client) Unlisten(topics ...string) {
	for _, topic := range topics {
		nonce := nonce.New()
		data, err := json.Marshal(MessageData{Type: Unlisten, Nonce: nonce, Data: TopicListenData{[]string{topic}, pubsub.token}})
		if err != nil {
			continue
		}
		pubsub.pendingUnlistens[nonce] = topic
		pubsub.Write(websocket.TextMessage, data)
	}
}

// Connect to the PubSub server.
func (pubsub *Client) Connect() (err error) {
	u := url.URL{Scheme: "wss", Host: IP}
	pubsub.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	pubsub.topics = []string{}
	pubsub.done = make(chan byte)
	pubsub.isConnected = true
	pubsub.isReconnecting = false
	go pubsub.reader()
	go pubsub.ticker()
	pubsub.Ping()
	return
}

// Write sends a message to the PubSub server.
func (pubsub *Client) Write(messageType int, data []byte) error {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()
	return pubsub.conn.WriteMessage(messageType, data)
}

// Ping the PubSub server.
func (pubsub *Client) Ping() error {
	data, err := json.Marshal(MessageData{Type: Ping})
	if err != nil {
		return err
	}
	pubsub.lastPing = time.Now()
	return pubsub.Write(websocket.TextMessage, data)
}

// Reconnect to the PubSub server.
func (pubsub *Client) Reconnect() {
	if !pubsub.isConnected {
		return
	}
	pubsub.isReconnecting = true
	pubsub.Close()
}

// Close the PubSub connection.
func (pubsub *Client) Close() {
	pubsub.Write(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	select {
	case <-pubsub.done:
	case <-time.After(time.Second):
		pubsub.conn.Close()
	}
}

// GetLatency returns the time.Duration for the last Ping.
func (pubsub *Client) GetLatency() time.Duration {
	return pubsub.latency
}

// IsConnected returns true if the PubSub client has an active connection.
func (pubsub *Client) IsConnected() bool {
	return pubsub.isConnected
}

// IsReconnecting returns true if the PubSub client is attempting to reconnect.
func (pubsub *Client) IsReconnecting() bool {
	return pubsub.isReconnecting
}

func (pubsub *Client) reader() {
	defer pubsub.conn.Close()
	for {
		msgType, message, err := pubsub.conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			break
		}
		data := &MessageData{}
		if err := json.Unmarshal(message, &data); err != nil {
			break
		}
		switch data.Type {
		case Pong:
			pubsub.latency = time.Since(pubsub.lastPing)
			pubsub.lastPing = time.Time{}
		case Reconnect:
			pubsub.Reconnect()
		case Response:
			pubsub.handleResponse(*data)
		case Message:
			pubsub.handleMessage(*data)
		}
	}
	pubsub.isConnected = false
	close(pubsub.done)
	for _, f := range pubsub.onDisconnect {
		f()
	}
}

func (pubsub *Client) ticker() {
	for {
		select {
		case <-pubsub.done:
			if pubsub.isReconnecting {
				pubsub.Connect()
			}
			break
		case <-time.After(time.Minute * 5):
			pubsub.Ping()
		}
	}
}

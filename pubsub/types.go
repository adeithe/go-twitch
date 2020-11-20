package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrShardTooManyTopics returned when a shard has attempted to join too many topics
	ErrShardTooManyTopics = errors.New("too many topics on shard")
	// ErrShardIDOutOfBounds returned when an invalid shard id is provided
	ErrShardIDOutOfBounds = errors.New("shard id out of bounds")
	// ErrNonceTimeout returned when the server doesnt respond to a nonced message in time
	ErrNonceTimeout = errors.New("nonced message timeout")

	// ErrBadMessage returned when the server receives an invalid message
	ErrBadMessage = errors.New("server received an invalid message")
	// ErrBadAuth returned when a topic doesnt have the permissions required
	ErrBadAuth = errors.New("bad authentication for topic")
	// ErrBadTopic returned when an invalid topic was requested
	ErrBadTopic = errors.New("invalid topic")
	// ErrServer returned when something went wrong on the servers end
	ErrServer = errors.New("something went wrong on the servers end")
	// ErrUnknown returned when the server sends back an error that wasnt handled by the reader
	ErrUnknown = errors.New("server sent back an unknown error")
)

// Packet stores data about a message sent to/from the PubSub server
type Packet struct {
	Type  MessageType  `json:"type"`
	Nonce string       `json:"nonce,omitempty"`
	Data  interface{}  `json:"data,omitempty"`
	Error MessageError `json:"error,omitempty"`
}

// MessageData stores data about a message packet
type MessageData struct {
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

// TopicData stores data about a topic
type TopicData struct {
	Topics []string `json:"topics"`
	Token  string   `json:"auth_token,omitempty"`
}

// MessageType stores the type provided in MessageData
type MessageType string

const (
	// Listen outgoing message type
	Listen MessageType = "LISTEN"
	// Unlisten outgoing message type
	Unlisten MessageType = "UNLISTEN"
	// Ping outgoing message type
	Ping MessageType = "PING"

	// Response incoming message type
	Response MessageType = "RESPONSE"
	// Message incoming message type
	Message MessageType = "MESSAGE"
	// Pong incoming message type
	Pong MessageType = "PONG"
	// Reconnect incoming message type
	Reconnect MessageType = "RECONNECT"
)

// MessageError stores the error provided in MessageData
type MessageError string

const (
	// BadMessage server received an invalid message
	BadMessage MessageError = "ERR_BADMESSAGE"
	// BadAuth provided token does not have required permissions
	BadAuth MessageError = "ERR_BADAUTH"
	// TooManyTopics attempted to listen to too many topics
	TooManyTopics MessageError = "ERR_TOO_MANY_TOPICS"
	// BadTopic provided topic is invalid
	BadTopic MessageError = "ERR_BADTOPIC"
	// InvalidTopic provided topic is invalid
	InvalidTopic MessageError = "Invalid Topic"
	// ServerError something went wrong on the servers side
	ServerError MessageError = "ERR_SERVER"
)

// ParseTopic returns a topic string with the provided arguments
func ParseTopic(str string, args ...interface{}) string {
	var params []string
	for _, arg := range args {
		params = append(params, fmt.Sprint(arg))
	}
	return fmt.Sprintf("%s.%s", str, strings.Join(params, "."))
}

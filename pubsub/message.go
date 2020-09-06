package pubsub

// MessageData stores data about a PubSub message.
type MessageData struct {
	Type  MessageType `json:"type"`
	Nonce string      `json:"nonce,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// MessageType string referring to the possible PubSub message types.
type MessageType string

const (
	Listen    MessageType = "LISTEN"
	Unlisten  MessageType = "UNLISTEN"
	Message   MessageType = "MESSAGE"
	Response  MessageType = "RESPONSE"
	Reconnect MessageType = "RECONNECT"
	Ping      MessageType = "PING"
	Pong      MessageType = "PONG"
)

// TopicListenData stores data about a listen message. Usually sent to the server by the client.
type TopicListenData struct {
	Topics []string `json:"topics"`
	Token  string   `json:"auth_token,omitempty"`
}

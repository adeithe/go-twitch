package pubsub

type MessageData struct {
	Type  MessageType `json:"type"`
	Nonce string      `json:"nonce,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type MessageType string

const (
	Listen    MessageType = "LISTEN"
	Unlisten              = "UNLISTEN"
	Message               = "MESSAGE"
	Response              = "RESPONSE"
	Reconnect             = "RECONNECT"
	Ping                  = "PING"
	Pong                  = "PONG"
)

type TopicListenData struct {
	Topics []string `json:"topics"`
	Token  string   `json:"auth_token,omitempty"`
}

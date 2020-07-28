package pubsub

import "encoding/json"

// OnDisconnect event called when the PubSub client gets disconnected from the server.
func (pubsub *Client) OnDisconnect(handler func()) {
	pubsub.onDisconnect = append(pubsub.onDisconnect, handler)
}

// OnTopicListen event called when the PubSub client listens to a topic.
func (pubsub *Client) OnTopicListen(handler func(topic string)) {
	pubsub.onTopicListen = append(pubsub.onTopicListen, handler)
}

// OnTopicUnlisten event called when the PubSub client unlistens a topic.
func (pubsub *Client) OnTopicUnlisten(handler func(topic string)) {
	pubsub.onTopicUnlisten = append(pubsub.onTopicUnlisten, handler)
}

// OnTopicResponseError event called when the PubSub client fails to listen or unlisten a topic.
func (pubsub *Client) OnTopicResponseError(handler func(topic string, err string)) {
	pubsub.onTopicResponseError = append(pubsub.onTopicResponseError, handler)
}

// OnMessage event called when the PubSub client receives a message about a listened topic.
func (pubsub *Client) OnMessage(handler func(topic string, data json.RawMessage)) {
	pubsub.onMessage = append(pubsub.onMessage, handler)
}

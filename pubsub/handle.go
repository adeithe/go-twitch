package pubsub

import (
	"encoding/json"
)

func (pubsub *Client) handleResponse(data MessageData) {
	if len(data.Error) <= 0 {
		if topic, ok := pubsub.pendingListens[data.Nonce]; ok {
			pubsub.topics = append(pubsub.topics, topic)
			delete(pubsub.pendingListens, data.Nonce)
			for _, f := range pubsub.onTopicListen {
				f(topic)
			}
		}
		if topic, ok := pubsub.pendingUnlistens[data.Nonce]; ok {
			var topics []string
			for _, listening := range pubsub.topics {
				if listening != topic {
					topics = append(topics, listening)
				}
			}
			pubsub.topics = topics
			delete(pubsub.pendingUnlistens, data.Nonce)
			for _, f := range pubsub.onTopicUnlisten {
				f(topic)
			}
		}
	} else {
		var topic string
		var ok bool
		if topic, ok = pubsub.pendingListens[data.Nonce]; !ok {
			topic = pubsub.pendingUnlistens[data.Nonce]
			delete(pubsub.pendingUnlistens, data.Nonce)
		}
		delete(pubsub.pendingListens, data.Nonce)
		for _, f := range pubsub.onTopicResponseError {
			f(topic, data.Error)
		}
	}
}

func (pubsub *Client) handleMessage(data MessageData) {
	type Message struct {
		Topic string          `json:"topic"`
		Data  json.RawMessage `json:"message"`
	}
	msg := &Message{}
	bytes, err := json.Marshal(data.Data)
	if err != nil {
		return
	}
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return
	}
	for _, f := range pubsub.onMessage {
		f(msg.Topic, msg.Data)
	}
}

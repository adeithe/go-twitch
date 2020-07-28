package twitch

import (
	"github.com/Adeithe/go-twitch/api"
	"github.com/Adeithe/go-twitch/irc"
	"github.com/Adeithe/go-twitch/pubsub"
)

// API provides tools for developing integrations with Twitch.
func API(clientID string) *api.Client {
	return api.New(clientID)
}

// IRC is the Twitch interface for chat functionality.
func IRC() *irc.Client {
	return irc.New()
}

// PubSub enables you to subscribe to a topic for updates.
func PubSub() *pubsub.Client {
	return pubsub.New()
}

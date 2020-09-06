# go-twitch [![GoDoc](https://godoc.org/github.com/Adeithe/go-twitch?status.svg)](https://godoc.org/github.com/Adeithe/go-twitch) [![Go Report Card](https://goreportcard.com/badge/github.com/Adeithe/go-twitch)](https://goreportcard.com/report/github.com/Adeithe/go-twitch)

A complete interface for Twitch services.

## Getting Started
```go
package main

import (
	"fmt"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/pubsub"
)

func main() {
	// Create an API client
	api := twitch.API("p0gch4mp101fy451do9uod1s1x9i4a")
	user := api.NewBearer("2gbdx6oar67tqtcmt49t3wpcgycthx")

	// Create a IRC client
	irc := twitch.IRC()
	irc.SetLogin("username", "oauth:123123123") // Skip this to login anonymously
	if err := irc.Connect(); err != nil {
		panic(err)
	}
	irc.Join("channel1", "channel2", "channel3")

	// Create a PubSub client
	ps := twitch.PubSub()
	if err := ps.Connect(); err != nil {
		panic(err)
	}
	ps.Listen(pubsub.ParseTopic(pubsub.ChatModeratorActions, 44322889))
}
```
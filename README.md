# go-twitch [![GoDoc](https://godoc.org/github.com/adeithe/go-twitch?status.svg)](https://godoc.org/github.com/adeithe/go-twitch) [![Go Report Card](https://goreportcard.com/badge/github.com/adeithe/go-twitch)](https://goreportcard.com/report/github.com/adeithe/go-twitch) [![CircleCI](https://circleci.com/gh/Adeithe/go-twitch/tree/master.svg?style=svg)](https://circleci.com/gh/Adeithe/go-twitch/tree/master)

A complete interface for Twitch services.

## Getting Started

### Installing

```sh
$ go get -u github.com/Adeithe/go-twith
```

### Usage

```go
package main

import (
	"fmt"

	twitch "github.com/Adeithe/go-twitch"
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
	ps.Listen("community-points-channel-v1", 44322889)
}
```

### Examples

Below is a short list of the available examples.

 - [Chat](https://github.com/Adeithe/go-twitch/tree/master/.examples/chat) - Connect to a Twitch channels chatroom
 - [PubSub](https://github.com/Adeithe/go-twitch/tree/master/.examples/pubsub) - Listen to various PubSub topics for a Twitch channel
 - [Get Users](https://github.com/Adeithe/go-twitch/tree/master/.examples/getusers) - Get up to 100 users by their username and see if they are live
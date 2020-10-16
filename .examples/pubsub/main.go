package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/pubsub"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	ps := twitch.PubSub()

	ps.OnTopicListen(func(topic string) {
		fmt.Printf("LISTEN %s\n", topic)
	})

	ps.OnTopicResponseError(func(topic string, err string) {
		fmt.Printf("LISTEN %s ERROR %s\n", topic, err)
	})

	ps.OnMessage(func(topic string, data json.RawMessage) {
		fmt.Printf("%s> %s\n", topic, data)
	})

	ps.OnDisconnect(func() {
		fmt.Println("Disconnected from PubSub!")
		sc <- syscall.SIGINT
	})

	if err := ps.Connect(); err != nil {
		panic(err)
	}

	ps.UseToken("2gbdx6oar67tqtcmt49t3wpcgycthx")
	ps.Listen(pubsub.ParseTopic(pubsub.ChatModeratorActions, 44322889))

	<-sc
	ps.Close()
}

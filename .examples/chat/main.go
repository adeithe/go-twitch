package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Adeithe/go-twitch"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	chat := twitch.IRC()
	chat.SetLogin("username", "oauth:123123123")
	chat.OnDisconnect(func() {
		fmt.Println("Disconnected from IRC!")
		sc <- syscall.SIGTERM
	})
	if err := chat.Connect(); err != nil {
		panic(err)
	}
	chat.Join("channel1")
	fmt.Println("Connected to IRC!")

	<-sc
	chat.Close()
}

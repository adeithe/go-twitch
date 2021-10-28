package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	writer := &irc.Conn{}
	writer.SetLogin("username", "oauth:123123123")
	if err := writer.Connect(ctx); err != nil {
		panic(err)
	}

	reader := irc.New()
	fmt.Println("Connected to IRC!")

	<-sc
	writer.Close()
}

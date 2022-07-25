package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Adeithe/go-twitch/irc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go listenForSignals(cancel)

	events := &irc.Events{
		Ready:      make(chan struct{}),
		RawMessage: make(chan *irc.Message),
	}

	conn := irc.New(events,
		irc.WithAuth("myusername", "oauth:yfvzjqb705z12hrhy1zkwa9xt7v662"),
		irc.WithTags(), irc.WithCommands(), irc.WithMembership(),
	)

	go handleEvents(ctx, conn, events)
	if err := conn.Connect(ctx); err != nil {
		panic(err)
	}
}

func handleEvents(ctx context.Context, conn *irc.Conn, events *irc.Events) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-events.Ready:
			_ = conn.SendRaw("JOIN #jtv")
		case msg := <-events.RawMessage:
			fmt.Println(msg.Raw)
		}
	}
}

func listenForSignals(cancel context.CancelFunc) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	cancel()
}

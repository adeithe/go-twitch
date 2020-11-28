package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	chat := twitch.IRC()

	chat.OnShardReconnect(func(shardID int) {
		fmt.Printf("Shard #%d reconnected\n", shardID)
	})

	chat.OnShardLatencyUpdate(func(shardID int, latency time.Duration) {
		fmt.Printf("Shard #%d has %dms ping\n", shardID, latency.Milliseconds())
	})

	chat.OnShardMessage(func(shardID int, msg irc.ChatMessage) {
		fmt.Printf("#%s %s: %s\n", msg.Channel, msg.Sender.DisplayName, msg.Message)
	})

	if err := chat.Join("dallas"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to IRC!")

	<-sc
	fmt.Println("Stopping...")
	chat.Close()
}

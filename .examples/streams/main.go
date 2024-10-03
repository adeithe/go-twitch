package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/adeithe/go-twitch/api"
)

var (
	clientID    string
	bearerToken string
)

func init() {
	flag.StringVar(&clientID, "client-id", "", "Twitch Client ID")
	flag.StringVar(&bearerToken, "token", "", "Twitch Bearer Token")
	flag.Parse()
}

func main() {
	client := api.New(clientID)

	var cursor string
	call := client.Streams.List().First(100)
	for {
		streams, err := call.After(cursor).Do(context.Background(), api.WithBearerToken(bearerToken))
		if err != nil {
			panic(err)
		}

		if len(streams.Data) == 0 {
			break
		}

		for _, stream := range streams.Data {
			fmt.Printf("%s is live! - Streaming %s to %d viewers\n", stream.UserLogin, stream.GameName, stream.ViewerCount)
		}
		cursor = streams.Cursor
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Adeithe/go-twitch/api"
	"github.com/Adeithe/go-twitch/api/kraken"
)

var reader *bufio.Reader

func main() {
	usernames := []string{}
	for i := 0; i < 100; i++ {
		fmt.Print("Get By Username: ")
		username := stdin()
		if username == "" {
			break
		}
		usernames = append(usernames, username)
	}
	client := api.Official.Kraken()
	users, err := client.GetUsers(kraken.UserOpts{
		Logins: usernames,
	})
	ids := []string{}
	for _, user := range users.Data {
		ids = append(ids, user.ID)
	}
	streams, err := client.GetStreams(kraken.StreamOpts{
		ChannelIDs: ids,
	})
	if err != nil {
		panic(err)
	}
	live := make(map[string]string)
	for _, stream := range streams.Data {
		live[fmt.Sprintf("%v", stream.Channel.ID)] = strings.ToUpper(stream.Type)
	}
	for i, user := range users.Data {
		live, _ := live[user.ID]
		fmt.Printf("[%d] %s (User ID: %s) %s\n", i, user.Login, user.ID, live)
	}
}

func stdin() string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	str, _ := reader.ReadString('\n')
	return strings.TrimSuffix(strings.TrimSuffix(str, "\r\n"), "\n")
}

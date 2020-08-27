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
	users, err := api.Official.Kraken().GetUsers(kraken.UserOpts{
		Logins: usernames,
	})
	if err != nil {
		panic(err)
	}
	for i, user := range users.Data {
		fmt.Printf("[%d] %s (User ID: %s)\n", i, user.Login, user.ID)
	}
}

func stdin() string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	str, _ := reader.ReadString('\n')
	return strings.TrimSuffix(strings.TrimSuffix(str, "\r\n"), "\n")
}

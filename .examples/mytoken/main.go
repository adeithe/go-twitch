package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Adeithe/go-twitch/api"
)

var reader *bufio.Reader

func main() {
	fmt.Print("Twitch Username: ")
	username := stdin()
	fmt.Print("Twitch Password: ")
	password := stdin()
	login, err := api.Official.Login(username, password)
	if err != nil {
		panic(err)
	}
	if login.GetErrorCode() == 3022 {
		fmt.Print("Twitch 2FA: ")
		code := stdin()
		if err := login.Verify(code); err != nil {
			panic(err)
		}
	}
	if login.GetErrorCode() != 0 {
		fmt.Printf("Failed: %s\n", login.GetError())
		return
	}
	fmt.Printf("Twitch Access Token: %s\n", login.GetAccessToken())
}

func stdin() string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	str, _ := reader.ReadString('\n')
	return strings.TrimSuffix(strings.TrimSuffix(str, "\r\n"), "\n")
}

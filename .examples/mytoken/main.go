package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Adeithe/go-twitch/api"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Twitch Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Twitch Password: ")
	password, _ := reader.ReadString('\n')
	login, err := api.Official.Login(clean(username), clean(password))
	if err != nil {
		panic(err)
	}
	if login.GetErrorCode() == 3022 {
		fmt.Print("Twitch 2FA: ")
		code, _ := reader.ReadString('\n')
		if err := login.Verify(clean(code)); err != nil {
			panic(err)
		}
	}
	if login.GetErrorCode() != 0 {
		fmt.Printf("Failed: %s\n", login.GetError())
		return
	}
	fmt.Printf("Twitch Access Token: %s\n", login.GetAccessToken())
}

func clean(str string) string {
	return strings.TrimSuffix(str, "\r\n")
}

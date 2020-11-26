package main

import (
	"fmt"
	"strings"

	"github.com/Adeithe/go-twitch/graphql"
)

func main() {
	gql := graphql.New()

	users, err := gql.GetUsersByID("44322889")
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		var str string
		str += fmt.Sprintf("ID: %s\n", user.ID)
		str += fmt.Sprintf("Login: %s\n", user.Login)
		str += fmt.Sprintf("Display Name: %s\n", user.DisplayName)
		str += fmt.Sprintf("Profile URL: %s\n", user.ChannelURL)
		str += fmt.Sprintf("Banner URL: %s\n", user.BannerImageURL)
		str += fmt.Sprintf("Offline Screen URL: %s\n", user.OfflineImageURL)
		str += fmt.Sprintf("Chat Color: %s\n", user.ChatColor)
		str += fmt.Sprintf("Description: %s\n", user.Description)
		str += fmt.Sprintf("Profile Views: %d\n", user.ProfileViewCount)
		if user.Hosting != nil {
			str += fmt.Sprintf("Hosting:\n")
			str += fmt.Sprintf("\tID: %s\n", user.Hosting.ID)
			str += fmt.Sprintf("\tLogin: %s\n", user.Hosting.DisplayName)
			str += fmt.Sprintf("\tDisplay Name: %s\n", user.Hosting.DisplayName)
			str += fmt.Sprintf("\tChannel URL: %s\n", user.Hosting.Channel.URL)
			str += fmt.Sprintf("\tStream:\n")
			for _, s := range sts(user.Hosting.Stream) {
				str += fmt.Sprintf("\t\t%s\n", s)
			}
			str += fmt.Sprintf("Is Affiliate: %t\n", user.Hosting.Roles.IsAffiliate)
			str += fmt.Sprintf("\tIs Partner: %t\n", user.Hosting.Roles.IsPartner)
			str += fmt.Sprintf("\tCreated At: %s\n", user.Hosting.CreatedAt)
			str += fmt.Sprintf("\tUpdated At: %s\n", user.Hosting.UpdatedAt)
		} else if user.Stream != nil {
			str += fmt.Sprintf("Streaming:\n")
			for _, s := range sts(*user.Stream) {
				str += fmt.Sprintf("\t%s\n", s)
			}
		}
		str += fmt.Sprintf("Has Prime: %t\n", user.HasPrime)
		str += fmt.Sprintf("Has Turbo: %t\n", user.HasTurbo)
		str += fmt.Sprintf("Is Affiliate: %t\n", user.Roles.IsAffiliate)
		str += fmt.Sprintf("Is Partner: %t\n", user.Roles.IsPartner)
		if user.DeletedAt.Unix() < user.CreatedAt.Unix() {
			str += fmt.Sprintf("Created At: %s\n", user.CreatedAt)
			str += fmt.Sprintf("Updated At: %s\n", user.UpdatedAt)
			fmt.Println(str)
			continue
		}
		str += fmt.Sprintf("Deleted At: %s\n", user.DeletedAt)
		fmt.Println(str)
	}
}

func sts(stream graphql.Stream) []string {
	var str string
	str += fmt.Sprintf("ID: %s\n", stream.ID)
	str += fmt.Sprintf("Title: %s\n", stream.Channel.BroadcastSettings.Title)
	str += fmt.Sprintf("Viewers Count: %d\n", stream.ViewersCount)

	str += fmt.Sprintf("Channel:\n")
	str += fmt.Sprintf("\tID: %s\n", stream.Channel.ID)
	str += fmt.Sprintf("\tLogin: %s\n", stream.Channel.Name)
	str += fmt.Sprintf("\tDisplay Name: %s\n", stream.Channel.DisplayName)
	str += fmt.Sprintf("\tLanguage: %s\n", stream.Channel.BroadcastSettings.Language)
	str += fmt.Sprintf("\tIs Mature: %t\n", stream.Channel.BroadcastSettings.IsMature)

	if stream.Game != nil {
		str += fmt.Sprintf("Game:\n")
		str += fmt.Sprintf("\tID: %s\n", stream.Game.ID)
		str += fmt.Sprintf("\tName: %s\n", stream.Game.Name)
		str += fmt.Sprintf("\tDisplay Name: %s\n", stream.Game.DisplayName)
		str += fmt.Sprintf("\tLive Channels: %d\n", stream.Game.BroadcastersCount)
		str += fmt.Sprintf("\tViewers: %d\n", stream.Game.ViewersCount)
		str += fmt.Sprintf("\tFollowers: %d\n", stream.Game.FollowersCount)
		str += fmt.Sprintf("\tPopularity: %d\n", stream.Game.PopularityScore)
		if stream.Game.GiantBombID != nil {
			str += fmt.Sprintf("\tGiantBomb ID: %s\n", stream.Game.GiantBombID)
		}
		if stream.Game.PrestoID != nil {
			str += fmt.Sprintf("\tPresto ID: %s\n", stream.Game.PrestoID)
		}
	}

	str += fmt.Sprintf("Broadcaster Software: %s\n", stream.BroadcasterSoftware)
	str += fmt.Sprintf("Average FPS: %.1f\n", stream.AverageFPS)
	str += fmt.Sprintf("Bitrate: %.1f\n", stream.Bitrate)
	str += fmt.Sprintf("Codec: %s\n", stream.Codec)
	str += fmt.Sprintf("Type: %s\n", stream.Type)
	str += fmt.Sprintf("Created At: %s\n", stream.CreatedAt)
	str += fmt.Sprintf("Updated At: %s", stream.UpdatedAt)
	return strings.Split(str, "\n")
}

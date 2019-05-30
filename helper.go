package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func errCheck(message string, err error) {
	if err != nil {
		fmt.Println("ERROR", message, err)
	}
}

func log(messages ...interface{}) {
	fmt.Println(messages...)
}

// PROD ENV WILL WATCH ALL CHANNELS EXCEPT HACK
// ELSE ONLY HACK WILL BE USED
func channelCheck(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	channel, _ := session.Channel(message.ChannelID)
	switch os.Getenv("ENV") {
	case "PROD":
		if channel.Name == "hack" {
			return false
		}
		return true
	default:
		if channel.Name != "hack" {
			return false
		}
		return true
	}
}

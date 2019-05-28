package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var botID string
var commandPrefix string

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	errCheck("could not connect to discord:", err)

	botID, err := session.User("@me")
	errCheck("ERROR could not get user info", err)

	fmt.Println("Logged in as", botID)

	session.AddHandler(onReady)

	err = session.Open()
	errCheck("could not open connection to Discord", err)
	defer session.Close()
	defer session.UpdateStatus(1, "")

	commandPrefix = "!"

	<-make(chan struct{})
}

func errCheck(message string, err error) {
	if err != nil {
		fmt.Println("ERROR", message, err)
	}
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) {
	session.UpdateStatus(0, "Introduction to Basics")
}

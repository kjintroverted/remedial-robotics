package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botID string
var commandPrefix string
var voteDuration time.Duration

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	errCheck("could not connect to discord:", err)

	botID, err := session.User("@me")
	errCheck("ERROR could not get user info", err)

	fmt.Println("Logged in as", botID)

	commandPrefix = "!"

	session.AddHandler(onReady)
	session.AddHandler(onCommand)

	err = session.Open()
	errCheck("could not open connection to Discord", err)
	defer session.Close()
	defer session.UpdateStatus(1, "")

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

func onCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	if len(message.Content) < 2 || message.Content[0:1] != "!" {
		return
	}

	command := strings.Split(message.Content, " ")[0][1:]
	switch command {
	case "vote":
		onVote(session, message)
	}
}

func onVote(session *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println("Let's put it to a vote:", message.Content)
	err := session.MessageReactionAdd(message.ChannelID, message.ID, "👍")
	errCheck("failed to react", err)
	err = session.MessageReactionAdd(message.ChannelID, message.ID, "👎")
	errCheck("failed to react", err)
	go checkVote(session, message)
}

func checkVote(session *discordgo.Session, message *discordgo.MessageCreate) {
	time.Sleep(10 * time.Second)
	users, err := session.MessageReactions(message.ChannelID, message.ID, "👍", 100)
	yea := len(users) - 1
	users, err = session.MessageReactions(message.ChannelID, message.ID, "👎", 100)
	ney := len(users) - 1
	score := strconv.Itoa(yea) + "-" + strconv.Itoa(ney)

	errCheck("getting +1:", err)
	text := strings.SplitN(message.Content, " ", 2)[1]
	var result string
	switch {
	case yea == 0 && ney == 0:
		result = "No one cares"
		break
	case yea == ney:
		result = "Its a tie, so idk..."
		break
	case yea > ney:
		result = "in favor of " + text
		break
	case yea < ney:
		result = text + " just isn't gonna work"
	}
	session.ChannelMessageSend(message.ChannelID, "The people have spoken..."+result+" ("+score+")")
}

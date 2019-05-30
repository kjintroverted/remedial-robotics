package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botID string
var commandPrefix string
var voteDuration time.Duration
var shortVoteDuration time.Duration
var deleteDuration time.Duration

func main() {
	// CONNECT TO SERVER
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	errCheck("could not connect to discord:", err)
	botID, err := session.User("@me")
	errCheck("ERROR could not get user info", err)
	log("Logged in as", botID)

	// SET SOME VARS
	commandPrefix = "-"
	voteDuration = 2 * time.Hour
	log("vote duration:", voteDuration)
	shortVoteDuration = 5 * time.Minute
	log("short vote duration:", shortVoteDuration)
	deleteDuration = 24 * time.Hour
	log("delete delay:", deleteDuration)

	// ADD HANDLERS
	session.AddHandler(onReady)
	session.AddHandler(onCommand)
	session.AddHandler(onHelp)

	// OPEN AND STAY OPEN
	err = session.Open()
	errCheck("could not open connection to Discord", err)
	defer session.Close()
	defer session.UpdateStatus(1, "")

	<-make(chan struct{})
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) {
	session.UpdateStatus(0, "Introduction to Basics")
}

func onHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	channel, _ := session.Channel(message.ChannelID)
	if channel.Name == "help" {
		log("Marked for delete:", message.Content)
		go func(session *discordgo.Session, message *discordgo.MessageCreate) {
			// WAIT
			time.Sleep(deleteDuration)
			// DELETE MESSAGE
			session.ChannelMessageDelete(message.ChannelID, message.ID)
		}(session, message)
	}
}

func onCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	if len(message.Content) < 2 || message.Content[0:1] != commandPrefix {
		return
	}

	command := strings.Split(message.Content, " ")[0][1:]
	switch command {
	case "vote":
		onVote(session, message, voteDuration)
		break
	case "quick-vote":
		onVote(session, message, shortVoteDuration)
		break
	default:
		return
	}
}

func onVote(session *discordgo.Session, message *discordgo.MessageCreate, duration time.Duration) {
	log("Let's put it to a vote:", message.Content)

	session.MessageReactionAdd(message.ChannelID, message.ID, "👍")
	session.MessageReactionAdd(message.ChannelID, message.ID, "👎")
	go checkVote(session, message, duration)
}

func checkVote(session *discordgo.Session, message *discordgo.MessageCreate, duration time.Duration) {
	// WAIT FOR VOTES
	time.Sleep(duration)
	log("Times up for", message.Content)

	// GET COUNTS
	users, _ := session.MessageReactions(message.ChannelID, message.ID, "👍", 100)
	yea := len(users) - 1
	users, _ = session.MessageReactions(message.ChannelID, message.ID, "👎", 100)
	ney := len(users) - 1
	score := strconv.Itoa(yea) + "-" + strconv.Itoa(ney)

	// GET VOTE CONTEXT
	text := strings.SplitN(message.Content, " ", 2)[1]

	// CREATE RESULT MESSAGE
	var result string
	switch {
	case yea == 0 && ney == 0:
		result = "No one cares about " + text
		break
	case yea == ney:
		result = "Its a tie for " + text + ", so idk..."
		break
	case yea > ney:
		result = "in favor of " + text
		break
	case yea < ney:
		result = text + " just isn't gonna work"
	}

	// SEND MESSAGE TO CHANNEL
	session.ChannelMessageSend(message.ChannelID, "The people have spoken..."+result+" ("+score+")")
}

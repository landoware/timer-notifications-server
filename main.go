package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token string

func main() {
	discord, err := discordgo.New("Bot" + Token)

	if err != nil {
		fmt.Println("Failed to create discord session, ", err)
		return
	}

	// Register the messageCreate function as the handler for MessageCreate events
	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsDirectMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("Failed to open connection, ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, userId string, message string) {
	channel, err := s.UserChannelCreate(userId)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		fmt.Println("error creating channel:", err)
		return
	}

	// Send it
	_, err = s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		fmt.Println("error sending DM message:", err)
	}
}

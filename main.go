package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	Token := ""
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Failed to create discord session, ", err)
		return
	}

	discord.Identify.Intents = discordgo.IntentsDirectMessages

	discord.AddHandler(onboard)

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

func onboard(s *discordgo.Session, e *discordgo.IntegrationCreate) {
	fmt.Println("IntegrationCreate")
	channel, err := s.UserChannelCreate(e.ID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		return
	}

	// Send the plugin setup message
	_, err = s.ChannelMessageSend(channel.ID, "Put this in the Discord ID field of the plugin: `"+e.ID+"`")
	if err != nil {
		fmt.Println("error sending DM message:", err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("TOKEN")

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Failed to create discord session, ", err)
		return
	}

	discord.Identify.Intents = discordgo.IntentsDirectMessages

	initCommand := discordgo.ApplicationCommand{
		Name:        "initialize",
		Description: "Set up notifications for farm runs and setup the RuneLite integration.",
	}

	_, err = discord.ApplicationCommandCreate("1511535737651335199", "", &initCommand)
	if err != nil {
		fmt.Println("Failed to create command, ", err)
		return
	}

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
	err = discord.Close()
	if err != nil {
		fmt.Println("Failed to close connection, ", err)
		return
	}
}

func onboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ApplicationCommandData().Name != "initialize" {
		return
	}

	var userID string
	if i.User != nil {
		userID = i.User.ID
	} else if i.Member != nil && i.Member.User != nil {
		userID = i.Member.User.ID
	} else {
		fmt.Println("error resolving interaction user")
		return
	}

	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		fmt.Println("error creating channel:", err)
		return
	}

	// Send the plugin setup message
	_, err = s.ChannelMessageSend(channel.ID, "Put this in the Discord ID field of the plugin: `"+userID+"`")
	if err != nil {
		fmt.Println("error sending DM message:", err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I sent you a DM with setup instructions.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("error responding to interaction:", err)
	}
}

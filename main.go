package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded: %v", err)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is required")
	}

	applicationID := os.Getenv("APPLICATION_ID")
	if applicationID == "" {
		log.Fatal("APPLICATION_ID is required")
	}

	commandGuildID := os.Getenv("COMMAND_GUILD_ID")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	scheduler := NewScheduler(discord)

	discord.Identify.Intents = discordgo.IntentsDirectMessages
	discord.AddHandler(newInteractionHandler(scheduler))

	commands := applicationCommands()
	if _, err := discord.ApplicationCommandBulkOverwrite(applicationID, "", commands); err != nil {
		log.Fatalf("Failed to register global commands: %v", err)
	}
	log.Printf("Registered %d global commands", len(commands))

	if commandGuildID != "" {
		if _, err := discord.ApplicationCommandBulkOverwrite(applicationID, commandGuildID, commands); err != nil {
			log.Fatalf("Failed to register guild commands for %s: %v", commandGuildID, err)
		}
		log.Printf("Registered %d guild commands for %s", len(commands), commandGuildID)
	}

	if err := discord.Open(); err != nil {
		log.Fatalf("Failed to open Discord connection: %v", err)
	}
	defer func() {
		if err := discord.Close(); err != nil {
			log.Printf("Failed to close Discord connection: %v", err)
		}
	}()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: NewAPIHandler(scheduler),
	}

	go func() {
		log.Printf("HTTP server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shut down HTTP server cleanly: %v", err)
	}
}

func applicationCommands() []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{{
		Name:        "initialize",
		Description: "Set up notifications for farm runs and setup the RuneLite integration.",
	}}

	for _, cropGroup := range sortedCropGroups() {
		command := &discordgo.ApplicationCommand{
			Name:        string(cropGroup),
			Description: fmt.Sprintf("Schedule a %s notification.", cropGroup.DisplayName()),
		}

		crops := cropsForGroup(cropGroup)
		if len(crops) > 1 {
			cropChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(crops))
			for _, crop := range crops {
				cropChoices = append(cropChoices, &discordgo.ApplicationCommandOptionChoice{
					Name:  crop.Name,
					Value: crop.Value,
				})
			}

			command.Options = []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "crop",
				Description: "Which crop you just planted.",
				Required:    cropOptionRequired(cropGroup),
				Choices:     cropChoices,
			}}
		}

		commands = append(commands, command)
	}

	return commands
}

func newInteractionHandler(scheduler *Scheduler) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		commandName := i.ApplicationCommandData().Name
		if commandName == "initialize" {
			handleInitializeCommand(s, i)
			return
		}

		cropGroup := CropGroup(commandName)
		if err := cropGroup.Validate(); err != nil {
			return
		}

		handleCropCommand(s, i, scheduler, cropGroup)
	}
}

func handleInitializeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID, err := interactionUserID(i)
	if err != nil {
		log.Printf("error resolving interaction user: %v", err)
		return
	}

	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		log.Printf("error creating channel: %v", err)
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "Put this in the Discord ID field of the plugin: `"+userID+"`")
	if err != nil {
		log.Printf("error sending DM message: %v", err)
		return
	}

	respondToInteraction(s, i, "I sent you a DM with setup instructions.", true)
}

func handleCropCommand(s *discordgo.Session, i *discordgo.InteractionCreate, scheduler *Scheduler, cropGroup CropGroup) {
	if i.GuildID != "" {
		respondToInteraction(s, i, "Use this command in a DM with the bot.", true)
		return
	}

	userID, err := interactionUserID(i)
	if err != nil {
		log.Printf("error resolving interaction user: %v", err)
		return
	}

	data := i.ApplicationCommandData()
	var (
		crop Crop
		ok   bool
	)
	if len(data.Options) == 0 {
		crop, ok = defaultCropForGroup(cropGroup)
	} else {
		crop, ok = cropForGroup(cropGroup, data.Options[0].StringValue())
	}
	if !ok {
		respondToInteraction(s, i, "That crop is not supported for this patch type.", true)
		return
	}

	minutes := int(crop.Duration / time.Minute)
	response, err := scheduler.Reschedule(NotificationRequest{
		UserID:          userID,
		CropGroup:       cropGroup,
		NotifyInMinutes: minutes,
		CropName:        crop.Name,
		CropValue:       crop.Value,
	})
	if err != nil {
		log.Printf("error scheduling notification from slash command: %v", err)
		respondToInteraction(s, i, "Failed to schedule that notification.", true)
		return
	}

	message := fmt.Sprintf("%s notification %s for %s at <t:%d:F>.", cropGroup.DisplayName(), response.Status, crop.Name, response.ScheduledFor.Unix())
	respondToInteraction(s, i, message, false)
}

func interactionUserID(i *discordgo.InteractionCreate) (string, error) {
	if i.User != nil {
		return i.User.ID, nil
	}
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID, nil
	}

	return "", fmt.Errorf("interaction user not found")
}

func respondToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string, ephemeral bool) {
	data := &discordgo.InteractionResponseData{Content: content}
	if ephemeral {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	}); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

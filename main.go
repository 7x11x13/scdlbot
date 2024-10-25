package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"regexp"

	"github.com/7x11x13/scdlbot/download"
	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = os.Getenv("BOT_TOKEN")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

var (
	dmPermission                   = false

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "scdl",
			Description: "Download song from SoundCloud as a Discord embed",
			DMPermission: &dmPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "SoundCloud track URL",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"scdl": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// defer interaction
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})

			handleError := func(err error) {
				log.Println("Error: ", err)
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Internal error",
				})
			}

			// validate url arg
			log.Printf("Validating args")
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			url := optionMap["url"].StringValue()
			valid, _ := regexp.MatchString("^https?://(www\\.|m\\.)?soundcloud\\.com/[^/]+/[^/]+(/[^/]+)?$", url)
			invalid, _ := regexp.MatchString("^https?://(www\\.|m\\.)?soundcloud\\.com/[^/]+/sets/[^/]+$", url)
			valid = valid && !invalid
			if !valid {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Invalid URL: " + url,
				})
				return
			}

			log.Printf("Creating temp dir")
			// create temp dir
			dname, err := os.MkdirTemp("", "scdlbot")
			if err != nil {
				handleError(err)
				return
			}

			defer os.RemoveAll(dname)

			log.Printf("Downloading file")
			// download file
			file, err := download.SoundCloud(dname, url)
			if err != nil {
				handleError(err)
				return
			}

			log.Printf("Uploading file")
			// upload as embed
			f, err := os.Open(*file)
			if err != nil {
				handleError(err)
				return
			}
			defer f.Close()
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Files: []*discordgo.File{
					{
						Name:   f.Name(),
						Reader: f,
					},
				},
			})
		},
	}
)

func init() {
	flag.Parse()
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

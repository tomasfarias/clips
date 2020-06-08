package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token        string
	ClientId     string
	ClientSecret string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.StringVar(&ClientId, "c", "", "Twitch client id")
	flag.StringVar(&ClientSecret, "s", "", "Twitch client secret")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("error creating Discord session, ", err)
	}

	dg.AddHandler(handleCommand)

	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection, ", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.Contains(m.Content, "!clips") {
		return
	}
	log.Printf("Got message %s", m.Content)
	command, err := ParseCommand(m.Content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "I need at least the name of a streamer to look for clips! Use \"!clips help\" for more info.")
		return
	}
	t := NewTwitchApi(ClientId, ClientSecret)

	broadcasters, err := t.GetBroadcastersByName([]string{command.Broadcaster})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Couldn't find a streamer named \""+command.Broadcaster+"\". Could you check the name and try again?")
		return
	}
	log.Printf("Command: %s", command)
	if command.Title == "" {
		// No filtering, should define some default behavior, latest? most popular?is
		clips, _ := t.GetClipsByBroadcasterId(broadcasters[0].Id, "", "", "", "", 100)

		if len(clips) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Couldn't find any \""+command.Broadcaster+"\" clips.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Here's your clip: "+clips[0].Url)
		return
	}

	targetClip := Clip{
		BroadcasterId: broadcasters[0].Id,
		Title:         command.Title,
	}

	result := t.FindClip(targetClip, matchTitle)
	s.ChannelMessageSend(m.ChannelID, "Found your clip: "+result.Url)

	// switch nResults := len(results); {
	// case nResults == 1:
	// case nResults > 1:
	//	s.ChannelMessageSend(m.ChannelID, "Found many clips")
	//}
	return
}

func matchTitle(clip1, clip2 Clip) bool {
	return clip1.Title == clip2.Title
}

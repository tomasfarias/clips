package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Token, ClientId, ClientSecret string

func main() {

	flag.StringVar(&Token, "t", "a-token", "Bot token")
	flag.StringVar(&ClientId, "c", "a-client-id", "Twitch client id")
	flag.StringVar(&ClientSecret, "s", "a-client-secret", "Twitch client secret")
	flag.Parse()

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

func handleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	help := `Search for Twitch clips.
Usage: !clips streamer "title" creator start_date end_date
Required arguments:
	- streamer: The name of the Twitch channel/streamer where to look for clips.
Optional arguments:
	- title: Find a clip with a specific title. **Must** be enclosed in double quotes.
	- creator: Filter by clips created by a specific user. If defined, **must** always come after streamer argument.
	- start_date: Look for a clip created from this date onwards. Defaults to **1 week ago**. Format as YYYY-MM-DD. Will make things run faster if used.
	- end_date: Look for a clip created before this date. Format as YYYY-MM-DD. Will make things run faster if used.`
	s.ChannelMessageSend(m.ChannelID, help)
	return
}

func handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.Contains(m.Content, "!clips") {
		return
	}
	log.Printf("Got message %s", m.Content)
	command, err := ParseCommand(m.Content)
	if command.Broadcaster == "help" {
		handleHelpCommand(s, m)
		return
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "I need at least the name of a streamer to look for clips! Use \"!clips help\" for more info.")
		return
	}
	t := NewTwitchApi(ClientId, ClientSecret, true)

	broadcasters, err := t.GetBroadcastersByName([]string{command.Broadcaster})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Couldn't find a streamer named \""+command.Broadcaster+"\". Could you check the name and try again?")
		return
	}
	log.Printf("Command: %s", command)
	targetClip := Clip{
		BroadcasterId: broadcasters[0].Id,
		Title:         command.Title,
		StartedAt:     command.StartedAt,
		EndedAt:       command.EndedAt,
		CreatorName:   command.Creator,
	}
	if targetClip.StartedAt.IsZero() {
		targetClip.StartedAt = time.Now().AddDate(0, 0, -7)
	}
	matchFunc := matchMany(matchTitle, matchCreator)

	var result Clip
	if targetClip.Title == "" || targetClip.CreatorName == "" {
		// There may be many clips with the same creator or title, so we look for the most popular one
		result = t.FindMostPopularClip(targetClip, matchFunc)
	} else {
		// Otherwise, we're looking for a specific clip
		result = t.FindClip(targetClip, matchFunc)
	}

	if result == targetClip {
		s.ChannelMessageSend(m.ChannelID, "I couldn't find a clip that matches your search.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Found your clip: "+result.Url)
	return
}

func matchMany(funcs ...func(Clip, Clip) bool) func(Clip, Clip) bool {
	return func(clip1, clip2 Clip) bool {
		result := true
		for _, f := range funcs {
			result = result && f(clip1, clip2)
		}

		return result
	}
}

func matchTitle(clip1, clip2 Clip) bool {
	return strings.Contains(strings.ToLower(clip1.Title), strings.ToLower(clip2.Title))
}

func matchCreator(clip1, clip2 Clip) bool {
	return strings.Contains(strings.ToLower(clip1.CreatorName), strings.ToLower(clip2.CreatorName))
}

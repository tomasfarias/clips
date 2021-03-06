package main

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Command represents a clips bot command
type Command struct {
	Broadcaster string
	Creator     string
	StartedAt   time.Time
	SubCommand  string
	EndedAt     time.Time
	Title       string
	Top         int
}

// ParseCommand parses a Discord message string to a Command
func ParseCommand(args string) (Command, error) {
	if !strings.HasPrefix(args, "!clips") {
		return Command{}, errors.New("command: invalid must start with \"!clips\"")
	}
	command := Command{}
	args = removeSubStrings(args, []string{"!clips"})

	start, end, stringDates := parseDates(args)
	if len(stringDates) > 0 {
		args = removeSubStrings(args, stringDates)
		command.StartedAt = start
		if len(stringDates) > 1 {
			command.EndedAt = end
		}
	}

	start, end, simpleStringDate := parseSimpleDate(args)
	if simpleStringDate != "" {
		args = removeSubStrings(args, []string{simpleStringDate})
		command.StartedAt = start
		command.EndedAt = end
	}

	title := strings.FieldsFunc(args, splitQuote)
	if title != nil && len(title) >= 2 && title[1] != "" {
		args = removeSubStrings(args, []string{"\"" + title[1] + "\""})
		command.Title = title[1]
	}

	words := strings.Fields(args)
	log.Printf("Parsing args: %s", words)
	if len(words) > 0 {
		switch potentialSubCommand := words[0]; {
		case potentialSubCommand == "help":
			command.SubCommand = potentialSubCommand
			words = words[1:]
		case strings.HasPrefix(potentialSubCommand, "top"):
			command.SubCommand = "top"
			n := potentialSubCommand[3:]
			if n == "" {
				command.Top = 10 // Default is top 10
			} else {
				top, err := strconv.Atoi(n)
				if err != nil {
					log.Printf("strconv error, : %s", err)
				}

				command.Top = top
			}
			words = words[1:]
		}
	}

	switch nParams := len(words); {
	case nParams == 0:
		return command, nil
	case nParams >= 2:
		command.Creator = words[1]
		fallthrough
	case nParams >= 1:
		command.Broadcaster = words[0]
	}

	return command, nil
}

func removeSubStrings(target string, toRemove []string) string {
	for _, remove := range toRemove {
		start := strings.Index(target, remove)
		end := start + len(remove)
		target = target[:start] + target[end:]
	}
	return target
}

func parseDates(args string) (time.Time, time.Time, []string) {
	regex, err := regexp.Compile(`[12][0-9]{3}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])`)
	if err != nil {
		log.Fatal("regex error: ", err)
	}
	matched := regex.FindAllString(args, 2)
	if len(matched) == 0 {
		return time.Time{}, time.Time{}, []string{}
	}

	start, err := time.Parse("2006-01-02", matched[0])
	if err != nil {
		log.Fatal("time error: ", err)
	}

	end := time.Time{}
	if len(matched) > 1 {
		end, err = time.Parse("2006-01-02", matched[1])
		if err != nil {
			log.Fatal("time error: ", err)
		}
	}

	return start, end, matched
}

func parseSimpleDate(args string) (time.Time, time.Time, string) {
	regex, err := regexp.Compile(`(?P<Number>\d+)(?P<Unit>d|m|y)`)
	if err != nil {
		log.Fatal("regex error: ", err)
	}
	matched := regex.FindStringSubmatch(args)
	if len(matched) == 0 {
		return time.Time{}, time.Time{}, ""
	}

	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()) // Twitch ignores everything after minute
	value, err := strconv.Atoi(matched[1])
	if err != nil {
		log.Fatal("strconv error: ", err)
	}

	switch matched[2] {
	case "d":
		return currentDate.AddDate(0, 0, -value), currentDate, matched[0]
	case "m":
		return currentDate.AddDate(0, -value, 0), currentDate, matched[0]
	case "y":
		return currentDate.AddDate(-value, 0, 0), currentDate, matched[0]
	}
	return time.Time{}, time.Time{}, ""
}

func splitQuote(r rune) bool {
	quotes := []rune{'"', '\''}

	for _, s := range quotes {
		if r == s {
			return true
		}
	}
	return false
}

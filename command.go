package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

type Command struct {
	Broadcaster string
	Creator     string
	StartedAt   time.Time
	EndedAt     time.Time
	Title       string
}

func ParseCommand(args string) (Command, error) {
	if !strings.Contains(args, "!clips") {
		return Command{}, errors.New("command: invalid must start with \"!clips\"")
	}
	command := Command{}

	stringDates := parseDates(args)
	dates := []time.Time{}
	for _, date := range stringDates {
		parsed, err := time.Parse("2006-01-02", date)
		if err != nil {
			return Command{}, err
		}
		dates = append(dates, parsed)
	}
	if len(stringDates) > 0 {
		args = removeSubStrings(args, stringDates)
		command.StartedAt = dates[0]
		command.EndedAt = dates[1]
	}

	title := strings.FieldsFunc(args, splitQuote)
	if title != nil && len(title) >= 2 && title[1] != "" {
		args = removeSubStrings(args, []string{"\"" + title[1] + "\""})
		command.Title = title[1]
	}

	words := strings.Fields(args)
	log.Printf("Parsing args: %s", words)
	switch nParams := len(words); {
	case nParams == 1:
		return Command{}, errors.New("command: no arguments passed")
	case nParams > 2:
		command.Creator = words[2]
		fallthrough
	case nParams > 1:
		command.Broadcaster = words[1]
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

func parseDates(args string) []string {
	regex, err := regexp.Compile("[12][0-9]{3}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])")
	if err != nil {
		log.Fatal("regex error: ", err)
	}
	matched := regex.FindAllString(args, -1)

	return matched
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

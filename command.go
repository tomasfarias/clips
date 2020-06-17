package main

import (
	"errors"
	"log"
	"regexp"
	"strconv"
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
	if !strings.HasPrefix(args, "!clips") {
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
		if len(stringDates) > 1 {
			command.EndedAt = dates[1]
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
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
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

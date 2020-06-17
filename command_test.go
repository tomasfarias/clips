package main

import (
    "testing"
    "time"
)

func TestParseCommand(t *testing.T) {
    inputCommand := "!clips Streamer \"Super funny clip!\" Creator 2020-05-30 2020-06-30"
    result, err := ParseCommand(inputCommand)
    if err != nil {
        t.Errorf("Got an error while parsing test command: %s", err)
    }

    if result.Broadcaster != "Streamer" {
        t.Errorf("Broadcaster not properly parsed: expected \"Streamer\" got %s", result.Broadcaster)
    }
    if result.Title != "Super funny clip!" {
        t.Errorf("Title not properly parsed: expected \"Super funny clip!\" got %s", result.Title)
    }
    if result.Creator != "Creator" {
        t.Errorf("Creator not properly parsed: expected \"Creator\" got %s", result.Creator)
    }
    
    started, _ := time.Parse("2006-01-02", "2020-05-30")
    if result.StartedAt != started {
        t.Errorf("StartedAt not properly parsed: expected \"%s\" got %s", started, result.StartedAt)
    }
    
    ended, _ := time.Parse("2006-01-02", "2020-06-30")
    if result.EndedAt != ended {
        t.Errorf("EndedAt not properly parsed: expected \"%s\" got %s", ended, result.EndedAt)
    }
}

func TestParseInvalidCommand(t *testing.T) {
    inputCommand := "This is not a valid clips command even if !clips is in it"
    result, err := ParseCommand(inputCommand)
    if err == nil {
        t.Errorf("Should have returned an error for invalid command, got: %v", result)
    }
}

func TestParseCommandOnlyTitle(t *testing.T) {
    inputCommand := "!clips Streamer \"Super funny clip!\""
    result, err := ParseCommand(inputCommand)
    if err != nil {
        t.Errorf("Got an error while parsing test command: %s", err)
    }

    if result.Broadcaster != "Streamer" {
        t.Errorf("Broadcaster not properly parsed: expected \"Streamer\" got %s", result.Broadcaster)
    }
    if result.Title != "Super funny clip!" {
        t.Errorf("Title not properly parsed: expected \"Super funny clip!\" got %s", result.Title)
    }
    if result.Creator != "" {
        t.Errorf("Creator not properly parsed: expected \"\" got %s", result.Creator)
    }
    
    emptyTime := time.Time{}
    if result.StartedAt != emptyTime {
        t.Errorf("StartedAt not properly parsed: expected \"%s\" got %s", emptyTime, result.StartedAt)
    }
    
    if result.EndedAt != emptyTime {
        t.Errorf("EndedAt not properly parsed: expected \"%s\" got %s", emptyTime, result.EndedAt)
    }
}

func TestParseCommandOnlyStreamer(t *testing.T) {
    inputCommand := "!clips Streamer"
    result, err := ParseCommand(inputCommand)
    if err != nil {
        t.Errorf("Got an error while parsing test command: %s", err)
    }

    if result.Broadcaster != "Streamer" {
        t.Errorf("Broadcaster not properly parsed: expected \"Streamer\" got %s", result.Broadcaster)
    }
    if result.Title != "" {
        t.Errorf("Title not properly parsed: expected \"\" got %s", result.Title)
    }
    if result.Creator != "" {
        t.Errorf("Creator not properly parsed: expected \"\" got %s", result.Creator)
    }

    emptyTime := time.Time{}
    if result.StartedAt != emptyTime {
        t.Errorf("StartedAt not properly parsed: expected \"%s\" got %s", emptyTime, result.StartedAt)
    }
    
    if result.EndedAt != emptyTime {
        t.Errorf("EndedAt not properly parsed: expected \"%s\" got %s", emptyTime, result.EndedAt)
    }
}

func TestParseCommandSimpleDays(t *testing.T) {
    inputCommand := "!clips Streamer \"Super funny clip!\" Creator 7d"
    result, err := ParseCommand(inputCommand)
    if err != nil {
        t.Errorf("Got an error while parsing test command: %s", err)
    }

    if result.Broadcaster != "Streamer" {
        t.Errorf("Broadcaster not properly parsed: expected \"Streamer\" got %s", result.Broadcaster)
    }
    if result.Title != "Super funny clip!" {
        t.Errorf("Title not properly parsed: expected \"Super funny clip!\" got %s", result.Title)
    }
    if result.Creator != "Creator" {
        t.Errorf("Creator not properly parsed: expected \"Creator\" got %s", result.Creator)
    }
    
    now := time.Now()
    currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    if result.StartedAt != currentDate.AddDate(0, 0, -7) {
        t.Errorf("StartedAt not properly parsed: expected \"%s\" got %s", currentDate.AddDate(0, 0, -7), result.StartedAt)
    }
    if result.EndedAt != currentDate {
        t.Errorf("EndedAt not properly parsed: expected \"%s\" got %s", currentDate, result.EndedAt)
    }
}

func TestParseCommandSimpleMonths(t *testing.T) {
    inputCommand := "!clips Streamer \"Super funny clip!\" Creator 1m"
    result, err := ParseCommand(inputCommand)
    if err != nil {
        t.Errorf("Got an error while parsing test command: %s", err)
    }

    if result.Broadcaster != "Streamer" {
        t.Errorf("Broadcaster not properly parsed: expected \"Streamer\" got %s", result.Broadcaster)
    }
    if result.Title != "Super funny clip!" {
        t.Errorf("Title not properly parsed: expected \"Super funny clip!\" got %s", result.Title)
    }
    if result.Creator != "Creator" {
        t.Errorf("Creator not properly parsed: expected \"Creator\" got %s", result.Creator)
    }
    
    now := time.Now()
    currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    if result.StartedAt != currentDate.AddDate(0, -1, 0) {
        t.Errorf("StartedAt not properly parsed: expected \"%s\" got %s", currentDate.AddDate(0, -1, 0), result.StartedAt)
    }
    if result.EndedAt != currentDate {
        t.Errorf("EndedAt not properly parsed: expected \"%s\" got %s", currentDate, result.EndedAt)
    }
}

func TestParseCommandSimpleYear(t *testing.T) {
    inputCommand := "!clips Streamer \"Super funny clip!\" Creator 2y"
    result, err := ParseCommand(inputCommand)
    if err != nil {
        t.Errorf("Got an error while parsing test command: %s", err)
    }

    if result.Broadcaster != "Streamer" {
        t.Errorf("Broadcaster not properly parsed: expected \"Streamer\" got %s", result.Broadcaster)
    }
    if result.Title != "Super funny clip!" {
        t.Errorf("Title not properly parsed: expected \"Super funny clip!\" got %s", result.Title)
    }
    if result.Creator != "Creator" {
        t.Errorf("Creator not properly parsed: expected \"Creator\" got %s", result.Creator)
    }
    
    now := time.Now()
    currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    if result.StartedAt != currentDate.AddDate(-2, 0, 0) {
        t.Errorf("StartedAt not properly parsed: expected \"%s\" got %s", currentDate.AddDate(-2, 0, 0), result.StartedAt)
    }
    if result.EndedAt != currentDate {
        t.Errorf("EndedAt not properly parsed: expected \"%s\" got %s", currentDate, result.EndedAt)
    }
}

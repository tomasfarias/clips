package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewTwitchApi(t *testing.T) {
	twitch := NewTwitchApi("client-id", "client-secret", false)

	if twitch.ClientId != "client-id" {
		t.Errorf("ClientId not properly set, expected \"client-id\" got %s", twitch.ClientId)
	}
	if twitch.ClientSecret != "client-secret" {
		t.Errorf("ClientSecret not properly set, expected \"client-secret\" got %s", twitch.ClientSecret)
	}
	if twitch.AccessToken != "" {
		t.Errorf("AccessToken not properly set, expected \"\" got %s", twitch.AccessToken)
	}

	baseUrl := url.URL{Scheme: "https", Host: "api.twitch.tv"}
	if twitch.BaseUrl != baseUrl {
		t.Errorf("BaseUrl not properly set, expected \"%v\" got %v", baseUrl, twitch.BaseUrl)
	}
	authUrl := url.URL{Scheme: "https", Host: "id.twitch.tv", Path: "/oauth2/token"}
	if twitch.AuthUrl != authUrl {
		t.Errorf("AuthUrl not properly set, expected \"%v\" got %v", authUrl, twitch.AuthUrl)
	}
}

func TestPrepareQuery(t *testing.T) {
	m := make(map[string]string)
	m["test_id"] = "some-id"
	m["test_secret"] = "some-secret"
	m["grant_type"] = "some-type"

	query := url.Values{}
	queryString := prepareQuery(query, m)
	expected := "grant_type=some-type&test_id=some-id&test_secret=some-secret"
	if queryString != expected {
		t.Errorf("Query string not properly encoded, expected \"%s\" got %s", expected, queryString)
	}
}

func TestPrepareRequest(t *testing.T) {
	twitch := NewTwitchApi("client-id", "client-secret", false)
	twitch.AccessToken = "some-token"
	req := twitch.prepareRequest("GET", "https://some.fancy/url")

	if req.Header.Get("Client-ID") != "client-id" {
		t.Errorf("Client-ID Header not properly set by prepareRequest, expected \"client-id\" got %s", req.Header.Get("Client-ID"))
	}
	if req.Header.Get("Authorization") != "Bearer some-token" {
		t.Errorf("Authorization Header not properly set by prepareRequest, expected \"Bearer some-token\" got %s", req.Header.Get("Authorization"))
	}
	if req.Method != "GET" {
		t.Errorf("Request method not properly set by prepareRequest, expected \"GET\" got %s", req.Method)
	}
	if req.Host != "some.fancy" {
		t.Errorf("Request Host not properly set by prepareRequest, expected \"some.fancy\" got %s", req.Host)
	}
}

func TestSetAuthToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(authHandler))

	twitch := NewTwitchApi("client-id", "client-secret", false)
	mockUrl, _ := url.Parse(ts.URL)
	twitch.AuthUrl = *mockUrl

	twitch.SetAuthToken()
	if twitch.AccessToken != "my-test-token" {
		t.Errorf("AccessToken not properly set by SetAuthToken, expected \"my-test-token\" got %s", twitch.AccessToken)
	}
	ts.Close()
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	res := TokenResponse{AccessToken: "my-test-token"}
	json.NewEncoder(w).Encode(res)
}

func TestGetBroadcastersByName(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(broadcastersHandler))

	twitch := NewTwitchApi("client-id", "client-secret", false)
	mockUrl, _ := url.Parse(ts.URL)
	twitch.BaseUrl = *mockUrl

	broadcasters, _ := twitch.GetBroadcastersByName([]string{"test-login"})
	if broadcasters[0].Id != "test-id" {
		t.Errorf("Broadcaster.Id not correctly returned by GetBroadcastersByName, expected \"test-id\" got %s", broadcasters[0].Id)
	}
	ts.Close()
}

func broadcastersHandler(w http.ResponseWriter, r *http.Request) {
	b := Broadcaster{Id: "test-id"}
	res := BroadcasterResponse{Data: []Broadcaster{b}}
	json.NewEncoder(w).Encode(res)
}

func TestGetClipsByBroadcasterId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(clipsHandler))

	twitch := NewTwitchApi("client-id", "client-secret", false)
	mockUrl, _ := url.Parse(ts.URL)
	twitch.BaseUrl = *mockUrl

	clips, _ := twitch.GetClipsByBroadcasterId("broadcaster", "", "", time.Time{}, time.Time{}, 100)
	if clips[0].Id != "test-id" {
		t.Errorf("Clip.Id not correctly returned by GetBroadcastersByName, expected \"test-id\" got %s", clips[0].Id)
	}
	if clips[0].BroadcasterId != "broadcaster_id=broadcaster&first=100" {
		t.Errorf("Clip.BroadcasterId not correctly returned by GetBroadcastersByName, expected \"broadcaster_id=broadcaster\" got %s", clips[0].BroadcasterId)
	}
	ts.Close()
}

func clipsHandler(w http.ResponseWriter, r *http.Request) {
	b := Clip{Id: "test-id", BroadcasterId: r.URL.RawQuery}
	res := ClipsResponse{Data: []Clip{b}}
	json.NewEncoder(w).Encode(res)
}

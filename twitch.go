package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ClipsResponse struct {
	Data       []Clip `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type Clip struct {
	Id              string `json:"id"`
	Url             string `json:"url"`
	EmbedUrl        string `json:"embed_url"`
	BroadcasterId   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	CreatorId       string `json:"creator_id"`
	CreatorName     string `json:"creator_name"`
	VideoId         string `json:"video_id"`
	GameId          string `json:"game_id"`
	Language        string `json:"language"`
	Title           string `json:"title"`
	ViewCount       int    `json:"view_count"`
	CreatedAt       string `json:"created_at"`
	ThumbnailUrl    string `json:"thumbnail_url"`
}

type BroadcasterResponse struct {
	Data []Broadcaster `json:"data"`
}

type Broadcaster struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

type twitchApi struct {
	ClientId     string
	ClientSecret string
	AccessToken  string
	BaseUrl      url.URL
	Client       *http.Client
}

func NewTwitchApi(clientId string, clientSecret string) twitchApi {
	t := twitchApi{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		BaseUrl:      url.URL{Scheme: "https", Host: "api.twitch.tv"},
		Client:       &http.Client{},
	}
	t.SetAuthToken()

	return t
}

func (t twitchApi) GetBroadcastersByName(broadcasterNames []string) ([]Broadcaster, error) {
	endpoint := t.BaseUrl
	endpoint.Path = "/helix/users"

	q := endpoint.Query()
	for _, broadcasterName := range broadcasterNames {
		q.Set("login", broadcasterName)
	}
	endpoint.RawQuery = q.Encode()

	req := t.prepareRequest("GET", endpoint.String())

	log.Printf("Request: %s", req)
	jsonResponse, err := t.Client.Do(req)
	if err != nil {
		log.Fatal("request failed: ", err)
	}

	resp := BroadcasterResponse{}
	json.NewDecoder(jsonResponse.Body).Decode(&resp)

	if len(resp.Data) == 0 {
		return nil, errors.New("twitch: no broadcasters found")
	}
	return resp.Data, nil
}

func (t twitchApi) prepareRequest(method string, endpoint string) *http.Request {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		log.Fatal("failed to create request: ", err)
	}

	req.Header.Add("Client-ID", t.ClientId)
	req.Header.Add("Authorization", "Bearer "+t.AccessToken)

	return req
}

func (t twitchApi) GetClipsByBroadcasterId(broadcasterId string, after string, before string, ended_at string, started_at string, first int) ([]Clip, string) {
	endpoint := t.BaseUrl
	endpoint.Path = "/helix/clips"

	q := endpoint.Query()
	q.Set("broadcaster_id", broadcasterId)
	q.Set("first", strconv.Itoa(first))
	// q.Set("started_at", started_at)
	if after != "" {
		q.Set("after", after)
	}
	endpoint.RawQuery = q.Encode()

	req := t.prepareRequest("GET", endpoint.String())
	log.Printf("Request: %s", req)

	jsonResponse, err := t.Client.Do(req)
	if err != nil {
		log.Fatal("request failed: ", err)
	}

	resp := ClipsResponse{}
	json.NewDecoder(jsonResponse.Body).Decode(&resp)

	return resp.Data, resp.Pagination.Cursor
}

func (t twitchApi) FindClip(targetClip Clip, matchFunc func(Clip, Clip) bool) Clip {
	clips, cursor := t.GetClipsByBroadcasterId(targetClip.BroadcasterId, "", "", "", "", 100) // Need to figure out how to pass dates

	for {
		if len(clips) == 0 {
			// return the same clip passed if nothing is found?
			return targetClip
		}

		for _, clip := range clips {
			if matchFunc(clip, targetClip) {
				return clip
			}
		}

		clips, cursor = t.GetClipsByBroadcasterId(targetClip.BroadcasterId, cursor, "", "", "", 100) // Same here
	}
}

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scopes       []string `json:"scopes"`
	TokenType    string   `json:"token_type"`
}

func (t *twitchApi) SetAuthToken() {
	endpoint := url.URL{Scheme: "https", Host: "id.twitch.tv", Path: "/oauth2/token"}
	q := endpoint.Query()
	q.Set("client_id", t.ClientId)
	q.Set("client_secret", t.ClientSecret)
	q.Set("grant_type", "client_credentials")
	endpoint.RawQuery = q.Encode()

	req := t.prepareRequest("POST", endpoint.String())

	jsonResponse, err := t.Client.Do(req)
	if err != nil {
		log.Fatal("request failed: ", err)
	}

	resp := TokenResponse{}
	json.NewDecoder(jsonResponse.Body).Decode(&resp)

	t.AccessToken = resp.AccessToken
}

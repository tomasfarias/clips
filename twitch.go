package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ClipsResponse represents a response from a request to Twitch's Get Clips
type ClipsResponse struct {
	Data       []Clip `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// Clip represents a Twitch clip
type Clip struct {
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	EmbedURL        string    `json:"embed_url"`
	BroadcasterID   string    `json:"broadcaster_id"`
	BroadcasterName string    `json:"broadcaster_name"`
	CreatorID       string    `json:"creator_id"`
	CreatorName     string    `json:"creator_name"`
	VideoID         string    `json:"video_id"`
	GameID          string    `json:"game_id"`
	Language        string    `json:"language"`
	Title           string    `json:"title"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       string    `json:"created_at"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	StartedAt       time.Time `json:",omitempty"`
	EndedAt         time.Time `json:",omitempty"`
}

// BroadcasterResponse represents a response from a request to Twitch's Get users
type BroadcasterResponse struct {
	Data []Broadcaster `json:"data"`
}

// Broadcaster represents a Twitch user, used for broadcasters
type Broadcaster struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

// TwitchAPI holds all configuration needed for a Twitch API connection
type TwitchAPI struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	BaseURL      url.URL
	AuthURL      url.URL
	Client       *http.Client
}

// NewTwitchAPI returns a new TwitchAPI after setting the access token
func NewTwitchAPI(clientID string, clientSecret string, setAuth bool) TwitchAPI {
	t := TwitchAPI{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      url.URL{Scheme: "https", Host: "api.twitch.tv"},
		AuthURL:      url.URL{Scheme: "https", Host: "id.twitch.tv", Path: "/oauth2/token"},
		Client:       &http.Client{},
	}

	if setAuth == true {
		t.SetAuthToken()
	}

	return t
}

// GetBroadcastersByName finds a Broadcaster with given names
func (t TwitchAPI) GetBroadcastersByName(broadcasterNames []string) ([]Broadcaster, error) {
	endpoint := t.BaseURL
	endpoint.Path = "/helix/users"

	q := endpoint.Query()
	for _, broadcasterName := range broadcasterNames {
		q.Set("login", broadcasterName)
	}
	endpoint.RawQuery = q.Encode()

	req := t.prepareRequest("GET", endpoint.String())

	log.Printf("Request: %v", req)
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

func (t TwitchAPI) prepareRequest(method string, endpoint string) *http.Request {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		log.Fatal("failed to create request: ", err)
	}

	req.Header.Add("Client-ID", t.ClientID)
	req.Header.Add("Authorization", "Bearer "+t.AccessToken)

	return req
}

// GetClipsByBroadcasterID finds clips from a given broadcaster
func (t TwitchAPI) GetClipsByBroadcasterID(broadcasterID string, after string, before string, endedAt time.Time, startedAt time.Time, first int) ([]Clip, string) {
	endpoint := t.BaseURL
	endpoint.Path = "/helix/clips"

	m := make(map[string]string)
	m["broadcaster_id"] = broadcasterID
	m["first"] = strconv.Itoa(first)
	if !startedAt.IsZero() {
		m["started_at"] = startedAt.Format(time.RFC3339)
	}
	if !endedAt.IsZero() {
		m["ended_at"] = endedAt.Format(time.RFC3339)
	}
	if after != "" {
		m["after"] = after
	}
	if before != "" {
		m["before"] = before
	}
	q := endpoint.Query()
	endpoint.RawQuery = prepareQuery(q, m)

	req := t.prepareRequest("GET", endpoint.String())
	log.Printf("Request: %v", req)

	jsonResponse, err := t.Client.Do(req)
	if err != nil {
		log.Fatal("request failed: ", err)
	}

	resp := ClipsResponse{}
	json.NewDecoder(jsonResponse.Body).Decode(&resp)

	return resp.Data, resp.Pagination.Cursor
}

// FindClip compares Twitch clips to targetClip using matchFunc
func (t TwitchAPI) FindClip(targetClip Clip, matchFunc func(Clip, Clip) bool) Clip {
	clips, cursor := t.GetClipsByBroadcasterID(targetClip.BroadcasterID, "", "", targetClip.EndedAt, targetClip.StartedAt, 100)

	for {
		if len(clips) == 0 || cursor == "" {
			// return the same clip passed if nothing is found?
			return targetClip
		}

		for _, clip := range clips {
			if matchFunc(clip, targetClip) {
				return clip
			}
		}

		clips, cursor = t.GetClipsByBroadcasterID(targetClip.BroadcasterID, cursor, "", targetClip.EndedAt, targetClip.StartedAt, 100)
	}
}

// FindMostPopularClip compares Twitch clips to targetClip using matchFunc and returns only the most popular
func (t TwitchAPI) FindMostPopularClip(targetClip Clip, matchFunc func(Clip, Clip) bool) Clip {
	clips, cursor := t.GetClipsByBroadcasterID(targetClip.BroadcasterID, "", "", targetClip.EndedAt, targetClip.StartedAt, 100)
	mostPopular := targetClip

	for {
		if len(clips) == 0 || cursor == "" {
			return mostPopular
		}

		for _, clip := range clips {
			if clip.ViewCount > mostPopular.ViewCount && matchFunc(clip, targetClip) {
				mostPopular = clip
			}
		}

		clips, cursor = t.GetClipsByBroadcasterID(targetClip.BroadcasterID, cursor, "", targetClip.EndedAt, targetClip.StartedAt, 100)
	}
}

// TokenResponse represents a response from the auth endpoint containing an access token
type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scopes       []string `json:"scopes"`
	TokenType    string   `json:"token_type"`
}

// SetAuthToken sets the AccessToken in TwitchAPI
func (t *TwitchAPI) SetAuthToken() {
	endpoint := t.AuthURL

	q := endpoint.Query()
	m := make(map[string]string)
	m["client_id"] = t.ClientID
	m["client_secret"] = t.ClientSecret
	m["grant_type"] = "client_credentials"

	endpoint.RawQuery = prepareQuery(q, m)
	req := t.prepareRequest("POST", endpoint.String())

	jsonResponse, err := t.Client.Do(req)
	if err != nil {
		log.Fatal("request failed: ", err)
	}

	resp := TokenResponse{}
	json.NewDecoder(jsonResponse.Body).Decode(&resp)

	t.AccessToken = resp.AccessToken
}

func prepareQuery(query url.Values, m map[string]string) string {
	for k, v := range m {
		query.Set(k, v)
	}
	return query.Encode()
}

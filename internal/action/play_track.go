package action

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
)

// PlayTrack is an action for getting a user's profile from Spotify
type PlayTrack struct {
    // Inputs
	AuthID string
    DeviceID string
    SpotifyURI string
	Client *http.Client

    // Fetched resources
	Token string
}

// NewPlayTrack constructs and returns a PlayTrack action
func NewPlayTrack(authID, deviceID, spotifyURI string) PlayTrack {
	return PlayTrack{
		AuthID: authID,
        DeviceID: deviceID,
        SpotifyURI: spotifyURI,
		Client: &http.Client{},
	}
}

// Fetch the data needed to process the request
func (a *PlayTrack) Fetch(ctx context.Context) error {
    // Construct request to get access token
	req, err := http.NewRequest(http.MethodGet, constant.URLAPIToken, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

    q := req.URL.Query()
    q.Add("auth_id", a.AuthID)
    req.URL.RawQuery = q.Encode()

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var token model.Token
	if err := json.Unmarshal(body, &token); err != nil {
		return err
	}

	a.Token = token.AccessToken
	return nil
}

// Execute the request
func (a *PlayTrack) Execute(ctx context.Context) error {
    // Construct the request
	reqBody := model.MapCreatePlayerRequest(a.SpotifyURI)
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, constant.URLSpotifyPlayTrack, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + a.Token)

    q := req.URL.Query()
    q.Add("device_id", a.DeviceID)
    req.URL.RawQuery = q.Encode()

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Save the output to the database
func (a *PlayTrack) Save(ctx context.Context) error {
	return nil
}

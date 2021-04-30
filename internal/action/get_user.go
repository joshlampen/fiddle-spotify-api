package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/format"
)

// GetUser is an action for getting a user's profile from Spotify
type GetUser struct {
    // Inputs
	AuthID string
	Client *http.Client

    // Fetched resources
	Token string

    // Outputs
	SpotifyResponse model.SpotifyUser
    DBResponse model.DBUser
}

// NewGetUser constructs and returns a GetUser action
func NewGetUser(authID string) GetUser {
	return GetUser{
		AuthID: authID,
		Client: &http.Client{},
	}
}

// Fetch the data needed to process the request
func (a *GetUser) Fetch(ctx context.Context) error {
    // Construct request to get access token
	req, err := http.NewRequest(http.MethodGet, format.Url(constant.URLAPIToken), nil)
	if err != nil {
		return fmt.Errorf("GetUser - could not create get token request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

    q := req.URL.Query()
    q.Add("auth_id", a.AuthID)
    req.URL.RawQuery = q.Encode()

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("GetUser - get token request failed: %w", err)
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("GetUser - failed to read get token response: %w", err)
	}

	var token model.Token
	if err := json.Unmarshal(body, &token); err != nil {
		return fmt.Errorf("GetUser - failed to unmarshal get token response: %w", err)
	}

	a.Token = token.AccessToken
	return nil
}

// Execute the request
func (a *GetUser) Execute(ctx context.Context) error {
	// Construct the request
	req, err := http.NewRequest(http.MethodGet, constant.URLSpotifyUserProfile, nil)
	if err != nil {
		return fmt.Errorf("GetUser - could not create get user request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + a.Token)

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("GetUser - get user request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("GetUser - failed to read get user response body: %w", err)
	}

	var jsonResp model.SpotifyUser
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return fmt.Errorf("GetUser - failed to unmarshal get user response body: %w", err)
	}

	a.SpotifyResponse = jsonResp
	return nil
}

// Save the output to the database
func (a *GetUser) Save(ctx context.Context) error {
	// Construct the request
	reqBody := model.MapCreateUserRequest(a.Token, a.SpotifyResponse)
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("GetUser - failed to marshal post user request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, format.Url(constant.URLAPIUsers), bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("GetUser - could not create post user request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("GetUser - post user request failed: %w", err)
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("GetUser - failed to read post user response: %w", err)
	}

	var user model.DBUser
	if err := json.Unmarshal(body, &user); err != nil {
		return fmt.Errorf("GetUser - failed to unmarshal post user response: %w", err)
	}

    a.DBResponse = user
	return nil
}

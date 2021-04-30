package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/format"
)

// GetToken is an action for getting an access token from Spotify
type GetToken struct {
    // Inputs
    AuthID string
    AuthCode string
	ClientID string
	ClientSecret string
	GrantType string
	RedirectURL string
	Client *http.Client

    // Outputs
	Response model.Token
}

// NewGetToken constructs and returns a GetToken action
func NewGetToken(authID, authCode, clientID, clientSecret, grantType, redirectURL string) GetToken {
	return GetToken{
        AuthID: authID,
        AuthCode: authCode,
		ClientID: clientID,
		ClientSecret: clientSecret,
		GrantType: grantType,
		RedirectURL: redirectURL,
		Client: &http.Client{},
	}
}

// Fetch the data needed to process the request
func (a *GetToken) Fetch(ctx context.Context) error {
	return nil
}

// Execute the request
func (a *GetToken) Execute(ctx context.Context) error {
	// Construct the request
	data := url.Values{}
	data.Set("client_id", a.ClientID)
	data.Set("client_secret", a.ClientSecret)
	data.Set("code", a.AuthCode)
	data.Set("grant_type", a.GrantType)
	data.Set("redirect_uri", a.RedirectURL)

	req, err := http.NewRequest(http.MethodPost, constant.URLSpotifyToken, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("GetToken - could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("GetToken - request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("GetToken - failed to read response body: %w", err)
	}
	var jsonResp model.Token
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return fmt.Errorf("GetToken - failed to unmarshal response body: %w", err)
	}

	a.Response = jsonResp
	return nil
}

// Save the output to the database
func (a *GetToken) Save(ctx context.Context) error {
    // Construct the request
    reqBody := model.MapCreateTokenRequest(a.AuthID, a.Response)
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("GetToken - failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, format.Url(constant.URLAPIToken), bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("GetToken - could not create post request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	_, err = a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("GetToken - post request failed: %w", err)
	}

	return nil
}

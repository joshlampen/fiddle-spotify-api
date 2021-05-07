package action

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/format"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/logger"
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
    logger := logger.NewLogger()

	// Construct the request
	data := url.Values{}
	data.Set("client_id", a.ClientID)
	data.Set("client_secret", a.ClientSecret)
	data.Set("code", a.AuthCode)
	data.Set("grant_type", a.GrantType)
	data.Set("redirect_uri", a.RedirectURL)

	req, err := http.NewRequest(http.MethodPost, constant.URLSpotifyToken, strings.NewReader(data.Encode()))
	if err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - failed to create get token request")
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - get token request failed")
		return err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - failed to read get token response body")
		return err
	}
	var jsonResp model.Token
	if err := json.Unmarshal(body, &jsonResp); err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - failed to unmarshal get token response body")
		return err
	}

	a.Response = jsonResp
	return nil
}

// Save the output to the database
func (a *GetToken) Save(ctx context.Context) error {
    logger := logger.NewLogger()

    // Construct the request
    reqBody := model.MapCreateTokenRequest(a.AuthID, a.Response)
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - failed to marshal post token request body")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, format.Url(constant.URLAPIToken), bytes.NewBuffer(jsonBody))
	if err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - failed to create post token request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	_, err = a.Client.Do(req)
	if err != nil {
        logger.Error().
            Err(err).
            Str("authID", a.AuthID).
            Msg("action.GetToken - post token request failed")
		return err
	}

	return nil
}

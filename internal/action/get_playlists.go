package action

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/format"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/logger"
)

// GetPlaylists is an action for getting a user's profile from Spotify
type GetPlaylists struct {
    // Inputs
    AuthID string
	UserID string
	SpotifyUserID string
	Client *http.Client

    // Fetched resources
	Token string

    // Outputs
	SpotifyResponse model.SpotifyPlaylists // response from Spotify
	DBResponse model.DBPlaylists // response from db
}

// NewGetPlaylists constructs and returns a GetPlaylists action
func NewGetPlaylists(authID, userID, spotifyUserID string) GetPlaylists {
	return GetPlaylists{
        AuthID: authID,
		UserID: userID,
		SpotifyUserID: spotifyUserID,
		Client: &http.Client{},
	}
}

// Fetch the data needed to process the request
func (a *GetPlaylists) Fetch(ctx context.Context) error {
    logger := logger.NewLogger()

    // Construct request to get access token
	req, err := http.NewRequest(http.MethodGet, format.Url(constant.URLAPIToken), nil)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to create get token request")
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
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - get token request failed")
		return err
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to read get token response")
		return err
	}

	var token model.Token
	if err := json.Unmarshal(body, &token); err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to unmarshal get token response")
		return err
	}

	a.Token = token.AccessToken
	return nil
}

// Execute the request
func (a *GetPlaylists) Execute(ctx context.Context) error {
    logger := logger.NewLogger()

	// Construct the request
	req, err := http.NewRequest(http.MethodGet, constant.URLSpotifyUserPlaylists, nil)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - could not create get playlists request")
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + a.Token)

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - get playlists request failed")
		return err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to read get playlists response body")
		return err
	}

	var jsonResp model.SpotifyPlaylists
	if err := json.Unmarshal(body, &jsonResp); err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to unmarshal get playlists response body")
		return err
	}

	a.SpotifyResponse = jsonResp
	return nil
}

// Save the output to the database
func (a *GetPlaylists) Save(ctx context.Context) error {
    logger := logger.NewLogger()

	// Construct the request
	reqBody := model.MapCreatePlaylistsRequest(a.UserID, a.SpotifyUserID, a.SpotifyResponse.Items)
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to marshal post playlists request body")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, format.Url(constant.URLAPIPlaylists), bytes.NewBuffer(jsonBody))
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to create post playlists request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - post playlists request failed")
		return err
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to read post playlists response")
		return err
	}

	var playlists model.DBPlaylists
	if err := json.Unmarshal(body, &playlists); err != nil {
        logger.Error().
            Err(err).
            Str("userID", a.UserID).
            Msg("action.GetPlaylists - failed to unmarshal post playlists response")
		return err
	}

	a.DBResponse = playlists
	return nil
}

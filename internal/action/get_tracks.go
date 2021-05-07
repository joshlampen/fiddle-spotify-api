package action

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/format"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/logger"
)

// GetTracks is an action for getting a playlist's tracks from Spotify
type GetTracks struct {
    // Inputs
    AuthID string
	PlaylistID string
	SpotifyPlaylistID string
	TotalTracks int
	Client *http.Client

    // Fetched resources
	Token string

    // Outputs
	SpotifyTracksResponse model.SpotifyPlaylistTracks
    SpotifyArtistsResponse map[string]model.SpotifyArtist // converted to map
    DBResponse model.DBTracks
}

// NewGetTracks constructs and returns a GetTracks action
func NewGetTracks(authID, playlistID, spotifyPlaylistID string, totalTracks int) GetTracks {
	return GetTracks{
        AuthID: authID,
		PlaylistID: playlistID,
		SpotifyPlaylistID: spotifyPlaylistID,
		TotalTracks: totalTracks,
		Client: &http.Client{},
	}
}

// Fetch the data needed to process the request
func (a *GetTracks) Fetch(ctx context.Context) error {
    logger := logger.NewLogger()

    // Construct request to get access token
	req, err := http.NewRequest(http.MethodGet, format.Url(constant.URLAPIToken), nil)
	if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to create get token request")
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
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - get token request failed")
		return err
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to read get token response")
		return err
	}

	var token model.Token
	if err := json.Unmarshal(body, &token); err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to unmarshal get token response")
		return err
	}

	a.Token = token.AccessToken
	return nil
}

// Execute the request
func (a *GetTracks) Execute(ctx context.Context) error {
    logger := logger.NewLogger()

    if err := getTracks(a); err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to run getTracks")
		return err
    }
    // if err := getTrackArtists(a); err != nil {
    //     return fmt.Errorf("GetTracks - failed to run getTrackArtists")
    // }

	return nil
}

// Save the output to the database
func (a *GetTracks) Save(ctx context.Context) error {
    logger := logger.NewLogger()

	// Construct the request
	reqBody, err := model.MapCreateTracksRequest(a.PlaylistID, a.SpotifyTracksResponse.Items, a.SpotifyArtistsResponse)
    if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to map post tracks request")
		return err
    }

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to marshal post tracks request body")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, format.Url(constant.URLAPITracks), bytes.NewBuffer(jsonBody))
	if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to create post tracks request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := a.Client.Do(req)
	if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - post tracks request failed")
		return err
	}

    // Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to read post tracks response")
		return err
	}

    var tracks model.DBTracks
	if err := json.Unmarshal(body, &tracks); err != nil {
        logger.Error().
            Err(err).
            Str("playlistID", a.PlaylistID).
            Msg("action.GetTracks - failed to unmarshal post tracks response")
		return err
	}

	a.DBResponse = tracks
	return nil
}

func getTracks(a *GetTracks) error {
    // Spotify accepts max 50 artists at a time
	// Construct requests until all artists have been requested
	remainingTracks := a.TotalTracks
	offset := 0
	for remainingTracks > 0 {
		var limit int
		if remainingTracks > 100 {
			limit = 100
		} else {
			limit = remainingTracks
		}
		// Construct the request
		req, err := http.NewRequest(http.MethodGet, constant.URLSpotifyPlaylists + "/" + a.SpotifyPlaylistID + "/tracks", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer " + a.Token)

		q := req.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		q.Add("offset", strconv.Itoa(offset))
		req.URL.RawQuery = q.Encode()

		// Do the request
		resp, err := a.Client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var jsonResp model.SpotifyPlaylistTracks
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			return err
		}

		if len(a.SpotifyTracksResponse.Items) == 0 {
			a.SpotifyTracksResponse = jsonResp
		} else {
			a.SpotifyTracksResponse.Items = append(a.SpotifyTracksResponse.Items, jsonResp.Items...)
		}

		offset += limit
		remainingTracks -= limit
	}

    return nil
}

// func getTrackArtists(a *GetTracks) error {
    // // Spotify accepts max 50 artists at a time
	// // Construct requests until all artists have been requested
	// artistIDs := model.MapGetArtistsRequest(a.SpotifyTracksResponse.Items)
	// for len(artistIDs) > 0 {
	// 	var reqIDs []string
	// 	if len(artistIDs) > 50 {
	// 		reqIDs = artistIDs[0:50]
	// 		artistIDs = artistIDs[50:]
	// 		} else {
	// 		reqIDs = artistIDs
	// 		artistIDs = []string{}
	// 	}

	// 	// Construct the request
	// 	req, err := http.NewRequest(http.MethodGet, constant.URLSpotifyArtists, nil)
	// 	if err != nil {
	// 		return fmt.Errorf("GetArtists - could not create get artists request: %w", err)
	// 	}
	// 	req.Header.Set("Accept", "application/json")
	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", "Bearer " + a.Token)

	// 	q := req.URL.Query()
	// 	q.Add("ids", strings.Join(reqIDs, ","))
	// 	req.URL.RawQuery = q.Encode()

	// 	// Do the request
	// 	resp, err := a.Client.Do(req)
	// 	if err != nil {
	// 		return fmt.Errorf("GetArtists - get artists request failed: %w", err)
	// 	}
	// 	defer resp.Body.Close()

	// 	// Read the response
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return fmt.Errorf("GetArtists - failed to read get artists response body: %w", err)
	// 	}

	// 	var jsonResp model.SpotifyArtists
	// 	if err := json.Unmarshal(body, &jsonResp); err != nil {
	// 		return fmt.Errorf("GetArtists - failed to unmarshal get artists response body: %w", err)
	// 	}

	// 	if a.SpotifyArtistsResponse == nil {
	// 		a.SpotifyArtistsResponse = make(map[string]model.SpotifyArtist)
	// 	}
	// 	for _, artist := range jsonResp.Artists {
	// 		a.SpotifyArtistsResponse[artist.ID] = artist
	// 	}
	// }

	// return nil
// }

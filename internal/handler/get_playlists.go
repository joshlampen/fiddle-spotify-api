package handler

import (
	"fmt"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/action"
	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	actionRunner "github.com/JoshLampen/fiddle/spotify-api/internal/utils/action"
	jsonWriter "github.com/JoshLampen/fiddle/spotify-api/internal/utils/json"
)

// GetPlaylists is an HTTP handler for getting a user's playlists from Spotify
func GetPlaylists(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")

    // Get query params
    authID := r.URL.Query().Get(constant.URLParamAuthID)
    userID := r.URL.Query().Get(constant.URLParamUserID)
    spotifyUserID := r.URL.Query().Get(constant.URLParamSpotifyUserID)

	playlists := action.NewGetPlaylists(authID, userID, spotifyUserID)
	if err := actionRunner.Run(r.Context(), &playlists); err != nil {
		jsonWriter.WriteError(
            w,
            fmt.Errorf("Failed to get playlists from Spotify: %w", err),
            http.StatusInternalServerError,
        )
		return
	}

    jsonWriter.WriteResponse(w, playlists.DBResponse)
}

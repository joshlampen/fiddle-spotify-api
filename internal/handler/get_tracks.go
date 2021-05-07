package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/action"
	"github.com/JoshLampen/fiddle/spotify-api/internal/model"
	actionRunner "github.com/JoshLampen/fiddle/spotify-api/internal/utils/action"
	jsonWriter "github.com/JoshLampen/fiddle/spotify-api/internal/utils/json"
)

// GetTracks is an HTTP handler for getting a playlist's tracks from Spotify
func GetTracks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")

    // Read the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        jsonWriter.WriteError(
            w,
            fmt.Errorf("Failed to read get tracks request: %w", err),
            http.StatusInternalServerError,
        )
		return
	}

	var playlists model.DBPlaylists
	if err := json.Unmarshal(body, &playlists); err != nil {
        fmt.Println("handler.GetTracks - failed to unmarshal request body:", err)
		return
	}

    var respBody model.DBTracks
	for _, playlist := range playlists.Items {
		tracks := action.NewGetTracks(
            playlists.AuthID,
			playlist.ID,
			playlist.SpotifyID,
			playlist.TotalTracks,
		)
		if err := actionRunner.Run(r.Context(), &tracks); err != nil {
            jsonWriter.WriteError(
                w,
                fmt.Errorf("Failed to get tracks from Spotify: %w", err),
                http.StatusInternalServerError,
            )
			return
		}

        respBody.Items = append(respBody.Items, tracks.DBResponse.Items...)
	}

    jsonWriter.WriteResponse(w, respBody)
}

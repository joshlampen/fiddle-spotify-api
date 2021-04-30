package handler

import (
	"fmt"
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/action"
	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	actionRunner "github.com/JoshLampen/fiddle/spotify-api/internal/utils/action"
)

// PlayTrack sends a play request to the Spotify player
func PlayTrack(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")

    // Get auth ID from url
    authID := r.URL.Query().Get(constant.URLParamAuthID)
    deviceID := r.URL.Query().Get(constant.URLParamDeviceID)
    spotifyURI := r.URL.Query().Get(constant.URLParamSpotifyURI)

	playTrack := action.NewPlayTrack(authID, deviceID, spotifyURI)
	if err := actionRunner.Run(r.Context(), &playTrack); err != nil {
		fmt.Println("handler.PlayTrack - failed to execute PlayTrack action:", err)
		return
	}

    // // Send a response
	// jsonBody, err := json.Marshal(user.DBResponse)
	// if err != nil {
	// 	fmt.Println("handler.GetUser - failed to marshal response body:", err)
	// 	return
	// }
	// w.Write(jsonBody)
}

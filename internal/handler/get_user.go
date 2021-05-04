package handler

import (
	"net/http"

	"github.com/JoshLampen/fiddle/spotify-api/internal/action"
	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	actionRunner "github.com/JoshLampen/fiddle/spotify-api/internal/utils/action"
	json "github.com/JoshLampen/fiddle/spotify-api/internal/utils/json"
)

// GetUser is an HTTP handler for getting a user's information from Spotify
func GetUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")

    // Get auth ID from url
    authID := r.URL.Query().Get(constant.URLParamAuthID)

	user := action.NewGetUser(authID)
	if err := actionRunner.Run(r.Context(), &user); err != nil {
		json.WriteError(w, err, http.StatusInternalServerError)
		return
	}

    json.WriteResponse(w, user.DBResponse)
}

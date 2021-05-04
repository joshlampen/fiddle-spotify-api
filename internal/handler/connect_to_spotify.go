package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	jsonWriter "github.com/JoshLampen/fiddle/spotify-api/internal/utils/json"
)

var authID string

// ConnectToSpotify is an HTTP handler for connecting to Spotify
func ConnectToSpotify(w http.ResponseWriter, r *http.Request) {
	// Get environment variables
	clientID := os.Getenv(constant.EnvVarClientID)
	redirectUrl := os.Getenv(constant.EnvVarRedirectURL)

    // Get auth ID from url
	authID = r.URL.Query().Get(constant.URLParamAuthID)

	spotifyAuthConfig := &oauth2.Config{
		ClientID: clientID,
		RedirectURL: redirectUrl,
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"playlist-read-private",
            "streaming",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL: constant.URLSpotifyAuth,
		},
	}

	// Create state token to protect against cross-site request forgery
	uuid, err := uuid.NewRandom()
	if err != nil {
        err := fmt.Errorf("Failed to generate uuid: %w", err)
		jsonWriter.WriteError(w, err, http.StatusInternalServerError)
	}
	csrf := uuid.String()

    url := spotifyAuthConfig.AuthCodeURL(csrf)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

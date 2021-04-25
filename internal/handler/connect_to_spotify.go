package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
)

var authID string

// ConnectToSpotify is an HTTP handler for connecting to Spotify
func ConnectToSpotify(w http.ResponseWriter, r *http.Request) {
	// Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		// If port is not defined, load local env file
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
			fmt.Println("handler.ConnectToSpotify - failed to load .env file:", err)
		}
	}
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
		fmt.Println("handler.ConnectToSpotify - failed to generate uuid:", err)
	}
	csrf := uuid.String()

    url := spotifyAuthConfig.AuthCodeURL(csrf)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

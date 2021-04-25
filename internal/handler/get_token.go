package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/spotify-api/internal/action"
	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	actionRunner "github.com/JoshLampen/fiddle/spotify-api/internal/utils/action"
)

const grantTypeAuthCode = "authorization_code"

const popupWindowSuccessHTML = `
	<html>
		<head>
			<script type="text/javascript">
				function close_popup(){window.close();}
			</script>
		</head>
		<body onLoad="setTimeout('close_popup()', 0)">
		</body>
	</html>`

// GetToken is an HTTP handler for getting a user's access token from Spotify
func GetToken(w http.ResponseWriter, r *http.Request) {
    // Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		// If port is not defined, load local env file
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
			fmt.Println("handler.ConnectToSpotify - failed to load .env file:", err)
		}
	}
	clientID := os.Getenv(constant.EnvVarClientID)
	clientSecret := os.Getenv(constant.EnvVarClientSecret)
	redirectUrl := os.Getenv(constant.EnvVarRedirectURL)

    // Get auth code from url
    authCode := r.URL.Query().Get(constant.URLParamCode)

    // Get auth ID from global constant (from connect_to_spotify.go)
    authID := authID

	// Get access token
	token := action.NewGetToken(
        authID,
        authCode,
		clientID,
		clientSecret,
		grantTypeAuthCode,
		redirectUrl,
	)
	if err := actionRunner.Run(r.Context(), &token); err != nil {
		fmt.Println("handler.GetToken - failed to execute GetToken action:", err)
		return
	}

    if _, err := fmt.Fprint(w, popupWindowSuccessHTML); err != nil {
		fmt.Println("handler.GetAuthCode - failed to write html:", err)
	}
}

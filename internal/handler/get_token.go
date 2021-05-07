package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JoshLampen/fiddle/spotify-api/internal/action"
	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	actionRunner "github.com/JoshLampen/fiddle/spotify-api/internal/utils/action"
	jsonWriter "github.com/JoshLampen/fiddle/spotify-api/internal/utils/json"
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
        jsonWriter.WriteError(
            w,
            fmt.Errorf("Failed to get token from Spotify: %w", err),
            http.StatusInternalServerError,
        )
		return
	}

    if _, err := fmt.Fprint(w, popupWindowSuccessHTML); err != nil {
        jsonWriter.WriteError(
            w,
            fmt.Errorf("Failed to write to popup window: %w", err),
            http.StatusInternalServerError,
        )
		return
	}
}

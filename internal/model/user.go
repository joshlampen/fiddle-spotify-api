package model

// SpotifyUser is the response body from requesting a user from Spotify
type SpotifyUser struct {
	ID string `json:"id"`
	DisplayName string `json:"display_name"`
	Email string `json:"email"`
	ExternalURLs SpotifyUserExternalURL `json:"external_urls"`
	Images []SpotifyUserImage `json:"images"`
}

// SpotifyUserImage is the image provided for the user
type SpotifyUserImage struct {
	URL string `json:"url"`
}

// SpotifyUserExternalURL contains the url for the user's Spotify account
type SpotifyUserExternalURL struct {
	Spotify string `json:"spotify"`
}

// DBUser is the response body from a DBCreateUserRequest
type DBUser struct {
	ID string `json:"id"`
	DisplayName string `json:"display_name"`
	Email string `json:"email"`
	SpotifyURL string `json:"spotify_url"`
	SpotifyImageURL string `json:"spotify_image_url"`
	SpotifyID string `json:"spotify_id"`
    Token string `json:"token"`
}

// MapCreateUserRequest maps a Spotify user response to a core API user request
func MapCreateUserRequest(token string, profile SpotifyUser) DBUser {
	var request DBUser
	request.DisplayName = profile.DisplayName
	request.Email = profile.Email
	request.SpotifyURL = profile.ExternalURLs.Spotify
	request.SpotifyImageURL = profile.Images[0].URL
	request.SpotifyID = profile.ID
    request.Token = token

	return request
}

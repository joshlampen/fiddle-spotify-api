package model

// Token is the request body for creating an access token in the database
type Token struct {
	AuthID       string `json:"auth_id"`
    AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
    CreatedAt    string `json:"created_at"`
}

// MapCreateTokenRequest maps a user's auth ID to a Spotify token
func MapCreateTokenRequest(authID string, token Token) Token {
	var request Token
    request.AuthID = authID
    request.AccessToken = token.AccessToken
    request.TokenType = token.TokenType
    request.Scope = token.Scope
    request.ExpiresIn = token.ExpiresIn
    request.RefreshToken = token.RefreshToken

	return request
}

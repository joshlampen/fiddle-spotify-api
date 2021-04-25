package model

// SpotifyPlaylists is the response body from requesting a user's playlists from Spotify
type SpotifyPlaylists struct {
	Items []SpotifyPlaylist `json:"items"`
}

// SpotifyPlaylist is a user's Spotify playlist
type SpotifyPlaylist struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ExternalURLs SpotifyPlaylistExternalURL `json:"external_urls"`
	Owner SpotifyPlaylistOwner `json:"owner"`
	Tracks SpotifyPlaylistTrackInfo `json:"tracks"`
}

// SpotifyPlaylistExternalURL contains the url for the user's playlist on Spotify
type SpotifyPlaylistExternalURL struct {
	Spotify string `json:"spotify"`
}

// SpotifyPlaylistOwner contains the Spotify user ID of the playlist owner
type SpotifyPlaylistOwner struct {
	ID string `json:"id"`
}

// SpotifyPlaylistTrackInfo contains the number of tracks for a playlist
type SpotifyPlaylistTrackInfo struct {
	Total int `json:"total"`
}

// DBPlaylists is the request body for creating a user's playlists in the database
type DBPlaylists struct {
    AuthID string `json:"auth_id"`
	UserID string `json:"user_id"`
	Items []DBPlaylist `json:"items"`
}

// DBPlaylist is a playlist within the database
type DBPlaylist struct {
	ID string `json:"id"`
	Name string `json:"name"`
	SpotifyURL string `json:"spotify_url"`
	SpotifyID string `json:"spotify_id"`
	TotalTracks int `json:"total_tracks"`
}

// MapCreatePlaylistsRequest maps a Spotify user playlists response to a core API user playlists request
func MapCreatePlaylistsRequest(userID, spotifyUserID string, playlists []SpotifyPlaylist) DBPlaylists {
	var request DBPlaylists
	request.UserID = userID

	for _, playlist := range playlists {
		// If the owner of the playlist is not the user we are fetching, ignore it
		if playlist.Owner.ID != spotifyUserID {
			continue
		}

		var pr DBPlaylist
		pr.Name = playlist.Name
		pr.SpotifyURL = playlist.ExternalURLs.Spotify
		pr.SpotifyID = playlist.ID
		pr.TotalTracks = playlist.Tracks.Total

		request.Items = append(request.Items, pr)
	}

	return request
}

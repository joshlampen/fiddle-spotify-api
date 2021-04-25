package model

import (
	"github.com/jmoiron/sqlx/types"
)

// SpotifyPlaylistTracks is the response body from requesting a playlist's tracks from Spotify
type SpotifyPlaylistTracks struct {
	Items []SpotifyPlaylistItem `json:"items"`
}

// SpotifyPlaylistItem is an item within a Spotify playlist
type SpotifyPlaylistItem struct {
	Track SpotifyTrack `json:"track"`
    AddedAt string `json:"added_at"`
}

// SpotifyTrack is a track within a SpotifyPlaylistItem
type SpotifyTrack struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Popularity int `json:"popularity"`
    Duration int `json:"duration_ms"`
    URI string `json:"uri"`
	ExternalURLs SpotifyTrackExternalURL `json:"external_urls"`
    Artists []SpotifyArtist `json:"artists"`
    Album SpotifyAlbum `json:"album"`

	// ArtistIDs []SpotifyArtistID `json:"artists"`
}

// // SpotifyArtistID is an artist ID assocaited with a SpotifyTrack
// type SpotifyArtistID struct {
// 	ID string `json:"id"`
// }

// SpotifyTrackExternalURL contains the url for the track on Spotify
type SpotifyTrackExternalURL struct {
	Spotify string `json:"spotify"`
}

// DBTracks is the request body for creating a user's tracks in the database
type DBTracks struct {
    AuthID string `json:"auth_id"`
	Items []DBTrack `json:"items"`
}

// DBTrack is a track within CreateTracksRequest
type DBTrack struct {
    ID string `json:"id"`
	Name string `json:"name"`
	Popularity int `json:"popularity"`
    Duration int `json:"duration_ms"`
    AddedAt string `json:"added_at"`
    SpotifyURI string `json:"spotify_uri"`
	SpotifyURL string `json:"spotify_url"`
	SpotifyID string `json:"spotify_id"`
    Artists types.JSONText `json:"artists"`
    Album types.JSONText `json:"album"`
	PlaylistID string `json:"playlist_id"`
}

// MapCreateTracksRequest maps a Spotify playlist tracks response to a core API tracks request
func MapCreateTracksRequest(playlistID string, items []SpotifyPlaylistItem, artists map[string]SpotifyArtist) (DBTracks, error) {
	var request DBTracks

	for _, item := range items {
		var reqItem DBTrack
		reqItem.Name = item.Track.Name
		reqItem.Popularity = item.Track.Popularity
        reqItem.Duration = item.Track.Duration
        reqItem.AddedAt = item.AddedAt
        reqItem.SpotifyURI = item.Track.URI
		reqItem.SpotifyURL = item.Track.ExternalURLs.Spotify
		reqItem.SpotifyID = item.Track.ID
		reqItem.PlaylistID = playlistID

        mappedArtists, err := mapArtists(item.Track.Artists)
        if err != nil {
            return DBTracks{}, err
        }
        reqItem.Artists = mappedArtists

        mappedAlbum, err := mapAlbum(item.Track.Album)
        if err != nil {
            return DBTracks{}, err
        }
        reqItem.Album = mappedAlbum

        // mappedArtists, err := mapArtists(item.Track.ArtistIDs, artists)
        // if err != nil {
        //     return DBTracks{}, err
        // }
        // reqItem.Artists = mappedArtists

		request.Items = append(request.Items, reqItem)
	}

	return request, nil
}

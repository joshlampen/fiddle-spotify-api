package model

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

// SpotifyAlbum is the album for a SpotifyPlaylistItem
type SpotifyAlbum struct {
    Name         string                  `json:"name"`
    ExternalURLs SpotifyAlbumExternalURL `json:"external_urls"`
    Images       []SpotifyAlbumImage     `json:"images"`
}

// SpotifyAlbumExternalURL contains the url for the album on Spotify
type SpotifyAlbumExternalURL struct {
	Spotify string `json:"spotify"`
}

// SpotifyAlbumImage is the image for a SpotifyAlbum
type SpotifyAlbumImage struct {
    URL string `json:"url"`
}

// DBAlbum is a representation of an album in the database
type DBAlbum struct {
    Name            string `json:"name"`
    SpotifyURL      string `json:"spotify_url"`
    SpotifyImageURL string `json:"spotify_image_url"`
}

func mapAlbum(album SpotifyAlbum) (types.JSONText, error) {
    var ma DBAlbum

    ma.Name = album.Name
    ma.SpotifyURL = album.ExternalURLs.Spotify
    ma.SpotifyImageURL = album.Images[0].URL

    jsonBytes, err := json.Marshal(ma)
    if err != nil {
        return nil, err
    }

    return jsonBytes, nil
}

package model

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

// SpotifyArtists is the response body from requesting artists from Spotify
type SpotifyArtists struct {
	Artists []SpotifyArtist `json:"artists"`
}

// SpotifyArtist is an artist from Spotify
type SpotifyArtist struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
    ExternalURLs SpotifyArtistExternalURL `json:"external_urls"`

	// Genres []string `json:"genres"`
}

// SpotifyArtistExternalURL contains the url for the artist on Spotify
type SpotifyArtistExternalURL struct {
	Spotify string `json:"spotify"`
}

// DBArtist is a representation of an artist in the database
type DBArtist struct {
    Name       string `json:"name"`
    SpotifyURL string `json:"spotify_url"`
    // SpotifyID string `json:"spotify_id"`
}

// // MapGetArtistsRequest maps a GetArtistsRequest
// func MapGetArtistsRequest(items []SpotifyPlaylistItem) []string {
//     // Insert the artist IDs from all tracks into a single slice
//     var trackArtistIDs []SpotifyArtistID
//     for _, item := range items {
//         trackArtistIDs = append(trackArtistIDs, item.Track.ArtistIDs...)
//     }

//     // Check for duplicates before returning
// 	var artistIDs []string
// 	for _, trackArtistID := range trackArtistIDs {
// 		if containsID(artistIDs, trackArtistID.ID) {
// 			continue
// 		}

// 		artistIDs = append(artistIDs, trackArtistID.ID)
// 	}

// 	return artistIDs
// }

// func containsID(slice []string, id string) bool {
//     for _, item := range slice {
//         if item == id {
//             return true
//         }
//     }
//     return false
// }

// func mapArtists(artistIDs []SpotifyArtistID, artists map[string]SpotifyArtist) (types.JSONText, error) {
// 	var mappedArtists []SpotifyArtist

// 	for _, artistID := range artistIDs {
// 		artist := artists[artistID.ID]
// 		mappedArtists = append(mappedArtists, artist)
// 	}

//     jsonBytes, err := json.Marshal(mappedArtists)
//     if err != nil {
//         return nil, err
//     }

// 	return jsonBytes, nil
// }

func mapArtists(artists []SpotifyArtist) (types.JSONText, error) {
    var mappedArtists []DBArtist

    for _, artist := range artists {
        var ma DBArtist
        ma.Name = artist.Name
        ma.SpotifyURL = artist.ExternalURLs.Spotify
        // ma.SpotifyID = artist.ID

        mappedArtists = append(mappedArtists, ma)
    }

    jsonBytes, err := json.Marshal(mappedArtists)
    if err != nil {
        return nil, err
    }

    return jsonBytes, nil
}

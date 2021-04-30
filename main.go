package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/JoshLampen/fiddle/spotify-api/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // localhost
	}

	r := mux.NewRouter()

	r.HandleFunc("/authorize", handler.ConnectToSpotify).Methods("GET", "OPTIONS")
	r.HandleFunc("/authorize/callback", handler.GetToken).Methods("GET", "OPTIONS")
	r.HandleFunc("/users", handler.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/playlists", handler.GetPlaylists).Methods("GET", "OPTIONS")
	r.HandleFunc("/tracks/get", handler.GetTracks).Methods("PUT", "OPTIONS")
	r.HandleFunc("/player", handler.PlayTrack).Methods("GET", "OPTIONS")

	http.ListenAndServe(":" + port, r)
}

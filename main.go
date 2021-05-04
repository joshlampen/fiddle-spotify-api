package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // localhost
        // If port is not defined, load local env file
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
			panic(fmt.Errorf("failed to load .env file: %w", err))
		}
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

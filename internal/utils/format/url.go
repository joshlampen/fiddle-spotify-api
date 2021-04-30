package format

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
)

func Url(endpoint string) (string) {
    // Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		// If port is not defined, load local env file
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
			fmt.Println("handler.ConnectToSpotify - failed to load .env file:", err)
		}
	}
	baseUrl := os.Getenv(constant.EnvVarAPIBaseURL)

    return baseUrl + endpoint
}

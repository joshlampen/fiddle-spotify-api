package format

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/JoshLampen/fiddle/spotify-api/internal/constant"
	"github.com/JoshLampen/fiddle/spotify-api/internal/utils/logger"
)

func Url(endpoint string) (string) {
    logger := logger.NewLogger()

    // Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		// If port is not defined, load local env file
		if err := godotenv.Load(constant.DotEnvFilePath); err != nil {
            logger.Error().Err(err).Msg("handler.ConnectToSpotify - failed to load .env file")
		}
	}
	baseUrl := os.Getenv(constant.EnvVarAPIBaseURL)

    return baseUrl + endpoint
}
